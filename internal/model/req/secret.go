package req

import corev1 "k8s.io/api/core/v1"

type Secret struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    []Item            `json:"labels"`
	Data      []Item            `json:"data"`
	Type      corev1.SecretType `json:"type"` // Opaque | kubernetes.io/dockerconfigjson
}
