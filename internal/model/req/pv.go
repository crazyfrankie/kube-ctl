package req

import corev1 "k8s.io/api/core/v1"

type PersistentVolume struct {
	Name          string                               `json:"name"`
	Labels        []Item                               `json:"labels"`
	Capacity      string                               `json:"capacity"`
	AccessModes   []corev1.PersistentVolumeAccessMode  `json:"accessModes"`
	ReclaimPolicy corev1.PersistentVolumeReclaimPolicy `json:"reclaimPolicy"`
	VolumeSource  VolumeSource                         `json:"volumeSource"`
}

type VolumeSource struct {
	Type string          `json:"type"`
	NFS  NFSVolumeSource `json:"nfs"`
}

type NFSVolumeSource struct {
	NfsPath   string `json:"nfsPath"`
	NfsServer string `json:"nfsServer"`
	ReadOnly  bool   `json:"readOnly"`
}
