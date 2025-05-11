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

type DaemonSetHandler struct {
	svc service.DaemonSetService
}

func NewDaemonSetHandler(svc service.DaemonSetService) *DaemonSetHandler {
	return &DaemonSetHandler{svc: svc}
}

func (h *DaemonSetHandler) RegisterRoute(r *gin.Engine) {
	daemonGroup := r.Group("api/daemonset")
	{
		daemonGroup.POST("", h.CreateOrUpdateDaemonSet())
		daemonGroup.DELETE("", h.DeleteDaemonSet())
		daemonGroup.GET("", h.GetDaemonSetDetail())
		daemonGroup.GET("list", h.GetDaemonSetList())
	}
}

// CreateOrUpdateDaemonSet
// @Summary 创建或更新 DaemonSet
// @Description 创建新的 DaemonSet 或更新已存在的 DaemonSet
// @Tags DaemonSet 管理
// @Accept json
// @Produce json
// @Param pod body req.DaemonSet true "DaemonSet 配置信息"
// @Success 200 {object} response.Response "操作成功，返回成功消息"
// @Failure 400 {object} response.Response "参数错误(code=20001)"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/daemonset [post]
func (h *DaemonSetHandler) CreateOrUpdateDaemonSet() gin.HandlerFunc {
	return func(c *gin.Context) {
		var createReq req.DaemonSet
		if err := c.ShouldBind(&createReq); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20001, "bind error "+err.Error()))
			return
		}

		err := h.svc.CreateOrUpdateDaemonSet(context.Background(), &createReq)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// DeleteDaemonSet
// @Summary 删除 DaemonSet
// @Description 删除指定命名空间下的指定DaemonSet
// @Tags DaemonSet 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "DaemonSet 名称"
// @Success 200 {object} response.Response "删除成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/daemonset [delete]
func (h *DaemonSetHandler) DeleteDaemonSet() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		err := h.svc.DeleteDaemonSet(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// GetDaemonSetDetail
// @Summary 获取DaemonSet详情
// @Description 获取指定命名空间下指定DaemonSet的详细信息
// @Tags DaemonSet 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "DaemonSet 名称"
// @Success 200 {object} response.Response{data=req.DaemonSet} "返回DaemonSet的详细信息"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/daemonset [get]
func (h *DaemonSetHandler) GetDaemonSetDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		res, err := h.svc.GetDaemonSetDetail(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		daemon := convert.DaemonSetConvertReq(res)

		response.SuccessWithData(c, daemon)
	}
}

// GetDaemonSetList
// @Summary 获取DaemonSet列表
// @Description 获取指定命名空间下的所有DaemonSet列表
// @Tags DaemonSet 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param keyword query string true "关键词"
// @Success 200 {object} response.Response{data=[]resp.DaemonSet} "返回DaemonSet列表"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/daemonset/list [get]
func (h *DaemonSetHandler) GetDaemonSetList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")
		keyword := c.Query("keyword")

		res, err := h.svc.GetDaemonSetList(context.Background(), ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		daemons := make([]resp.DaemonSet, 0, len(res))
		for _, d := range res {
			if strings.Contains(d.Name, keyword) {
				daemons = append(daemons, convert.DaemonSetConvertResp(&d))
			}
		}

		response.SuccessWithData(c, daemons)
	}
}
