package k8s

import (
	"context"
	"fmt"
	"net/http"

	"github.com/crazyfrankie/gem/gerrors"
	"github.com/gin-gonic/gin"

	"github.com/crazyfrankie/kube-ctl/api/model/req"
	"github.com/crazyfrankie/kube-ctl/api/model/resp"
	"github.com/crazyfrankie/kube-ctl/pkg/convert"
	"github.com/crazyfrankie/kube-ctl/pkg/response"
	"github.com/crazyfrankie/kube-ctl/pkg/validate"
	"github.com/crazyfrankie/kube-ctl/service"
)

type PodHandler struct {
	svc service.PodService
}

func NewPodHandler(svc service.PodService) *PodHandler {
	return &PodHandler{svc: svc}
}

func (p *PodHandler) RegisterRoute(r *gin.Engine) {
	podGroup := r.Group("api/pod")
	{
		podGroup.POST("", p.CreateOrUpdatePod())
		podGroup.POST("search", p.SearchPod())
		podGroup.GET("namespace", p.GetNameSpace())
		podGroup.GET("detail", p.GetPod())
		podGroup.GET("list", p.GetPodList())
		podGroup.DELETE("", p.DeletePod())
	}
}

// GetNameSpace
// @Summary 获取命名空间列表
// @Description 获取所有可用的命名空间列表
// @Tags Pod管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]resp.Namespace} "返回命名空间列表，每个命名空间包含名称、创建时间和状态"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/pod/namespace [get]
func (p *PodHandler) GetNameSpace() gin.HandlerFunc {
	return func(c *gin.Context) {
		items, err := p.svc.GetNamespace(context.Background())
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		ns := make([]resp.Namespace, 0, len(items))
		for _, i := range items {
			ns = append(ns, resp.Namespace{
				Name:       i.Name,
				CreateTime: i.CreationTimestamp.Unix(),
				Status:     string(i.Status.Phase),
			})
		}

		response.SuccessWithData(c, ns)
	}
}

// CreateOrUpdatePod
// @Summary 创建或更新Pod
// @Description 创建新的Pod或更新已存在的Pod
// @Tags Pod管理
// @Accept json
// @Produce json
// @Param pod body req.Pod true "Pod配置信息"
// @Success 200 {object} response.Response "操作成功，返回成功消息"
// @Failure 400 {object} response.Response "参数错误(code=20001)或验证错误(code=20002)"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/pod [post]
func (p *PodHandler) CreateOrUpdatePod() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqPod req.Pod
		if err := c.ShouldBind(&reqPod); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20001, "bind error "+err.Error()))
			return
		}

		vd := &validate.PodValidate{}
		err := vd.Validate(&reqPod)
		if err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20002, "validate pod err: "+err.Error()))
			return
		}

		pod := convert.PodReqConvert(&reqPod)

		err = p.svc.CreateOrUpdatePod(context.Background(), pod)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.SuccessWithMsg(c, fmt.Sprintf("Pod[namespace=%s,name=%s] action success", pod.Namespace, pod.Name))
	}
}

// GetPod
// @Summary 获取Pod详情
// @Description 获取指定命名空间下指定Pod的详细信息
// @Tags Pod管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "Pod名称"
// @Success 200 {object} response.Response{data=req.Pod} "返回Pod的详细信息，包含基础信息、卷配置、网络配置、初始化容器和主容器配置"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/pod/detail [get]
func (p *PodHandler) GetPod() gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Query("namespace")
		name := c.Query("name")

		detail, err := p.svc.GetPod(context.Background(), namespace, name)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err)
			return
		}

		pod := convert.PodConvertReq(detail)

		response.SuccessWithData(c, pod)
	}
}

// GetPodList
// @Summary 获取Pod列表
// @Description 获取指定命名空间下的所有Pod列表
// @Tags Pod管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Success 200 {object} response.Response{data=[]resp.PodListItem} "返回Pod列表，每个Pod包含名称、就绪状态、运行状态、重启次数、运行时长、IP和所在节点"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/pod/list [get]
func (p *PodHandler) GetPodList() gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Query("namespace")

		items, err := p.svc.GetPodList(context.Background(), namespace)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		pods := make([]resp.PodListItem, 0, len(items))
		for _, i := range items {
			pods = append(pods, convert.PodListConvertResp(i))
		}

		response.SuccessWithData(c, pods)
	}
}

// DeletePod
// @Summary 删除Pod
// @Description 删除指定命名空间下的指定Pod
// @Tags Pod管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "Pod名称"
// @Success 200 {object} response.Response "删除成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/pod [delete]
func (p *PodHandler) DeletePod() gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Query("namespace")
		name := c.Query("name")

		err := p.svc.DeletePod(context.Background(), namespace, name)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}

// SearchPod
// @Summary 搜索Pod
// @Description 搜索指定命名空间下的指定Pod
// @Tags Pod管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param name query string true "Pod名称"
// @Success 200 {object} response.Response{data=resp.PodListItem} "搜索成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/pod/search [post]
func (p *PodHandler) SearchPod() gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Query("namespace")
		name := c.Query("name")

		res, err := p.svc.SearchPod(context.Background(), namespace, name)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		pod := convert.PodListConvertResp(*res)

		response.SuccessWithData(c, pod)
	}
}
