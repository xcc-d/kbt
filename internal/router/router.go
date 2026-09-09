package router

import (
	"kbt/internal/client"
	"kbt/internal/handler"
	"kbt/internal/service"

	"github.com/gin-gonic/gin"
)

type Client struct {
	K8sClient *client.K8sClient
}

func Setup(c *Client) *gin.Engine {
	r := gin.Default()

	helloSvc := service.NewHelloService()
	helloHand := handler.NewHelloHandler(helloSvc)

	nsSvc := service.NewNamespaceService(c.K8sClient)
	nsHand := handler.NewNamespaceHandler(nsSvc)

	depSvc := service.NewDeploymentService(c.K8sClient)
	depHand := handler.NewDeploymentHandler(depSvc)

	r.GET("/hello", helloHand.Hello)

	corev1 := r.Group("/api/v1")
	{
		corev1.GET("/namespaces", nsHand.List)
		corev1.GET("/namespaces/:namespace", nsHand.Get)
		corev1.POST("/namespaces", nsHand.Create)
		corev1.DELETE("/namespaces/:namespace", nsHand.Delete)

		corev1.GET("/namespaces/:namespace/deployments", depHand.List)
		corev1.GET("/namespaces/:namespace/deployments/:deployment", depHand.Get)
		corev1.POST("/namespaces/:namespace/deployments", depHand.Create)
		corev1.DELETE("/namespaces/:namespace/deployments/:deployment", depHand.Delete)
	}

	return r
}
