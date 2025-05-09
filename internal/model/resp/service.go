package resp

import corev1 "k8s.io/api/core/v1"

type Service struct {
	Name       string             `json:"name"`
	Namespace  string             `json:"namespace"`
	Type       corev1.ServiceType `json:"type"`
	ClusterIP  string             `json:"clusterIP"`
	ExternalIP []string           `json:"externalIP"`
	Age        int64              `json:"age"`
}
