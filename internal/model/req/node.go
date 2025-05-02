package req

import corev1 "k8s.io/api/core/v1"

type UpdateLabelReq struct {
	Name   string `json:"name"`
	Labels []Item `json:"labels"`
}

type UpdateTaintReq struct {
	Name   string         `json:"name"`
	Taints []corev1.Taint `json:"taints"`
}
