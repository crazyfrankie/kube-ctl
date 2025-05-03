package resp

import corev1 "k8s.io/api/core/v1"

type Secret struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	DataNum   int               `json:"dataNum"`
	Age       int64             `json:"age"`
	Type      corev1.SecretType `json:"type"`
}

type SecretDetail struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	DataNum   int               `json:"dataNum"`
	Age       int64             `json:"age"`
	Type      corev1.SecretType `json:"type"`
	Labels    []Item            `json:"labels"`
	Data      []Item            `json:"data"`
}
