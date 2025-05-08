package req

import (
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
)

type StorageClass struct {
	Name                 string                               `json:"name"`
	Namespace            string                               `json:"namespace"`
	Labels               []Item                               `json:"labels"`
	Provisioner          string                               `json:"provisioner"`
	Parameters           []Item                               `json:"parameters"`
	ReclaimPolicy        corev1.PersistentVolumeReclaimPolicy `json:"reclaimPolicy"` // Delete | Retain
	MountOptions         []string                             `json:"mountOptions"`
	AllowVolumeExpansion bool                                 `json:"allowVolumeExpansion"` // whether permit expand
	VolumeBindingMode    storagev1.VolumeBindingMode          `json:"volumeBindingMode"`
}
