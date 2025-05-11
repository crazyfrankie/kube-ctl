package k8s

import (
	"context"
	"github.com/crazyfrankie/gem/gerrors"
	"github.com/crazyfrankie/kube-ctl/internal/model/convert"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/internal/service"
	"github.com/crazyfrankie/kube-ctl/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type DeploymentHandler struct {
	svc service.DeploymentService
}

func NewDeploymentHandler(svc service.DeploymentService) *DeploymentHandler {
	return &DeploymentHandler{svc: svc}
}

func (h *DeploymentHandler) RegisterRoute(r *gin.Engine) {
	deploymentGroup := r.Group("api/deployment")
	{
		deploymentGroup.POST("", h.CreateOrUpdateDeployment())
		deploymentGroup.DELETE("", h.DeleteDeployment())
		deploymentGroup.GET("detail", h.GetDeploymentDetail())
		deploymentGroup.GET("list", h.GetDeploymentList())
	}
}

// CreateOrUpdateDeployment
// @Summary 创建或更新 Deployment
// @Description 创建新的 Deployment 或更新已存在的 Deployment
// @Tags Deployment 管理
// @Accept json
// @Produce json
// @Param pod body req.Deployment true "Deployment 配置信息"
// @Success 200 {object} response.Response "操作成功，返回成功消息"
// @Failure 400 {object} response.Response "参数错误(code=20001)"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/deployment [post]
func (h *DeploymentHandler) CreateOrUpdateDeployment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var createReq req.Deployment
		if err := c.ShouldBind(&createReq); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20001, "bind error "+err.Error()))
			return
		}

		err := h.svc.CreateOrUpdateDeployment(context.Background(), &createReq)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// DeleteDeployment
// @Summary 删除 Deployment
// @Description 删除指定命名空间下的指定Deployment
// @Tags Deployment 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "Deployment 名称"
// @Success 200 {object} response.Response "删除成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/deployment [delete]
func (h *DeploymentHandler) DeleteDeployment() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		err := h.svc.DeleteDeployment(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// GetDeploymentDetail
// @Summary 获取Deployment详情
// @Description 获取指定命名空间下指定Deployment的详细信息
// @Tags Deployment 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "Deployment 名称"
// @Success 200 {object} response.Response{data=req.Deployment} "返回Deployment的详细信息"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/deployment [get]
func (h *DeploymentHandler) GetDeploymentDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		res, err := h.svc.GetDeploymentDetail(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		deploy := convert.DeploymentConvertReq(res)

		response.SuccessWithData(c, deploy)
	}
}

// GetDeploymentList
// @Summary 获取Deployment列表
// @Description 获取指定命名空间下的所有Deployment列表
// @Tags Deployment 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param keyword query string true "关键词"
// @Success 200 {object} response.Response{data=[]resp.Deployment} "返回Deployment列表"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/deployment/list [get]
func (h *DeploymentHandler) GetDeploymentList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")
		keyword := c.Query("keyword")

		res, err := h.svc.GetDeploymentList(context.Background(), ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		deploys := make([]resp.Deployment, 0, len(res))
		for _, d := range res {
			if strings.Contains(d.Name, keyword) {
				deploys = append(deploys, convert.DeploymentConvertResp(d))
			}
		}

		response.SuccessWithData(c, deploys)
	}
}
