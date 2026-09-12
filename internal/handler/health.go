package handler

import (
	"kbt/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	Client *model.Client
}

func NewHealthHandler(client *model.Client) *HealthHandler {
	return &HealthHandler{Client: client}
}

func (h *HealthHandler) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *HealthHandler) Readyz(c *gin.Context) {
	if h.Client == nil || h.Client.K8sClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}
