/*
 */
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	serverGracePeriod = 30 * time.Second
)

func newWebServer(addr string, logger *log.Logger) *http.Server {
	router := http.NewServeMux()
	router.HandleFunc("/webhook", handleGithubWebhook)

	return &http.Server{
		Addr:         addr,
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}
}

func gracefulShutdown(server *http.Server, logger *log.Logger, quit <-chan os.Signal, done chan<- interface{}) {
	<-quit
	logger.Println("server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), serverGracePeriod)
	defer func() {
		// Close any other client connections you have here.
		cancel()
	}()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("error while shutting down the server: %v\n", err)
	}
	close(done)
}
