package convert

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/pkg/utils"
)

func DeploymentReqConvert(req *req.Deployment) *appsv1.Deployment {
	pod := PodReqConvert(&req.Template)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels:    utils.ReqItemToMap(req.Labels),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &req.Replicas,
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

func DeploymentConvertReq(deploy *appsv1.Deployment) req.Deployment {
	var replicas int32
	if deploy.Spec.Replicas != nil {
		replicas = *deploy.Spec.Replicas
	}
	return req.Deployment{
		Name:      deploy.Name,
		Namespace: deploy.Namespace,
		Labels:    utils.ReqMapToItem(deploy.Labels),
		Replicas:  replicas,
		Selector:  utils.ReqMapToItem(deploy.Spec.Selector.MatchLabels),
		Template: *PodConvertReq(&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Labels: deploy.Spec.Template.Labels,
			},
			Spec: deploy.Spec.Template.Spec,
		}),
	}
}

func DeploymentConvertResp(deployment *appsv1.Deployment) resp.Deployment {
	var replicas int32
	if deployment.Spec.Replicas != nil {
		replicas = *deployment.Spec.Replicas
	}
	return resp.Deployment{
		Name:      deployment.Name,
		Namespace: deployment.Namespace,
		Replicas:  replicas,
		Ready:     deployment.Status.ReadyReplicas,
		UpToDate:  deployment.Status.UpdatedReplicas,
		Available: deployment.Status.AvailableReplicas,
		Age:       deployment.CreationTimestamp.Unix(),
	}
}
