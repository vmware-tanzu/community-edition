package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	github "github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

const (
	runnerOnline string = "online"
)

// Errors
var (
	// ErrRunnerNotFound Runner not found
	ErrRunnerNotFound = errors.New("Runner not found")
	// ErrRunnerOffline Runner is offline
	ErrRunnerOffline = errors.New("Runner is offline")
)

// get github client
func getGitHubClientWithEnvToken() (*github.Client, error) {
	var token string
	if v := os.Getenv("GITHUB_TOKEN"); v != "" {
		token = v
	}

	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	// new GitHub client
	client := github.NewClient(tc)

	return client, nil
}

func createRunnerToken(client *github.Client) (string, error) {
	if client == nil {
		err := ErrClientInvalid
		fmt.Printf("Client == nil. Err: %v\n", err)
		return "", err
	}

	token, _, err := client.Actions.CreateRegistrationToken(context.Background(), "vmware-tanzu", "community-edition")
	if err != nil {
		fmt.Printf("Actions.CreateRegistrationToken returned Err: %v\n", err)
		return "", err
	}

	fmt.Printf("Runner token created successfully\n")
	return *token.Token, nil
}

func getGitHubRunner(client *github.Client, runnerName string) (*github.Runner, error) {
	if client == nil {
		err := ErrClientInvalid
		fmt.Printf("Client == nil. Err: %v\n", err)
		return nil, err
	}

	opts := &github.ListOptions{}
	runners, _, err := client.Actions.ListRunners(context.Background(), "vmware-tanzu", "community-edition", opts)
	if err != nil {
		fmt.Printf("Actions.ListRunners failed. Err: %v\n", err)
		return nil, err
	}

	if runners.TotalCount == 0 {
		fmt.Printf("No runners found...")
		return nil, ErrRunnerNotFound
	}

	for _, runner := range runners.Runners {
		fmt.Printf("Runner: %s\n", *runner.Name)

		if !strings.EqualFold(runnerName, *runner.Name) {
			continue
		}

		fmt.Printf("Runner found! ID: %d\n", *runner.ID)
		return runner, nil
	}

	fmt.Printf("Runner is OFFLINE...")
	return nil, ErrRunnerOffline
}

func deleteGitHubRunnerByName(client *github.Client, runnerName string) error {
	if client == nil {
		err := ErrClientInvalid
		fmt.Printf("Client == nil. Err: %v\n", err)
		return err
	}
	if runnerName == "" {
		err := ErrClientInvalid
		fmt.Printf("runnerName is empty. Err: %v\n", err)
		return err
	}

	runner, err := getGitHubRunner(client, runnerName)
	if err != nil {
		fmt.Printf("getGitHubRunner failed. Err: %v\n", err)
		return err
	}

	return deleteGitHubRunnerByID(client, *runner.ID)
}

func deleteGitHubRunnerByID(client *github.Client, runnerID int64) error {
	if client == nil {
		err := ErrClientInvalid
		fmt.Printf("Client == nil. Err: %v\n", err)
		return err
	}

	_, err := client.Actions.RemoveRunner(context.Background(), "vmware-tanzu", "community-edition", runnerID)
	if err != nil {
		fmt.Printf("Actions.RemoveRunner failed. Err: %v\n", err)
		return err
	}

	fmt.Printf("Runner has been deleted!")
	return nil
}
