package convert

import (
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/pkg/utils"
)

func ServiceAccountReqConvert(req *req.ServiceAccount) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels:    utils.ReqItemToMap(req.Labels),
		},
	}
}

func RoleReqConvert(req *req.Role) *rbacv1.Role {
	return &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels:    utils.ReqItemToMap(req.Labels),
		},
		Rules: req.Rules,
	}
}

func ClusterRoleReqConvert(req *req.Role) *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:   req.Name,
			Labels: utils.ReqItemToMap(req.Labels),
		},
		Rules: req.Rules,
	}
}

func RoleBindingReqConvert(req *req.RoleBinding) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels:    utils.ReqItemToMap(req.Labels),
		},
		Subjects: getReqRoleBindingSubjects(req.Subjects),
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     req.RoleRef,
		},
	}
}

func ClusterRoleBindingReqConvert(req *req.RoleBinding) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   req.Name,
			Labels: utils.ReqItemToMap(req.Labels),
		},
		Subjects: getReqRoleBindingSubjects(req.Subjects),
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     req.RoleRef,
		},
	}
}

func getReqRoleBindingSubjects(subjects []req.Subject) []rbacv1.Subject {
	res := make([]rbacv1.Subject, 0, len(subjects))
	for _, s := range subjects {
		res = append(res, rbacv1.Subject{
			Kind:      s.Kind,
			Name:      s.Name,
			Namespace: s.Namespace,
		})
	}

	return res
}

func RoleConvertReq(role *rbacv1.Role) req.Role {
	return req.Role{
		Name:      role.Name,
		Namespace: role.Namespace,
		Labels:    utils.ReqMapToItem(role.Labels),
		Rules:     role.Rules,
	}
}

func ClusterRoleConvertReq(role *rbacv1.ClusterRole) req.Role {
	return req.Role{
		Name:      role.Name,
		Namespace: role.Namespace,
		Labels:    utils.ReqMapToItem(role.Labels),
		Rules:     role.Rules,
	}
}

func RoleBindingConvertReq(role *rbacv1.RoleBinding) req.RoleBinding {
	return req.RoleBinding{
		Name:      role.Name,
		Namespace: role.Namespace,
		Labels:    utils.ReqMapToItem(role.Labels),
		RoleRef:   role.RoleRef.Name,
		Subjects:  getRoleBindingSubjectsReq(role.Subjects),
	}
}

func ClusterRoleBindingConvertReq(role *rbacv1.ClusterRoleBinding) req.RoleBinding {
	return req.RoleBinding{
		Name:      role.Name,
		Namespace: role.Namespace,
		Labels:    utils.ReqMapToItem(role.Labels),
		RoleRef:   role.RoleRef.Name,
		Subjects:  getRoleBindingSubjectsReq(role.Subjects),
	}
}

func getRoleBindingSubjectsReq(subjects []rbacv1.Subject) []req.Subject {
	res := make([]req.Subject, 0, len(subjects))
	for _, s := range subjects {
		res = append(res, req.Subject{
			Name:      s.Name,
			Namespace: s.Namespace,
			Kind:      s.Kind,
		})
	}

	return res
}

func RoleConvertResp(role *rbacv1.Role) resp.Role {
	return resp.Role{
		Name:      role.Name,
		Namespace: role.Namespace,
		Age:       role.CreationTimestamp.Unix(),
	}
}

func ClusterRoleConvertResp(role *rbacv1.ClusterRole) resp.Role {
	return resp.Role{
		Name:      role.Name,
		Namespace: role.Namespace,
		Age:       role.CreationTimestamp.Unix(),
	}
}

func RoleBindingConvertResp(role *rbacv1.RoleBinding) resp.RoleBinding {
	return resp.RoleBinding{
		Name:      role.Name,
		Namespace: role.Namespace,
		Age:       role.CreationTimestamp.Unix(),
	}
}

func ClusterRoleBindingConvertResp(role *rbacv1.ClusterRoleBinding) resp.RoleBinding {
	return resp.RoleBinding{
		Name:      role.Name,
		Namespace: role.Namespace,
		Age:       role.CreationTimestamp.Unix(),
	}
}
