package convert

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/pkg/utils"
)

func JobReqConvert(req *req.Job) *batchv1.Job {
	pod := PodReqConvert(&req.Template)
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels:    utils.ReqItemToMap(req.Labels),
		},
		Spec: batchv1.JobSpec{
			Completions: &req.Completions,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: pod.ObjectMeta,
				Spec:       pod.Spec,
			},
		},
	}
}

func JobConvertReq(job *batchv1.Job) req.Job {
	return req.Job{
		Name:        job.Name,
		Namespace:   job.Namespace,
		Labels:      utils.ReqMapToItem(job.Labels),
		Completions: *job.Spec.Completions,
		Template: *PodConvertReq(&corev1.Pod{
			ObjectMeta: job.Spec.Template.ObjectMeta,
			Spec:       job.Spec.Template.Spec,
		}),
	}
}

func JobConvertResp(job *batchv1.Job) resp.Job {
	res := resp.Job{
		Name:        job.Name,
		Namespace:   job.Namespace,
		Completions: *job.Spec.Completions,
		Succeeded:   job.Status.Succeeded,
		Age:         job.CreationTimestamp.Unix(),
	}
	if job.Status.CompletionTime != nil && job.Status.StartTime != nil {
		res.Duration = job.Status.CompletionTime.Unix() - job.Status.StartTime.Unix()
	}

	return res
}
