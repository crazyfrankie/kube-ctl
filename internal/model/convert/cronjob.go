package convert

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/pkg/utils"
)

func CronJobReqConvert(req *req.CronJob) *batchv1.CronJob {
	pod := PodReqConvert(&req.Template)

	return &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels:    utils.ReqItemToMap(req.Labels),
		},
		Spec: batchv1.CronJobSpec{
			Schedule:                   req.Schedule,
			Suspend:                    &req.Suspend,
			ConcurrencyPolicy:          req.ConcurrencyPolicy,
			FailedJobsHistoryLimit:     &req.FailedJobsHistoryLimit,
			SuccessfulJobsHistoryLimit: &req.SuccessfulJobsHistoryLimit,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					BackoffLimit: &req.JobBase.BackoffLimit,
					Completions:  &req.JobBase.Completions,
					Template: corev1.PodTemplateSpec{
						ObjectMeta: pod.ObjectMeta,
						Spec:       pod.Spec,
					},
				},
			},
		},
	}
}

func CronJobConvertReq(cron *batchv1.CronJob) req.CronJob {
	return req.CronJob{
		Name:                       cron.Name,
		Namespace:                  cron.Namespace,
		Labels:                     utils.ReqMapToItem(cron.Labels),
		Schedule:                   cron.Spec.Schedule,
		Suspend:                    *cron.Spec.Suspend,
		ConcurrencyPolicy:          cron.Spec.ConcurrencyPolicy,
		SuccessfulJobsHistoryLimit: *cron.Spec.SuccessfulJobsHistoryLimit,
		FailedJobsHistoryLimit:     *cron.Spec.FailedJobsHistoryLimit,
		JobBase: req.JobBase{
			Completions:  *cron.Spec.JobTemplate.Spec.Completions,
			BackoffLimit: *cron.Spec.JobTemplate.Spec.BackoffLimit,
		},
		Template: *PodConvertReq(&corev1.Pod{
			ObjectMeta: cron.Spec.JobTemplate.Spec.Template.ObjectMeta,
			Spec:       cron.Spec.JobTemplate.Spec.Template.Spec,
		}),
	}
}

func CronJobConvertResp(cron *batchv1.CronJob) resp.CronJob {
	res := resp.CronJob{
		Name:      cron.Name,
		Namespace: cron.Namespace,
		Schedule:  cron.Spec.Schedule,
		Suspend:   *cron.Spec.Suspend,
		Active:    len(cron.Status.Active),
		Age:       cron.CreationTimestamp.Unix(),
	}
	if cron.Status.LastScheduleTime != nil {
		res.LastScheduleTime = cron.Status.LastScheduleTime.Unix()
	}

	return res
}
