package convert

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/pkg/utils"
)

func CMReqConvert(req *req.ConfigMap) *corev1.ConfigMap {

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:              req.Name,
			Namespace:         req.Namespace,
			CreationTimestamp: metav1.Time{Time: time.Now()},
			Labels:            utils.ReqItemToMap(req.Labels),
		},
		Data: utils.ReqItemToMap(req.Data),
	}
}

func CMConvertListResp(cm *corev1.ConfigMap) resp.ConfigMap {
	return resp.ConfigMap{
		Name:      cm.Name,
		Namespace: cm.Namespace,
		DataNum:   len(cm.Data),
		Age:       cm.CreationTimestamp.Time.Unix(),
	}
}

func CMConvertDetailResp(cm *corev1.ConfigMap) resp.ConfigMapDetail {
	return resp.ConfigMapDetail{
		Name:      cm.Name,
		Namespace: cm.Namespace,
		DataNum:   len(cm.Data),
		Age:       cm.CreationTimestamp.Time.Unix(),
		Data:      utils.ResMapToItem(cm.Data),
		Labels:    utils.ResMapToItem(cm.Labels),
	}
}
