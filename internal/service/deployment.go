package service

import (
	"context"
	apperr "kbt/pkg/errors"

	"kbt/internal/client"

	appv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type DeploymentService struct {
	K8sClient *client.K8sClient
}

func NewDeploymentService(k8sClient *client.K8sClient) *DeploymentService {
	return &DeploymentService{K8sClient: k8sClient}
}

func (d *DeploymentService) DeploymentList(ctx context.Context, namespace string) ([]string, error) {
	return d.K8sClient.DeploymentList(ctx, namespace)
}

func (d *DeploymentService) DeploymentGet(ctx context.Context, namespace, name string) (*appv1.Deployment, error) {
	return d.K8sClient.DeploymentGet(ctx, namespace, name)
}

func (d *DeploymentService) DeploymentCreate(ctx context.Context, yamlFile []byte) error {
	var dep appv1.Deployment
	if err := yaml.Unmarshal(yamlFile, &dep); err != nil {
		return apperr.NewBadRequest("invalid yaml")
	}

	if dep.Namespace == "" || dep.Name == "" {
		return apperr.NewBadRequest("namespace or name is empty")
	}

	return d.K8sClient.DeploymentCreate(ctx, dep)
}

func (d *DeploymentService) DeploymentDelete(ctx context.Context, namespace, name string) error {
	return d.K8sClient.DeploymentDelete(ctx, namespace, name)
}
