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

type CronJobHandler struct {
	svc service.CronJobService
}

func NewCronJobHandler(svc service.CronJobService) *CronJobHandler {
	return &CronJobHandler{svc: svc}
}

func (h *CronJobHandler) RegisterRoute(r *gin.Engine) {
	cronCronJobGroup := r.Group("api/cronCronJob")
	{
		cronCronJobGroup.POST("", h.CreateOrUpdateCronJob())
		cronCronJobGroup.DELETE("", h.DeleteCronJob())
		cronCronJobGroup.GET("", h.GetCronJobDetail())
		cronCronJobGroup.GET("list", h.GetCronJobList())
	}
}

// CreateOrUpdateCronJob
// @Summary 创建或更新 CronJob
// @Description 创建新的 CronJob 或更新已存在的 CronJob
// @Tags CronJob 管理
// @Accept json
// @Produce json
// @Param pod body req.CronJob true "CronJob 配置信息"
// @Success 200 {object} response.Response "操作成功，返回成功消息"
// @Failure 400 {object} response.Response "参数错误(code=20001)"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/cronjob [post]
func (h *CronJobHandler) CreateOrUpdateCronJob() gin.HandlerFunc {
	return func(c *gin.Context) {
		var creatReq req.CronJob
		if err := c.ShouldBind(&creatReq); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20001, "bind error "+err.Error()))
			return
		}

		err := h.svc.CreateOrUpdateCronJob(context.Background(), &creatReq)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// DeleteCronJob
// @Summary 删除 CronJob
// @Description 删除指定命名空间下的指定CronJob
// @Tags CronJob 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "CronJob 名称"
// @Success 200 {object} response.Response "删除成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/cronjob [delete]
func (h *CronJobHandler) DeleteCronJob() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		err := h.svc.DeleteCronJob(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// GetCronJobDetail
// @Summary 获取CronJob详情
// @Description 获取指定命名空间下指定CronJob的详细信息
// @Tags CronJob 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "CronJob 名称"
// @Success 200 {object} response.Response{data=req.CronJob} "返回CronJob的详细信息"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/cronjob [get]
func (h *CronJobHandler) GetCronJobDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		res, err := h.svc.GetCronJobDetail(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		cronJob := convert.CronJobConvertReq(res)

		response.SuccessWithData(c, cronJob)
	}
}

// GetCronJobList
// @Summary 获取 CronJob 列表
// @Description 获取指定命名空间下的所有CronJob列表
// @Tags CronJob 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param keyword query string true "关键词"
// @Success 200 {object} response.Response{data=[]resp.CronJob} "返回CronJob列表"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/cronjob/list [get]
func (h *CronJobHandler) GetCronJobList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")
		keyword := c.Query("keyword")

		res, err := h.svc.GetCronJobList(context.Background(), ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		cronJobs := make([]resp.CronJob, 0, len(res))
		for _, j := range res {
			if strings.Contains(j.Name, keyword) {
				cronJobs = append(cronJobs, convert.CronJobConvertResp(&j))
			}
		}

		response.SuccessWithData(c, cronJobs)
	}
}
