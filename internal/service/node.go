package service

import (
	"context"

	"github.com/bytedance/sonic"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"

	"github.com/crazyfrankie/kube-ctl/internal/model/req"
)

type NodeService interface {
	NodeList(ctx context.Context) ([]corev1.Node, error)
	NodeDetail(ctx context.Context, name string) (*corev1.Node, error)
	UpdateNodeLabel(ctx context.Context, req req.UpdateLabelReq) error
	UpdateNodeTaints(ctx context.Context, req req.UpdateTaintReq) error
	GetNodePods(ctx context.Context, namespace string, nodeName string) ([]corev1.Pod, error)
}

type nodeService struct {
	clientSet *kubernetes.Clientset
}

func NewNodeService(cs *kubernetes.Clientset) NodeService {
	return &nodeService{clientSet: cs}
}

func (s *nodeService) NodeList(ctx context.Context) ([]corev1.Node, error) {
	res, err := s.clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return res.Items, nil
}

func (s *nodeService) NodeDetail(ctx context.Context, name string) (*corev1.Node, error) {
	res, err := s.clientSet.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *nodeService) UpdateNodeLabel(ctx context.Context, req req.UpdateLabelReq) error {
	labels := make(map[string]string)
	for _, l := range req.Labels {
		labels[l.Key] = l.Value
	}

	labels["$patch"] = "replace"
	update := map[string]any{
		"metadata": map[string]any{
			"labels": labels,
		},
	}

	data, err := sonic.Marshal(&update)
	if err != nil {
		return err
	}

	_, err = s.clientSet.CoreV1().Nodes().Patch(ctx, req.Name, types.StrategicMergePatchType, data, metav1.PatchOptions{})

	return err
}

func (s *nodeService) UpdateNodeTaints(ctx context.Context, req req.UpdateTaintReq) error {
	taints := map[string]any{
		"spec": map[string]any{
			"taints": req.Taints,
		},
	}

	data, err := sonic.Marshal(taints)
	if err != nil {
		return err
	}

	_, err = s.clientSet.CoreV1().Nodes().Patch(ctx, req.Name, types.StrategicMergePatchType, data, metav1.PatchOptions{})

	return err
}

func (s *nodeService) GetNodePods(ctx context.Context, namespace string, nodeName string) ([]corev1.Pod, error) {
	pods, err := s.clientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	res := make([]corev1.Pod, 0, len(pods.Items))
	for _, i := range pods.Items {
		if i.Spec.NodeName == nodeName {
			res = append(res, i)
		}
	}

	return res, nil
}
