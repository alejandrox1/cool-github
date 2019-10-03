/*
- Allow server to finish any open request. Handle work in progress with a max timeout.
- Do not keep any connections that finish alive.

    Container orchestration systems will ususally send a SIGTERM and allow the container
to performa a graceful shutdown. After this period, if the process is still running then
a SIGKILL will be sent.
*/
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	listenAddr = "0.0.0.0:8080"
)

var (
	githubSecretsConfig string
	branchPolicyConfig  string

	// Github secrets.
	githubAccessToken   string
	githubWebhookSecret string

	// Github branch rules.
	actingBranchPolicy *BranchPolicy
)

func parseFlags() {
	flag.StringVar(&githubSecretsConfig, "githubSecrets", "", "configuration file with access token and webhook secret")
	flag.StringVar(&branchPolicyConfig, "branchPolicy", "branch_policy.yaml", "coniguration file for branch protection policy")
	flag.Parse()
}

func parseConfigs() error {
	githubToken, webhookSecret, err := readGithubSecrets(githubSecretsConfig)
	if err != nil {
		return err
	}

	// Set the global variables.
	githubAccessToken = githubToken
	githubWebhookSecret = webhookSecret

	bp, err := readBranchProtectionConfig(branchPolicyConfig)
	if err != nil {
		log.Printf("error reading branch policy config: %+v\n", err)
	}
	actingBranchPolicy = bp
	log.Println(actingBranchPolicy)
	return nil
}

func main() {
	parseFlags()
	if err := parseConfigs(); err != nil {
		panic(err)
	}

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
