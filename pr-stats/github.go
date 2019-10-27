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
		fmt.Println("Environment variable GITHUB_AUTH_TOKEN was not set...")
		fmt.Printf("GitHub access token: ")
		token, err := terminal.ReadPassword(syscall.Stdin)
		fmt.Printf("\n")
		if err != nil {
			panic(err)
		}
		githubToken = string(token)
	}

	return githubToken
}

func newGithubClient() (*github.Client, context.Context, context.CancelFunc) {
	githubToken := getGithubToken()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)

	ctx, cancel := context.WithCancel(context.Background())
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), ctx, cancel
}
