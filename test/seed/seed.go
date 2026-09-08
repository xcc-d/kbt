package seed

import (
	"kbt/internal/client"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// NamespaceNames 预置的 Namespace 名称列表
var NamespaceNames = []string{
	"default",
	"kube-system",
	"kube-public",
	"dev",
	"prod",
}

// NamespaceObjects 返回预置的 Namespace 对象列表（供 fake.NewSimpleClientset 使用）
func NamespaceObjects() []*corev1.Namespace {
	objs := make([]*corev1.Namespace, 0, len(NamespaceNames))
	for _, name := range NamespaceNames {
		objs = append(objs, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		})
	}
	return objs
}

// NewFakeK8sClient 返回预置了 Namespace 数据的 fake client。
func NewFakeK8sClient() *client.K8sClient {
	objs := make([]runtime.Object, 0, len(NamespaceNames))
	for _, ns := range NamespaceObjects() {
		objs = append(objs, ns)
	}
	return client.NewFakeK8sClientWithObjects(objs...)
}
