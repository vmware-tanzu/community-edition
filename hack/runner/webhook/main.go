// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"net/http"
	"os"
	"path/filepath"

	webhook "github.com/go-playground/webhooks/v6/github"
	klog "k8s.io/klog/v2"
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

func initLogging() {
	// init klog
	klog.InitFlags(nil)

	// flags
	err := flag.Set("v", "2")
	if err != nil {
		panic(err)
	}
	err = flag.Set("logtostderr", "false")
	if err != nil {
		panic(err)
	}
	err = flag.Set("alsologtostderr", "true")
	if err != nil {
		panic(err)
	}

	exec, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exec = filepath.Dir(exec)
	logFile := filepath.Join(exec, "webhook.log")
	err = flag.Set("log_file", logFile)
	if err != nil {
		panic(err)
	}

	flag.Parse()
}

func main() {
	// init logging
	initLogging()

	// envvars
	var port string
	if v := os.Getenv("LISTEN_PORT"); v != "" {
		port = v
	} else {
		port = defaultListenPort
	}
	klog.Infof("Listening on: %s\n", port)

	var githubSecret string
	if v := os.Getenv("GITHUB_WEBHOOK_SECRET"); v != "" {
		githubSecret = v
	}

	// flags check
	if githubSecret == "" {
		klog.Errorf("Must supply GITHUB_WEBHOOK_SECRET\n")
		return
	}

	// set up GH webhook
	hook1, _ := webhook.New(webhook.Options.Secret(githubSecret))

	http.HandleFunc(webhookPath, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook1.Parse(r, webhook.PingEvent, webhook.PullRequestEvent, webhook.WorkflowRunEvent)
		if err != nil {
			if err == webhook.ErrEventNotFound {
				klog.Errorf("Received event we weren't interested in. %v\n", err)
				http.Error(w, err.Error(), 404)
				return
			}
		}

		switch payloadType := payload.(type) {
		case webhook.PingPayload:
			ping := payload.(webhook.PingPayload)
			handlePing(&ping)

		case webhook.PullRequestPayload:
			pullRequest := payload.(webhook.PullRequestPayload)
			handlePullRequest(&pullRequest)

		case webhook.WorkflowJobPayload:
			workflowJob := payload.(webhook.WorkflowJobPayload)
			handleWorkflowJob(&workflowJob)

		case webhook.WorkflowRunPayload:
			workflowRun := payload.(webhook.WorkflowRunPayload)
			handleWorkflowRun(&workflowRun)

		default:
			klog.Errorf("Unsupported Request Type. Type: %v, Dump: %v\n", payloadType, payload)
		}
	})

	// generic version check
	http.HandleFunc(versionPath, func(w http.ResponseWriter, r *http.Request) {
		version := Version{
			Version: versionStr,
		}

		encoder := json.NewEncoder(w)
		if encoder != nil {
			err := encoder.Encode(version)
			if err != nil {
				klog.Errorf("Encode Version failed. Err: %v\n", err)
			}
		}
	})

	klog.Infof("Starting server...\n\n")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		klog.Errorf("ListenAndServe failed. Err: %v\n", err)
	}
}
