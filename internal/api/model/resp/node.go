package resp

import corev1 "k8s.io/api/core/v1"

type NodeListItem struct {
	Name             string `json:"name"`
	Status           string `json:"status"`
	Age              int64  `json:"age"`
	Version          string `json:"version"` // kubelet version
	KernelVersion    string `json:"kernelVersion"`
	InternalIP       string `json:"internalIP"`
	ExternalIP       string `json:"externalIP"`
	OSImage          string `json:"OSImage"`
	ContainerRuntime string `json:"containerRuntime"`
}

type NodeDetail struct {
	Name             string         `json:"name"`
	Status           string         `json:"status"`
	Age              int64          `json:"age"`
	Version          string         `json:"version"` // kubelet version
	KernelVersion    string         `json:"kernelVersion"`
	InternalIP       string         `json:"internalIP"`
	ExternalIP       string         `json:"externalIP"`
	OSImage          string         `json:"OSImage"`
	ContainerRuntime string         `json:"containerRuntime"`
	Labels           []Item         `json:"labels"`
	Taints           []corev1.Taint `json:"taints"`
}

type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
