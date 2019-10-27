# PR statistics

Welcome!
In this page we create a small program for analyzing the frequency of labels in
pull requests (PRs) in the
[kubernetes/kubernetes](https://github.com/kubernetes/kubernetes)
repo.

## Usage

This is a very simple program, you need only to have Go installed in your
system and a GitHub access token.

To run this program you can either pass the GitHub access token at runtime or
you can set it in the `GITHUB_AUTH_TOKEN` environment variable.
When the program is running you will be prompted for what to do, i.e.,
```
$ go build && ./pr-stats
Environment variable GITHUB_AUTH_TOKEN was not set...
GitHub access token:
```

## Configuration
Currently, this programm will gothrough all PRs that were created more than 3
months ago and will go through these starting with the oldest one.

If you want to modify this behaviour then you will have to modify the
parameters for the `PRStats` struct in the `newPRStats()` method in
[pr-stats.go](./pr-stats.go):
```go
func newPRStats() *PRStats {
  return &PRStats{
    ...
    // List all open PRs and sort them based on when they were created
    // beginning with the ones that were created first.
    PROpt: &github.PullRequestListOptions{
      State:       "open",
      Sort:        "created",
      Direction:   "asc",
      ...
    },
    //Only look at PR that are created before 3 months ago today.
    CreatedNMonthsAgo: 3,
    // Labels to remove from map.
    LabelsToIgnore: []string{`¯\_(ツ)_/¯`},
    ...
  }
}
```

## Output
Running this program will produce output of this sort:
```
$ go build && ./pr-stats
Environment variable GITHUB_AUTH_TOKEN was not set...
GitHub access token:
Currently at page 0 out of 11. Latest date observed was 2019-03-12 22:55:02 +0000 UTC
Currently at page 2 out of 11. Latest date observed was 2019-05-22 13:22:41 +0000 UTC
Currently at page 3 out of 11. Latest date observed was 2019-06-20 09:46:43 +0000 UTC
Currently at page 4 out of 11. Latest date observed was 2019-07-23 09:34:33 +0000 UTC
Currently at page 5 out of 11. Latest date observed was 2019-07-26 18:13:11 +0000 UTC
Currently at page 6 out of 11. Latest date observed was 2019-07-26 18:13:11 +0000 UTC
Currently at page 7 out of 11. Latest date observed was 2019-07-26 18:13:11 +0000 UTC
Currently at page 8 out of 11. Latest date observed was 2019-07-26 18:13:11 +0000 UTC
Currently at page 9 out of 11. Latest date observed was 2019-07-26 18:13:11 +0000 UTC
Currently at page 10 out of 11. Latest date observed was 2019-07-26 18:13:11 +0000 UTC
Currently at page 11 out of 0. Latest date observed was 2019-07-26 18:13:11 +0000 UTC
app is shutting down...
- 408   cncf-cla: yes
- 257   release-note-none
- 228   needs-priority
- 183   ok-to-test
- 149   kind/cleanup
- 130   release-note
- 126   kind/bug
- 111   size/XS
- 110   sig/node
- 96    sig/api-machinery
- 94    area/kubelet
- 94    lifecycle/stale
- 92    size/L
- 91    size/M
- 82    kind/feature
- 80    sig/testing
- 79    size/S
- 76    sig/cli
- 73    area/test
- 73    lgtm
- 72    needs-rebase
- 70    area/kubectl
- 70    sig/apps
- 68    lifecycle/rotten
- 65    priority/backlog
- 56    needs-ok-to-test
- 56    priority/important-soon
- 53    do-not-merge/hold
- 46    sig/cluster-lifecycle
- 44    sig/storage
- 43    kind/api-change
- 38    priority/important-longterm
- 37    area/apiserver
- 34    sig/scheduling
- 32    sig/network
- 30    approved
- 29    needs-kind
- 28    do-not-merge/work-in-progress
- 27    do-not-merge/release-note-label-needed
- 25    size/XL
- 20    sig/cloud-provider
- 18    sig/auth
- 17    kind/documentation
- 15    area/dependency
- 14    area/kubeadm
- 14    size/XXL
- 13    lifecycle/frozen
- 11    needs-sig
- 11    sig/architecture
- 9     area/cloudprovider
- 9     sig/instrumentation
- 9     sig/release
- 9     sig/windows
- 8     area/code-generation
- 7     cncf-cla: no
- 6     area/provider/aws
- 6     area/provider/gcp
- 6     area/release-eng
- 5     api-review
- 5     area/conformance
- 5     area/ipvs
- 5     sig/autoscaling
- 5     sig/scalability
- 4     kind/design
- 4     kind/failing-test
- 3     area/e2e-test-framework
- 3     do-not-merge/contains-merge-commits
- 3     kind/flake
- 2     area/provider/vmware
- 2     do-not-merge/cherry-pick-not-approved
- 2     wg/apply
- 1     area/app-lifecycle
- 1     area/client-libraries
- 1     area/code-organization
- 1     area/kubelet-api
- 1     area/provider/azure
- 1     do-not-merge
- 1     do-not-merge/invalid-commit-message
- 1     priority/awaiting-more-evidence
- 1     priority/critical-urgent
- 1     release-note-action-required
- 1     wg/component-standard
```
