package resp

import corev1 "k8s.io/api/core/v1"

type PersistentVolumeClaim struct {
	Name                 string                              `json:"name"`
	Namespace            string                              `json:"namespace"`
	Status               corev1.PersistentVolumeClaimPhase   `json:"status"`
	Volume               string                              `json:"volume"` // PV name
	Capacity             string                              `json:"capacity"`
	AccessModes          []corev1.PersistentVolumeAccessMode `json:"accessModes"`
	VolumeAttributeClass string                              `json:"volumeAttributeClass"`
	Age                  int64                               `json:"age"`
}
