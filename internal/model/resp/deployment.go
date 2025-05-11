package resp

type Deployment struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Replicas  int32  `json:"replicas"`
	Ready     int32  `json:"ready"`     // Ready 字段表示 Deployment 中正在运行的 Pod 副本的数量
	UpToDate  int32  `json:"upToDate"`  // UpToDate 字段表示与 Deployment 所期望的副本数相比，有多少个 Pod 副本是最新的
	Available int32  `json:"available"` // Available 字段表示 Deployment 中可用的 Pod 副本数
	Age       int64  `json:"age"`
}
