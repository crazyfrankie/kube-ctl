package convert

import (
	"encoding/base64"
	"time"
	"unsafe"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"github.com/crazyfrankie/kube-ctl/pkg/utils"
)

func SecretReqConvert(req *req.Secret) *corev1.Secret {
	data := make(map[string]string)
	for _, i := range req.Data {
		val := base64.StdEncoding.EncodeToString(unsafeToBytes(i.Value))
		data[i.Key] = val
	}
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:              req.Name,
			Namespace:         req.Namespace,
			Labels:            utils.ReqItemToMap(req.Labels),
			CreationTimestamp: metav1.Time{Time: time.Now()},
		},
		StringData: data,
		Type:       req.Type,
	}
}

func unsafeToBytes(s string) []byte {
	sh := (*[2]uintptr)(unsafe.Pointer(&s))

	res := [3]uintptr{sh[0], sh[1], sh[1]}

	return *(*[]byte)(unsafe.Pointer(&res))
}

func SecretConvertListResp(s *corev1.Secret) resp.Secret {
	return resp.Secret{
		Name:      s.Name,
		Namespace: s.Namespace,
		Type:      s.Type,
		DataNum:   len(s.Data),
		Age:       int64(time.Now().Sub(s.CreationTimestamp.Time).Seconds()),
	}
}

func SecretConvertDetailResp(s *corev1.Secret) resp.SecretDetail {
	return resp.SecretDetail{
		Name:      s.Name,
		Namespace: s.Namespace,
		DataNum:   len(s.Data),
		Age:       s.CreationTimestamp.Time.Unix(),
		Type:      s.Type,
		Labels:    utils.ResMapToItem(s.Labels),
		Data:      utils.ResMapToItem(s.StringData),
	}
}
