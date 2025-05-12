package service

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/crazyfrankie/kube-ctl/internal/model/convert"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
)

type DaemonSetService interface {
	CreateOrUpdateDaemonSet(ctx context.Context, req *req.DaemonSet) error
	DeleteDaemonSet(ctx context.Context, name string, namespace string) error
	GetDaemonSetDetail(ctx context.Context, name string, namespace string) (*appsv1.DaemonSet, error)
	GetDaemonSetList(ctx context.Context, namespace string) ([]appsv1.DaemonSet, error)
}

type daemonSetService struct {
	clientSet *kubernetes.Clientset
}

func NewDaemonSetService(cs *kubernetes.Clientset) DaemonSetService {
	return &daemonSetService{clientSet: cs}
}

func (s *daemonSetService) CreateOrUpdateDaemonSet(ctx context.Context, req *req.DaemonSet) error {
	daemon := convert.DaemonSetReqConvert(req)

	if _, err := s.clientSet.AppsV1().DaemonSets(daemon.Namespace).Get(ctx, daemon.Name, metav1.GetOptions{}); err == nil {
		_, err := s.clientSet.AppsV1().DaemonSets(daemon.Namespace).Update(ctx, daemon, metav1.UpdateOptions{})

		return err
	}

	_, err := s.clientSet.AppsV1().DaemonSets(daemon.Namespace).Create(ctx, daemon, metav1.CreateOptions{})

	return err
}

func (s *daemonSetService) DeleteDaemonSet(ctx context.Context, name string, namespace string) error {
	return s.clientSet.AppsV1().DaemonSets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (s *daemonSetService) GetDaemonSetDetail(ctx context.Context, name string, namespace string) (*appsv1.DaemonSet, error) {
	res, err := s.clientSet.AppsV1().DaemonSets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *daemonSetService) GetDaemonSetList(ctx context.Context, namespace string) ([]appsv1.DaemonSet, error) {
	res, err := s.clientSet.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return res.Items, nil
}
