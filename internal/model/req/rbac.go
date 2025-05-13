package req

import rbacv1 "k8s.io/api/rbac/v1"

type ServiceAccount struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Labels    []Item `json:"labels"`
}

type Role struct {
	Name      string              `json:"name"`
	Namespace string              `json:"namespace"` // Namespace == "" ? ClusterRole : Role
	Labels    []Item              `json:"labels"`
	Rules     []rbacv1.PolicyRule `json:"rules"`
}

type RoleBinding struct {
	Name      string           `json:"name"`
	Namespace string           `json:"namespace"` // Namespace == "" ?  ClusterRoleBinding: RoleBinding
	Labels    []Item           `json:"labels"`
	RoleRef   string           `json:"roleRef"`
	Subjects  []ServiceAccount `json:"subjects"`
}
