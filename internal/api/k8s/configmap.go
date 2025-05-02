package k8s

import (
	"context"
	"net/http"

	"github.com/crazyfrankie/gem/gerrors"
	"github.com/gin-gonic/gin"

	"github.com/crazyfrankie/kube-ctl/internal/model/convert"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/internal/service"
	"github.com/crazyfrankie/kube-ctl/pkg/response"
)

type ConfigMapHandler struct {
	svc service.ConfigMapService
}

func NewConfigMapHandler(svc service.ConfigMapService) *ConfigMapHandler {
	return &ConfigMapHandler{svc: svc}
}

func (h *ConfigMapHandler) RegisterRoute(r *gin.Engine) {
	cmGroup := r.Group("api/configmap")
	{
		cmGroup.POST("", h.CreateOrUpdateConfigMap())
		cmGroup.GET("", h.GetConfigMap())
		cmGroup.GET("list", h.GetConfigMapList())
		cmGroup.DELETE("", h.DeleteConfigMap())
	}
}

// CreateOrUpdateConfigMap
// @Summary 创建或更新 ConfigMap
// @Description 创建新的 ConfigMap 或更新已存在的 ConfigMap
// @Tags ConfigMap管理
// @Accept json
// @Produce json
// @Param pod body req.ConfigMap true "ConfigMap 配置信息"
// @Success 200 {object} response.Response "操作成功，返回成功消息"
// @Failure 400 {object} response.Response "参数错误(code=20001)"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/configmap [post]
func (h *ConfigMapHandler) CreateOrUpdateConfigMap() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cmReq req.ConfigMap
		if err := c.ShouldBind(&cmReq); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20001, "bind error "+err.Error()))
			return
		}

		err := h.svc.CreateOrUpdateConfigMap(context.Background(), &cmReq)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// GetConfigMap
// @Summary 获取ConfigMap详情
// @Description 获取指定命名空间下指定ConfigMap的详细信息
// @Tags ConfigMap管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "ConfigMap名称"
// @Success 200 {object} response.Response{data=resp.ConfigMapDetail} "返回ConfigMap的详细信息"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/configmap [get]
func (h *ConfigMapHandler) GetConfigMap() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")
		name := c.Query("name")

		res, err := h.svc.GetConfigMap(context.Background(), ns, name)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		cm := convert.CMConvertDetailResp(res)

		response.SuccessWithData(c, cm)
	}
}

// GetConfigMapList
// @Summary 获取ConfigMap列表
// @Description 获取指定命名空间下的所有ConfigMap列表
// @Tags ConfigMap管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Success 200 {object} response.Response{data=[]resp.ConfigMap} "返回ConfigMap列表"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/configmap/list [get]
func (h *ConfigMapHandler) GetConfigMapList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")

		res, err := h.svc.GetConfigMapList(context.Background(), ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		cms := make([]resp.ConfigMap, 0, len(res))
		for _, i := range res {
			cms = append(cms, convert.CMConvertListResp(&i))
		}

		response.SuccessWithData(c, cms)
	}
}

// DeleteConfigMap
// @Summary 删除ConfigMap
// @Description 删除指定命名空间下的指定ConfigMap
// @Tags ConfigMap管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "ConfigMap名称"
// @Success 200 {object} response.Response "删除成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/configmap [delete]
func (h *ConfigMapHandler) DeleteConfigMap() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")
		name := c.Query("name")

		err := h.svc.DeleteConfigMap(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}
