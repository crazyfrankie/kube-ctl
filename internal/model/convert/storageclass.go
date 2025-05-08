package convert

import (
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/pkg/utils"
)

func StorageClassReqConvert(req *req.StorageClass) *storagev1.StorageClass {
	return &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels:    utils.ReqItemToMap(req.Labels),
		},
		Provisioner:          req.Provisioner,
		MountOptions:         req.MountOptions,
		ReclaimPolicy:        &req.ReclaimPolicy,
		AllowVolumeExpansion: &req.AllowVolumeExpansion,
		VolumeBindingMode:    &req.VolumeBindingMode,
		Parameters:           utils.ReqItemToMap(req.Parameters),
	}
}

func StorageClassConvertResp(sc *storagev1.StorageClass) resp.StorageClass {
	return resp.StorageClass{
		Name:                 sc.Name,
		Labels:               utils.ResMapToItem(sc.Labels),
		MountOptions:         sc.MountOptions,
		Parameters:           utils.ResMapToItem(sc.Parameters),
		Provisioner:          sc.Provisioner,
		ReclaimPolicy:        *sc.ReclaimPolicy,
		VolumeBindingMode:    *sc.VolumeBindingMode,
		AllowVolumeExpansion: *sc.AllowVolumeExpansion,
		Age:                  sc.CreationTimestamp.Unix(),
	}
}
