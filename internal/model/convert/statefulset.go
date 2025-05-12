package convert

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/pkg/utils"
)

func StatefulSetReqConvert(req *req.StatefulSet) *appsv1.StatefulSet {
	pod := PodReqConvert(&req.Template)

	vct := make([]corev1.PersistentVolumeClaim, 0, len(req.VolumeClaimTemplates))
	for _, i := range req.VolumeClaimTemplates {
		vct = append(vct, *PVCReqConvert(&i))
	}
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels:    utils.ReqItemToMap(req.Labels),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &req.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: utils.ReqItemToMap(req.Selector),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: pod.ObjectMeta,
				Spec:       pod.Spec,
			},
			ServiceName:          req.ServiceName,
			VolumeClaimTemplates: vct,
		},
	}
}

func StatefulSetConvertReq(state *appsv1.StatefulSet) req.StatefulSet {
	var replicas int32
	if state.Spec.Replicas != nil {
		replicas = *state.Spec.Replicas
	}

	vct := make([]req.PersistentVolumeClaim, 0, len(state.Spec.VolumeClaimTemplates))
	for _, i := range state.Spec.VolumeClaimTemplates {
		vct = append(vct, req.PersistentVolumeClaim{
			Name:             i.Name,
			AccessModes:      i.Spec.AccessModes,
			Capacity:         i.Spec.Resources.Requests.Storage().String(),
			StorageClassName: *i.Spec.StorageClassName,
		})
	}

	return req.StatefulSet{
		Name:      state.Name,
		Namespace: state.Namespace,
		Labels:    utils.ReqMapToItem(state.Labels),
		Replicas:  replicas,
		Selector:  utils.ReqMapToItem(state.Spec.Selector.MatchLabels),
		Template: *PodConvertReq(&corev1.Pod{
			ObjectMeta: state.Spec.Template.ObjectMeta,
			Spec:       state.Spec.Template.Spec,
		}),
		VolumeClaimTemplates: vct,
		ServiceName:          state.Spec.ServiceName,
	}
}

func StatefulSetConvertResp(state *appsv1.StatefulSet) resp.StatefulSet {
	return resp.StatefulSet{
		Name:      state.Name,
		Namespace: state.Namespace,
		Ready:     state.Status.ReadyReplicas,
		Replicas:  state.Status.Replicas,
		Age:       state.CreationTimestamp.Unix(),
	}
}
