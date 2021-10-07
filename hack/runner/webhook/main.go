package main

import (
	/*
		"context"
		"encoding/json"
		"errors"
		"fmt"
		"net/http"
		"os"
		"strings"

		webhook "github.com/go-playground/webhooks/v6/github"
		github "github.com/google/go-github/github"
		"golang.org/x/oauth2"

		"github.com/aws/aws-sdk-go/aws"
		"github.com/aws/aws-sdk-go/aws/credentials"
		"github.com/aws/aws-sdk-go/aws/session"
		ec2 "github.com/aws/aws-sdk-go/service/ec2"
	*/

	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	webhook "github.com/go-playground/webhooks/v6/github"
)

/*
	Must start this service with the follow parameters defined:

	LISTEN_PORT

	GITHUB_TOKEN
	GITHUB_WEBHOOK_SECRET

	AWS_REGION
	AWS_ACCESS_KEY_ID
	AWS_SECRET_ACCESS_KEY
	AWS_AMI_ID
*/

const (
	webhookPath string = "/community-edition"
	versionPath string = "/version"

	defaultListenPort string = "8080"
	versionStr        string = "v0.0.1"
)

// Errors
var (
	// ErrClientInvalid client is not initialized
	ErrClientInvalid = errors.New("client is not initialized")
)

type Version struct {
	Version string `json:"version"`
}

func main() {
	// flags
	var port string
	if v := os.Getenv("LISTEN_PORT"); v != "" {
		port = v
	} else {
		port = defaultListenPort
	}
	fmt.Printf("Listening on: %s\n", port)

	var githubSecret string
	if v := os.Getenv("GITHUB_WEBHOOK_SECRET"); v != "" {
		githubSecret = v
	}

	// flags check
	if githubSecret == "" {
		fmt.Printf("Must supply GITHUB_WEBHOOK_SECRET\n")
		return
	}

	// set up GH webhook
	hook1, _ := webhook.New(webhook.Options.Secret(githubSecret))

	http.HandleFunc(webhookPath, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook1.Parse(r, webhook.WorkflowJobEvent)
		if err != nil {
			if err == webhook.ErrEventNotFound {
				fmt.Printf("Received event we weren't interested in. %v\n", err)
				http.Error(w, err.Error(), 404)
				return
			}
		}

		switch payload.(type) {

		case webhook.WorkflowJobPayload:

			workflowJob := payload.(webhook.WorkflowJobPayload)

			// Dump event
			fmt.Printf("%+v", workflowJob)

			// if one of these events... create a runner!
			if strings.EqualFold(workflowJob.Action, "queued") {
				fmt.Printf("Workflow is queued.  ID: %d, Name: %s",
					workflowJob.WorkflowJob.ID, workflowJob.WorkflowJob.Name)
				err = createRunner()
				if err != nil {
					fmt.Printf("createRunner failed. Err: %v\n", err)
				}
			} else if strings.EqualFold(workflowJob.Action, "in_progress") {
				fmt.Printf("Workflow is in_progress.  ID: %d, Name: %s",
					workflowJob.WorkflowJob.ID, workflowJob.WorkflowJob.Name)
			} else if strings.EqualFold(workflowJob.Action, "completed") {
				fmt.Printf("Workflow is completed.  ID: %d, Name: %s",
					workflowJob.WorkflowJob.ID, workflowJob.WorkflowJob.Name)
				// for _, tag := range workflowJob.WorkflowJob.

				// err = deleteRunner(workflowJob.)
				// if err != nil {

				// }
			}

		default:
			fmt.Printf("Unsupported Request Type. Dump: %v", payload)

		}
	})

	// generic version check
	http.HandleFunc(versionPath, func(w http.ResponseWriter, r *http.Request) {
		version := Version{
			Version: versionStr,
		}

		json.NewEncoder(w).Encode(version)
	})

	fmt.Printf("Starting server...")
	http.ListenAndServe(":"+port, nil)
}
