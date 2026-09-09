package router

import (
	"kbt/internal/client"
	"kbt/internal/handler"
	"kbt/internal/middleware"
	"kbt/internal/service"

	"github.com/gin-gonic/gin"
)

type Client struct {
	K8sClient *client.K8sClient
}

func Setup(c *Client, authCfg middleware.AuthConfig) *gin.Engine {
	r := gin.Default()

	helloSvc := service.NewHelloService()
	helloHand := handler.NewHelloHandler(helloSvc)

	nsSvc := service.NewNamespaceService(c.K8sClient)
	nsHand := handler.NewNamespaceHandler(nsSvc)

	depSvc := service.NewDeploymentService(c.K8sClient)
	depHand := handler.NewDeploymentHandler(depSvc)

	auth, err := middleware.NewAuthMiddleware(authCfg)
	if err != nil {
		panic("failed to init auth middleware: " + err.Error())
	}

	r.GET("/hello", helloHand.Hello)

	corev1 := r.Group("/api/v1")
	corev1.Use(auth.Authenticate())
	{
		corev1.GET("/namespaces", nsHand.List)
		corev1.GET("/namespaces/:namespace", nsHand.Get)

		corev1.GET("/namespaces/:namespace/deployments", depHand.List)
		corev1.GET("/namespaces/:namespace/deployments/:deployment", depHand.Get)

		write := corev1.Group("")
		write.Use(auth.IsAdmin())
		{
			write.POST("/namespaces", nsHand.Create)
			write.DELETE("/namespaces/:namespace", nsHand.Delete)

			write.POST("/namespaces/:namespace/deployments", depHand.Create)
			write.DELETE("/namespaces/:namespace/deployments/:deployment", depHand.Delete)
		}

	}

	return r
}
