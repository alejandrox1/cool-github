package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v28/github"
)

type PRStats struct {
	client        *github.Client
	ctx           context.Context
	CtxCancelFunc context.CancelFunc

	// User-specified parameters.
	PROpt             *github.PullRequestListOptions
	CreatedNMonthsAgo int
	CreatedBefore     time.Time
	LabelsToIgnore    []string

	// Parameters for analyze label frequency.
	LabelFreq        map[string]int
	OrderedLabelFreq []kv
}

func newPRStats() *PRStats {
	client, ctx, ctxCancel := newGithubClient()
	return &PRStats{
		client:        client,
		ctx:           ctx,
		CtxCancelFunc: ctxCancel,

		// List all open PRs and sort them based on when they were created
		// beginning with the ones that were created first.
		PROpt: &github.PullRequestListOptions{
			State:       "open",
			Sort:        "created",
			Direction:   "asc",
			ListOptions: github.ListOptions{PerPage: 100},
		},
		//Only look at PR that are created before 3 months ago today.
		CreatedNMonthsAgo: 3,
		// Labels to remove from map.
		LabelsToIgnore: []string{`¯\_(ツ)_/¯`},

		LabelFreq:        make(map[string]int),
		OrderedLabelFreq: []kv{},
	}
}

func (p *PRStats) initCreatedBefore() {
	today := time.Now()
	p.CreatedBefore = today.AddDate(0, -p.CreatedNMonthsAgo, 0)
}

func (p *PRStats) AnalyzePRLabelFreq(owner, repo string) error {
	p.initCreatedBefore()

	var lastCreatedAt time.Time
	for {
		// Paginate through all open k/k PRs.
		prs, resp, err := p.client.PullRequests.List(p.ctx, owner, repo, p.PROpt)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				fmt.Printf("hit rate limit: %v\n", err)
				time.Sleep(1 * time.Minute)
				continue
			}
			return err
		}

		// Go through each PR...
		for _, pr := range prs {
			// IF PR was created before the 'CreatedBefore' date...
			if (*pr.CreatedAt).Before(p.CreatedBefore) {
				// Then count each label in the PR.
				for _, label := range pr.Labels {
					p.LabelFreq[*label.Name]++
				}
				lastCreatedAt = *pr.CreatedAt
			}
		}

		fmt.Printf("Currently at page %d out of %d. Latest date observed was %v\n", p.PROpt.Page, resp.LastPage, lastCreatedAt)
		if resp.NextPage == 0 {
			break
		}
		p.PROpt.Page = resp.NextPage
	}

	return nil
}

func (p *PRStats) SummarizeResults() {
	for _, label := range p.LabelsToIgnore {
		if _, ok := p.LabelFreq[label]; ok {
			delete(p.LabelFreq, label)
		}
	}
	p.OrderedLabelFreq = orderMap(p.LabelFreq)

	for _, kv := range p.OrderedLabelFreq {
		val := strconv.Itoa(kv.Value)
		fmt.Printf("- %s %s %s\n", val, strings.Repeat(" ", 4-len(val)), kv.Key)
	}
}
