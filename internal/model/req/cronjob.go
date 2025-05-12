package req

import batchv1 "k8s.io/api/batch/v1"

type CronJob struct {
	Name                       string                    `json:"name"`
	Namespace                  string                    `json:"namespace"`
	Labels                     []Item                    `json:"labels"`
	Schedule                   string                    `json:"schedule"`          // cron 表达式
	Suspend                    bool                      `json:"suspend"`           // 是否暂停 cronjob
	ConcurrencyPolicy          batchv1.ConcurrencyPolicy `json:"concurrencyPolicy"` // 并发策略
	SuccessfulJobsHistoryLimit int32                     `json:"successfulJobsHistoryLimit"`
	FailedJobsHistoryLimit     int32                     `json:"failedJobsHistoryLimit"`
	JobBase                    JobBase                   `json:"jobBase"`
	Template                   Pod                       `json:"template"`
}

type JobBase struct {
	Completions  int32 `json:"completions"`
	BackoffLimit int32 `json:"backoffLimit"`
}
