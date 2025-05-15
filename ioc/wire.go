//go:build wireinject

package ioc

import (
	"fmt"
	"github.com/crazyfrankie/kube-ctl/internal/metrics"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	promapi "github.com/prometheus/client_golang/api"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/crazyfrankie/kube-ctl/conf"
	"github.com/crazyfrankie/kube-ctl/docs"
	"github.com/crazyfrankie/kube-ctl/internal/api/k8s"
	"github.com/crazyfrankie/kube-ctl/internal/api/mw"
	"github.com/crazyfrankie/kube-ctl/internal/service"
)

type App struct {
	Engine  *gin.Engine
	Metrics *metrics.MetricsHandler
}

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

func InitKubernetesWithDiscovery() *kubernetes.Clientset {
	if isInCluster() {
		cfg, err := rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
		clientSet, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			panic(err)
		}

		return clientSet
	} else {
		return InitKubernetes()
	}
}

func isInCluster() bool {
	tokenFile := "/var/run/secrets/kubernetes.io/serviceaccount/token"

	_, err := os.Stat(tokenFile)
	if err != nil {
		return false
	}

	return true
}

func InitPromAPI() promv1.API {
	client, err := promapi.NewClient(promapi.Config{
		Address: fmt.Sprintf("%s://%s", conf.GetConf().Prom.Scheme, conf.GetConf().Prom.Host),
	})
	if err != nil {
		panic(err)
	}

	return promv1.NewAPI(client)
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
	job *k8s.JobHandler, cron *k8s.CronJobHandler,
	rbac *k8s.RbacHandler, metrics *k8s.MetricsHandler) *gin.Engine {
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
	rbac.RegisterRoute(srv)
	metrics.RegisterRoute(srv)

	srv.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	docs.SwaggerInfo.BasePath = "/api"

	return srv
}

func InitApp() *App {
	wire.Build(
		InitMws,
		InitKubernetesWithDiscovery,
		InitPromAPI,

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
		service.NewRbacService,
		service.NewMetricsService,
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
		k8s.NewRbacHandler,
		k8s.NewMetricsHandler,

		InitGin,
		metrics.NewMetricsHandler,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
