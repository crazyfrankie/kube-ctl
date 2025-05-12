package req

type Job struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Labels      []Item `json:"labels"`
	Completions int32  `json:"completions"` // Job 的 Pod 副本数，全部副本数运行成功，才能代表job运行成功
	Template    Pod    `json:"template"`
}
