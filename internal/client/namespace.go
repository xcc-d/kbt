package client

import (
	"context"
	apperr "kbt/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *K8sClient) NamespaceList(ctx context.Context) ([]string, error) {
	list, err := c.ClientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, apperr.FromError(err)
	}
	names := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		names = append(names, item.Name)
	}
	return names, nil
}

func (c *K8sClient) NamespaceGet(ctx context.Context, name string) (*corev1.Namespace, error) {
	ns, err := c.ClientSet.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, apperr.FromError(err)
	}

	return ns, nil
}

func (c *K8sClient) NamespaceCreate(ctx context.Context, ns *corev1.Namespace) (*corev1.Namespace, error) {
	created, err := c.ClientSet.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if err != nil {
		return nil, apperr.FromError(err)
	}
	return created, nil
}

func (c *K8sClient) NamespaceDelete(ctx context.Context, name string) error {
	err := c.ClientSet.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return apperr.FromError(err)
	}
	return nil
}
