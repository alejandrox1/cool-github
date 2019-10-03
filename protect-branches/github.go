/*
   Define and read the configuration file containing a GitHub access token and
   a webhook secret.
   Here we also define other utility functions for interacting with the GitHub
   API.
*/
package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

// githubSecrets represents the config file detailing the access token and
// webhook secret. This is REQUIRED for the application to run.
type githubSecrets struct {
	Token         string `yaml:"token"`
	WebhookSecret string `yaml:"webhookSecret"`
}

// readGithubSecrets reads a configuration file and returns the access token
// and the webhook secret.
func readGithubSecrets(config string) (string, string, error) {
	yamlFile, err := ioutil.ReadFile(config)
	if err != nil {
		return "", "", err
	}

	gs := &githubSecrets{}
	if err := yaml.Unmarshal(yamlFile, gs); err != nil {
		return "", "", err
	}

	if gs != nil && gs.Token != "" && gs.WebhookSecret != "" {
		return gs.Token, gs.WebhookSecret, nil
	} else {
		return "", "", fmt.Errorf("one or more values were missing from %s", config)
	}
}

// newGithubClient creates an instance of the GitHub API client.
func newGithubClient() (*github.Client, context.Context) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), ctx
}
