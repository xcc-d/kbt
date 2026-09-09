package main

import (
	"context"
	"kbt/internal/middleware"
	"kbt/internal/router"
	"kbt/test/seed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 本地开发用 fake client，并预置 seed 数据（列表接口有数据可查）
	k8sClient := seed.NewFakeK8sClient()
	authCfg := middleware.AuthConfig{
		Issuer:   getEnv("KEYCLOAK_ISSUER", "http://localhost:8080/realms/kbt"),
		ClientID: getEnv("KEYCLOAK_CLIENT_ID", "kbt-backend"),
	}

	r := router.Setup(
		&router.Client{K8sClient: k8sClient},
		authCfg)

	addr := ":11313"

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Println("Starting server... ")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting gracefully")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
