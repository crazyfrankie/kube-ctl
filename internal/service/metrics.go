package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/pkg/utils"
)

type MetricsService interface {
	GetClusterBaseInfo(ctx context.Context) ([]resp.MetricsItem, error)
	GetClusterResource(ctx context.Context) ([]resp.MetricsItem, error)
	GetClusterUsage(ctx context.Context) ([]resp.MetricsItem, error)
	GetClusterUsageRange(ctx context.Context) ([]resp.MetricsItem, error)
}

type metricsService struct {
	clientSet *kubernetes.Clientset
	promApi   promv1.API
}

func NewMetricsService(cs *kubernetes.Clientset, promApi promv1.API) MetricsService {
	return &metricsService{clientSet: cs, promApi: promApi}
}

func (s *metricsService) GetClusterBaseInfo(ctx context.Context) ([]resp.MetricsItem, error) {
	metrics := make([]resp.MetricsItem, 0, 3)

	// cluster info
	metrics = append(metrics, resp.MetricsItem{
		Title: "Cluster",
		Value: "Kubernetes",
	})

	// version info
	ver, err := s.clientSet.ServerVersion()
	if err != nil {
		return nil, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "Kubernetes Version",
		Value: fmt.Sprintf("%s.%s", ver.Major, ver.Minor),
	})

	// node info
	list, err := s.clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "Nodes",
		Value: strconv.Itoa(len(list.Items)),
	})

	// cluster init time
	var ctime time.Time
	for _, item := range list.Items {
		for k, _ := range item.Labels {
			if k == "node-role.kubernetes.io/control-plane" {
				if ctime.IsZero() {
					ctime = item.CreationTimestamp.Time
				}
				if !ctime.IsZero() && ctime.Sub(item.CreationTimestamp.Time) > 0 {
					ctime = item.CreationTimestamp.Time
				}
			}
		}
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "Created",
		Value: utils.FormatTime(time.Now().Sub(ctime)),
	})

	for i, m := range metrics {
		metrics[i].Color = utils.GenerateHashBaseRGB(m.Title)
	}

	return metrics, nil
}

func (s *metricsService) GetClusterResource(ctx context.Context) ([]resp.MetricsItem, error) {
	metrics := make([]resp.MetricsItem, 0, 18)

	// namespace
	ns, err := s.clientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "Namespace",
		Value: strconv.Itoa(len(ns.Items)),
	})
	// pod
	pods, err := s.clientSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "Pod",
		Value: strconv.Itoa(len(pods.Items)),
	})
	// configMap
	cm, err := s.clientSet.CoreV1().ConfigMaps("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "ConfigMap",
		Value: strconv.Itoa(len(cm.Items)),
	})
	// secret
	se, err := s.clientSet.CoreV1().Secrets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "Secret",
		Value: strconv.Itoa(len(se.Items)),
	})
	// pv
	pv, err := s.clientSet.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "PersistentVolume",
		Value: strconv.Itoa(len(pv.Items)),
	})
	// pvc
	pvc, err := s.clientSet.CoreV1().PersistentVolumeClaims("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "PersistentVolumeClaim",
		Value: strconv.Itoa(len(pvc.Items)),
	})
	// storageClass
	sc, err := s.clientSet.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "StorageClass",
		Value: strconv.Itoa(len(sc.Items)),
	})
	// service
	svc, err := s.clientSet.CoreV1().Services("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "Service",
		Value: strconv.Itoa(len(svc.Items)),
	})
	// ingress
	ingress, err := s.clientSet.NetworkingV1().Ingresses("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "Ingress",
		Value: strconv.Itoa(len(ingress.Items)),
	})
	// deployment
	deploy, err := s.clientSet.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "Deployment",
		Value: strconv.Itoa(len(deploy.Items)),
	})
	// daemonSet
	daemon, err := s.clientSet.AppsV1().DaemonSets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "DaemonSet",
		Value: strconv.Itoa(len(daemon.Items)),
	})
	// statefulSet
	state, err := s.clientSet.AppsV1().StatefulSets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "StatefulSet",
		Value: strconv.Itoa(len(state.Items)),
	})
	// job
	job, err := s.clientSet.BatchV1().Jobs("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "Job",
		Value: strconv.Itoa(len(job.Items)),
	})
	// cronJob
	cron, err := s.clientSet.BatchV1().CronJobs("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "CronJob",
		Value: strconv.Itoa(len(cron.Items)),
	})
	// serviceAccount
	sa, err := s.clientSet.CoreV1().ServiceAccounts("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "ServiceAccount",
		Value: strconv.Itoa(len(sa.Items)),
	})
	// roles
	roles, err := s.clientSet.RbacV1().Roles("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "Role",
		Value: strconv.Itoa(len(roles.Items)),
	})
	// clusterRole
	cr, err := s.clientSet.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "ClusterRole",
		Value: strconv.Itoa(len(cr.Items)),
	})
	// roleBinding
	rb, err := s.clientSet.RbacV1().RoleBindings("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "RoleBinding",
		Value: strconv.Itoa(len(rb.Items)),
	})
	// clusterRoleBinding
	crb, err := s.clientSet.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "ClusterRoleBinding",
		Value: strconv.Itoa(len(crb.Items)),
	})

	for i, metric := range metrics {
		metrics[i].Color = utils.GenerateHashBaseRGB(metric.Title)
	}

	return metrics, nil
}

func (s *metricsService) GetClusterUsage(ctx context.Context) ([]resp.MetricsItem, error) {
	metrics := make([]resp.MetricsItem, 0, 3)

	url := "/apis/metrics.k8s.io/v1beta1/nodes"

	raw, err := s.clientSet.RESTClient().Get().AbsPath(url).DoRaw(ctx)
	if err != nil {
		return metrics, nil
	}

	var nodeMetrics resp.NodeMetricsList
	err = sonic.Unmarshal(raw, &nodeMetrics)
	if err != nil {
		return metrics, err
	}

	nodes, err := s.clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}

	pods, err := s.clientSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return metrics, err
	}

	var cpuUsage, cpuTotal int64
	var memUsage, memTotal int64
	var podUsage, podTotal int64
	podUsage = int64(len(pods.Items))
	count := 0
	for i, item := range nodes.Items {
		if len(nodeMetrics.Items) != count {
			cpuUsage += nodeMetrics.Items[i].Usage.Cpu().Value()
			memUsage += nodeMetrics.Items[i].Usage.Memory().Value()
			count++
		}
		cpuTotal += item.Status.Capacity.Cpu().Value()
		memTotal += item.Status.Capacity.Memory().Value()
		podTotal += item.Status.Capacity.Pods().Value()
	}

	metrics = append(metrics, resp.MetricsItem{
		Title: "Pod proportion",
		Value: fmt.Sprintf("%.2f", (float64(podUsage)/float64(podTotal))*100),
	})
	metrics = append(metrics, resp.MetricsItem{
		Title: "CPU proportion",
		Value: fmt.Sprintf("%.2f", (float64(cpuUsage)/float64(cpuTotal))*100),
		Label: "cluster_cpu",
	})
	metrics = append(metrics, resp.MetricsItem{
		Title: "Memory proportion",
		Value: fmt.Sprintf("%.2f", (float64(memUsage)/float64(memTotal))*100),
		Label: "cluster_mem",
	})

	return metrics, nil
}

func (s *metricsService) GetClusterUsageRange(ctx context.Context) ([]resp.MetricsItem, error) {
	metrics := make([]resp.MetricsItem, 0, 2)

	cpu, err := s.getMetricsFromProm("cluster_cpu")
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "CPU changing trend",
		Value: cpu,
	})
	mem, err := s.getMetricsFromProm("cluster_mem")
	if err != nil {
		return metrics, err
	}
	metrics = append(metrics, resp.MetricsItem{
		Title: "Memory changing trend",
		Value: mem,
	})

	return metrics, nil
}

func (s *metricsService) getMetricsFromProm(metricName string) (string, error) {
	resultMap := make(map[string][]string)
	now := time.Now()
	start, end := now.Add(-time.Hour*24), now
	rg := promv1.Range{
		Start: start,
		End:   end,
		Step:  5 * time.Minute,
	}

	queryRange, _, err := s.promApi.QueryRange(context.TODO(), metricName, rg)
	if err != nil {
		return "", err
	}
	matrix := queryRange.(model.Matrix)
	if len(matrix) == 0 {
		err = fmt.Errorf("prometheus query data is null")
		return "", err
	}

	beijingLoc, _ := time.LoadLocation("Asia/Shanghai")
	x := make([]string, 0)
	y := make([]string, 0)
	for _, value := range matrix[0].Values {
		beijingTime := value.Timestamp.Time().In(beijingLoc)
		format := beijingTime.Format("15:04")
		x = append(x, format)
		y = append(y, value.Value.String())
	}
	resultMap["x"] = x
	resultMap["y"] = y
	raw, _ := sonic.Marshal(resultMap)
	return string(raw), nil
}
