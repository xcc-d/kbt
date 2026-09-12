package service

import (
	"context"
	"kbt/internal/audit"
	apperr "kbt/pkg/errors"
	"kbt/pkg/utils"

	"kbt/internal/client"

	"go.uber.org/zap"
	appv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type DeploymentService struct {
	K8sClient *client.K8sClient
	Recorder  *audit.Recorder
}

func NewDeploymentService(k8sClient *client.K8sClient, recorder *audit.Recorder) *DeploymentService {
	return &DeploymentService{K8sClient: k8sClient, Recorder: recorder}
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
		utils.Biz(ctx, "create", "deployment", "").Fail(err)
		return apperr.NewBadRequest("invalid yaml")
	}

	if dep.Namespace == "" || dep.Name == "" {
		utils.Biz(ctx, "create", "deployment", dep.Name).Fail(apperr.NewBadRequest("namespace or name is empty"))
		return apperr.NewBadRequest("namespace or name is empty")
	}

	if err := d.K8sClient.DeploymentCreate(ctx, dep); err != nil {
		utils.Biz(ctx, "create", "deployment", dep.Name).
			WithExtra(zap.String("namespace", dep.Namespace)).
			Fail(err)
		return err
	}

	utils.Biz(ctx, "create", "deployment", dep.Name).
		WithExtra(zap.String("namespace", dep.Namespace)).
		Success()
	return nil
}

func (d *DeploymentService) DeploymentDelete(ctx context.Context, namespace, name string) error {
	if err := d.K8sClient.DeploymentDelete(ctx, namespace, name); err != nil {
		utils.Biz(ctx, "delete", "deployment", name).
			WithExtra(zap.String("namespace", namespace)).
			Fail(err)
		return err
	}

	utils.Biz(ctx, "delete", "deployment", name).
		WithExtra(zap.String("namespace", namespace)).
		Success()
	return nil
}

func (d *DeploymentService) DeploymentPatch(ctx context.Context, namespace, name string, path []byte) (*appv1.Deployment, error) {
	before, err := d.K8sClient.DeploymentGet(ctx, namespace, name)
	if err != nil {
		utils.Biz(ctx, "update", "deployment", name).
			WithExtra(zap.String("namespace", namespace)).
			Fail(err)
		return nil, err
	}

	after, err := d.K8sClient.DeploymentPatch(ctx, namespace, name, path)
	if err != nil {
		d.Recorder.Record(&audit.Entry{
			Operator:     utils.UserFromCtx(ctx),
			Action:       "update",
			Resource:     "deployment",
			ResourceName: name,
			Namespace:    namespace,
			Result:       "failed",
			ErrorMsg:     err.Error(),
			Before:       before,
			After:        nil,
			ClientIP:     utils.ClientIPFromCtx(ctx),
			UserAgent:    utils.UserAgentFromCtx(ctx),
		})
		utils.Biz(ctx, "update", "deployment", name).
			WithExtra(zap.String("namespace", namespace)).
			Fail(err)
		return nil, err
	}

	d.Recorder.Record(&audit.Entry{
		Operator:     utils.UserFromCtx(ctx),
		Action:       "update",
		Resource:     "deployment",
		ResourceName: name,
		Namespace:    namespace,
		Result:       "success",
		Before:       before,
		After:        after,
		ClientIP:     utils.ClientIPFromCtx(ctx),
		UserAgent:    utils.UserAgentFromCtx(ctx),
	})

	utils.Biz(ctx, "update", "deployment", name).
		WithExtra(zap.String("namespace", namespace)).
		Success()
	return after, nil

}
