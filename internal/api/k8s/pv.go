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

type PVHandler struct {
	svc service.PVService
}

func NewPVHandler(svc service.PVService) *PVHandler {
	return &PVHandler{svc: svc}
}

func (h *PVHandler) RegisterRoute(r *gin.Engine) {
	pvGroup := r.Group("api/pv")
	{
		pvGroup.POST("", h.CreatePV())
		pvGroup.DELETE("", h.DeletePV())
		pvGroup.GET("", h.GetPVList())
	}
}

// CreatePV
// @Summary 创建 PV
// @Description 创建一个 PersistentVolume 存储空间
// @Tags PV 管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "创建 PV 成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/pv [post]
func (h *PVHandler) CreatePV() gin.HandlerFunc {
	return func(c *gin.Context) {
		var createReq req.PersistentVolume
		if err := c.Bind(&createReq); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20000, "bind error "+err.Error()))
			return
		}

		err := h.svc.CreatePV(context.Background(), &createReq)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// DeletePV
// @Summary 删除 PV
// @Description 删除一个 PersistentVolume 存储空间
// @Tags PV 管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "删除 PV 成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/pv [delete]
func (h *PVHandler) DeletePV() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")

		err := h.svc.DeletePV(context.Background(), name)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// GetPVList
// @Summary 获取 PV 列表
// @Description 获取所有 PersistentVolume 存储空间的信息
// @Tags PV 管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]resp.PersistentVolumeItem} "获取成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/pv [get]
func (h *PVHandler) GetPVList() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := h.svc.GetPVList(context.Background())
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		pvs := make([]resp.PersistentVolumeItem, 0, len(res))
		for _, i := range res {
			pvs = append(pvs, convert.PVConvertResp(&i))
		}

		response.SuccessWithData(c, pvs)
	}
}
