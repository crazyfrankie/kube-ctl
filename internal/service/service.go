package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/crazyfrankie/kube-ctl/internal/model/convert"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
)

type SvcService interface {
	CreateOrUpdateService(ctx context.Context, req *req.Service) error
	DeleteService(ctx context.Context, name string, namespace string) error
	GetServiceDetail(ctx context.Context, name string, namespace string) (*corev1.Service, error)
	GetServiceList(ctx context.Context, namespace string) ([]corev1.Service, error)
}

type svcService struct {
	clientSet *kubernetes.Clientset
}

func NewServiceService(cs *kubernetes.Clientset) SvcService {
	return &svcService{clientSet: cs}
}

func (s *svcService) CreateOrUpdateService(ctx context.Context, req *req.Service) error {
	svc := convert.ServiceReqConvert(req)

	if _, err := s.clientSet.CoreV1().Services(svc.Namespace).Get(ctx, svc.Name, metav1.GetOptions{}); err == nil {
		_, err := s.clientSet.CoreV1().Services(svc.Namespace).Update(ctx, svc, metav1.UpdateOptions{})

		return err
	}

	_, err := s.clientSet.CoreV1().Services(svc.Namespace).Create(ctx, svc, metav1.CreateOptions{})

	return err
}

func (s *svcService) DeleteService(ctx context.Context, name string, namespace string) error {
	return s.clientSet.CoreV1().Services(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (s *svcService) GetServiceDetail(ctx context.Context, name string, namespace string) (*corev1.Service, error) {
	res, err := s.clientSet.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *svcService) GetServiceList(ctx context.Context, namespace string) ([]corev1.Service, error) {
	res, err := s.clientSet.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return res.Items, nil
}
