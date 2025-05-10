package req

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type IngressRouteSpec struct {
	EntryPoints []string `json:"entryPoints"`
	Routes      []struct {
		Kind     string `json:"kind"`
		Match    string `json:"match"`
		Services []struct {
			Name string `json:"name"`
			Port int32  `json:"port"`
		} `json:"services"`
		Middlewares []struct {
			Name string `json:"name"`
		} `json:"middlewares"`
	} `json:"routes"`
	Tls *struct {
		SecretName string `json:"secretName"`
	} `json:"tls"`
}

type Middleware struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ObjectMeta `json:"metadata"`
}

type IngressRoute struct {
	Name             string `json:"name"`
	Namespace        string `json:"namespace"`
	Labels           []Item `json:"labels"`
	IngressRouteSpec `json:"ingressRouteSpec"`
}
