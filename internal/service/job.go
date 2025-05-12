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
	"k8s.io/utils/pointer"

	"github.com/crazyfrankie/kube-ctl/internal/model/convert"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
)

type JobService interface {
	CreateOrUpdateJob(ctx context.Context, req *req.Job) error
	DeleteJob(ctx context.Context, name string, namespace string) error
	GetJobDetail(ctx context.Context, name string, namespace string) (*batchv1.Job, error)
	GetJobList(ctx context.Context, namespace string) ([]batchv1.Job, error)
}

type jobService struct {
	clientSet *kubernetes.Clientset
}

func NewJobService(cs *kubernetes.Clientset) JobService {
	return &jobService{clientSet: cs}
}

func (s *jobService) CreateOrUpdateJob(ctx context.Context, req *req.Job) error {
	job := convert.JobReqConvert(req)

	if exists, err := s.clientSet.BatchV1().Jobs(job.Namespace).Get(ctx, job.Name, metav1.GetOptions{}); err == nil {
		// 校验
		jobCp := *job
		newName := jobCp.Name + "-validate"
		jobCp.Name = newName
		_, err := s.clientSet.BatchV1().Jobs(job.Namespace).Create(ctx, &jobCp, metav1.CreateOptions{
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
		for k, v := range exists.Spec.Template.Labels {
			podSelector = append(podSelector, fmt.Sprintf("%s=%s", k, v))
		}
		// 启动监听
		watcher, err := s.clientSet.BatchV1().Jobs(job.Namespace).Watch(ctx, metav1.ListOptions{
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
						// 删除关联的 pod
						if list, err := s.clientSet.CoreV1().Pods(exists.Namespace).List(ctx, metav1.ListOptions{
							LabelSelector: strings.Join(podSelector, ","),
						}); err == nil {
							for _, i := range list.Items {
								// delete pod
								background := metav1.DeletePropagationBackground
								err = s.clientSet.CoreV1().Pods(i.Namespace).Delete(ctx, i.Name, metav1.DeleteOptions{
									GracePeriodSeconds: pointer.Int64(0),
									PropagationPolicy:  &background,
								})
							}
						}
						_, err = s.clientSet.BatchV1().Jobs(job.Namespace).Create(ctx, job, metav1.CreateOptions{})
						notify <- err
						return
					}
				case <-time.After(time.Second * 5):
					notify <- errors.New("update job timeout")
					return
				}
			}
		}()
		// 删除 Job
		background := metav1.DeletePropagationForeground
		err = s.clientSet.CoreV1().Pods(job.Namespace).Delete(ctx, exists.Name, metav1.DeleteOptions{
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

	_, err := s.clientSet.BatchV1().Jobs(job.Namespace).Create(ctx, job, metav1.CreateOptions{})

	return err
}

func (s *jobService) DeleteJob(ctx context.Context, name string, namespace string) error {
	job, err := s.clientSet.BatchV1().Jobs(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	var labelSelector []string
	for k, v := range job.Labels {
		labelSelector = append(labelSelector, fmt.Sprintf("%s=%s", k, v))
	}
	watcher, err := s.clientSet.BatchV1().Jobs(namespace).Watch(ctx, metav1.ListOptions{
		LabelSelector: strings.Join(labelSelector, ","),
	})
	if err != nil {
		return err
	}
	var podLabelSelector []string
	for k, v := range job.Spec.Template.Labels {
		podLabelSelector = append(podLabelSelector, fmt.Sprintf("%s=%s", k, v))
	}

	err = s.clientSet.BatchV1().Jobs(namespace).Delete(ctx, name, metav1.DeleteOptions{})
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
					// 删除关联 Pod
					if list, err := s.clientSet.CoreV1().Pods(job.Namespace).
						List(ctx, metav1.ListOptions{
							LabelSelector: strings.Join(podLabelSelector, ","),
						}); err == nil {
						//清理 job 关联的 Pod
						for _, i := range list.Items {
							// delete pod
							background := metav1.DeletePropagationBackground
							err = s.clientSet.CoreV1().Pods(i.Namespace).Delete(ctx, i.Name, metav1.DeleteOptions{
								GracePeriodSeconds: pointer.Int64(0),
								PropagationPolicy:  &background,
							})
						}
					}
					notify <- nil
					return
				}
			case <-time.After(5 * time.Second):
				notify <- errors.New("删除Job超时")
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

func (s *jobService) GetJobDetail(ctx context.Context, name string, namespace string) (*batchv1.Job, error) {
	res, err := s.clientSet.BatchV1().Jobs(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *jobService) GetJobList(ctx context.Context, namespace string) ([]batchv1.Job, error) {
	res, err := s.clientSet.BatchV1().Jobs(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return res.Items, nil
}
