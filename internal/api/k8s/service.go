package k8s

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/crazyfrankie/gem/gerrors"
	"github.com/crazyfrankie/kube-ctl/internal/model/convert"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/internal/service"
	"github.com/crazyfrankie/kube-ctl/pkg/response"
)

type ServiceHandler struct {
	svc service.SvcService
}

func NewServiceHandler(svc service.SvcService) *ServiceHandler {
	return &ServiceHandler{svc: svc}
}

func (h *ServiceHandler) RegisterRoute(r *gin.Engine) {
	serviceGroup := r.Group("api/service")
	{
		serviceGroup.POST("", h.CreateOrUpdateService())
		serviceGroup.DELETE("", h.DeleteService())
		serviceGroup.GET("", h.GetServiceDetail())
		serviceGroup.GET("list", h.GetServiceList())
	}
}

// CreateOrUpdateService
// @Summary 创建或更新 Service
// @Description 创建新的 Service 或更新已存在的 Service
// @Tags Service 管理
// @Accept json
// @Produce json
// @Param pod body req.Service true "Service 配置信息"
// @Success 200 {object} response.Response "操作成功，返回成功消息"
// @Failure 400 {object} response.Response "参数错误(code=20001)"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/service [post]
func (h *ServiceHandler) CreateOrUpdateService() gin.HandlerFunc {
	return func(c *gin.Context) {
		var createReq req.Service
		if err := c.ShouldBind(&createReq); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20001, "bind error "+err.Error()))
			return
		}

		err := h.svc.CreateOrUpdateService(context.Background(), &createReq)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// DeleteService
// @Summary 删除 Service
// @Description 删除一个 Service
// @Tags Service 管理
// @Accept json
// @Produce json
// @Param name query string true "Service 名称"
// @Param namespace query string true "命名空间"
// @Success 200 {object} response.Response "删除 Service 成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/service [delete]
func (h *ServiceHandler) DeleteService() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		err := h.svc.DeleteService(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// GetServiceDetail
// @Summary 获取Service详情
// @Description 获取指定命名空间下指定Service的详细信息
// @Tags Service 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "Service 名称"
// @Success 200 {object} response.Response{data=req.Service} "返回Service的详细信息"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/service [get]
func (h *ServiceHandler) GetServiceDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		res, err := h.svc.GetServiceDetail(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		svc := convert.ServiceConvertReq(res)

		response.SuccessWithData(c, svc)
	}
}

// GetServiceList
// @Summary 获取Service列表
// @Description 获取指定命名空间下指定Service的列表
// @Tags Service 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param keyword query string true "关键词"
// @Success 200 {object} response.Response{data=[]resp.Service} "返回Service的列表"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/service/list [get]
func (h *ServiceHandler) GetServiceList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")
		keyword := c.Query("keyword")

		res, err := h.svc.GetServiceList(context.Background(), ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		svcs := make([]resp.Service, 0, len(res))
		for _, i := range res {
			if strings.Contains(i.Name, keyword) {
				svcs = append(svcs, convert.ServiceConvertResp(&i))
			}
		}

		response.SuccessWithData(c, svcs)
	}
}
