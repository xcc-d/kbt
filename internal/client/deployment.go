package client

import (
	"context"

	apperr "kbt/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *K8sClient) DeploymentList(ctx context.Context, namespace string) ([]string, error) {
	list, err := c.ClientSet.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, apperr.FromError(err)
	}
	names := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		names = append(names, item.Name)
	}
	return names, nil
}

func (c *K8sClient) DeploymentGet(ctx context.Context, namespace, name string) (*appsv1.Deployment, error) {
	get, err := c.ClientSet.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, apperr.FromError(err)
	}
	return get, nil
}

func (c *K8sClient) DeploymentCreate(ctx context.Context, dep appsv1.Deployment) error {
	_, err := c.ClientSet.AppsV1().Deployments(dep.Namespace).Create(ctx, &dep, metav1.CreateOptions{})
	if err != nil {
		return apperr.FromError(err)
	}
	return nil
}

func (c *K8sClient) DeploymentDelete(ctx context.Context, namespace, name string) error {
	err := c.ClientSet.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return apperr.FromError(err)
	}
	return nil
}
