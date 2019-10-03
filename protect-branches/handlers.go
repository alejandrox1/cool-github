/*
 */
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/v28/github"
)

func handleGithubWebhook(w http.ResponseWriter, r *http.Request) {
	logger := log.New(os.Stdout, "githubWebhook: ", log.LstdFlags)

	payload, err := github.ValidatePayload(r, webhookSecret)
	if err != nil {
		logger.Printf("error validating request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		logger.Printf("error parsing webhook: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch event := event.(type) {
	case *github.RepositoryEvent:
		if event.Action != nil && *event.Action == "created" {
			logger.Printf("adding branch default branch protection policy for %s/%s\n", *event.Repo.Owner.Login, *event.Repo.Name)
			githubBranchPolicy := NewGithubRepoPolicy()
			githubBranchPolicy.createBranchProtection(*event.Repo.Owner.Login, *event.Repo.Name, "master")
		}
	}
}
