package resp

import corev1 "k8s.io/api/core/v1"

type PersistentVolumeItem struct {
	Name                 string                               `json:"name"`
	Labels               []Item                               `json:"labels"`
	Capacity             string                               `json:"capacity"`
	AccessModes          []corev1.PersistentVolumeAccessMode  `json:"accessModes"`
	ReclaimPolicy        corev1.PersistentVolumeReclaimPolicy `json:"reclaimPolicy"`
	Status               corev1.PersistentVolumePhase         `json:"status"`
	Claim                string                               `json:"claim"` // bind for pvc
	StorageClass         string                               `json:"storageClass"`
	VolumeAttributeClass string                               `json:"volumeAttributeClass"`
	Reason               string                               `json:"reason"`
	Age                  int64                                `json:"age"`
}
