package metrics

import (
	"context"
	"strconv"

	"github.com/crazyfrankie/kube-ctl/internal/service"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricsHandler struct {
	clusterCpu prometheus.Gauge
	clusterMem prometheus.Gauge
	svc        service.MetricsService
}

func NewMetricsHandler(svc service.MetricsService) *MetricsHandler {
	return &MetricsHandler{
		clusterCpu: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "cluster_cpu",
				Help: "collector cluster cpu info",
			}),
		clusterMem: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "cluster_mem",
				Help: "collector cluster memory info",
			}),
		svc: svc,
	}
}

func (h *MetricsHandler) Describe(descs chan<- *prometheus.Desc) {
	h.clusterCpu.Describe(descs)
	h.clusterMem.Describe(descs)
}

func (h *MetricsHandler) Collect(metrics chan<- prometheus.Metric) {
	usageArr, err := h.svc.GetClusterUsage(context.Background())
	if err != nil {
		return
	}
	for _, item := range usageArr {
		switch item.Label {
		case "cluster_cpu":
			newValue, _ := strconv.ParseFloat(item.Value, 64)
			h.clusterCpu.Set(newValue)
			h.clusterCpu.Collect(metrics)
		case "cluster_mem":
			newValue, _ := strconv.ParseFloat(item.Value, 64)
			h.clusterMem.Set(newValue)
			h.clusterMem.Collect(metrics)
		}
	}
}
