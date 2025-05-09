package req

import (
	networkingv1 "k8s.io/api/networking/v1"
)

type Ingress struct {
	Name      string        `json:"name"`
	Namespace string        `json:"namespace"`
	Labels    []Item        `json:"labels"`
	Rules     []IngressRule `json:"rules"`
}

type IngressRule struct {
	Host  string                        `json:"host"`
	Value networkingv1.IngressRuleValue `json:"value"`
}
