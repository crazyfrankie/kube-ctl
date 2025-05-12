package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"

	"github.com/crazyfrankie/kube-ctl/internal/model/convert"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
)

type CronJobService interface {
	CreateOrUpdateCronJob(ctx context.Context, req *req.CronJob) error
	DeleteCronJob(ctx context.Context, name string, namespace string) error
	GetCronJobDetail(ctx context.Context, name string, namespace string) (*batchv1.CronJob, error)
	GetCronJobList(ctx context.Context, namespace string) ([]batchv1.CronJob, error)
}

func NewCronJobService(cs *kubernetes.Clientset) CronJobService {
	return &cronJobService{clientSet: cs}
}

type cronJobService struct {
	clientSet *kubernetes.Clientset
}

func (s *cronJobService) CreateOrUpdateCronJob(ctx context.Context, req *req.CronJob) error {
	cron := convert.CronJobReqConvert(req)

	if exists, err := s.clientSet.BatchV1().CronJobs(cron.Namespace).Get(ctx, cron.Name, metav1.GetOptions{}); err == nil {
		// 校验
		cronCp := *cron
		newName := cronCp.Name + "-validate"
		cronCp.Name = newName
		_, err := s.clientSet.BatchV1().CronJobs(cron.Namespace).Create(ctx, &cronCp, metav1.CreateOptions{
			DryRun: []string{metav1.DryRunAll},
		})
		if err != nil {
			return err
		}
		// 获取监听标签
		var labelSelector []string
		for k, v := range exists.Labels {
			labelSelector = append(labelSelector, fmt.Sprintf("%s=%s", k, v))
		}
		var podSelector []string
		for k, v := range exists.Spec.JobTemplate.Spec.Template.Labels {
			podSelector = append(podSelector, fmt.Sprintf("%s=%s", k, v))
		}
		// 启动监听
		watcher, err := s.clientSet.BatchV1().CronJobs(cron.Namespace).Watch(ctx, metav1.ListOptions{
			LabelSelector: strings.Join(labelSelector, ","),
		})
		if err != nil {
			return err
		}
		notify := make(chan error)

		// 异步删除 Pod
		go func() {
			for {
				select {
				case e := <-watcher.ResultChan():
					switch e.Type {
					case watch.Deleted:
						_, err = s.clientSet.BatchV1().CronJobs(cron.Namespace).Create(ctx, cron, metav1.CreateOptions{})
						notify <- err
						return
					}
				case <-time.After(time.Second * 5):
					notify <- errors.New("update cronjob timeout")
					return
				}
			}
		}()
		// 删除 cronjob
		background := metav1.DeletePropagationForeground
		err = s.clientSet.CoreV1().Pods(cron.Namespace).Delete(ctx, exists.Name, metav1.DeleteOptions{
			PropagationPolicy: &background,
		})
		if err != nil {
			return err
		}
		//监听删除后重新创建的信息
		select {
		case errx := <-notify:
			if errx != nil {
				return errx
			}
		}
	}

	_, err := s.clientSet.BatchV1().CronJobs(cron.Namespace).Create(ctx, cron, metav1.CreateOptions{})

	return err
}

func (s *cronJobService) DeleteCronJob(ctx context.Context, name string, namespace string) error {
	job, err := s.clientSet.BatchV1().CronJobs(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	var labelSelector []string
	for k, v := range job.Labels {
		labelSelector = append(labelSelector, fmt.Sprintf("%s=%s", k, v))
	}
	watcher, err := s.clientSet.BatchV1().CronJobs(namespace).Watch(ctx, metav1.ListOptions{
		LabelSelector: strings.Join(labelSelector, ","),
	})
	if err != nil {
		return err
	}
	var podLabelSelector []string
	for k, v := range job.Spec.JobTemplate.Spec.Template.Labels {
		podLabelSelector = append(podLabelSelector, fmt.Sprintf("%s=%s", k, v))
	}

	err = s.clientSet.BatchV1().CronJobs(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	notify := make(chan error)
	go func() {
		for {
			select {
			case e := <-watcher.ResultChan():
				switch e.Type {
				case watch.Deleted:
					notify <- nil
					return
				}
			case <-time.After(5 * time.Second):
				notify <- errors.New("delete CronJob timeout")
				return
			}
		}
	}()
	select {
	case errx := <-notify:
		if errx != nil {
			return errx
		}
	}

	return nil
}

func (s *cronJobService) GetCronJobDetail(ctx context.Context, name string, namespace string) (*batchv1.CronJob, error) {
	res, err := s.clientSet.BatchV1().CronJobs(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *cronJobService) GetCronJobList(ctx context.Context, namespace string) ([]batchv1.CronJob, error) {
	res, err := s.clientSet.BatchV1().CronJobs(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return res.Items, nil
}
