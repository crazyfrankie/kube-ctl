package resp

type PodListItem struct {
	Name     string `json:"name"`
	Ready    string `json:"ready"`    // 0/1 | 1/1
	Status   string `json:"status"`   // Running | Error
	Restarts int32  `json:"restarts"` // Number of restarts
	Age      int64  `json:"age"`      // Runtime
	IP       string `json:"ip"`       // Pod id
	Node     string `json:"node"`     // Which Node the Pod is dispatched to
}
