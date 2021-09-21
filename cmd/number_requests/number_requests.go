package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sergey-kurenkov/test_number_requests/internal/http_server"
)

const (
	defaultApiAddress = ":8888"
)

var apiAddress string

func getVar(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	apiAddress = getVar("API_ADDRESS", defaultApiAddress)
}

func main() {
	log.Print("start")
	app := http_server.NewGetNumberRequestsApp(60 * time.Second)

	server := &http.Server{Addr: apiAddress, Handler: app.Handler()}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	log.Print("wait to be stopped")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-stop

	log.Print("stopping")
	app.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("ListenAndServe(): %v", err)
	}

	log.Print("stopped")
}
