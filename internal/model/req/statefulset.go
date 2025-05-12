package req

type StatefulSet struct {
	Name                 string                  `json:"name"`
	Namespace            string                  `json:"namespace"`
	Labels               []Item                  `json:"labels"`
	Replicas             int32                   `json:"replicas"`
	Selector             []Item                  `json:"selector"`
	Template             Pod                     `json:"template"`
	VolumeClaimTemplates []PersistentVolumeClaim `json:"volumeClaimTemplates"`
	ServiceName          string                  `json:"serviceName"`
}
