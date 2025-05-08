package service

import (
	"context"

	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/crazyfrankie/kube-ctl/internal/model/convert"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
)

type StorageClassService interface {
	CreateStorageClass(ctx context.Context, req *req.StorageClass) error
	DeleteStorageClass(ctx context.Context, name string) error
	GetStorageClassList(ctx context.Context) ([]storagev1.StorageClass, error)
}

type storageClassService struct {
	clientSet *kubernetes.Clientset
}

func NewStorageClassService(cs *kubernetes.Clientset) StorageClassService {
	return &storageClassService{clientSet: cs}
}

func (s *storageClassService) CreateStorageClass(ctx context.Context, req *req.StorageClass) error {
	sc := convert.StorageClassReqConvert(req)
	_, err := s.clientSet.StorageV1().StorageClasses().Create(ctx, sc, metav1.CreateOptions{})

	return err
}

func (s *storageClassService) DeleteStorageClass(ctx context.Context, name string) error {
	return s.clientSet.StorageV1().StorageClasses().Delete(ctx, name, metav1.DeleteOptions{})
}

func (s *storageClassService) GetStorageClassList(ctx context.Context) ([]storagev1.StorageClass, error) {
	list, err := s.clientSet.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}
