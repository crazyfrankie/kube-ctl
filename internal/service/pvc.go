package service

import (
	"context"
	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/crazyfrankie/kube-ctl/internal/model/convert"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
)

type PVCService interface {
	CreatePVC(ctx context.Context, req *req.PersistentVolumeClaim) error
	DeletePVC(ctx context.Context, name string, namespace string) error
	GetPVCList(ctx context.Context, namespace string) ([]corev1.PersistentVolumeClaim, error)
}

type pvcService struct {
	clientSet *kubernetes.Clientset
}

func NewPVCService(cs *kubernetes.Clientset) PVCService {
	return &pvcService{clientSet: cs}
}

func (s *pvcService) CreatePVC(ctx context.Context, req *req.PersistentVolumeClaim) error {
	pvc := convert.PVCReqConvert(req)
	_, err := s.clientSet.CoreV1().PersistentVolumeClaims(pvc.Namespace).Create(ctx, pvc, metav1.CreateOptions{})

	return err
}

func (s *pvcService) DeletePVC(ctx context.Context, name string, namespace string) error {
	return s.clientSet.CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (s *pvcService) GetPVCList(ctx context.Context, namespace string) ([]corev1.PersistentVolumeClaim, error) {
	list, err := s.clientSet.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	
	return list.Items, nil
}
