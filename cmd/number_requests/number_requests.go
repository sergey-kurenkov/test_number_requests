package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sergey-kurenkov/test_number_requests/internal/app"
)

const (
	defaultAPIAddress = ":8888"
)

func getVar(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Print("start")

	const defaultInterval time.Duration = 60 * time.Second
	application := app.NewGetNumberRequestsApp(defaultInterval)

	apiAddress := getVar("API_ADDRESS", defaultAPIAddress)
	server := &http.Server{Addr: apiAddress, Handler: application.Handler()}

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	log.Print("wait to be stopped")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-stop

	log.Print("stopping")
	application.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("ListenAndServe(): %v", err) // nolint:gocritic
	}

	log.Print("stopped")
}
