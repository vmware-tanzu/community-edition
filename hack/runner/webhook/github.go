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
	runnerIdle   string = "idle"
	runnerOnline string = "online"
	runnerActive string = "active"
	// runnerOffline string = "offline"
)

// Errors
var (
	// ErrClientCreateFailed client create failed
	ErrClientCreateFailed = errors.New("client create failed")

	// ErrRunnerOffline Runner is offline
	ErrRunnerOffline = errors.New("runner is offline")

	// ErrRunnerIsBusy Runner is busy
	ErrRunnerIsBusy = errors.New("runner is busy")
)

// GitHub object
type GitHub struct {
	client *github.Client
}

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

// NewGitHub generates a new GH object
func NewGitHub() (*GitHub, error) {
	client, err := getGitHubClientWithEnvToken()
	if err != nil {
		klog.Errorf("getGitHubClientWithEnvToken failed. Err: %v\n", err)
		return nil, err
	}

	mygh := &GitHub{
		client: client,
	}
	return mygh, nil
}

// CreateRunnerToken create a token
func (g *GitHub) CreateRunnerToken() (string, error) {
	token, _, err := g.client.Actions.CreateRegistrationToken(context.Background(), "vmware-tanzu", "community-edition")
	if err != nil {
		klog.Errorf("Actions.CreateRegistrationToken returned Err: %v\n", err)
		return "", err
	}

	klog.Infof("Runner token created successfully\n")
	return *token.Token, nil
}

// GetGitHubRunners get all runner
func (g *GitHub) GetGitHubRunners() (*github.Runners, error) {
	opts := &github.ListOptions{}
	runners, _, err := g.client.Actions.ListRunners(context.Background(), "vmware-tanzu", "community-edition", opts)
	if err != nil {
		klog.Errorf("Actions.ListRunners failed. Err: %v\n", err)
		return nil, err
	}

	klog.Infof("GetGitHubRunners succeeded!\n")
	return runners, nil
}

// GetGitHubRunner get runner
func (g *GitHub) GetGitHubRunner(runnerName string) (*github.Runner, error) {
	klog.Infof("getGitHubRunner(%s)\n", runnerName)

	opts := &github.ListOptions{}
	runners, _, err := g.client.Actions.ListRunners(context.Background(), "vmware-tanzu", "community-edition", opts)
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

// DeleteGitHubRunnerByName delete by name
func (g *GitHub) DeleteGitHubRunnerByName(runnerName string) error {
	klog.Infof("deleteGitHubRunnerByName(%s)\n", runnerName)

	if runnerName == "" {
		err := ErrClientInvalid
		klog.Errorf("runnerName is empty. Err: %v\n", err)
		return err
	}

	runner, err := g.GetGitHubRunner(runnerName)
	if err != nil {
		klog.Errorf("getGitHubRunner failed. Err: %v\n", err)
		return err
	}

	if *runner.Busy {
		klog.Infof("Runner %s is working on another job\n", runnerName)
		return ErrRunnerIsBusy
	}

	return g.DeleteGitHubRunnerByID(*runner.ID)
}

// DeleteGitHubRunnerByID delete by ID
func (g *GitHub) DeleteGitHubRunnerByID(runnerID int64) error {
	klog.Infof("deleteGitHubRunnerByID(%d)\n", runnerID)

	_, err := g.client.Actions.RemoveRunner(context.Background(), "vmware-tanzu", "community-edition", runnerID)
	if err != nil {
		klog.Errorf("Actions.RemoveRunner failed. Err: %v\n", err)
		return err
	}

	klog.Infof("Runner has been deleted!")
	return nil
}
