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

type SecretHandler struct {
	svc service.SecretService
}

func NewSecretHandler(svc service.SecretService) *SecretHandler {
	return &SecretHandler{svc: svc}
}

func (h *SecretHandler) RegisterRoute(r *gin.Engine) {
	secretGroup := r.Group("api/secret")
	{
		secretGroup.POST("", h.CreateOrUpdateSecret())
		secretGroup.GET("", h.GetSecret())
		secretGroup.GET("list", h.GetSecretList())
		secretGroup.DELETE("", h.DeleteSecret())
	}
}

// CreateOrUpdateSecret
// @Summary 创建或更新 Secret
// @Description 创建新的 Secret 或更新已存在的 Secret
// @Tags Secret管理
// @Accept json
// @Produce json
// @Param pod body req.Secret true "Secret 配置信息"
// @Success 200 {object} response.Response "操作成功，返回成功消息"
// @Failure 400 {object} response.Response "参数错误(code=20001)"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/secret [post]
func (h *SecretHandler) CreateOrUpdateSecret() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cmReq req.Secret
		if err := c.ShouldBind(&cmReq); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20001, "bind error "+err.Error()))
			return
		}

		err := h.svc.CreateOrUpdateSecret(context.Background(), &cmReq)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// GetSecret
// @Summary 获取Secret详情
// @Description 获取指定命名空间下指定Secret的详细信息
// @Tags Secret管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "Secret名称"
// @Success 200 {object} response.Response{data=resp.SecretDetail} "返回Secret的详细信息"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/secret [get]
func (h *SecretHandler) GetSecret() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")
		name := c.Query("name")

		res, err := h.svc.GetSecret(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		se := convert.SecretConvertDetailResp(res)

		response.SuccessWithData(c, se)
	}
}

// GetSecretList
// @Summary 获取Secret列表
// @Description 获取指定命名空间下的所有Secret列表
// @Tags Secret管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Success 200 {object} response.Response{data=[]resp.Secret} "返回Secret列表"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/secret/list [get]
func (h *SecretHandler) GetSecretList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")

		res, err := h.svc.GetSecretList(context.Background(), ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		secrets := make([]resp.Secret, 0, len(res))
		for _, i := range res {
			secrets = append(secrets, convert.SecretConvertListResp(&i))
		}

		response.SuccessWithData(c, secrets)
	}
}

// DeleteSecret
// @Summary 删除Secret
// @Description 删除指定命名空间下的指定Secret
// @Tags Secret管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "Secret名称"
// @Success 200 {object} response.Response "删除成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/secret [delete]
func (h *SecretHandler) DeleteSecret() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")
		name := c.Query("name")

		err := h.svc.DeleteSecret(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}
