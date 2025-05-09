package convert

import (
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/pkg/utils"
)

func ServiceReqConvert(req *req.Service) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels:    utils.ReqItemToMap(req.Labels),
		},
		Spec: corev1.ServiceSpec{
			Ports:    getServicePorts(req.Ports),
			Selector: utils.ReqItemToMap(req.Selector),
			Type:     req.Type,
		},
	}
}

func getServicePorts(ports []req.ServicePort) []corev1.ServicePort {
	res := make([]corev1.ServicePort, 0, len(ports))
	for _, p := range ports {
		res = append(res, corev1.ServicePort{
			Name:       p.Name,
			Port:       p.Port,
			TargetPort: intstr.FromInt32(p.TargetPort),
			NodePort:   p.NodePort,
		})
	}

	return res
}

func ServiceConvertReq(svc *corev1.Service) req.Service {
	return req.Service{
		Name:      svc.Name,
		Namespace: svc.Namespace,
		Labels:    utils.ReqMapToItem(svc.Labels),
		Type:      svc.Spec.Type,
		Selector:  utils.ReqMapToItem(svc.Spec.Selector),
		Ports:     getReqServicePorts(svc.Spec.Ports),
	}
}

func getReqServicePorts(ports []corev1.ServicePort) []req.ServicePort {
	res := make([]req.ServicePort, 0, len(ports))
	for _, p := range ports {
		res = append(res, req.ServicePort{
			Name:       p.Name,
			Port:       p.Port,
			TargetPort: p.TargetPort.IntVal,
			NodePort:   p.NodePort,
		})
	}

	return res
}

func ServiceConvertResp(svc *corev1.Service) resp.Service {
	return resp.Service{
		Name:       svc.Name,
		Namespace:  svc.Namespace,
		Type:       svc.Spec.Type,
		ClusterIP:  svc.Spec.ClusterIP,
		ExternalIP: svc.Spec.ExternalIPs,
		Age:        svc.CreationTimestamp.Unix(),
	}
}
