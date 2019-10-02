/*
 */
package main

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/google/go-github/v28/github"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/oauth2"
)

var (
	githubTokenEnv = "GITHUB_AUTH_TOKEN"
)

func getGithubToken() string {
	var githubToken string

	githubToken = os.Getenv(githubTokenEnv)
	if githubToken == "" {
		fmt.Printf("GitHub access token: ")
		token, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			panic(err)
		}
		githubToken = string(token)
	}

	return githubToken
}

func newGithubClient() (*github.Client, context.Context) {
	githubToken := getGithubToken()
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), ctx
}
