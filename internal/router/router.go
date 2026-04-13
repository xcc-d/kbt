package router

import (
	"kbt/internal/handler"
	"kbt/internal/service"
	"net/http"
)

func Setup() *http.ServeMux {
	helloSvc := service.NewHelloService()
	helloHand := handler.NewHelloHandler(helloSvc)

	mux := http.NewServeMux()
	mux.Handle("/hello", helloHand)

	return mux
}
