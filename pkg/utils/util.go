package utils

import (
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
)

func ReqItemToMap(req []req.Item) map[string]string {
	res := make(map[string]string)
	for _, i := range req {
		res[i.Key] = i.Value
	}

	return res
}

func ReqMapToItem(ma map[string]string) []req.Item {
	res := make([]req.Item, 0, len(ma))
	for k, v := range ma {
		res = append(res, req.Item{
			Key:   k,
			Value: v,
		})
	}

	return res
}

func ResMapToItem(ma map[string]string) []resp.Item {
	res := make([]resp.Item, 0, len(ma))
	for k, v := range ma {
		res = append(res, resp.Item{
			Key:   k,
			Value: v,
		})
	}

	return res
}
