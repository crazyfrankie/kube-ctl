package req

type DaemonSet struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Labels    []Item `json:"labels"`
	Selector  []Item `json:"selector"`
	Template  Pod    `json:"template"`
}
