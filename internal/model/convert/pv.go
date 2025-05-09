package convert

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/pkg/utils"
)

const (
	VolumeTypeNFS = "nfs"
)

func PVReqConvert(req *req.PersistentVolume) *corev1.PersistentVolume {
	var volumeSource corev1.PersistentVolumeSource
	switch req.VolumeSource.Type {
	case VolumeTypeNFS:
		volumeSource.NFS = &corev1.NFSVolumeSource{
			Server:   req.VolumeSource.NFS.NfsServer,
			Path:     req.VolumeSource.NFS.NfsPath,
			ReadOnly: req.VolumeSource.NFS.ReadOnly,
		}
		// TODO other type
	default:
		panic("unsupported volume type")
	}

	return &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:   req.Name,
			Labels: utils.ReqItemToMap(req.Labels),
		},
		Spec: corev1.PersistentVolumeSpec{
			AccessModes: req.AccessModes,
			Capacity: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceStorage: resource.MustParse(req.Capacity),
			},
			PersistentVolumeReclaimPolicy: req.ReclaimPolicy,
			PersistentVolumeSource:        volumeSource,
		},
	}
}

func PVConvertResp(pv *corev1.PersistentVolume) resp.PersistentVolumeItem {
	var attributeName, claim string
	if pv.Spec.VolumeAttributesClassName != nil {
		attributeName = *pv.Spec.VolumeAttributesClassName
	} else {
		attributeName = "<unset>"
	}
	if pv.Spec.ClaimRef != nil {
		claim = pv.Spec.ClaimRef.Name
	}
	return resp.PersistentVolumeItem{
		Name:                 pv.Name,
		Labels:               utils.ResMapToItem(pv.Labels),
		Capacity:             pv.Spec.Capacity.Storage().String(),
		AccessModes:          pv.Spec.AccessModes,
		ReclaimPolicy:        pv.Spec.PersistentVolumeReclaimPolicy,
		Status:               pv.Status.Phase,
		Claim:                claim,
		StorageClass:         pv.Spec.StorageClassName,
		VolumeAttributeClass: attributeName,
		Reason:               pv.Status.Reason,
		Age:                  pv.CreationTimestamp.Time.Unix(),
	}
}
