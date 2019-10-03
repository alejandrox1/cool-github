/*
 */
package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/google/go-github/v28/github"
)

type BranchProtectionRule struct {
	*github.Protection
	*github.SignaturesProtectedBranch `json:"required_signatures"`
}

type GithubRepoPolicy struct {
	Client  *github.Client
	Context context.Context
	Logger  *log.Logger
}

func NewGithubRepoPolicy() *GithubRepoPolicy {
	client, ctx := newGithubClient()
	logger := log.New(os.Stdout, "githubBranchPolicy: ", log.LstdFlags)

	return &GithubRepoPolicy{
		Client:  client,
		Context: ctx,
		Logger:  logger,
	}
}

func (g *GithubRepoPolicy) Printf(format string, a ...interface{}) {
	g.Logger.Printf(format, a...)
}

func (g *GithubRepoPolicy) createBranchProtection(owner, repo, branch string) {
	protectionRequest := &github.ProtectionRequest{
		RequiredStatusChecks: &github.RequiredStatusChecks{
			Strict:   true,
			Contexts: []string{},
		},
		RequiredPullRequestReviews: &github.PullRequestReviewsEnforcementRequest{
			DismissStaleReviews:          true,
			RequiredApprovingReviewCount: 2,
		},
		EnforceAdmins: true,
		Restrictions:  nil,
	}

	protection, resp, err := g.Client.Repositories.UpdateBranchProtection(g.Context, owner, repo, "master", protectionRequest)
	if (resp.StatusCode < 200 || resp.StatusCode > 299) || err != nil {
		g.Printf("updateBranch protection returned status code: %v and err: %v\n", resp.StatusCode, err)
		return
	}

	signatureProtection, resp, err := g.Client.Repositories.RequireSignaturesOnProtectedBranch(g.Context, owner, repo, "master")
	if (resp.StatusCode < 200 || resp.StatusCode > 299) || err != nil {
		g.Printf("require signature returned status code: %v and err: %v\n", resp.StatusCode)
		return
	}

	if protection != nil {
		branchProtectionRule := BranchProtectionRule{protection, signatureProtection}
		jsonProtection, err := json.MarshalIndent(branchProtectionRule, "", "\t")
		if err != nil {
			g.Printf("error marshaling protection: %v\n", err)
			return
		}
		g.Printf("Protection added: %s\n", string(jsonProtection))
	}
}
