package resp

type CronJob struct {
	Name             string `json:"name"`
	Namespace        string `json:"namespace"`
	Schedule         string `json:"schedule"`
	Suspend          bool   `json:"suspend"`
	Active           int    `json:"active"`
	LastScheduleTime int64  `json:"lastScheduleTime"`
	Age              int64  `json:"age"`
}
