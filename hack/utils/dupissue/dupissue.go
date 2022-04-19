// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	// TemplateFullPathFilename filename
	TemplateFullPathFilename string = "issue.template"

	// DefaultIssueTemplate the default string
	DefaultIssueTemplate string = "This issue was originally filed in the TCE repo, but on a quick inspection, we think this issue probably is best suited for the Tanzu Framework repo.\n\nThe original issue filed in TCE: https://github.com/vmware-tanzu/community-edition/issues/%d\n\n\nThis is a copy and paste from that issue originally filed by @%s:\n\n%s"
)

func getClientWithEnvToken() (*github.Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return client, nil
}

func getIssue(client *github.Client, issueID int) (*github.Issue, error) {
	issueGH, _, err := client.Issues.Get(context.Background(), "vmware-tanzu", "community-edition", issueID)
	if err != nil {
		return nil, err
	}

	return issueGH, nil
}

func createIssue(client *github.Client, template string, issue *github.Issue) error {
	newBody := fmt.Sprintf(template, *issue.Number, *issue.User.Login, *issue.Body)

	issueRequest := &github.IssueRequest{
		Title: issue.Title,
		Body:  &newBody,
	}
	_, _, err := client.Issues.Create(context.Background(), "vmware-tanzu", "tanzu-framework", issueRequest)
	if err != nil {
		return err
	}

	return nil
}

func readTemplate() (string, error) {
	template, err := os.ReadFile(TemplateFullPathFilename)
	if err == nil {
		fmt.Printf("Using template from local directory\n")
		return string(template), nil
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	templateFromHomeDir := fmt.Sprintf("%s/%s", homedir, TemplateFullPathFilename)
	template, err = os.ReadFile(templateFromHomeDir)
	if err == nil {
		fmt.Printf("Using template from homedir\n")
		return string(template), nil
	}

	fmt.Printf("Using default template\n")
	return DefaultIssueTemplate, nil
}

func fileIssue(issueID int) error {
	client, err := getClientWithEnvToken()
	if err != nil {
		fmt.Printf("getClientWithEnvToken returned error: %v\n", err)
		return err
	}

	template, err := readTemplate()
	if err != nil {
		fmt.Printf("readTemplate failed. Err: %v\n", err)
		return err
	}

	issue, err := getIssue(client, issueID)
	if err != nil {
		fmt.Printf("getIssue failed: %v\n", err)
		return err
	}

	err = createIssue(client, template, issue)
	if err != nil {
		fmt.Printf("createIssue failed: %v\n", err)
		return err
	}

	return nil
}

func main() {
	var issueID int
	flag.IntVar(&issueID, "issue", 0, "The issue to duplicate to Tanzu Framework")
	flag.Parse()

	if issueID == 0 {
		fmt.Printf("A issue must be provided\n")
		return
	}

	err := fileIssue(issueID)
	if err != nil {
		return
	}

	fmt.Printf("Succeeded\n")
}
