package k8s

import (
	"context"
	"fmt"
	
	"github.com/crazyfrankie/kube-ctl/pkg/response"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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
		podGroup.GET("list", p.GetPodList())
	}
}

func (p *PodHandler) GetPodList() gin.HandlerFunc {
	return func(c *gin.Context) {
		pods, err := p.clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		for _, p := range pods.Items {
			fmt.Println(p.Name)
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
