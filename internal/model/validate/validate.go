package validate

import (
	"errors"
	"fmt"

	"github.com/crazyfrankie/kube-ctl/conf"
	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/pkg/consts"
)

func PodValidate(pod *req.Pod) error {
	// Check sum required
	if pod.Base.Name == "" {
		return errors.New("pod name is necessary")
	}
	if len(pod.Containers) == 0 {
		return errors.New("pod containers is necessary")
	}

	// Non-required setting defaults
	if len(pod.InitContainers) > 0 {
		for i, c := range pod.InitContainers {
			if c.Name == "" {
				return errors.New(fmt.Sprintf("pod init containers: %d , name is necessary", i))
			}
			if c.Image == "" {
				return errors.New(fmt.Sprintf("pod init containers: %d , image is necessary", i))
			}
			if c.ImagePullPolicy == "" {
				pod.InitContainers[i].ImagePullPolicy = consts.ImagePullPolicyIfNotPresent
			}
		}
	}

	if len(pod.Containers) > 0 {
		for i, c := range pod.Containers {
			if c.Name == "" {
				return errors.New(fmt.Sprintf("pod containers: %d , name is necessary", i))
			}
			if c.Image == "" {
				return errors.New(fmt.Sprintf("pod containers: %d , image is necessary", i))
			}
			if c.ImagePullPolicy == "" {
				pod.Containers[i].ImagePullPolicy = consts.ImagePullPolicyIfNotPresent
			}
		}
	}

	if pod.Base.RestartPolicy == "" {
		pod.Base.RestartPolicy = consts.RestartPolicyAlways
	}

	return nil
}

func StorageClassValidate(sc *req.StorageClass) error {
	expectedProv := conf.GetConf().StorageClass.Provisioner
	flag := false
	for _, i := range expectedProv {
		if i != sc.Provisioner {
			flag = true
			break
		}
	}
	if flag {
		return nil
	}

	return errors.New("unsupported provisioner type")
}
