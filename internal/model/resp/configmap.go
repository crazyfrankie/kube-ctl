package resp

type ConfigMap struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	DataNum   int    `json:"dataNum"`
	Age       int64  `json:"age"`
}

type ConfigMapDetail struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	DataNum   int    `json:"dataNum"`
	Age       int64  `json:"age"`
	Labels    []Item `json:"labels"`
	Data      []Item `json:"data"`
}
