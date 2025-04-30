package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NodeService interface {
	NodeList(ctx context.Context) ([]corev1.Node, error)
	NodeDetail(ctx context.Context, name string) (*corev1.Node, error)
}

type nodeService struct {
	clientSet *kubernetes.Clientset
}

func NewNodeService(cs *kubernetes.Clientset) NodeService {
	return &nodeService{clientSet: cs}
}

func (s *nodeService) NodeList(ctx context.Context) ([]corev1.Node, error) {
	res, err := s.clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return res.Items, nil
}

func (s *nodeService) NodeDetail(ctx context.Context, name string) (*corev1.Node, error) {
	res, err := s.clientSet.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}
