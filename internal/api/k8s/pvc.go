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

type PVCHandler struct {
	svc service.PVCService
}

func NewPVCHandler(svc service.PVCService) *PVCHandler {
	return &PVCHandler{svc: svc}
}

func (h *PVCHandler) RegisterRoute(r *gin.Engine) {
	pvcGroup := r.Group("api/pvc")
	{
		pvcGroup.POST("", h.CreatePVC())
		pvcGroup.DELETE("", h.DeletePVC())
		pvcGroup.GET("", h.GetPVCList())
	}
}

// CreatePVC
// @Summary 创建 PVC
// @Description 创建一个 PersistentVolumeClaim 声明用户的存储需求
// @Tags PVC 管理
// @Accept json
// @Produce json
// @Param pod body req.PersistentVolumeClaim true "PVC 信息"
// @Success 200 {object} response.Response "创建 PVC 成功"
// @Failure 400 {object} response.Response "参数错误(code=20000)"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/pvc [post]
func (h *PVCHandler) CreatePVC() gin.HandlerFunc {
	return func(c *gin.Context) {
		var createReq req.PersistentVolumeClaim
		if err := c.ShouldBind(&createReq); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20000, "bind error "+err.Error()))
			return
		}

		err := h.svc.CreatePVC(context.Background(), &createReq)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// DeletePVC
// @Summary 删除 PVC
// @Description 删除一个 PersistentVolumeClaim 存储声明
// @Tags PVC 管理
// @Accept json
// @Produce json
// @Param name query string true "PVC 名称"
// @Success 200 {object} response.Response "删除 PVC 成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/pvc [delete]
func (h *PVCHandler) DeletePVC() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		err := h.svc.DeletePVC(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// GetPVCList
// @Summary 获取 PVC 列表
// @Description 获取所有 PersistentVolumeClaim 存储声明信息
// @Tags PVC 管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]resp.PersistentVolumeClaim} "获取成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/pvc [get]
func (h *PVCHandler) GetPVCList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")

		res, err := h.svc.GetPVCList(context.Background(), ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		pvcs := make([]resp.PersistentVolumeClaim, 0, len(res))
		for _, i := range res {
			pvcs = append(pvcs, convert.PVCRespConvert(&i))
		}

		response.SuccessWithData(c, pvcs)
	}
}
