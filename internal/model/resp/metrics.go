package resp

import (
	"time"
	
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MetricsItem struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Color string `json:"color"`
	Label string `json:"label"`
}

type NodeMetricsList struct {
	Kind       string            `json:"kind"`
	ApiVersion string            `json:"apiVersion"`
	Metadata   metav1.ObjectMeta `json:"metadata"`
	Items      []NodeMetric      `json:"items"`
}

type NodeMetric struct {
	Metadata  metav1.ObjectMeta   `json:"metadata"`
	Timestamp time.Time           `json:"timestamp"`
	Window    string              `json:"window"`
	Usage     corev1.ResourceList `json:"usage"`
}
