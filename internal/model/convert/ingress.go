package convert

import (
	"strings"
	
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/pkg/utils"
)

func IngressReqConvert(req *req.Ingress) *networkingv1.Ingress {
	rules := make([]networkingv1.IngressRule, 0, len(req.Rules))
	for _, r := range req.Rules {
		rules = append(rules, networkingv1.IngressRule{
			Host:             r.Host,
			IngressRuleValue: r.Value,
		})
	}
	return &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels:    utils.ReqItemToMap(req.Labels),
		},
		Spec: networkingv1.IngressSpec{
			Rules: rules,
		},
	}
}

func IngressConvertReq(ingress *networkingv1.Ingress) req.Ingress {
	rules := make([]req.IngressRule, 0, len(ingress.Spec.Rules))
	for _, r := range ingress.Spec.Rules {
		rules = append(rules, req.IngressRule{
			Host:  r.Host,
			Value: r.IngressRuleValue,
		})
	}
	return req.Ingress{
		Name:      ingress.Name,
		Namespace: ingress.Namespace,
		Labels:    utils.ReqMapToItem(ingress.Labels),
		Rules:     rules,
	}
}

func IngressConvertResp(ingress *networkingv1.Ingress) resp.Ingress {
	var class string
	if ingress.Spec.IngressClassName != nil {
		class = *ingress.Spec.IngressClassName
	} else {
		class = "<none>"
	}
	var hosts string
	host := make([]string, 0, len(ingress.Spec.Rules))
	for _, r := range ingress.Spec.Rules {
		host = append(host, r.Host)
	}
	hosts = strings.Join(host, ",")
	return resp.Ingress{
		Name:      ingress.Name,
		Namespace: ingress.Namespace,
		Class:     class,
		Hosts:     hosts,
		Age:       ingress.CreationTimestamp.Unix(),
	}
}
