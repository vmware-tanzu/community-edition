// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	"os"
	"strings"

	klog "k8s.io/klog/v2"

	github "github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

const (
	runnerOnline string = "online"
)

// Errors
var (
	// ErrClientCreateFailed client create failed
	ErrClientCreateFailed = errors.New("client create failed")

	// ErrRunnerOffline Runner is offline
	ErrRunnerOffline = errors.New("runner is offline")
)

// get github client
func getGitHubClientWithEnvToken() (*github.Client, error) {
	token := os.Getenv("GITHUB_TOKEN")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	if tc == nil {
		klog.Errorf("oauth2.NewClient failed\n")
		return nil, ErrClientCreateFailed
	}

	// new GitHub client
	client := github.NewClient(tc)
	if client == nil {
		klog.Errorf("github.NewClient failed\n")
		return nil, ErrClientCreateFailed
	}

	return client, nil
}

func createRunnerToken(client *github.Client) (string, error) {
	if client == nil {
		err := ErrClientInvalid
		klog.Errorf("Client == nil. Err: %v\n", err)
		return "", err
	}

	token, _, err := client.Actions.CreateRegistrationToken(context.Background(), "vmware-tanzu", "community-edition")
	if err != nil {
		klog.Errorf("Actions.CreateRegistrationToken returned Err: %v\n", err)
		return "", err
	}

	klog.Infof("Runner token created successfully\n")
	return *token.Token, nil
}

func getGitHubRunner(client *github.Client, runnerName string) (*github.Runner, error) {
	klog.Infof("getGitHubRunner(%s)\n", runnerName)

	if client == nil {
		err := ErrClientInvalid
		klog.Errorf("Client == nil. Err: %v\n", err)
		return nil, err
	}

	opts := &github.ListOptions{}
	runners, _, err := client.Actions.ListRunners(context.Background(), "vmware-tanzu", "community-edition", opts)
	if err != nil {
		klog.Errorf("Actions.ListRunners failed. Err: %v\n", err)
		return nil, err
	}

	if runners.TotalCount > 0 {
		for _, runner := range runners.Runners {
			klog.V(4).Infof("Runner: %s\n", *runner.Name)

			if !strings.EqualFold(runnerName, *runner.Name) {
				continue
			}

			klog.Infof("Runner found! ID: %d\n", *runner.ID)
			return runner, nil
		}
	}

	klog.Infof("Runner is OFFLINE...\n")
	return nil, ErrRunnerOffline
}

func deleteGitHubRunnerByName(client *github.Client, runnerName string) error {
	klog.Infof("deleteGitHubRunnerByName(%s)\n", runnerName)

	if client == nil {
		err := ErrClientInvalid
		klog.Errorf("Client == nil. Err: %v\n", err)
		return err
	}
	if runnerName == "" {
		err := ErrClientInvalid
		klog.Errorf("runnerName is empty. Err: %v\n", err)
		return err
	}

	runner, err := getGitHubRunner(client, runnerName)
	if err != nil {
		klog.Errorf("getGitHubRunner failed. Err: %v\n", err)
		return err
	}

	return deleteGitHubRunnerByID(client, *runner.ID)
}

func deleteGitHubRunnerByID(client *github.Client, runnerID int64) error {
	klog.Infof("deleteGitHubRunnerByID(%d)\n", runnerID)

	if client == nil {
		err := ErrClientInvalid
		klog.Errorf("Client == nil. Err: %v\n", err)
		return err
	}

	_, err := client.Actions.RemoveRunner(context.Background(), "vmware-tanzu", "community-edition", runnerID)
	if err != nil {
		klog.Errorf("Actions.RemoveRunner failed. Err: %v\n", err)
		return err
	}

	klog.Infof("Runner has been deleted!")
	return nil
}
