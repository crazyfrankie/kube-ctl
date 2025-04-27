package k8s

import (
	"context"
	"fmt"

	"github.com/crazyfrankie/gem/gerrors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/crazyfrankie/kube-ctl/api/model/req"
	"github.com/crazyfrankie/kube-ctl/api/model/resp"
	"github.com/crazyfrankie/kube-ctl/pkg/convert"
	"github.com/crazyfrankie/kube-ctl/pkg/response"
	"github.com/crazyfrankie/kube-ctl/pkg/validate"
	"github.com/crazyfrankie/kube-ctl/service"
	"github.com/crazyfrankie/kube-ctl/docs"
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
		podGroup.GET("namespace", p.GetNameSpace())
		podGroup.GET("detail", p.GetPod())
		podGroup.GET("list", p.GetPodList())
		podGroup.DELETE("", p.DeletePod())
	}
}

// @Summary Get all namespaces
// @Description Get list of all Kubernetes namespaces
// @Tags namespace
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]resp.Namespace} "Success"
// @Failure 400 {object} response.Response "Error"
// @Router /pod/namespace [get]
func (p *PodHandler) GetNameSpace() gin.HandlerFunc {
	return func(c *gin.Context) {
		items, err := p.svc.GetNamespace(context.Background())
		if err != nil {
			response.Error(c, gerrors.NewBizError(30000, err.Error()))
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

func (p *PodHandler) CreateOrUpdatePod() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqPod req.Pod
		if err := c.Bind(&reqPod); err != nil {
			response.Error(c, gerrors.NewBizError(20001, "bind error "+err.Error()))
			return
		}

		vd := &validate.PodValidate{}
		err := vd.Validate(&reqPod)
		if err != nil {
			response.Error(c, gerrors.NewBizError(20002, "validate pod err: "+err.Error()))
			return
		}

		pod := convert.PodReqConvert(&reqPod)

		err = p.svc.CreateOrUpdatePod(context.Background(), pod)
		if err != nil {
			response.Error(c, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.SuccessWithMsg(c, fmt.Sprintf("Pod[namespace=%s,name=%s] action success", pod.Namespace, pod.Name))
	}
}

func (p *PodHandler) GetPod() gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Query("namespace")
		name := c.Query("name")

		detail, err := p.svc.GetPod(context.Background(), namespace, name)
		if err != nil {
			response.Error(c, err)
			return
		}

		pod := convert.PodConvertReq(detail)

		response.SuccessWithData(c, pod)
	}
}

func (p *PodHandler) GetPodList() gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Query("namespace")

		items, err := p.svc.GetPodList(context.Background(), namespace)
		if err != nil {
			response.Error(c, gerrors.NewBizError(30000, err.Error()))
			return
		}

		pods := make([]resp.PodListItem, 0, len(items))
		for _, i := range items {
			pods = append(pods, convert.PodListConvertResp(i))
		}

		response.SuccessWithData(c, pods)
	}
}

func (p *PodHandler) DeletePod() gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Query("namespace")
		name := c.Query("name")

		err := p.svc.DeletePod(context.Background(), namespace, name)
		if err != nil {
			response.Error(c, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.Success(c)
	}
}
