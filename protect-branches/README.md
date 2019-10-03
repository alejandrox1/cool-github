# Protecting Branches Upon Creation

Welcome!

If you find yourself here then it means that you want to use a solution for
managing your GitHub organization's repositories.

The project in this repository will allow you to stand up a web server that will
perform an action whenever a new repository in your GitHub organization's page
is created.
Once a repository is created, this program will create a branch protection
policy (rule).
Whenever a branch protection policy is created an issue will also be created in
the given repository.
The issue will detail the branch protection policy that was created and it will
assign the issue to a predefined set of users.

## Usage
### Configuration
In order to run this service you will need two configuration files.
One in which you specify a GitHub access token and a webhook secret.
For more information about webhook secrets, please see
[Securing your webhooks](https://developer.github.com/webhooks/securing/).
These parameters are **REQUIRED** and should be specified in a YAML format as
follows:

`github_secrets.yaml`
```yaml
token: <access token goes here>
webhookSecret: <webhook secret here>
```

You can also specify a branch protection policy in YAML, we include an example
in this repo, see [`branch_policy.yaml`](./branch_policy.yaml).
This configuration file has the following structure:

```yaml
# Users who will be assigned to the issue detailing the created policy.
notifyUsers: ["alejandrox1"]
# Enforce all configured restrictions for administrators..
enforceAdmins: true
# Require signed commits.
requireSignatures: true
requiredPullRequestReviews:
  # New reviewable commits will dismiss pull request review approvals.
  dismissStaleReviews: true
  # Require this number of approving reviews.
  requiredApprovingReviewCount: 2
requireStatusChecks:
  # Require status checks to pass before merging.
  strict: true
  # Specify contexts, [statuses](https://developer.github.com/v3/repos/statuses/).
  contexts: []
```

The default branch protection policy follows the example given in
[`branch_policy.yaml`](./branch_policy.yaml).

### Running This Service

To run this service, you have to specify the configuration file containg the
token access and the webhook services and can optionally modify the branch
protection policy as follows:

```bash
./protect-branches -githubSecrets ~/.github/github_secrets.yaml -branchPolicy branch_policy.yaml
```

See `--help` for more info.
```bash
./protect-branches --help
Usage of ./protect-branches:
  -branchPolicy string
    	coniguration file for branch protection policy (default "branch_policy.yaml")
  -githubSecrets string
    	configuration file with access token and webhook secret
```

## Developers

If you are interested in building upon this project then keep on reading...
:rocket:

This application is composed of a webserver that listens for events from a
GitHub webhook.
When you configure a webhook for an organization you will have the opportunity
to specify that you want to listen for events whenever a repository is created.

In order to test this functionality, we recommend the use of ngrok as a reverse
proxy.
When running ngrok, you will obtain a HTTP and HTTPS urls that will be
reachable from the internet.
These urls will redirect the traffic to your localhost where you can run this
service.
Furthermore, ngrok has an endpoint for inspecting request,
http://127.0.0.1:4040/inspect/http.
In this endpoint you can replay request and iteratively develop your
application from your machine.

SO let's get started!

### ngrok

Please follow the oficial
[ngrok installation instructions](https://ngrok.com/download) to get it working
on your system.

Once you have done that, you can start ngrok:
```bash
ngrok http 8080
```

Note that we are redicrecting all the outside traffic to `localhost:8080`.
If you want to change the port then you will have to change the `listenAddr`
variable in [`main.go`](./main.go).

### Configuring Your Webhook

Now that you have ngrok running, save the https url that is provided for you,
we will need it to configure the webhook.

Go to GitHub, to your organization's page.
Click on setting, and go to webhooks.
In there, click on "add webhook".

In the "payload URL", copy the url that was given to you by ngrok and add
`/webhook` as a path (this is the default route specified in
[server.go](./server.go).
You should have something like `https://8fe9a5i6.ngrok.io/webhook`.
Below that, create some string for the secret (i.e., "my-secret-string).
Make sure you save this, this is your webhook secret and you have to have it
int eh configuration file along with your GitHub access token.

Below that, scroll to "Which events would you like to trigger this webhook?".
Click on "Let me select individual events.".
And now click on "Repositories".
We are telling GitHub that we want to be notified when a repository is created,
deleted, archived, unarchived, publicized, privatized, edited, renamed, or
transferred.
And save!

### Run This Program

AT this point, you can build and run this program:
```branch
export GO111MODULE=on; go build -o protect-branches \
    && ./protect-branches -githubSecrets ~/.github/github_secrets.yaml -branchPolicy branch_policy.yaml
```

You will see the log of events on the console.
And if you want, remeber you can replay request by going to
http://127.0.0.1:4040/inspect/http.
