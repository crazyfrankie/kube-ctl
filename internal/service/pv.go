package service

import (
	"context"
	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/crazyfrankie/kube-ctl/internal/model/convert"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
)

type PVService interface {
	CreatePV(ctx context.Context, req *req.PersistentVolume) error
	DeletePV(ctx context.Context, name string) error
	GetPVList(ctx context.Context) ([]corev1.PersistentVolume, error)
}

type pvService struct {
	clientSet *kubernetes.Clientset
}

func NewPVService(cs *kubernetes.Clientset) PVService {
	return &pvService{clientSet: cs}
}

func (s *pvService) CreatePV(ctx context.Context, req *req.PersistentVolume) error {
	pv := convert.PVReqConvert(req)
	_, err := s.clientSet.CoreV1().PersistentVolumes().Create(ctx, pv, metav1.CreateOptions{})

	return err
}

func (s *pvService) DeletePV(ctx context.Context, name string) error {
	return s.clientSet.CoreV1().PersistentVolumes().Delete(ctx, name, metav1.DeleteOptions{})
}

func (s *pvService) GetPVList(ctx context.Context) ([]corev1.PersistentVolume, error) {
	volume, err := s.clientSet.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return volume.Items, nil
}
