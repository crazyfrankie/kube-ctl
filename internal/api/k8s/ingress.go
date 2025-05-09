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

type IngressHandler struct {
	svc service.IngressService
}

func NewIngressHandler(svc service.IngressService) *IngressHandler {
	return &IngressHandler{svc: svc}
}

func (h *IngressHandler) RegisterRoute(r *gin.Engine) {
	ingressGroup := r.Group("api/ingress")
	{
		ingressGroup.POST("", h.CreateOrUpdateIngress())
		ingressGroup.DELETE("", h.DeleteIngress())
		ingressGroup.GET("detail", h.GetIngressDetail())
		ingressGroup.GET("list", h.GetIngressList())
	}
}

// CreateOrUpdateIngress
// @Summary 创建或更新 Ingress
// @Description 创建新的 Ingress 或更新已存在的 Ingress
// @Tags Ingress 管理
// @Accept json
// @Produce json
// @Param pod body req.Ingress true "Ingress 配置信息"
// @Success 200 {object} response.Response "操作成功，返回成功消息"
// @Failure 400 {object} response.Response "参数错误(code=20001)"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/ingress [post]
func (h *IngressHandler) CreateOrUpdateIngress() gin.HandlerFunc {
	return func(c *gin.Context) {
		var createReq req.Ingress
		if err := c.ShouldBind(&createReq); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20001, "bind error "+err.Error()))
			return
		}

		err := h.svc.CreateOrUpdateIngress(context.Background(), &createReq)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// DeleteIngress
// @Summary 删除 Ingress
// @Description 删除一个 Ingress
// @Tags Ingress 管理
// @Accept json
// @Produce json
// @Param name query string true "Ingress 名称"
// @Param namespace query string true "命名空间"
// @Success 200 {object} response.Response "删除 Ingress 成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/ingress [delete]
func (h *IngressHandler) DeleteIngress() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		err := h.svc.DeleteIngress(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// GetIngressDetail
// @Summary 获取Ingress详情
// @Description 获取指定命名空间下指定Ingress的详细信息
// @Tags Ingress 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "Ingress 名称"
// @Success 200 {object} response.Response{data=req.Ingress} "返回Ingress的详细信息"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/ingress/detail [get]
func (h *IngressHandler) GetIngressDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		res, err := h.svc.GetIngressDetail(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		ingress := convert.IngressConvertReq(res)

		response.SuccessWithData(c, ingress)
	}
}

// GetIngressList
// @Summary 获取Ingress列表
// @Description 获取指定命名空间下指定Ingress的列表
// @Tags Ingress 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param keyword query string true "关键词"
// @Success 200 {object} response.Response{data=[]resp.Ingress} "返回Ingress的列表"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/ingress/list [get]
func (h *IngressHandler) GetIngressList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")
		keyword := c.Query("keyword")

		res, err := h.svc.GetIngressList(context.Background(), ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		ingresses := make([]resp.Ingress, 0, len(res))
		for _, i := range res {
			if strings.Contains(i.Name, keyword) {
				ingresses = append(ingresses, convert.IngressConvertResp(&i))
			}
		}

		response.SuccessWithData(c, ingresses)
	}
}
