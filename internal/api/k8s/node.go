package k8s

import (
	"context"
	"net/http"
	"strings"

	"github.com/crazyfrankie/gem/gerrors"
	"github.com/gin-gonic/gin"

	"github.com/crazyfrankie/kube-ctl/internal/api/model/resp"
	"github.com/crazyfrankie/kube-ctl/internal/service"
	"github.com/crazyfrankie/kube-ctl/pkg/convert"
	"github.com/crazyfrankie/kube-ctl/pkg/response"
)

type NodeHandler struct {
	svc service.NodeService
}

func NewNodeHandler(svc service.NodeService) *NodeHandler {
	return &NodeHandler{svc: svc}
}

func (n *NodeHandler) RegisterRoute(r *gin.Engine) {
	nodeGroup := r.Group("api/node")
	{
		nodeGroup.GET("list", n.NodeList())
		nodeGroup.GET("detail", n.NodeDetail())
	}
}

// NodeList
// @Summary 获取 Node 列表
// @Description 获取集群中所有 Node 信息
// @Tags Node管理
// @Accept json
// @Produce json
// @Param keyword query string true "node 关键词"
// @Success 200 {object} response.Response{data=[]resp.NodeListItem} "获取成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
func (n *NodeHandler) NodeList() gin.HandlerFunc {
	return func(c *gin.Context) {
		keyword := c.Query("keyword")

		list, err := n.svc.NodeList(context.Background())
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		res := make([]resp.NodeListItem, 0, len(list))
		for _, n := range list {
			if strings.Contains(n.Name, keyword) {
				res = append(res, convert.NodeListItemConvertResp(n))
			}
		}

		response.SuccessWithData(c, res)
	}
}

// NodeDetail
// @Summary 获取 Node 详情
// @Description 获取集群中单个 Node 信息
// @Tags Node管理
// @Accept json
// @Produce json
// @Param keyword query string true "node name"
// @Success 200 {object} response.Response{data=resp.NodeDetail} "获取成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
func (n *NodeHandler) NodeDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")

		res, err := n.svc.NodeDetail(context.Background(), name)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		node := convert.NodeDetailConvertResp(res)

		response.SuccessWithData(c, node)
	}
}
