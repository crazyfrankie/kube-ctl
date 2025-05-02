package k8s

import (
	"context"
	"net/http"
	"strings"

	"github.com/crazyfrankie/gem/gerrors"
	"github.com/gin-gonic/gin"

	"github.com/crazyfrankie/kube-ctl/internal/model/convert"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/internal/service"
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
		nodeGroup.PUT("label", n.UpdateNodeLabel())
		nodeGroup.PUT("taint", n.UpdateNodeTaint())
		nodeGroup.GET("pods", n.GetNodePods())
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
// @Router /api/node/list [get]
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
// @Param name query string true "node name"
// @Success 200 {object} response.Response{data=resp.NodeDetail} "获取成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/node/detail [get]
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

// UpdateNodeLabel
// @Summary Node 添加标签
// @Description 为单个 Node 添加 label
// @Tags Node管理
// @Accept json
// @Produce json
// @Param node body req.UpdateLabelReq true "node name and labels"
// @Success 200 {object} response.Response "更新成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/node/label [put]
func (n *NodeHandler) UpdateNodeLabel() gin.HandlerFunc {
	return func(c *gin.Context) {
		var updateReq req.UpdateLabelReq
		if err := c.ShouldBind(&updateReq); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20001, "bind error "+err.Error()))
			return
		}

		err := n.svc.UpdateNodeLabel(context.Background(), updateReq)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.SuccessWithMsg(c, "update label successfully")
	}
}

// UpdateNodeTaint
// @Summary Node 污点更新
// @Description 集群中某个节点的污点更新
// @Tags Node管理
// @Accept json
// @Produce json
// @Param node body req.UpdateTaintReq true "node name and taints"
// @Success 200 {object} response.Response "更新成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/node/taint [put]
func (n *NodeHandler) UpdateNodeTaint() gin.HandlerFunc {
	return func(c *gin.Context) {
		var updateReq req.UpdateTaintReq
		if err := c.Bind(&updateReq); err != nil {
			response.Error(c, http.StatusBadRequest, gerrors.NewBizError(20001, "bind error "+err.Error()))
			return
		}

		err := n.svc.UpdateNodeTaints(context.Background(), updateReq)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		response.SuccessWithMsg(c, "update taints successfully")
	}
}

// GetNodePods
// @Summary Node 下的所有 Pod
// @Description 集群中某个节点下的所有 Pod
// @Tags Node管理
// @Accept json
// @Produce json
// @Param namespace query string true "命名空间"
// @Param node query string true "Node 名称"
// @Success 200 {object} response.Response{data=[]resp.PodListItem} "查询成功"
// @Failure 500 {object} response.Response "系统错误(code=30000)"
// @Router /api/node/pods [get]
func (n *NodeHandler) GetNodePods() gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Query("namespace")
		nodeName := c.Query("node")

		res, err := n.svc.GetNodePods(context.Background(), namespace, nodeName)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, gerrors.NewBizError(30000, err.Error()))
			return
		}

		pods := make([]resp.PodListItem, 0, len(res))
		for _, p := range res {
			pods = append(pods, convert.PodListConvertResp(p))
		}

		response.SuccessWithData(c, pods)
	}
}
