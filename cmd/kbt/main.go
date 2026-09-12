package main

import (
	"context"
	"errors"
	"kbt/internal/audit"
	"kbt/internal/client"
	"kbt/internal/middleware"
	"kbt/internal/model"
	"kbt/internal/router"
	"kbt/pkg/utils"
	"kbt/test/seed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	utils.InitLogger("level", "json")
	// 本地开发用 fake client，并预置 seed 数据（列表接口有数据可查）
	k8sClient := seed.NewFakeK8sClient()
	authCfg := middleware.AuthConfig{
		Issuer:   getEnv("KEYCLOAK_ISSUER", "http://localhost:8080/realms/kbt"),
		ClientID: getEnv("KEYCLOAK_CLIENT_ID", "kbt-backend"),
	}
	gormDB, err := client.NewMysqlClient(client.DbConfig{
		Host:                  getEnv("MYSQL_HOST", "127.0.0.1"),
		Port:                  3306,
		Username:              getEnv("MYSQL_USER", "root"),
		Password:              getEnv("MYSQL_PASSWORD", "root123"),
		Database:              getEnv("MYSQL_DATABASE", "kbt"),
		MaxIdleConnections:    5,
		MaxOpenConnections:    20,
		MaxConnectionLifeTime: time.Hour,
	})
	if err != nil {
		utils.L().Error("connect mysql failed:", zap.Error(err))
	}

	if err := gormDB.AutoMigrate(&model.AuditLog{}); err != nil {
		utils.L().Error("auto migrate failed: ", zap.Error(err))
	}

	secretKey, err := audit.LoadSecretKey()
	if err != nil {
		utils.L().Error("load secret key failed:", zap.Error(err))
	}

	recorder := audit.NewRecorder(gormDB, secretKey, 1024)

	r := router.Setup(
		&model.Client{K8sClient: k8sClient,
			DB: gormDB,
		},
		recorder,
		authCfg)

	addr := ":18080"

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		utils.L().Info("Starting server... ", zap.String("addr:", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			utils.L().Error("Server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.L().Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	utils.L().Info("Server exiting gracefully")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
