package req

type ConfigMap struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Labels    []Item `json:"labels"`
	Data      []Item `json:"data"`
}
