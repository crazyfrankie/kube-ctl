package req

import corev1 "k8s.io/api/core/v1"

type PersistentVolumeClaim struct {
	Name             string                              `json:"name"`
	Namespace        string                              `json:"namespace"`
	Labels           []Item                              `json:"labels"`
	AccessModes      []corev1.PersistentVolumeAccessMode `json:"accessModes"`
	Capacity         string                              `json:"capacity"`
	Selector         []Item                              `json:"selector"`
	StorageClassName string                              `json:"storageClassName"`
}
