package resp

import (
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
)

type StorageClass struct {
	Name                 string                               `json:"name"`
	Labels               []Item                               `json:"labels"`
	Provisioner          string                               `json:"provisioner"`
	ReclaimPolicy        corev1.PersistentVolumeReclaimPolicy `json:"reclaimPolicy"`
	MountOptions         []string                             `json:"mountOptions"`
	Parameters           []Item                               `json:"parameters"`
	VolumeBindingMode    storagev1.VolumeBindingMode          `json:"volumeBindingMode"`
	AllowVolumeExpansion bool                                 `json:"allowVolumeExpansion"`
	Age                  int64                                `json:"age"`
}
