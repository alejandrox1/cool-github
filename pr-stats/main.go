package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func gracefulShutdown(cancel context.CancelFunc, quit <-chan os.Signal, done chan<- interface{}) {
	<-quit
	fmt.Println("App is shutting down...")
	defer cancel()
	close(done)
}

func main() {
	prStats := newPRStats()

	quit := make(chan os.Signal, 1)
	done := make(chan interface{}, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go gracefulShutdown(prStats.CtxCancelFunc, quit, done)
	prStats.AnalyzePRLabelFreq("kubernetes", "kubernetes")

	close(quit)
	<-done
	prStats.SummarizeResults()
}
