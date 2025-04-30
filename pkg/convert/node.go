package convert

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/crazyfrankie/kube-ctl/internal/api/model/resp"
)

func NodeListItemConvertResp(node corev1.Node) resp.NodeListItem {
	return resp.NodeListItem{
		Name:             node.Name,
		Age:              node.CreationTimestamp.Unix(),
		Version:          node.Status.NodeInfo.KubeletVersion,
		OSImage:          node.Status.NodeInfo.OSImage,
		KernelVersion:    node.Status.NodeInfo.KernelVersion,
		ContainerRuntime: node.Status.NodeInfo.ContainerRuntimeVersion,
		Status:           getNodeStatus(node.Status.Conditions, corev1.NodeReady, corev1.ConditionTrue),
		InternalIP:       getNodeIP(node.Status.Addresses, corev1.NodeInternalIP),
		ExternalIP:       getNodeIP(node.Status.Addresses, corev1.NodeExternalIP),
	}
}

func getNodeStatus(conditions []corev1.NodeCondition, typ corev1.NodeConditionType, sta corev1.ConditionStatus) string {
	for _, cd := range conditions {
		if cd.Type == typ && cd.Status == sta {
			return "Ready"
		}
	}

	return "NotReady"
}

func getNodeIP(adds []corev1.NodeAddress, addrTyp corev1.NodeAddressType) string {
	for _, ad := range adds {
		if ad.Type == addrTyp {
			return ad.Address
		}
	}

	return "<none>"
}

func NodeDetailConvertResp(node *corev1.Node) resp.NodeDetail {
	return resp.NodeDetail{
		Name:             node.Name,
		Age:              node.CreationTimestamp.Unix(),
		Version:          node.Status.NodeInfo.KubeletVersion,
		OSImage:          node.Status.NodeInfo.OSImage,
		KernelVersion:    node.Status.NodeInfo.KernelVersion,
		ContainerRuntime: node.Status.NodeInfo.ContainerRuntimeVersion,
		Status:           getNodeStatus(node.Status.Conditions, corev1.NodeReady, corev1.ConditionTrue),
		InternalIP:       getNodeIP(node.Status.Addresses, corev1.NodeInternalIP),
		ExternalIP:       getNodeIP(node.Status.Addresses, corev1.NodeExternalIP),
		Labels:           getNodeDetailLabels(node.Labels),
		Taints:           node.Spec.Taints,
	}
}

func getNodeDetailLabels(labels map[string]string) []resp.Item {
	res := make([]resp.Item, 0, len(labels))
	for k, v := range labels {
		res = append(res, resp.Item{
			Key:   k,
			Value: v,
		})
	}

	return res
}
