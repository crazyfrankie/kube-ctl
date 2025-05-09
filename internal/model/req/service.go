package req

import corev1 "k8s.io/api/core/v1"

type Service struct {
	Name      string             `json:"name"`
	Namespace string             `json:"namespace"`
	Labels    []Item             `json:"labels"`
	Type      corev1.ServiceType `json:"type"`
	Selector  []Item             `json:"selector"`
	Ports     []ServicePort      `json:"ports"`
}

type ServicePort struct {
	Name       string `json:"name"`
	Port       int32  `json:"port"`
	TargetPort int32  `json:"targetPort"`
	NodePort   int32  `json:"nodePort"`
}
