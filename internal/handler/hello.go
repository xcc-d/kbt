package handler

import (
	"kbt/internal/service"

	"github.com/gin-gonic/gin"
)

type HelloHandler struct {
	svc *service.HelloService
}

func NewHelloHandler(svc *service.HelloService) *HelloHandler {
	return &HelloHandler{svc: svc}
}

func (h *HelloHandler) Hello(c *gin.Context) {
	name := c.Query("name")
	msg := h.svc.SayHello(name)
	Success(c, msg)
}
