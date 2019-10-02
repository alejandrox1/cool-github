/*
- Allow server to finish any open request. Handle work in progress with a max timeout.
- Do not keep any connections that finish alive.

    Container orchestration systems will ususally send a SIGTERM and allow the container
to performa a graceful shutdown. After this period, if the process is still running then
a SIGKILL will be sent.
*/
package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	listenAddr = "0.0.0.0:8080"
)

func main() {
	quit := make(chan os.Signal, 1)
	done := make(chan interface{}, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger := log.New(os.Stdout, "gitserver: ", log.LstdFlags)

	server := newWebServer(listenAddr, logger)

	go gracefulShutdown(server, logger, quit, done)

	logger.Printf("Server listening on: %s\n", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("error trying to listen on %v %v\n", listenAddr, err)
	}

	<-done
	logger.Println("server stopped")
}