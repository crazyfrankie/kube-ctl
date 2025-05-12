//go:build wireinject

package ioc

import (
	"github.com/crazyfrankie/kube-ctl/docs"
	"github.com/crazyfrankie/kube-ctl/internal/api/k8s"
	"github.com/crazyfrankie/kube-ctl/internal/api/mw"
	"github.com/crazyfrankie/kube-ctl/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

func InitGin(mws []gin.HandlerFunc, pod *k8s.PodHandler, node *k8s.NodeHandler,
	configmap *k8s.ConfigMapHandler, secret *k8s.SecretHandler, pv *k8s.PVHandler,
	pvc *k8s.PVCHandler, storage *k8s.StorageClassHandler,
	svc *k8s.ServiceHandler, ingress *k8s.IngressHandler,
	igRoute *k8s.IngressRouteHandler, deployment *k8s.DeploymentHandler,
	daemon *k8s.DaemonSetHandler, stateful *k8s.StatefulSetHandler,
	job *k8s.JobHandler, cron *k8s.CronJobHandler) *gin.Engine {
	srv := gin.Default()
	srv.Use(mws...)

	pod.RegisterRoute(srv)
	node.RegisterRoute(srv)
	configmap.RegisterRoute(srv)
	secret.RegisterRoute(srv)
	pv.RegisterRoute(srv)
	pvc.RegisterRoute(srv)
	storage.RegisterRoute(srv)
	svc.RegisterRoute(srv)
	ingress.RegisterRoute(srv)
	igRoute.RegisterRoute(srv)
	deployment.RegisterRoute(srv)
	daemon.RegisterRoute(srv)
	stateful.RegisterRoute(srv)
	job.RegisterRoute(srv)
	cron.RegisterRoute(srv)

	srv.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	docs.SwaggerInfo.BasePath = "/api"

	return srv
}

func InitServer() *gin.Engine {
	wire.Build(
		InitMws,
		InitKubernetes,

		service.NewPodService,
		service.NewNodeService,
		service.NewConfigMapService,
		service.NewSecretService,
		service.NewPVService,
		service.NewPVCService,
		service.NewStorageClassService,
		service.NewServiceService,
		service.NewIngressService,
		service.NewIngressRouteService,
		service.NewDeploymentService,
		service.NewDaemonSetService,
		service.NewStatefulSetService,
		service.NewJobService,
		service.NewCronJobService,
		k8s.NewPodHandler,
		k8s.NewNodeHandler,
		k8s.NewConfigMapHandler,
		k8s.NewSecretHandler,
		k8s.NewPVHandler,
		k8s.NewPVCHandler,
		k8s.NewStorageClassHandler,
		k8s.NewServiceHandler,
		k8s.NewIngressHandler,
		k8s.NewIngressRouteHandler,
		k8s.NewDeploymentHandler,
		k8s.NewDaemonSetHandler,
		k8s.NewStatefulSetHandler,
		k8s.NewJobHandler,
		k8s.NewCronJobHandler,

		InitGin,
	)
	return new(gin.Engine)
}
