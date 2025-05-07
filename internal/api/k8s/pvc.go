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
