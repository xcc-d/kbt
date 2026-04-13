package main

import (
	"context"
	"kbt/internal/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mux := router.Setup()
	addr := ":11313"

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		log.Println("Starting server... ")
		if err := http.ListenAndServe(addr, mux); err != nil {
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
