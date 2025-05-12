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

type JobHandler struct {
	svc service.JobService
}

func NewJobHandler(svc service.JobService) *JobHandler {
	return &JobHandler{svc: svc}
}

func (h *JobHandler) RegisterRoute(r *gin.Engine) {
	jobGroup := r.Group("api/job")
	{
		jobGroup.POST("", h.CreateOrUpdateJob())
		jobGroup.DELETE("", h.DeleteJob())
		jobGroup.GET("", h.GetJobDetail())
		jobGroup.GET("list", h.GetJobList())
	}
}

// CreateOrUpdateJob
// @Summary 创建或更新 Job
// @Description 创建新的 Job 或更新已存在的 Job
// @Tags Job 管理
// @Accept json
// @Produce json
// @Param pod body req.Job true "Job 配置信息"
// @Success 200 {object} response.Response "操作成功，返回成功消息"
// @Failure 400 {object} response.Response "参数错误(code=20001)"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/job [post]
func (h *JobHandler) CreateOrUpdateJob() gin.HandlerFunc {
	return func(c *gin.Context) {
		var creatReq req.Job
		if err := c.ShouldBind(&creatReq); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20001, "bind error "+err.Error()))
			return
		}

		err := h.svc.CreateOrUpdateJob(context.Background(), &creatReq)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// DeleteJob
// @Summary 删除 Job
// @Description 删除指定命名空间下的指定Job
// @Tags Job 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "Job 名称"
// @Success 200 {object} response.Response "删除成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/job [delete]
func (h *JobHandler) DeleteJob() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		err := h.svc.DeleteJob(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// GetJobDetail
// @Summary 获取Job详情
// @Description 获取指定命名空间下指定Job的详细信息
// @Tags Job 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "Job 名称"
// @Success 200 {object} response.Response{data=req.Job} "返回Job的详细信息"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/job [get]
func (h *JobHandler) GetJobDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		ns := c.Query("namespace")

		res, err := h.svc.GetJobDetail(context.Background(), name, ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		job := convert.JobConvertReq(res)

		response.SuccessWithData(c, job)
	}
}

// GetJobList
// @Summary 获取Job列表
// @Description 获取指定命名空间下的所有Job列表
// @Tags Job 管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param keyword query string true "关键词"
// @Success 200 {object} response.Response{data=[]resp.Job} "返回Job列表"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/job/list [get]
func (h *JobHandler) GetJobList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")
		keyword := c.Query("keyword")

		res, err := h.svc.GetJobList(context.Background(), ns)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		jobs := make([]resp.Job, 0, len(res))
		for _, j := range res {
			if strings.Contains(j.Name, keyword) {
				jobs = append(jobs, convert.JobConvertResp(&j))
			}
		}

		response.SuccessWithData(c, jobs)
	}
}
