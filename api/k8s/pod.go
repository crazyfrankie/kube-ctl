package k8s

import (
	"context"
	es "errors"
	"fmt"
	"github.com/crazyfrankie/gem/gerrors"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"strings"
	"time"

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
		podGroup.POST("", p.CreateOrUpdatePod())
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

		if get, err := p.clientSet.CoreV1().Pods(pod.Namespace).
			Get(c.Request.Context(), pod.Name, metav1.GetOptions{}); err == nil {
			// Verify that the parameters are legal
			cPod := *pod
			cPod.Name = cPod.Name + "-validate"
			_, err := p.clientSet.CoreV1().Pods(cPod.Namespace).Create(c.Request.Context(),
				&cPod, metav1.CreateOptions{DryRun: []string{metav1.DryRunAll}})
			if err != nil {
				response.Error(c, err)
				return
			}

			// Delete the Pod
			err = p.clientSet.CoreV1().Pods(pod.Namespace).Delete(c.Request.Context(), pod.Name, metav1.DeleteOptions{})
			if err != nil {
				response.Error(c, err)
				return
			}

			// Listen for deletion events
			labels := make([]string, 0, len(get.Labels))
			for k, v := range get.Labels {
				labels = append(labels, fmt.Sprintf("%s=%s", k, v))
			}

			ctx, cancel := context.WithTimeout(c.Request.Context(), 50*time.Second)
			defer cancel()
			ch, err := p.clientSet.CoreV1().Pods(pod.Namespace).Watch(ctx, metav1.ListOptions{
				LabelSelector: strings.Join(labels, ","),
			})
			if err != nil {
				response.Error(c, err)
				return
			}

			for event := range ch.ResultChan() {
				chPod := event.Object.(*corev1.Pod)

				// Fast paths, some Pods may be deleted quickly,
				// causing the listener to not start yet and subsequently keep blocking.
				// Query if the event has been deleted,
				// if it has been deleted, then you don't need to listen to the delete event.
				if _, err := p.clientSet.CoreV1().Pods(pod.Namespace).
					Get(c.Request.Context(), pod.Name, metav1.GetOptions{}); errors.IsNotFound(err) {
					// Delete successful, create new Pod
					newPod, err := p.clientSet.CoreV1().Pods(pod.Namespace).Create(c.Request.Context(),
						pod, metav1.CreateOptions{})
					if err != nil {
						msg := es.New(fmt.Sprintf("failed update pod, name: %s, %s", newPod.Name, err.Error()))
						response.Error(c, msg)
						return
					} else {
						response.SuccessWithMsg(c, fmt.Sprintf("Pod[namespace=%s,name=%s] updated success", pod.Namespace, pod.Name))
						return
					}
				}

				switch event.Type {
				case watch.Deleted:
					if chPod.Name != pod.Name {
						continue
					}

					// Delete successful, create new Pod
					newPod, err := p.clientSet.CoreV1().Pods(pod.Namespace).Create(c.Request.Context(),
						pod, metav1.CreateOptions{})
					if err != nil {
						msg := es.New(fmt.Sprintf("failed update pod, name: %s, %s", newPod.Name, err.Error()))
						response.Error(c, msg)
						return
					} else {
						response.SuccessWithMsg(c, fmt.Sprintf("Pod[namespace=%s,name=%s] updated success", pod.Namespace, pod.Name))
						return
					}
				}
			}

			select {
			case <-ctx.Done():
				response.Error(c, fmt.Errorf("timeout waiting for pod deletion"))
				return
			default:
			}
		}

		_, err = p.clientSet.CoreV1().Pods(pod.Namespace).Create(c.Request.Context(),
			pod, metav1.CreateOptions{})
		if err != nil {
			msg := es.New(fmt.Sprintf("failed create pod, name: %s, %s", pod.Name, err.Error()))
			response.Error(c, msg)
			return
		}

		response.SuccessWithMsg(c, fmt.Sprintf("Pod[namespace=%s,name=%s] created success", pod.Namespace, pod.Name))
	}
}

func (p *PodHandler) GetPodList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ns := c.Query("namespace")

		pods, err := p.clientSet.CoreV1().Pods(ns).List(c.Request.Context(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

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

		response.SuccessWithData(c, pods.Items)
	}
}
