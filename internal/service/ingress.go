package service

import (
	"context"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/crazyfrankie/kube-ctl/internal/model/convert"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IngressService interface {
	CreateOrUpdateIngress(ctx context.Context, req *req.Ingress) error
	DeleteIngress(ctx context.Context, name string, namespace string) error
	GetIngressDetail(ctx context.Context, name string, namespace string) (*networkingv1.Ingress, error)
	GetIngressList(ctx context.Context, namespace string) ([]networkingv1.Ingress, error)
}

type ingressService struct {
	clientSet *kubernetes.Clientset
}

func NewIngressService(cs *kubernetes.Clientset) IngressService {
	return &ingressService{clientSet: cs}
}

func (s *ingressService) CreateOrUpdateIngress(ctx context.Context, req *req.Ingress) error {
	ingress := convert.IngressReqConvert(req)

	if exists, err := s.clientSet.NetworkingV1().Ingresses(ingress.Namespace).Get(ctx, ingress.Name, metav1.GetOptions{}); err == nil {
		exists.Spec = ingress.Spec
		_, err := s.clientSet.NetworkingV1().Ingresses(ingress.Namespace).Update(ctx, exists, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	_, err := s.clientSet.NetworkingV1().Ingresses(ingress.Namespace).Create(ctx, ingress, metav1.CreateOptions{})

	return err
}

func (s *ingressService) DeleteIngress(ctx context.Context, name string, namespace string) error {
	return s.clientSet.NetworkingV1().Ingresses(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (s *ingressService) GetIngressDetail(ctx context.Context, name string, namespace string) (*networkingv1.Ingress, error) {
	res, err := s.clientSet.NetworkingV1().Ingresses(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ingressService) GetIngressList(ctx context.Context, namespace string) ([]networkingv1.Ingress, error) {
	res, err := s.clientSet.NetworkingV1().Ingresses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return res.Items, nil
}
