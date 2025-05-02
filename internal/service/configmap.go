package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/crazyfrankie/kube-ctl/internal/model/convert"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
)

type ConfigMapService interface {
	CreateOrUpdateConfigMap(ctx context.Context, req *req.ConfigMap) error
	GetConfigMap(ctx context.Context, name string, namespace string) (*corev1.ConfigMap, error)
	GetConfigMapList(ctx context.Context, namespace string) ([]corev1.ConfigMap, error)
	DeleteConfigMap(ctx context.Context, name string, namespace string) error
}

type configMapService struct {
	clientSet *kubernetes.Clientset
}

func NewConfigMapService(cs *kubernetes.Clientset) ConfigMapService {
	return &configMapService{clientSet: cs}
}

func (s *configMapService) CreateOrUpdateConfigMap(ctx context.Context, req *req.ConfigMap) error {
	cm := convert.CMReqConvert(req)

	if _, err := s.clientSet.CoreV1().ConfigMaps(cm.Namespace).Get(ctx, cm.Name, metav1.GetOptions{}); err == nil {
		_, err := s.clientSet.CoreV1().ConfigMaps(cm.Namespace).Update(ctx, cm, metav1.UpdateOptions{})
		if err != nil {
			return err
		}

		return nil
	}

	_, err := s.clientSet.CoreV1().ConfigMaps(cm.Namespace).Create(ctx, cm, metav1.CreateOptions{})

	return err
}

func (s *configMapService) GetConfigMap(ctx context.Context, name string, namespace string) (*corev1.ConfigMap, error) {
	res, err := s.clientSet.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *configMapService) GetConfigMapList(ctx context.Context, namespace string) ([]corev1.ConfigMap, error) {
	res, err := s.clientSet.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return res.Items, nil
}

func (s *configMapService) DeleteConfigMap(ctx context.Context, name string, namespace string) error {
	return s.clientSet.CoreV1().ConfigMaps(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
