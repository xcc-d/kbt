package handler

import (
	"kbt/internal/service"
	"net/http"
)

type HelloHandler struct {
	svc *service.HelloService
}

func NewHelloHandler(svc *service.HelloService) *HelloHandler {
	return &HelloHandler{svc: svc}
}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	msg := h.svc.SayHello(name)
	_, err := w.Write([]byte(msg))
	if err != nil {
		return
	}
}
