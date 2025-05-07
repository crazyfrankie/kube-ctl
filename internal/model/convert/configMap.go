package convert

import (
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crazyfrankie/kube-ctl/internal/model/req"
)

func CMReqConvert(req *req.ConfigMap) *corev1.ConfigMap {
	labels := make(map[string]string)
	for _, i := range req.Labels {
		labels[i.Key] = i.Value
	}
	data := make(map[string]string)
	for _, i := range req.Data {
		data[i.Key] = i.Value
	}
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:              req.Name,
			Namespace:         req.Namespace,
			CreationTimestamp: metav1.Time{Time: time.Now()},
			Labels:            labels,
		},
		Data: data,
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
	labels := make([]resp.Item, 0, len(cm.Labels))
	for k, v := range cm.Labels {
		labels = append(labels, resp.Item{
			Key:   k,
			Value: v,
		})
	}
	data := make([]resp.Item, 0, len(cm.Data))
	for k, v := range cm.Labels {
		data = append(data, resp.Item{
			Key:   k,
			Value: v,
		})
	}

	return resp.ConfigMapDetail{
		Name:      cm.Name,
		Namespace: cm.Namespace,
		DataNum:   len(cm.Data),
		Age:       cm.CreationTimestamp.Time.Unix(),
		Data:      data,
		Labels:    labels,
	}
}
