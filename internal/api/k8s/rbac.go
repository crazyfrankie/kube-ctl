package k8s

import (
	"context"
	"net/http"
	"strings"

	"github.com/crazyfrankie/gem/gerrors"
	"github.com/gin-gonic/gin"

	"github.com/crazyfrankie/kube-ctl/internal/model/convert"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/internal/service"
	"github.com/crazyfrankie/kube-ctl/pkg/response"
)

type RbacHandler struct {
	svc service.RbacService
}

func NewRbacHandler(svc service.RbacService) *RbacHandler {
	return &RbacHandler{svc: svc}
}

func (h *RbacHandler) RegisterRoute(r *gin.Engine) {
	rbacGroup := r.Group("api/rbac")
	{
		saGroup := rbacGroup.Group("sa")
		{
			saGroup.POST("", h.CreateServiceAccount())
			saGroup.DELETE("", h.DeleteServiceAccount())
			saGroup.GET("", h.GetServiceAccountList())
		}
		roleGroup := rbacGroup.Group("role")
		{
			roleGroup.POST("", h.CreateOrUpdateRole())
			roleGroup.DELETE("", h.DeleteRole())
			roleGroup.GET("", h.GetRoleDetail())
			roleGroup.GET("list", h.GetRoleList())
		}
		rbGroup := rbacGroup.Group("rb")
		{
			rbGroup.POST("", h.CreateOrUpdateRoleBinding())
			rbGroup.DELETE("", h.DeleteRoleBinding())
			rbGroup.GET("", h.GetRoleBindingDetail())
			rbGroup.GET("list", h.GetRoleBindingList())
		}
	}
}

// CreateServiceAccount
// @Summary 创建 ServiceAccount
// @Description 创建一个 ServiceAccount 以进行用户认证
// @Tags RBAC 管理
// @Accept json
// @Produce json
// @Param serviceAccount body req.ServiceAccount true "ServiceAccount 信息"
// @Success 200 {object} response.Response "创建 SA 成功"
// @Failure 400 {object} response.Response "参数错误(code=20001)"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/rbac/sa [post]
func (h *RbacHandler) CreateServiceAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var createReq req.ServiceAccount
		if err := c.ShouldBind(&createReq); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20001, "bind error "+err.Error()))
			return
		}

		err := h.svc.CreateServiceAccount(context.Background(), &createReq)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// DeleteServiceAccount
// @Summary 删除 ServiceAccount
// @Description 删除 ServiceAccount 用户访问集群账号
// @Tags RBAC 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "ServiceAccount 名称"
// @Success 200 {object} response.Response "删除 Role 成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/rbac/sa [delete]
func (h *RbacHandler) DeleteServiceAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		err := h.svc.DeleteServiceAccount(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// GetServiceAccountList
// @Summary 获取 ServiceAccount 列表
// @Description 获取所有 ServiceAccount 列表
// @Tags RBAC 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param keyword query string true "关键词"
// @Success 200 {object} response.Response{data=[]resp.ServiceAccount} "返回 ServiceAccount 列表"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/rbac/sa [get]
func (h *RbacHandler) GetServiceAccountList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")
		keyword := c.Query("keyword")

		res, err := h.svc.GetServiceAccountList(context.Background(), ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		sas := make([]resp.ServiceAccount, 0, len(res))
		for _, i := range res {
			if strings.Contains(i.Name, keyword) {
				sas = append(sas, resp.ServiceAccount{
					Name:      i.Name,
					Namespace: i.Namespace,
					Age:       i.CreationTimestamp.Unix(),
				})
			}
		}

		response.SuccessWithData(c, sas)
	}
}

// CreateOrUpdateRole
// @Summary 创建或更新 Role | ClusterRole
// @Description 创建或更新 Role | ClusterRole 作为用户访问规则
// @Tags RBAC 管理
// @Accept json
// @Produce json
// @Param role body req.Role true "Role | ClusterRole 规则信息"
// @Success 200 {object} response.Response "创建或更新 Role 成功"
// @Failure 400 {object} response.Response "参数错误(code=20001)"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/rbac/role [post]
func (h *RbacHandler) CreateOrUpdateRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		var createReq req.Role
		if err := c.ShouldBind(&createReq); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20001, "bind error "+err.Error()))
			return
		}

		err := h.svc.CreateOrUpdateRole(context.Background(), &createReq)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// DeleteRole
// @Summary 删除 Role | ClusterRole
// @Description 删除 Role | ClusterRole 用户访问规则
// @Tags RBAC 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "Role 名称"
// @Success 200 {object} response.Response "删除 Role 成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/rbac/role [delete]
func (h *RbacHandler) DeleteRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		err := h.svc.DeleteRole(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// GetRoleDetail
// @Summary 获取 Role 详情
// @Description 获取 Role 的详细信息
// @Tags RBAC 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "Role 名称"
// @Success 200 {object} response.Response{data=resp.Role} "返回Role的详细信息"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/rbac/role [get]
func (h *RbacHandler) GetRoleDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		role, cluster, err := h.svc.GetRoleDetail(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		var res req.Role
		if ns == "" {
			res = convert.ClusterRoleConvertReq(cluster)
		} else {
			res = convert.RoleConvertReq(role)
		}

		response.SuccessWithData(c, res)
	}
}

// GetRoleList
// @Summary 获取 Role 列表
// @Description 获取所有 Role 列表
// @Tags RBAC 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param keyword query string true "关键词"
// @Success 200 {object} response.Response{data=[]resp.Role} "返回 Role 列表"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/rbac/role/list [get]
func (h *RbacHandler) GetRoleList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")
		keyword := c.Query("keyword")

		roles, clusters, err := h.svc.GetRoleList(context.Background(), ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		var n int
		if ns == "" {
			n = len(clusters)
			res := make([]resp.Role, 0, n)
			for _, i := range clusters {
				if strings.Contains(i.Name, keyword) {
					res = append(res, convert.ClusterRoleConvertResp(&i))
				}
			}

			response.SuccessWithData(c, res)
		} else {
			n = len(roles)
			res := make([]resp.Role, 0, n)
			for _, i := range roles {
				if strings.Contains(i.Name, keyword) {
					res = append(res, convert.RoleConvertResp(&i))
				}
			}

			response.SuccessWithData(c, res)
		}
	}
}

// CreateOrUpdateRoleBinding
// @Summary 创建或更新 RoleBinding | ClusterRoleBinding
// @Description 创建或更新 RoleBinding | ClusterRoleBinding 将用户与规则绑定
// @Tags RBAC 管理
// @Accept json
// @Produce json
// @Param roleBinding body req.RoleBinding true "RoleBinding | ClusterRoleBinding 用户与规则绑定信息"
// @Success 200 {object} response.Response "创建或更新 RB 成功"
// @Failure 400 {object} response.Response "参数错误(code=20001)"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/rbac/rb [post]
func (h *RbacHandler) CreateOrUpdateRoleBinding() gin.HandlerFunc {
	return func(c *gin.Context) {
		var createReq req.RoleBinding
		if err := c.ShouldBind(&createReq); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20001, "bind error "+err.Error()))
			return
		}

		err := h.svc.CreateOrUpdateRoleBinding(context.Background(), &createReq)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// DeleteRoleBinding
// @Summary 删除 RoleBinding | ClusterRoleBinding
// @Description 删除 RoleBinding | ClusterRoleBinding 规则与用户绑定信息
// @Tags RBAC 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "RoleBinding 名称"
// @Success 200 {object} response.Response "删除 RoleBinding 成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/rbac/rb [delete]
func (h *RbacHandler) DeleteRoleBinding() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		err := h.svc.DeleteRoleBinding(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// GetRoleBindingDetail
// @Summary 获取 RoleBinding 详情
// @Description 获取 RoleBinding 的详细信息
// @Tags RBAC 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "Role 名称"
// @Success 200 {object} response.Response{data=resp.RoleBinding} "返回 RoleBinding 的详细信息"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/rbac/rb [get]
func (h *RbacHandler) GetRoleBindingDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		role, cluster, err := h.svc.GetRoleBindingDetail(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		var res req.RoleBinding
		if ns == "" {
			res = convert.ClusterRoleBindingConvertReq(cluster)
		} else {
			res = convert.RoleBindingConvertReq(role)
		}

		response.SuccessWithData(c, res)
	}
}

// GetRoleBindingList
// @Summary 获取 RoleBinding 列表
// @Description 获取所有 RoleBinding 列表
// @Tags RBAC 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param keyword query string true "关键词"
// @Success 200 {object} response.Response{data=[]resp.RoleBinding} "返回 Role 列表"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/rbac/rb/list [get]
func (h *RbacHandler) GetRoleBindingList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")
		keyword := c.Query("keyword")

		roles, clusters, err := h.svc.GetRoleBindingList(context.Background(), ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		var n int
		if ns == "" {
			n = len(clusters)
			res := make([]resp.RoleBinding, 0, n)
			for _, i := range clusters {
				if strings.Contains(i.Name, keyword) {
					res = append(res, convert.ClusterRoleBindingConvertResp(&i))
				}
			}

			response.SuccessWithData(c, res)
		} else {
			n = len(roles)
			res := make([]resp.RoleBinding, 0, n)
			for _, i := range roles {
				if strings.Contains(i.Name, keyword) {
					res = append(res, convert.RoleBindingConvertResp(&i))
				}
			}

			response.SuccessWithData(c, res)
		}
	}
}
