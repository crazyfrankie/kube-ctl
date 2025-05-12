package service

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/crazyfrankie/kube-ctl/internal/model/convert"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
)

type StatefulSetService interface {
	CreateOrUpdateStatefulSet(ctx context.Context, req *req.StatefulSet) error
	DeleteStatefulSet(ctx context.Context, name string, namespace string) error
	GetStatefulSetDetail(ctx context.Context, name string, namespace string) (*appsv1.StatefulSet, error)
	GetStatefulSetList(ctx context.Context, namespace string) ([]appsv1.StatefulSet, error)
}

type statefulSetService struct {
	clientSet *kubernetes.Clientset
}

func NewStatefulSetService(cs *kubernetes.Clientset) StatefulSetService {
	return &statefulSetService{clientSet: cs}
}

func (s *statefulSetService) CreateOrUpdateStatefulSet(ctx context.Context, req *req.StatefulSet) error {
	stateful := convert.StatefulSetReqConvert(req)

	if exists, err := s.clientSet.AppsV1().StatefulSets(stateful.Namespace).Get(ctx, stateful.Name, metav1.GetOptions{}); err != nil {
		exists.Spec = stateful.Spec
		_, err := s.clientSet.AppsV1().StatefulSets(stateful.Namespace).Create(ctx, exists, metav1.CreateOptions{})

		return err
	}

	_, err := s.clientSet.AppsV1().StatefulSets(stateful.Namespace).Update(ctx, stateful, metav1.UpdateOptions{})

	return err
}

func (s *statefulSetService) DeleteStatefulSet(ctx context.Context, name string, namespace string) error {
	return s.clientSet.AppsV1().StatefulSets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (s *statefulSetService) GetStatefulSetDetail(ctx context.Context, name string, namespace string) (*appsv1.StatefulSet, error) {
	res, err := s.clientSet.AppsV1().StatefulSets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *statefulSetService) GetStatefulSetList(ctx context.Context, namespace string) ([]appsv1.StatefulSet, error) {
	res, err := s.clientSet.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return res.Items, nil
}
