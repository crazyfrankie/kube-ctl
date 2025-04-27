package service

import (
	"context"
	es "errors"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type PodService interface {
	CreateOrUpdatePod(ctx context.Context, pod *corev1.Pod) error
	GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error)
	GetPodList(ctx context.Context, namespace string) ([]corev1.Pod, error)
	DeletePod(ctx context.Context, namespace string, name string) error
	GetNamespace(ctx context.Context) ([]corev1.Namespace, error)
}

type podService struct {
	clientSet *kubernetes.Clientset
}

func NewPodService(cs *kubernetes.Clientset) PodService {
	return &podService{clientSet: cs}
}

func (s *podService) CreateOrUpdatePod(ctx context.Context, pod *corev1.Pod) error {
	if get, err := s.clientSet.CoreV1().Pods(pod.Namespace).
		Get(ctx, pod.Name, metav1.GetOptions{}); err == nil {
		// Verify that the parameters are legal
		cPod := *pod
		cPod.Name = cPod.Name + "-validate"
		_, err := s.clientSet.CoreV1().Pods(cPod.Namespace).Create(ctx,
			&cPod, metav1.CreateOptions{DryRun: []string{metav1.DryRunAll}})
		if err != nil {
			return err
		}

		// Delete the Pod
		err = s.clientSet.CoreV1().Pods(pod.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}

		// Listen for deletion events
		labels := make([]string, 0, len(get.Labels))
		for k, v := range get.Labels {
			labels = append(labels, fmt.Sprintf("%s=%s", k, v))
		}

		ch, err := s.clientSet.CoreV1().Pods(pod.Namespace).Watch(ctx, metav1.ListOptions{
			LabelSelector: strings.Join(labels, ","),
		})
		if err != nil {
			return err
		}

		for event := range ch.ResultChan() {
			chPod := event.Object.(*corev1.Pod)

			// Fast paths, some Pods may be deleted quickly,
			// causing the listener to not start yet and subsequently keep blocking.
			// Query if the event has been deleted,
			// if it has been deleted, then you don't need to listen to the delete event.
			if _, err := s.clientSet.CoreV1().Pods(pod.Namespace).
				Get(ctx, pod.Name, metav1.GetOptions{}); errors.IsNotFound(err) {
				// Delete successful, create new Pod
				newPod, err := s.clientSet.CoreV1().Pods(pod.Namespace).Create(ctx,
					pod, metav1.CreateOptions{})
				if err != nil {
					return es.New(fmt.Sprintf("failed update pod, name: %s, %s", newPod.Name, err.Error()))
				} else {
					return nil
				}
			}

			switch event.Type {
			case watch.Deleted:
				if chPod.Name != pod.Name {
					continue
				}

				// Delete successful, create new Pod
				newPod, err := s.clientSet.CoreV1().Pods(pod.Namespace).Create(ctx,
					pod, metav1.CreateOptions{})
				if err != nil {
					return es.New(fmt.Sprintf("failed update pod, name: %s, %s", newPod.Name, err.Error()))
				} else {
					return nil
				}
			}
		}
	}

	_, err := s.clientSet.CoreV1().Pods(pod.Namespace).Create(ctx,
		pod, metav1.CreateOptions{})
	if err != nil {
		return es.New(fmt.Sprintf("failed create pod, name: %s, %s", pod.Name, err.Error()))
	}

	return nil
}

func (s *podService) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	pod, err := s.clientSet.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func (s *podService) GetPodList(ctx context.Context, namespace string) ([]corev1.Pod, error) {
	pods, err := s.clientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	return pods.Items, nil
}

func (s *podService) DeletePod(ctx context.Context, namespace string, name string) error {
	return s.clientSet.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (s *podService) GetNamespace(ctx context.Context) ([]corev1.Namespace, error) {
	list, err := s.clientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}
