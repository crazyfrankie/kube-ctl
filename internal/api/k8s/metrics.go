package k8s

import (
	"context"
	"net/http"

	"github.com/crazyfrankie/gem/gerrors"
	"github.com/gin-gonic/gin"

	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/internal/service"
	"github.com/crazyfrankie/kube-ctl/pkg/response"
)

type MetricsHandler struct {
	svc service.MetricsService
}

func NewMetricsHandler(svc service.MetricsService) *MetricsHandler {
	return &MetricsHandler{svc: svc}
}

func (h *MetricsHandler) RegisterRoute(r *gin.Engine) {
	r.GET("api/dashboard", h.GetDashBoard())
}

// GetDashBoard
// @Summary 获取集群Metrics信息
// @Description 获取集群Metrics信息构建Dashboard界面
// @Tags Metrics 管理
// @Accept json
// @Produce json
// @Failure 500 {object} response.Response{data=map[string][]resp.MetricsItem} "获取集群信息失败"
// @Success 200 {object} response.Response "获取集群信息成功"
// @Router /api/dashboard [get]
func (h *MetricsHandler) GetDashBoard() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := make(map[string][]resp.MetricsItem, 4)
		info, err := h.svc.GetClusterBaseInfo(context.Background())
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}
		res["cluster"] = info
		resources, err := h.svc.GetClusterResource(context.Background())
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}
		res["resources"] = resources

		usage, _ := h.svc.GetClusterUsage(context.Background())
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}
		res["usage"] = usage

		usageRange, _ := h.svc.GetClusterUsageRange(context.Background())
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}
		res["usageRange"] = usageRange

		response.SuccessWithData(c, res)
	}
}
