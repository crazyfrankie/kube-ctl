package convert

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/pkg/utils"
)

func PVCReqConvert(req *req.PersistentVolumeClaim) *corev1.PersistentVolumeClaim {
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels:    utils.ReqItemToMap(req.Labels),
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(req.Capacity),
				},
			},
			AccessModes: req.AccessModes,
			Selector: &metav1.LabelSelector{
				MatchLabels: utils.ReqItemToMap(req.Selector),
			},
			StorageClassName: &req.StorageClassName,
		},
	}
}

func PVCRespConvert(pvc *corev1.PersistentVolumeClaim) resp.PersistentVolumeClaim {
	var attributeName string
	if pvc.Spec.VolumeAttributesClassName != nil {
		attributeName = *pvc.Spec.VolumeAttributesClassName
	} else {
		attributeName = "<unset>"
	}

	return resp.PersistentVolumeClaim{
		Name:                 pvc.Name,
		Namespace:            pvc.Namespace,
		Status:               pvc.Status.Phase,
		Volume:               pvc.Spec.VolumeName,
		Capacity:             pvc.Spec.Resources.Requests.Storage().String(),
		AccessModes:          pvc.Spec.AccessModes,
		VolumeAttributeClass: attributeName,
		Age:                  pvc.CreationTimestamp.Unix(),
	}
}
