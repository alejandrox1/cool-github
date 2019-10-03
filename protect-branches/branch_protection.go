/*
   Implement the logic for creating a branch protection policy. We use the
   term policy as the logic here enables the user to create a protection rule
   with more specifications (i.e., require signed commits.
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v28/github"
)

// BranchProtectionPolicy
type BranchProtectionPolicy struct {
	*github.Protection
	*github.SignaturesProtectedBranch `json:"required_signatures"`
}

// GithubRepoPolicy serves as the agent for performing API calls.
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

// createBranchProtection creates a branch protection policy for the specified
// branch on the given repo.
func (g *GithubRepoPolicy) createBranchProtection(owner, repo, branch string) error {
	// Create a branch protection request. Populate it with the default values
	// from actingBranchPolicy.
	protectionRequest := &github.ProtectionRequest{
		RequiredStatusChecks: &github.RequiredStatusChecks{
			Strict:   actingBranchPolicy.RequireStatusChecks.Strict,
			Contexts: actingBranchPolicy.RequireStatusChecks.Contexts,
		},
		RequiredPullRequestReviews: &github.PullRequestReviewsEnforcementRequest{
			DismissStaleReviews:          actingBranchPolicy.RequiredPullRequestReviews.DismissStaleReviews,
			RequiredApprovingReviewCount: actingBranchPolicy.RequiredPullRequestReviews.RequiredApprovingReviewCount,
		},
		EnforceAdmins: actingBranchPolicy.EnforceAdmins,
		Restrictions:  nil,
	}
	protection, resp, err := g.Client.Repositories.UpdateBranchProtection(g.Context, owner, repo, "master", protectionRequest)
	if (resp.StatusCode < 200 || resp.StatusCode > 299) || err != nil {
		g.Printf("updateBranch protection returned status code: %v and err: %v\n", resp.StatusCode, err)
		return err
	}

	// If requiering signed commits, then make a request to do so.
	signatureProtection := &github.SignaturesProtectedBranch{}
	if actingBranchPolicy.RequireSignatures {
		signatureProtection, resp, err = g.Client.Repositories.RequireSignaturesOnProtectedBranch(g.Context, owner, repo, "master")
		if (resp.StatusCode < 200 || resp.StatusCode > 299) || err != nil {
			g.Printf("require signature returned status code: %v and err: %v\n", resp.StatusCode)
			return err
		}
	}

	// Create a string summary of the branch protection policy. This will be
	// included in the body of an issue.
	var protectionSummary string
	if protection != nil {
		branchProtectionRule := BranchProtectionPolicy{protection, signatureProtection}
		jsonProtection, err := json.MarshalIndent(branchProtectionRule, "", "\t")
		if err != nil {
			g.Printf("error marshaling protection: %v\n", err)
			return err
		}
		protectionSummary = fmt.Sprintf("Protection added:\n```json\n%s\n```\n", string(jsonProtection))
	}
	g.Printf(protectionSummary)

	// Create an issue with a summary of the created branch protection policy
	// and assign the issue to the desired set of users/
	issueRequest := &github.IssueRequest{
		Title:     github.String("Branch protection rules added to master branch"),
		Body:      github.String(protectionSummary),
		Assignees: &actingBranchPolicy.NotifyUsers,
	}
	issue, resp, err := g.Client.Issues.Create(g.Context, owner, repo, issueRequest)
	if (resp.StatusCode < 200 || resp.StatusCode > 299) || err != nil {
		g.Printf("creating issue returned status code: %v and err: %v\n", resp.StatusCode, err)
		return err
	}

	g.Printf("created notification issue: %s\n", *issue.HTMLURL)
	return nil
}
