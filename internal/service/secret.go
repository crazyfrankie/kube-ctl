package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/crazyfrankie/kube-ctl/internal/model/convert"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
)

type SecretService interface {
	CreateOrUpdateSecret(ctx context.Context, req *req.Secret) error
	GetSecret(ctx context.Context, name string, namespace string) (*corev1.Secret, error)
	GetSecretList(ctx context.Context, namespace string) ([]corev1.Secret, error)
	DeleteSecret(ctx context.Context, name string, namespace string) error
}

type secretService struct {
	clientSet *kubernetes.Clientset
}

func NewSecretService(cs *kubernetes.Clientset) SecretService {
	return &secretService{clientSet: cs}
}

func (s *secretService) CreateOrUpdateSecret(ctx context.Context, req *req.Secret) error {
	secret := convert.SecretReqConvert(req)
	if _, err := s.clientSet.CoreV1().Secrets(secret.Namespace).Get(ctx, secret.Name, metav1.GetOptions{}); err == nil {
		_, err := s.clientSet.CoreV1().Secrets(secret.Namespace).Update(ctx, secret, metav1.UpdateOptions{})
		if err != nil {
			return err
		}

		return nil
	}

	_, err := s.clientSet.CoreV1().Secrets(secret.Namespace).Create(ctx, secret, metav1.CreateOptions{})

	return err
}

func (s *secretService) GetSecret(ctx context.Context, name string, namespace string) (*corev1.Secret, error) {
	res, err := s.clientSet.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *secretService) GetSecretList(ctx context.Context, namespace string) ([]corev1.Secret, error) {
	res, err := s.clientSet.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return res.Items, nil
}

func (s *secretService) DeleteSecret(ctx context.Context, name string, namespace string) error {
	return s.clientSet.CoreV1().Secrets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
