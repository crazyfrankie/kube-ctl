package convert

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/pkg/utils"
)

func DaemonSetReqConvert(req *req.DaemonSet) *appsv1.DaemonSet {
	pod := PodReqConvert(&req.Template)

	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels:    utils.ReqItemToMap(req.Labels),
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: utils.ReqItemToMap(req.Selector),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: pod.ObjectMeta,
				Spec:       pod.Spec,
			},
		},
	}
}

func DaemonSetConvertReq(daemon *appsv1.DaemonSet) req.DaemonSet {
	return req.DaemonSet{
		Name:      daemon.Name,
		Namespace: daemon.Namespace,
		Labels:    utils.ReqMapToItem(daemon.Labels),
		Selector:  utils.ReqMapToItem(daemon.Spec.Selector.MatchLabels),
		Template: *PodConvertReq(&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Labels: daemon.Spec.Template.Labels,
			},
			Spec: daemon.Spec.Template.Spec,
		}),
	}
}

func DaemonSetConvertResp(daemon *appsv1.DaemonSet) resp.DaemonSet {
	return resp.DaemonSet{
		Name:      daemon.Name,
		Namespace: daemon.Namespace,
		Desired:   daemon.Status.DesiredNumberScheduled,
		Current:   daemon.Status.CurrentNumberScheduled,
		Ready:     daemon.Status.NumberReady,
		UpToDate:  daemon.Status.UpdatedNumberScheduled,
		Available: daemon.Status.NumberAvailable,
		Age:       daemon.CreationTimestamp.Unix(),
	}
}
