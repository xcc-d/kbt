package client

import (
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8sClient struct {
	ClientSet kubernetes.Interface
}

func NewK8sClient() (*K8sClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeConfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			return nil, err
		}
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &K8sClient{ClientSet: clientSet}, nil

}

func NewFakeK8sClient() *K8sClient {
	return &K8sClient{ClientSet: fake.NewSimpleClientset()}
}

// NewFakeK8sClientWithObjects 返回预置了初始对象的 fake client（用于本地开发/测试）
func NewFakeK8sClientWithObjects(objects ...runtime.Object) *K8sClient {
	return &K8sClient{ClientSet: fake.NewSimpleClientset(objects...)}
}
