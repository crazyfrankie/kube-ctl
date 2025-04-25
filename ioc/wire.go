//go:build wireinject

package ioc

import (
	"github.com/crazyfrankie/kube-ctl/api/k8s"
	"github.com/crazyfrankie/kube-ctl/api/mw"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func InitKubernetes() *kubernetes.Clientset {
	kubeConfig := ".kube/config"
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientSet
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientSet
}

func InitMws() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		mw.CORS(),
	}
}

func InitGin(pod *k8s.PodHandler, mws []gin.HandlerFunc) *gin.Engine {
	srv := gin.Default()
	srv.Use(mws...)
	pod.RegisterRoute(srv)

	return srv
}

func InitServer() *gin.Engine {
	wire.Build(
		InitMws,
		InitKubernetes,

		k8s.NewPodHandler,

		InitGin,
	)
	return new(gin.Engine)
}
