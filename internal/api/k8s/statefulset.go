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

type StatefulSetHandler struct {
	svc service.StatefulSetService
}

func NewStatefulSetHandler(svc service.StatefulSetService) *StatefulSetHandler {
	return &StatefulSetHandler{svc: svc}
}

func (h *StatefulSetHandler) RegisterRoute(r *gin.Engine) {
	statefulGroup := r.Group("api/statefulset")
	{
		statefulGroup.POST("", h.CreateOrUpdateStatefulSet())
		statefulGroup.DELETE("", h.DeleteStatefulSet())
		statefulGroup.GET("", h.GetStatefulSetDetail())
		statefulGroup.GET("list", h.GetStatefulSetList())
	}
}

// CreateOrUpdateStatefulSet
// @Summary 创建或更新 StatefulSet
// @Description 创建新的 StatefulSet 或更新已存在的 StatefulSet
// @Tags StatefulSet 管理
// @Accept json
// @Produce json
// @Param pod body req.StatefulSet true "StatefulSet 配置信息"
// @Success 200 {object} response.Response "操作成功，返回成功消息"
// @Failure 400 {object} response.Response "参数错误(code=20001)"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/statefulset [post]
func (h *StatefulSetHandler) CreateOrUpdateStatefulSet() gin.HandlerFunc {
	return func(c *gin.Context) {
		var createReq req.StatefulSet
		if err := c.ShouldBind(&createReq); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20001, "bind error "+err.Error()))
			return
		}

		err := h.svc.CreateOrUpdateStatefulSet(context.Background(), &createReq)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// DeleteStatefulSet
// @Summary 删除 StatefulSet
// @Description 删除指定命名空间下的指定 StatefulSet
// @Tags StatefulSet 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "StatefulSet 名称"
// @Success 200 {object} response.Response "删除成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/statefulset [delete]
func (h *StatefulSetHandler) DeleteStatefulSet() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		err := h.svc.DeleteStatefulSet(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// GetStatefulSetDetail
// @Summary 获取StatefulSet详情
// @Description 获取指定命名空间下指定StatefulSet的详细信息
// @Tags StatefulSet 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "StatefulSet 名称"
// @Success 200 {object} response.Response{data=req.StatefulSet} "返回StatefulSet的详细信息"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/statefulset [get]
func (h *StatefulSetHandler) GetStatefulSetDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		res, err := h.svc.GetStatefulSetDetail(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		state := convert.StatefulSetConvertReq(res)

		response.SuccessWithData(c, state)
	}
}

// GetStatefulSetList
// @Summary 获取StatefulSet列表
// @Description 获取指定命名空间下的所有StatefulSet列表
// @Tags StatefulSet 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param keyword query string true "关键词"
// @Success 200 {object} response.Response{data=[]resp.StatefulSet} "返回StatefulSet列表"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/statefulset/list [get]
func (h *StatefulSetHandler) GetStatefulSetList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")
		keyword := c.Query("keyword")

		res, err := h.svc.GetStatefulSetList(context.Background(), ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		daemons := make([]resp.StatefulSet, 0, len(res))
		for _, d := range res {
			if strings.Contains(d.Name, keyword) {
				daemons = append(daemons, convert.StatefulSetConvertResp(&d))
			}
		}

		response.SuccessWithData(c, daemons)
	}
}
