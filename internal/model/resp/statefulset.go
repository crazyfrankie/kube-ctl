package resp

type StatefulSet struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Ready     int32  `json:"ready"`
	Replicas  int32  `json:"replicas"`
	Age       int64  `json:"age"`
}
