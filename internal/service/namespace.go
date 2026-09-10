package service

import (
	"context"
	"kbt/internal/client"
	"kbt/internal/model"
	"kbt/pkg/utils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceService struct {
	k8sClient *client.K8sClient
}

func NewNamespaceService(k8sClient *client.K8sClient) *NamespaceService {
	return &NamespaceService{k8sClient: k8sClient}
}

func (n *NamespaceService) NamespaceList(ctx context.Context) ([]string, error) {
	return n.k8sClient.NamespaceList(ctx)
}

func (n *NamespaceService) NamespaceGet(ctx context.Context, name string) (*corev1.Namespace, error) {
	return n.k8sClient.NamespaceGet(ctx, name)
}

func (n *NamespaceService) NamespaceCreate(ctx context.Context, req *model.Namespace) (*corev1.Namespace, error) {
	name := req.Name
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: name},
	}

	created, err := n.k8sClient.NamespaceCreate(ctx, ns)
	if err != nil {
		utils.Biz(ctx, "create", "namespace", name).Fail(err)
		return nil, err
	}

	utils.Biz(ctx, "create", "namespace", name).Success()
	return created, nil
}

func (n *NamespaceService) NamespaceDelete(ctx context.Context, name string) error {
	if err := n.k8sClient.NamespaceDelete(ctx, name); err != nil {
		utils.Biz(ctx, "delete", "namespace", name).Fail(err)
		return err
	}

	utils.Biz(ctx, "delete", "namespace", name).Success()
	return nil
}
