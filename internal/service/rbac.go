package service

import (
	"context"
	corev1 "k8s.io/api/core/v1"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/crazyfrankie/kube-ctl/internal/model/convert"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
)

type RbacService interface {
	ServiceAccoutService
	RoleService
	RBService
}

type ServiceAccoutService interface {
	CreateServiceAccount(ctx context.Context, req *req.ServiceAccount) error
	DeleteServiceAccount(ctx context.Context, name string, namespace string) error
	GetServiceAccountList(ctx context.Context, namespace string) ([]corev1.ServiceAccount, error)
}

type RoleService interface {
	CreateOrUpdateRole(ctx context.Context, req *req.Role) error
	DeleteRole(ctx context.Context, name string, namespace string) error
	GetRoleDetail(ctx context.Context, name string, namespace string) (*rbacv1.Role, *rbacv1.ClusterRole, error)
	GetRoleList(ctx context.Context, namespace string) ([]rbacv1.Role, []rbacv1.ClusterRole, error)
}

type RBService interface {
	CreateOrUpdateRoleBinding(ctx context.Context, req *req.RoleBinding) error
	DeleteRoleBinding(ctx context.Context, name string, namespace string) error
	GetRoleBindingDetail(ctx context.Context, name string, namespace string) (*rbacv1.RoleBinding, *rbacv1.ClusterRoleBinding, error)
	GetRoleBindingList(ctx context.Context, namespace string) ([]rbacv1.RoleBinding, []rbacv1.ClusterRoleBinding, error)
}

type rbacService struct {
	clientSet *kubernetes.Clientset
}

func NewRbacService(cs *kubernetes.Clientset) RbacService {
	return &rbacService{clientSet: cs}
}

func (s *rbacService) CreateServiceAccount(ctx context.Context, req *req.ServiceAccount) error {
	sa := convert.ServiceAccountReqConvert(req)

	_, err := s.clientSet.CoreV1().ServiceAccounts(req.Namespace).Create(ctx, sa, metav1.CreateOptions{})

	return err
}

func (s *rbacService) DeleteServiceAccount(ctx context.Context, name string, namespace string) error {
	return s.clientSet.CoreV1().ServiceAccounts(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (s *rbacService) GetServiceAccountList(ctx context.Context, namespace string) ([]corev1.ServiceAccount, error) {
	res, err := s.clientSet.CoreV1().ServiceAccounts(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return res.Items, nil
}

func (s *rbacService) CreateOrUpdateRole(ctx context.Context, req *req.Role) error {
	// Namespace == "" ? ClusterRole : Role

	// create or update ClusterRole
	if req.Namespace == "" {
		clusterRole := convert.ClusterRoleReqConvert(req)
		if exists, err := s.clientSet.RbacV1().ClusterRoles().Get(ctx, clusterRole.Name, metav1.GetOptions{}); err == nil {
			exists.ObjectMeta.Labels = clusterRole.Labels
			exists.Rules = clusterRole.Rules

			_, err := s.clientSet.RbacV1().ClusterRoles().Update(ctx, exists, metav1.UpdateOptions{})

			return err
		}

		_, err := s.clientSet.RbacV1().ClusterRoles().Create(ctx, clusterRole, metav1.CreateOptions{})

		return err
	}

	// create or update Role
	role := convert.RoleReqConvert(req)
	if exists, err := s.clientSet.RbacV1().Roles(role.Namespace).Get(ctx, role.Name, metav1.GetOptions{}); err == nil {
		exists.ObjectMeta.Labels = role.Labels
		exists.Rules = role.Rules

		_, err := s.clientSet.RbacV1().Roles(role.Namespace).Update(ctx, exists, metav1.UpdateOptions{})

		return err
	}
	_, err := s.clientSet.RbacV1().Roles(req.Namespace).Create(ctx, role, metav1.CreateOptions{})

	return err
}

func (s *rbacService) CreateOrUpdateRoleBinding(ctx context.Context, req *req.RoleBinding) error {
	// Namespace == "" ? ClusterRoleBinding : RoleBinding

	// create ClusterRoleBinding
	if req.Namespace == "" {
		clusterRb := convert.ClusterRoleBindingReqConvert(req)
		if exists, err := s.clientSet.RbacV1().ClusterRoleBindings().Get(ctx, clusterRb.Name, metav1.GetOptions{}); err == nil {
			exists.ObjectMeta.Labels = clusterRb.Labels
			exists.Subjects = clusterRb.Subjects
			exists.RoleRef = clusterRb.RoleRef

			_, err := s.clientSet.RbacV1().ClusterRoleBindings().Update(ctx, exists, metav1.UpdateOptions{})

			return err
		}
		_, err := s.clientSet.RbacV1().ClusterRoleBindings().Create(ctx, clusterRb, metav1.CreateOptions{})

		return err
	}

	// create RoleBinding
	rb := convert.RoleBindingReqConvert(req)

	if exists, err := s.clientSet.RbacV1().RoleBindings(rb.Namespace).Get(ctx, rb.Name, metav1.GetOptions{}); err == nil {
		exists.ObjectMeta.Labels = rb.Labels
		exists.Subjects = rb.Subjects
		exists.RoleRef = rb.RoleRef

		_, err := s.clientSet.RbacV1().RoleBindings(rb.Namespace).Update(ctx, exists, metav1.UpdateOptions{})

		return err
	}
	_, err := s.clientSet.RbacV1().RoleBindings(req.Namespace).Create(ctx, rb, metav1.CreateOptions{})

	return err
}

func (s *rbacService) DeleteRole(ctx context.Context, name string, namespace string) error {
	if namespace == "" {
		return s.clientSet.RbacV1().ClusterRoles().Delete(ctx, name, metav1.DeleteOptions{})
	}

	return s.clientSet.RbacV1().Roles(name).Delete(ctx, name, metav1.DeleteOptions{})
}

func (s *rbacService) GetRoleDetail(ctx context.Context, name string, namespace string) (*rbacv1.Role, *rbacv1.ClusterRole, error) {
	if namespace == "" {
		res, err := s.clientSet.RbacV1().ClusterRoles().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return nil, nil, err
		}

		return nil, res, nil
	}

	res, err := s.clientSet.RbacV1().Roles(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	return res, nil, nil
}

func (s *rbacService) GetRoleList(ctx context.Context, namespace string) ([]rbacv1.Role, []rbacv1.ClusterRole, error) {
	if namespace == "" {
		res, err := s.clientSet.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, nil, err
		}

		return nil, res.Items, nil
	}

	res, err := s.clientSet.RbacV1().Roles(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}

	return res.Items, nil, nil
}

func (s *rbacService) DeleteRoleBinding(ctx context.Context, name string, namespace string) error {
	if namespace == "" {
		return s.clientSet.RbacV1().ClusterRoleBindings().Delete(ctx, name, metav1.DeleteOptions{})
	}

	return s.clientSet.RbacV1().RoleBindings(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (s *rbacService) GetRoleBindingDetail(ctx context.Context, name string, namespace string) (*rbacv1.RoleBinding, *rbacv1.ClusterRoleBinding, error) {
	if namespace == "" {
		res, err := s.clientSet.RbacV1().ClusterRoleBindings().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return nil, nil, err
		}

		return nil, res, nil
	}

	res, err := s.clientSet.RbacV1().RoleBindings(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	return res, nil, nil
}

func (s *rbacService) GetRoleBindingList(ctx context.Context, namespace string) ([]rbacv1.RoleBinding, []rbacv1.ClusterRoleBinding, error) {
	if namespace == "" {
		res, err := s.clientSet.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, nil, err
		}

		return nil, res.Items, nil
	}

	res, err := s.clientSet.RbacV1().RoleBindings(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}

	return res.Items, nil, nil
}
