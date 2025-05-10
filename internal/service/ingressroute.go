package service

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/crazyfrankie/kube-ctl/internal/model/req"
	"github.com/crazyfrankie/kube-ctl/pkg/utils"
)

type IngressRoute struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ObjectMeta    `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec            req.IngressRouteSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

type IngressRouteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []IngressRoute `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type MiddlewareList struct {
	Items           []req.Middleware `json:"items"`
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
}

type IngressRouteService interface {
	CreateOrUpdateIngressRoute(ctx context.Context, req *req.IngressRoute) error
	DeleteIngressRoute(ctx context.Context, name string, namespace string) error
	GetIngressRouteDetail(ctx context.Context, name string, namespace string) (*IngressRoute, error)
	GetIngressRouteList(ctx context.Context, namespace string) (*IngressRouteList, error)
	GetIngressRouteMws(ctx context.Context, namespace string) ([]string, error)
}

type ingressRouteService struct {
	clientSet *kubernetes.Clientset
}

func NewIngressRouteService(cs *kubernetes.Clientset) IngressRouteService {
	return &ingressRouteService{clientSet: cs}
}

func (s *ingressRouteService) CreateOrUpdateIngressRoute(ctx context.Context, request *req.IngressRoute) error {
	ig := &IngressRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       "traefik.io/v1alpha1",
			APIVersion: "IngressRoute",
		},
		Metadata: metav1.ObjectMeta{
			Name:      request.Name,
			Namespace: request.Namespace,
			Labels:    utils.ReqItemToMap(request.Labels),
		},
		Spec: request.IngressRouteSpec,
	}
	url := fmt.Sprintf("/apis/traefik.io/v1alpha1/namespaces/%s/ingressroutes/%s", ig.Metadata.Namespace, ig.Metadata.Name)

	raw, _ := sonic.Marshal(ig)
	if _, err := s.clientSet.NetworkingV1().RESTClient().Get().AbsPath(url).DoRaw(ctx); err == nil {
		_, err = s.clientSet.NetworkingV1().RESTClient().Put().AbsPath(url).Body(raw).DoRaw(ctx)

		return err
	}

	_, err := s.clientSet.NetworkingV1().RESTClient().Post().AbsPath(url).Body(raw).DoRaw(ctx)

	return err
}

func (s *ingressRouteService) DeleteIngressRoute(ctx context.Context, name string, namespace string) error {
	url := fmt.Sprintf("/apis/traefik.io/v1alpha1/namespaces/%s/ingressroutes/%s", namespace, name)
	_, err := s.clientSet.NetworkingV1().RESTClient().Delete().AbsPath(url).DoRaw(ctx)

	return err
}

func (s *ingressRouteService) GetIngressRouteDetail(ctx context.Context, name string, namespace string) (*IngressRoute, error) {
	url := fmt.Sprintf("/apis/traefik.io/v1alpha1/namespaces/%s/ingressroutes/%s", namespace, name)
	raw, err := s.clientSet.NetworkingV1().RESTClient().Get().AbsPath(url).DoRaw(ctx)
	if err != nil {
		return nil, err
	}

	var res IngressRoute
	if err = sonic.Unmarshal(raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (s *ingressRouteService) GetIngressRouteList(ctx context.Context, namespace string) (*IngressRouteList, error) {
	url := fmt.Sprintf("/apis/traefik.io/v1alpha1/namespaces/%s/ingressroutes", namespace)
	raw, err := s.clientSet.NetworkingV1().RESTClient().Get().AbsPath(url).DoRaw(ctx)
	if err != nil {
		return nil, err
	}

	var res IngressRouteList
	if err = sonic.Unmarshal(raw, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (s *ingressRouteService) GetIngressRouteMws(ctx context.Context, namespace string) ([]string, error) {
	url := fmt.Sprintf("/apis/traefik.io/v1alpha1/namespaces/%s/middlewares", namespace)

	raw, err := s.clientSet.RESTClient().Get().AbsPath(url).DoRaw(ctx)
	mws := make([]string, 0)
	var middlewareList MiddlewareList
	err = sonic.Unmarshal(raw, &middlewareList)
	if err != nil {
		return nil, err
	}
	for _, item := range middlewareList.Items {
		mws = append(mws, item.Metadata.Name)
	}

	return mws, nil
}
