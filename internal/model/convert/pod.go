package convert

import (
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
)

const (
	HTTPProbe = "http"
	TCPProbe  = "tcp"
	EXECProbe = "exec"

	EMPTYDIR = "emptyDir"

	ScheduleNodeName     = "nodeName"
	ScheduleNodeSelector = "nodeSelector"
	ScheduleNodeAffinity = "nodeAffinity"
	ScheduleNodeAny      = "nodeAny"

	RefConfigMap = "configMap"
	RefSecret    = "secret"
)

// PodReqConvert convert req.Pod to corev1.Pod
func PodReqConvert(req *req.Pod) *corev1.Pod {
	// get node scheduling
	affinity, selector, nodeName := getPodNodeScheduling(req.NodeScheduling)

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:              req.Base.Name,
			Namespace:         req.Base.Namespace,
			CreationTimestamp: metav1.Time{Time: time.Now()},
			Labels:            getPodLabels(req.Base.Labels),
		},
		Spec: corev1.PodSpec{
			Tolerations:    req.Tolerations,
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
			NodeName:      nodeName,
			NodeSelector:  selector,
			Affinity:      affinity,
		},
	}
}

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
			Env:             getPodEnvVar(c.Env),
			EnvFrom:         getPodEnvVarFrom(c.EnvsFrom),
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

func getPodEnvVar(items []req.EnvVar) []corev1.EnvVar {
	envs := make([]corev1.EnvVar, 0, len(items))
	for _, i := range items {
		env := corev1.EnvVar{Name: i.Name}
		switch i.Type {
		case RefConfigMap:
			env.ValueFrom = &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: i.RefName},
					Key:                  i.Value,
				},
			}
		case RefSecret:
			env.ValueFrom = &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: i.RefName},
					Key:                  i.Value,
				},
			}
		default:
			env.Value = i.Value
		}
		envs = append(envs)
	}

	return envs
}

func getPodEnvVarFrom(items []req.EnvVarFromResource) []corev1.EnvFromSource {
	envs := make([]corev1.EnvFromSource, 0, len(items))
	for _, i := range items {
		env := corev1.EnvFromSource{Prefix: i.Prefix}
		switch i.RefType {
		case RefConfigMap:
			env.ConfigMapRef = &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: i.Name},
			}
		case RefSecret:
			env.SecretRef = &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: i.Name},
			}
		}
		envs = append(envs, env)
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

func getPodNodeScheduling(sch req.NodeScheduling) (*corev1.Affinity, map[string]string, string) {
	switch sch.Type {
	case ScheduleNodeName:
		return nil, nil, sch.NodeName
	case ScheduleNodeSelector:
		res := make(map[string]string)
		for _, i := range sch.NodeSelector {
			res[i.Key] = i.Value
		}
		return nil, res, ""
	case ScheduleNodeAffinity:
		expr := sch.NodeAffinity
		matchExpr := make([]corev1.NodeSelectorRequirement, 0, len(expr))
		for _, e := range expr {
			matchExpr = append(matchExpr, corev1.NodeSelectorRequirement{
				Key:      e.Key,
				Operator: e.Operator,
				Values:   strings.Split(e.Value, ","),
			})
		}
		res := &corev1.Affinity{
			NodeAffinity: &corev1.NodeAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
					NodeSelectorTerms: []corev1.NodeSelectorTerm{
						{
							MatchExpressions: matchExpr,
						},
					}},
			},
		}
		return res, nil, ""
	case ScheduleNodeAny:
		// do nothing
		return nil, nil, ""
	default:
		// do nothing
		return nil, nil, ""
	}
}

// PodConvertReq corev1.Pod convert to req.Pod
func PodConvertReq(pod *corev1.Pod) *req.Pod {
	volume, volumeMap := getReqVolume(pod.Spec.Volumes)
	return &req.Pod{
		Base:           getReqBase(pod),
		Network:        getReqNetwork(pod),
		Volume:         volume,
		InitContainers: getReqContainers(pod.Spec.InitContainers, volumeMap),
		Containers:     getReqContainers(pod.Spec.Containers, volumeMap),
		Tolerations:    pod.Spec.Tolerations,
		NodeScheduling: getReqPodNodeScheduling(pod),
	}
}

func getReqBase(pod *corev1.Pod) req.Base {
	return req.Base{
		Name:          pod.Name,
		Labels:        getReqLabels(pod.Labels),
		Namespace:     pod.Namespace,
		RestartPolicy: string(pod.Spec.RestartPolicy),
	}
}

func getReqLabels(data map[string]string) []req.Item {
	res := make([]req.Item, 0, len(data))
	for k, v := range data {
		res = append(res, req.Item{
			Key:   k,
			Value: v,
		})
	}

	return res
}

func getReqNetwork(pod *corev1.Pod) req.Network {
	return req.Network{
		HostNetwork: pod.Spec.HostNetwork,
		HostName:    pod.Spec.NodeName,
		DnsPolicy:   string(pod.Spec.DNSPolicy),
		DnsConfig:   getReqDNSConfig(pod.Spec.DNSConfig),
		HostAliases: getReqHostAliases(pod.Spec.HostAliases),
	}
}

func getReqDNSConfig(dns *corev1.PodDNSConfig) req.DnsConfig {
	if dns == nil {
		return req.DnsConfig{}
	} else {
		return req.DnsConfig{Nameservers: dns.Nameservers}
	}
}

func getReqHostAliases(host []corev1.HostAlias) []req.Item {
	res := make([]req.Item, 0, len(host))
	for _, i := range host {
		res = append(res, req.Item{
			Key:   i.IP,
			Value: strings.Join(i.Hostnames, ","),
		})
	}

	return res
}

func getReqVolume(volumes []corev1.Volume) ([]req.Volume, map[string]string) {
	res := make([]req.Volume, 0, len(volumes))
	volumeMap := make(map[string]string)
	for _, v := range volumes {
		if v.EmptyDir == nil {
			continue
		}
		volumeMap[v.Name] = ""
		res = append(res, req.Volume{
			Name: v.Name,
			Type: EMPTYDIR,
		})
	}

	return res, volumeMap
}

func getReqContainers(containers []corev1.Container, volumeMap map[string]string) []req.Container {
	res := make([]req.Container, 0, len(containers))
	for _, c := range containers {
		var privileged bool
		if c.SecurityContext != nil && c.SecurityContext.Privileged != nil {
			privileged = *c.SecurityContext.Privileged
		}
		res = append(res, req.Container{
			Name:            c.Name,
			Image:           c.Image,
			ImagePullPolicy: string(c.ImagePullPolicy),
			Tty:             c.TTY,
			Ports:           getReqContainerPort(c.Ports),
			WorkingDir:      c.WorkingDir,
			Command:         c.Command,
			Args:            c.Args,
			Env:             getReqContainerEnv(c.Env),
			EnvsFrom:        getReqContainerEnvVarFrom(c.EnvFrom),
			Privileged:      privileged,
			Resources:       getReqContainerResource(&c.Resources),
			VolumeMounts:    getReqContainerVolumeMount(c.VolumeMounts, volumeMap),
			StartUpProbe:    getReqContainerProbe(c.StartupProbe),
			LivenessProbe:   getReqContainerProbe(c.LivenessProbe),
			ReadinessProbe:  getReqContainerProbe(c.ReadinessProbe),
		})
	}

	return res
}

func getReqContainerPort(ports []corev1.ContainerPort) []req.ContainerPort {
	res := make([]req.ContainerPort, 0, len(ports))
	for _, i := range ports {
		res = append(res, req.ContainerPort{
			Name:          i.Name,
			ContainerPort: i.ContainerPort,
			HostPort:      i.HostPort,
		})
	}

	return res
}

func getReqContainerEnv(envs []corev1.EnvVar) []req.EnvVar {
	res := make([]req.EnvVar, 0, len(envs))
	for _, i := range envs {
		env := req.EnvVar{Name: i.Name}
		if i.ValueFrom != nil {
			if i.ValueFrom.ConfigMapKeyRef != nil {
				env.RefName = i.ValueFrom.ConfigMapKeyRef.LocalObjectReference.Name
				env.Value = i.ValueFrom.ConfigMapKeyRef.Key
				env.Type = RefConfigMap
			}
			if i.ValueFrom.SecretKeyRef != nil {
				env.RefName = i.ValueFrom.SecretKeyRef.LocalObjectReference.Name
				env.Value = i.ValueFrom.SecretKeyRef.Key
				env.Type = RefSecret
			}
		} else {
			env.Value = i.Value
		}

		res = append(res, env)
	}

	return res
}

func getReqContainerEnvVarFrom(envs []corev1.EnvFromSource) []req.EnvVarFromResource {
	res := make([]req.EnvVarFromResource, 0, len(envs))
	for _, i := range envs {
		env := req.EnvVarFromResource{Prefix: i.Prefix}
		if i.ConfigMapRef != nil {
			env.RefType = RefConfigMap
			env.Name = i.ConfigMapRef.Name
		}
		if i.SecretRef != nil {
			env.RefType = RefSecret
			env.Name = i.SecretRef.Name
		}
		res = append(res, env)
	}

	return res
}

func getReqContainerResource(resource *corev1.ResourceRequirements) req.Resource {
	if resource == nil {
		return req.Resource{}
	}

	return req.Resource{
		Enable:      true,
		MemoryReq:   int32(resource.Requests.Memory().Value()),
		MemoryLimit: int32(resource.Limits.Memory().Value()),
		CPUReq:      int32(resource.Requests.Cpu().Value()),
		CPULimit:    int32(resource.Limits.Cpu().Value()),
	}
}

func getReqContainerVolumeMount(vm []corev1.VolumeMount, volumeMap map[string]string) []req.VolumeMount {
	res := make([]req.VolumeMount, 0, len(vm))
	for _, i := range vm {
		// Filter by non-emptyDir
		if _, ok := volumeMap[i.Name]; ok {
			res = append(res, req.VolumeMount{
				MountName: i.Name,
				MountPath: i.MountPath,
				ReadOnly:  i.ReadOnly,
			})
		}
	}

	return res
}

func getReqContainerProbe(probe *corev1.Probe) req.ContainerProbe {
	if probe != nil {
		if probe.HTTPGet != nil {
			return req.ContainerProbe{
				Enable: true,
				Type:   HTTPProbe,
				HttpGet: req.ProbeHTTPGet{
					Scheme:  string(probe.HTTPGet.Scheme),
					Host:    probe.HTTPGet.Host,
					Path:    probe.HTTPGet.Path,
					Port:    probe.HTTPGet.Port.IntVal,
					Headers: getReqProbeHTTPHeaders(probe.HTTPGet.HTTPHeaders),
				},
				ProbeTime: req.ProbeTime{
					InitialDelaySeconds: probe.InitialDelaySeconds,
					PeriodSeconds:       probe.PeriodSeconds,
					TimeoutSeconds:      probe.TimeoutSeconds,
					SuccessThreshold:    probe.SuccessThreshold,
					FailureThreshold:    probe.FailureThreshold,
				},
			}
		}
		if probe.TCPSocket != nil {
			return req.ContainerProbe{
				Enable: true,
				Type:   TCPProbe,
				TcpSocket: req.ProbeTcpSocket{
					Host: probe.TCPSocket.Host,
					Port: probe.TCPSocket.Port.IntVal,
				},
				ProbeTime: req.ProbeTime{
					InitialDelaySeconds: probe.InitialDelaySeconds,
					PeriodSeconds:       probe.PeriodSeconds,
					TimeoutSeconds:      probe.TimeoutSeconds,
					SuccessThreshold:    probe.SuccessThreshold,
					FailureThreshold:    probe.FailureThreshold,
				},
			}
		}
		if probe.Exec != nil {
			return req.ContainerProbe{
				Enable:  true,
				Type:    EXECProbe,
				Command: req.ProbeCommand{Command: probe.Exec.Command},
				ProbeTime: req.ProbeTime{
					InitialDelaySeconds: probe.InitialDelaySeconds,
					PeriodSeconds:       probe.PeriodSeconds,
					TimeoutSeconds:      probe.TimeoutSeconds,
					SuccessThreshold:    probe.SuccessThreshold,
					FailureThreshold:    probe.FailureThreshold,
				},
			}
		}
	}

	return req.ContainerProbe{}
}

func getReqProbeHTTPHeaders(header []corev1.HTTPHeader) []req.Item {
	res := make([]req.Item, 0, len(header))
	for _, i := range header {
		res = append(res, req.Item{
			Key:   i.Name,
			Value: i.Value,
		})
	}

	return res
}

func getReqPodNodeScheduling(pod *corev1.Pod) req.NodeScheduling {
	scheduling := req.NodeScheduling{Type: ScheduleNodeAny}
	if pod.Spec.NodeSelector != nil {
		scheduling.Type = ScheduleNodeSelector
		res := make([]req.Item, 0, len(pod.Spec.NodeSelector))
		for k, v := range pod.Spec.NodeSelector {
			res = append(res, req.Item{
				Key:   k,
				Value: v,
			})
		}
		scheduling.NodeSelector = res
	}

	if pod.Spec.Affinity != nil && pod.Spec.Affinity.NodeAffinity != nil {
		// Hard affinity scheduling by default
		scheduling.Type = ScheduleNodeAffinity
		term := pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0]
		res := make([]req.NodeAffinityTermExpressions, 0, len(term.MatchExpressions))
		for _, e := range term.MatchExpressions {
			res = append(res, req.NodeAffinityTermExpressions{
				Key:      e.Key,
				Value:    strings.Join(e.Values, ","),
				Operator: e.Operator,
			})
		}
		scheduling.NodeAffinity = res
	}

	if pod.Spec.NodeName != "" {
		scheduling.Type = ScheduleNodeName
		scheduling.NodeName = pod.Spec.NodeName
	}

	return scheduling
}

func PodListConvertResp(pod corev1.Pod) resp.PodListItem {
	var total, ready int
	var restart int32
	for _, c := range pod.Status.ContainerStatuses {
		if c.Ready {
			ready++
		}
		restart += c.RestartCount
		total++
	}

	var podStatus string
	if pod.Status.Phase != "Running" {
		podStatus = "Error"
	} else {
		podStatus = "Running"
	}

	return resp.PodListItem{
		Name:     pod.Name,
		Ready:    fmt.Sprintf("%d/%d", ready, total),
		Status:   podStatus,
		Restarts: restart,
		Age:      pod.CreationTimestamp.Unix(),
		IP:       pod.Status.PodIP,
		Node:     pod.Spec.NodeName,
	}
}
