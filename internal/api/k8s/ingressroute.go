package k8s

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/crazyfrankie/gem/gerrors"
	"github.com/gin-gonic/gin"

	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/internal/service"
	"github.com/crazyfrankie/kube-ctl/pkg/response"
	"github.com/crazyfrankie/kube-ctl/pkg/utils"
)

type IngressRouteHandler struct {
	svc service.IngressRouteService
}

func NewIngressRouteHandler(svc service.IngressRouteService) *IngressRouteHandler {
	return &IngressRouteHandler{svc: svc}
}

func (h *IngressRouteHandler) RegisterRoute(r *gin.Engine) {
	irGroup := r.Group("api/ingroute")
	{
		irGroup.POST("", h.CreateOrUpdateIngresRoute())
		irGroup.DELETE("", h.DeleteIngressRoute())
		irGroup.GET("", h.GetIngressRouteDetail())
		irGroup.GET("list", h.GetIngressRouteList())
		irGroup.GET("mws", h.GetIngressRouteMws())
	}
}

// CreateOrUpdateIngresRoute
// @Summary 创建或更新 IngressRoute
// @Description 创建新的 IngressRoute 或更新已存在的 IngressRoute
// @Tags IngressRoute 管理
// @Accept json
// @Produce json
// @Param pod body req.IngressRoute true "IngressRoute 配置信息"
// @Success 200 {object} response.Response "操作成功，返回成功消息"
// @Failure 400 {object} response.Response "参数错误(code=20001)"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/ingroute [post]
func (h *IngressRouteHandler) CreateOrUpdateIngresRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		var createReq req.IngressRoute
		if err := c.ShouldBind(&createReq); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20001, "bind error "+err.Error()))
			return
		}

		err := h.svc.CreateOrUpdateIngressRoute(context.Background(), &createReq)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// DeleteIngressRoute
// @Summary 删除 IngressRoute
// @Description 删除一个 IngressRoute
// @Tags IngressRoute 管理
// @Accept json
// @Produce json
// @Param name query string true "IngressRoute 名称"
// @Param namespace query string true "命名空间"
// @Success 200 {object} response.Response "删除 IngressRoute 成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/ingroute [delete]
func (h *IngressRouteHandler) DeleteIngressRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		err := h.svc.DeleteIngressRoute(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// GetIngressRouteDetail
// @Summary 获取IngresRoute详情
// @Description 获取指定命名空间下指定IngressRoute的详细信息
// @Tags IngressRoute 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "IngressRoute 名称"
// @Success 200 {object} response.Response{data=req.IngressRoute} "返回IngressRoute的详细信息"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/ingroute [get]
func (h *IngressRouteHandler) GetIngressRouteDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		res, err := h.svc.GetIngressRouteDetail(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.SuccessWithData(c, req.IngressRoute{
			Name:             res.Metadata.Name,
			Namespace:        res.Metadata.Namespace,
			Labels:           utils.ReqMapToItem(res.Metadata.Labels),
			IngressRouteSpec: res.Spec,
		})
	}
}

// GetIngressRouteList
// @Summary 获取IngressRoute列表
// @Description 获取指定命名空间下指定IngressRoute的列表
// @Tags IngressRoute 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param keyword query string true "关键词"
// @Success 200 {object} response.Response{data=[]resp.IngressRoute} "返回IngressRoute的列表"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/ingroute/list [get]
func (h *IngressRouteHandler) GetIngressRouteList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")
		keyword := c.Query("keyword")

		res, err := h.svc.GetIngressRouteList(context.Background(), ns)
		if err != nil {
			if errors.Is(err, service.ErrNoResource) {
				response.Success(c)
				return
			}
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		igs := make([]resp.IngressRoute, 0, len(res.Items))
		for _, i := range res.Items {
			if strings.Contains(i.Metadata.Name, keyword) {
				igs = append(igs, resp.IngressRoute{
					Name:      i.Metadata.Name,
					Namespace: i.Metadata.Namespace,
					Age:       i.Metadata.CreationTimestamp.Unix(),
				})
			}
		}

		response.SuccessWithData(c, igs)
	}
}

// GetIngressRouteMws
// @Summary 获取IngressRoute的Middlewares列表
// @Description 获取指定命名空间下指定IngressRoute的Middlewares列表
// @Tags IngressRoute 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Success 200 {object} response.Response{data=[]string} "返回IngressRoute的Middlewares列表"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/ingroute/mws [get]
func (h *IngressRouteHandler) GetIngressRouteMws() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")

		res, err := h.svc.GetIngressRouteMws(context.Background(), ns)
		if err != nil {
			if errors.Is(err, service.ErrNoResource) {
				response.Success(c)
				return
			}
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.SuccessWithData(c, res)
	}
}
