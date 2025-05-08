package k8s

import (
	"context"
	"github.com/crazyfrankie/kube-ctl/internal/model/validate"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/crazyfrankie/gem/gerrors"
	"github.com/crazyfrankie/kube-ctl/internal/model/convert"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/internal/service"
	"github.com/crazyfrankie/kube-ctl/pkg/response"
)

type StorageClassHandler struct {
	svc service.StorageClassService
}

func NewStorageClassHandler(svc service.StorageClassService) *StorageClassHandler {
	return &StorageClassHandler{svc: svc}
}

func (h *StorageClassHandler) RegisterRoute(r *gin.Engine) {
	scGroup := r.Group("api/storage")
	{
		scGroup.POST("", h.CreateStorageClass())
		scGroup.DELETE("", h.DeleteStorageClass())
		scGroup.GET("", h.GetStorageClassList())
	}
}

// CreateStorageClass
// @Summary 创建 StorageClass
// @Description 创建一个 StorageClas 存储类
// @Tags StorageClass 管理
// @Accept json
// @Produce json
// @Param pod body req.StorageClass true "StorageClass 信息"
// @Success 200 {object} response.Response "创建 StorageClass 成功"
// @Failure 400 {object} response.Response "参数错误(code=20000或20001)"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/storage [post]
func (h *StorageClassHandler) CreateStorageClass() gin.HandlerFunc {
	return func(c *gin.Context) {
		var createReq req.StorageClass
		if err := c.ShouldBind(&createReq); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20000, "bind error "+err.Error()))
			return
		}

		err := validate.StorageClassValidate(&createReq)
		if err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20001, err.Error()))
			return
		}

		err = h.svc.CreateStorageClass(context.Background(), &createReq)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// DeleteStorageClass
// @Summary 删除 StorageClass
// @Description 删除一个 StorageClass 存储类
// @Tags StorageClass 管理
// @Accept json
// @Produce json
// @Param name query string true "StorageClass 名称"
// @Success 200 {object} response.Response "删除 StorageClass 成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/storage [delete]
func (h *StorageClassHandler) DeleteStorageClass() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")

		err := h.svc.DeleteStorageClass(context.Background(), name)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// GetStorageClassList
// @Summary 获取 StorageClass 列表
// @Description 获取所有 StorageClass 存储类信息
// @Tags PVC 管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]resp.StorageClass} "获取成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/storage [get]
func (h *StorageClassHandler) GetStorageClassList() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := h.svc.GetStorageClassList(context.Background())
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		scs := make([]resp.StorageClass, 0, len(res))
		for _, i := range res {
			scs = append(scs, convert.StorageClassConvertResp(&i))
		}

		response.SuccessWithData(c, scs)
	}
}
