package service

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/crazyfrankie/kube-ctl/internal/model/convert"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
)

type DeploymentService interface {
	CreateOrUpdateDeployment(ctx context.Context, req *req.Deployment) error
	DeleteDeployment(ctx context.Context, name string, namespace string) error
	GetDeploymentDetail(ctx context.Context, name string, namespace string) (*appsv1.Deployment, error)
	GetDeploymentList(ctx context.Context, namespace string) ([]appsv1.Deployment, error)
}

type deploymentService struct {
	clientSet *kubernetes.Clientset
}

func NewDeploymentService(cs *kubernetes.Clientset) DeploymentService {
	return &deploymentService{clientSet: cs}
}

func (s *deploymentService) CreateOrUpdateDeployment(ctx context.Context, req *req.Deployment) error {
	deployment := convert.DeploymentReqConvert(req)

	if _, err := s.clientSet.AppsV1().Deployments(deployment.Namespace).Get(ctx, deployment.Name, metav1.GetOptions{}); err == nil {
		_, err := s.clientSet.AppsV1().Deployments(deployment.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})

		return err
	}

	_, err := s.clientSet.AppsV1().Deployments(deployment.Namespace).Create(ctx, deployment, metav1.CreateOptions{})

	return err
}

func (s *deploymentService) DeleteDeployment(ctx context.Context, name string, namespace string) error {
	return s.clientSet.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (s *deploymentService) GetDeploymentDetail(ctx context.Context, name string, namespace string) (*appsv1.Deployment, error) {
	res, err := s.clientSet.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *deploymentService) GetDeploymentList(ctx context.Context, namespace string) ([]appsv1.Deployment, error) {
	res, err := s.clientSet.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return res.Items, nil
}
