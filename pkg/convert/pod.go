package convert

import (
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/crazyfrankie/kube-ctl/api/model/req"
)

const (
	HTTPProbe = "http"
	TCPProbe  = "tcp"
	EXECProbe = "exec"

	EMPTYDIR = "emptyDir"
)

func PodReqConvert(req *req.Pod) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:              req.Base.Name,
			Namespace:         req.Base.Namespace,
			CreationTimestamp: metav1.Time{Time: time.Now()},
			Labels:            getPodLabels(req.Base.Labels),
		},
		Spec: corev1.PodSpec{
			Volumes:        getPodVolumes(req.Volume),
			InitContainers: getPodContainers(req.InitContainers),
			Containers:     getPodContainers(req.Containers),
			HostAliases:    getPodHostAliases(req.Network.HostAliases),
			Hostname:       req.Network.HostName,
			DNSConfig: &corev1.PodDNSConfig{
				Nameservers: req.Network.DnsConfig.Nameservers,
			},
			DNSPolicy:     corev1.DNSPolicy(req.Network.DnsPolicy),
			RestartPolicy: corev1.RestartPolicy(req.Base.RestartPolicy),
		},
	}
}

// getPodLabels converts the labels in the request to labels of type map[string]string as needed by client-go.
func getPodLabels(items []req.Item) map[string]string {
	res := make(map[string]string, len(items))
	for _, i := range items {
		res[i.Key] = i.Value
	}

	return res
}

func getPodVolumes(vms []req.Volume) []corev1.Volume {
	res := make([]corev1.Volume, 0, len(vms))
	for _, i := range vms {
		if i.Type != EMPTYDIR {
			continue
		}
		source := corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		}
		res = append(res, corev1.Volume{
			Name:         i.Name,
			VolumeSource: source,
		})
	}

	return res
}

func getPodContainers(cs []req.Container) []corev1.Container {
	res := make([]corev1.Container, 0, len(cs))

	for _, c := range cs {
		res = append(res, corev1.Container{
			Name:            c.Name,
			Image:           c.Image,
			Command:         c.Command,
			Args:            c.Args,
			WorkingDir:      c.WorkingDir,
			TTY:             c.Tty,
			Ports:           getPodPorts(c.Ports),
			Env:             getPodEnv(c.Env),
			ImagePullPolicy: corev1.PullPolicy(c.ImagePullPolicy),
			SecurityContext: &corev1.SecurityContext{
				Privileged: &c.Privileged,
			},
			Resources:      getPodContainerResource(c.Resources),
			VolumeMounts:   getPodContainerVolumeMounts(c.VolumeMounts),
			StartupProbe:   getPodContainerProbe(c.StartUpProbe),
			LivenessProbe:  getPodContainerProbe(c.LivenessProbe),
			ReadinessProbe: getPodContainerProbe(c.ReadinessProbe),
		})
	}

	return res
}

func getPodPorts(ports []req.ContainerPort) []corev1.ContainerPort {
	res := make([]corev1.ContainerPort, 0, len(ports))
	for _, i := range ports {
		res = append(res, corev1.ContainerPort{
			Name:          i.Name,
			HostPort:      i.HostPort,
			ContainerPort: i.ContainerPort,
		})
	}

	return res
}

func getPodEnv(items []req.Item) []corev1.EnvVar {
	envs := make([]corev1.EnvVar, 0, len(items))
	for _, i := range items {
		envs = append(envs, corev1.EnvVar{
			Name:  i.Key,
			Value: i.Value,
		})
	}

	return envs
}

func getPodContainerProbe(probe req.ContainerProbe) *corev1.Probe {
	if !probe.Enable {
		return nil
	}

	var hdl corev1.ProbeHandler
	switch probe.Type {
	case HTTPProbe:
		header := make([]corev1.HTTPHeader, 0, len(probe.HttpGet.Headers))
		for _, h := range probe.HttpGet.Headers {
			header = append(header, corev1.HTTPHeader{
				Name:  h.Key,
				Value: h.Value,
			})
		}
		hdl.HTTPGet = &corev1.HTTPGetAction{
			Path: probe.HttpGet.Path,
			Port: intstr.IntOrString{
				IntVal: probe.HttpGet.Port,
			},
			Host:        probe.HttpGet.Host,
			Scheme:      corev1.URIScheme(probe.HttpGet.Scheme),
			HTTPHeaders: header,
		}
	case TCPProbe:
		hdl.TCPSocket = &corev1.TCPSocketAction{
			Port: intstr.IntOrString{
				IntVal: probe.TcpSocket.Port,
			},
			Host: probe.TcpSocket.Host,
		}
	case EXECProbe:
		hdl.Exec = &corev1.ExecAction{
			Command: probe.Command.Command,
		}
	}

	return &corev1.Probe{
		ProbeHandler:        hdl,
		InitialDelaySeconds: probe.InitialDelaySeconds,
		TimeoutSeconds:      probe.TimeoutSeconds,
		PeriodSeconds:       probe.PeriodSeconds,
		SuccessThreshold:    probe.SuccessThreshold,
		FailureThreshold:    probe.FailureThreshold,
	}
}

func getPodContainerVolumeMounts(vms []req.VolumeMount) []corev1.VolumeMount {
	res := make([]corev1.VolumeMount, 0, len(vms))
	for _, vm := range vms {
		res = append(res, corev1.VolumeMount{
			Name:      vm.MountName,
			ReadOnly:  vm.ReadOnly,
			MountPath: vm.MountPath,
		})
	}

	return res
}

func getPodContainerResource(rs req.Resource) corev1.ResourceRequirements {
	var res corev1.ResourceRequirements
	if !rs.Enable {
		return res
	}
	res.Limits = corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse(fmt.Sprintf("%dm", rs.CPULimit)),
		corev1.ResourceMemory: resource.MustParse(fmt.Sprintf("%dMi", rs.MemoryLimit)),
	}
	res.Requests = corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse(fmt.Sprintf("%dm", rs.CPUReq)),
		corev1.ResourceMemory: resource.MustParse(fmt.Sprintf("%dMi", rs.MemoryReq)),
	}

	return res
}

func getPodHostAliases(as []req.Item) []corev1.HostAlias {
	res := make([]corev1.HostAlias, 0, len(as))
	for _, i := range as {
		res = append(res, corev1.HostAlias{
			IP:        i.Key,
			Hostnames: strings.Split(i.Value, ","),
		})
	}

	return res
}
