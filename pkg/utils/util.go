package utils

import (
	"fmt"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/internal/model/resp"
	"hash/fnv"
	"strconv"
	"time"
	"unsafe"
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

func GenerateHashBaseRGB(s string) string {
	hash := hashString(s)
	r, g, b := hashToRGB(hash)

	return strconv.Itoa(r) + "," + strconv.Itoa(g) + "," + strconv.Itoa(b) + ","
}

func FormatTime(d time.Duration) string {
	seconds := int(d.Seconds())
	minutes := int(d.Minutes())
	hours := int(d.Hours())
	days := hours / 24
	years := days / 365

	switch {
	case years > 0:
		return fmt.Sprintf("%d years", years)
	case days > 0:
		return fmt.Sprintf("%d days", days)
	case hours > 0:
		return fmt.Sprintf("%d hours", hours)
	case minutes > 0:
		return fmt.Sprintf("%d minutes", minutes)
	default:
		return fmt.Sprintf("%d seconds", seconds)
	}
}

func hashString(s string) uint32 {
	h := fnv.New32a()
	h.Write(unsafeStringToByte(s))

	return h.Sum32()
}

func hashToRGB(h uint32) (int, int, int) {
	r := int(h) & 0xFF
	g := int((h >> 8) & 0xFF)
	b := int((h >> 16) & 0xFF)

	return r, g, b
}

func unsafeStringToByte(s string) []byte {
	sh := (*[2]uintptr)(unsafe.Pointer(&s))

	res := [3]uintptr{sh[0], sh[1], sh[1]}

	return *(*[]byte)(unsafe.Pointer(&res))
}
