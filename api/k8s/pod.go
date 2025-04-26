package k8s

import (
	"context"
	es "errors"
	"fmt"

	"github.com/crazyfrankie/gem/gerrors"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/crazyfrankie/kube-ctl/api/model/req"
	"github.com/crazyfrankie/kube-ctl/api/model/resp"
	"github.com/crazyfrankie/kube-ctl/pkg/convert"
	"github.com/crazyfrankie/kube-ctl/pkg/response"
	"github.com/crazyfrankie/kube-ctl/pkg/validate"
)

type PodHandler struct {
	clientSet *kubernetes.Clientset
}

func NewPodHandler(cs *kubernetes.Clientset) *PodHandler {
	return &PodHandler{clientSet: cs}
}

func (p *PodHandler) RegisterRoute(r *gin.Engine) {
	podGroup := r.Group("api/pod")
	{
		podGroup.POST("", p.CreatePod())
		podGroup.GET("namespace", p.GetNameSpace())
		podGroup.GET("list", p.GetPodList())
	}
}

func (p *PodHandler) GetNameSpace() gin.HandlerFunc {
	return func(c *gin.Context) {
		list, err := p.clientSet.CoreV1().Namespaces().List(c.Request.Context(), metav1.ListOptions{})
		if err != nil {
			response.Error(c, err)
			return
		}

		ns := make([]resp.Namespace, 0, len(list.Items))
		for _, i := range list.Items {
			ns = append(ns, resp.Namespace{
				Name:       i.Name,
				CreateTime: i.CreationTimestamp.Unix(),
				Status:     string(i.Status.Phase),
			})
		}

		response.SuccessWithData(c, "OK", ns)
	}
}

func (p *PodHandler) CreatePod() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pod req.Pod
		if err := c.Bind(&pod); err != nil {
			response.Error(c, gerrors.NewBizError(20001, "bind error "+err.Error()))
			return
		}

		vd := &validate.PodValidate{}
		err := vd.Validate(&pod)
		if err != nil {
			response.Error(c, gerrors.NewBizError(20002, "validate pod err: "+err.Error()))
			return
		}

		pd, err := p.clientSet.CoreV1().Pods(pod.Base.Namespace).Create(c.Request.Context(),
			convert.PodReqConvert(&pod), metav1.CreateOptions{})
		if err != nil {
			msg := es.New(fmt.Sprintf("failed create pod, name: %s, %s", pod.Base.Name, err.Error()))
			response.Error(c, msg)
			return
		}

		response.SuccessWithData(c, fmt.Sprintf("Pod[%s-%s] created success", pod.Base.Name, pod.Base.Name), pd)
	}
}

func (p *PodHandler) GetPodList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")

		pods, err := p.clientSet.CoreV1().Pods(ns).List(c.Request.Context(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		namespace := "default"
		pod := "example-xxxxx"
		_, err = p.clientSet.CoreV1().Pods(namespace).Get(context.TODO(), pod, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %s in namespace %s: %v\n",
				pod, namespace, statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
		}

		response.SuccessWithData(c, "OK", pods.Items)
	}
}
