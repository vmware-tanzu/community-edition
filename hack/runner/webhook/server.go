// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	klog "k8s.io/klog/v2"

	ec2 "github.com/aws/aws-sdk-go/service/ec2"
	webhook "github.com/go-playground/webhooks/v6/github"
	github "github.com/google/go-github/v39/github"
)

const (
	workflowBuildStaging    string = "Build - Create Dev/Staging"
	workflowCheckMain       string = "Check - Main (All tests)"
	workflowCheckDockerMgmt string = "Check - Management Docker Cluster"
	workflowCheckDockerSa   string = "Check - Standalone Docker Cluster"
	workflowReleaseFake     string = "Release - Create Test/Fake Release"
	workflowReleaseNonGa    string = "Release - Create GA"
	workflowReleaseGa       string = "Release - Create RC, Beta, Alpha"

	defaultSleepHeadStart       int64 = 30
	defaultSleepBetweenPoll     int64 = 10
	defaultNumOfTimesToPoll     int   = 30
	defaultmustHaveStatusBefore int   = 10
	defaultNumOfTimesToRetry    int   = 3
)

// Errors
var (
	// ErrCreateAndConnectRunner failed to create and connect the runner
	ErrCreateAndConnectRunner = errors.New("failed to create and connect the runner")
)

func createOnlineRunner(ghClient *github.Client, ec2Client *ec2.EC2, uniqueID string) error {
	token, err := createRunnerToken(ghClient)
	if err != nil {
		klog.Errorf("getGitHubClientWithEnvToken failed. Err: %v\n", err)
		return err
	}

	instanceID, err := createEc2Runner(ec2Client, uniqueID, token)
	if err != nil {
		klog.Errorf("createEc2Runner failed. Err: %v\n", err)
		return err
	}

	klog.Infof("Giving head start...\n")
	time.Sleep(time.Duration(defaultSleepHeadStart) * time.Second)

	succeeded := false
	for i := 0; i < defaultNumOfTimesToPoll; i++ {
		runner, err := getGitHubRunner(ghClient, uniqueID)
		if err == nil {
			klog.Infof("Status: %s\n", *runner.Status)
			if !strings.EqualFold(*runner.Status, runnerOnline) && i > defaultmustHaveStatusBefore {
				klog.Infof("The node should have already returned some status... retry\n")
				break
			}

			if strings.EqualFold(*runner.Status, runnerOnline) {
				klog.Infof("Succeeded!\n")
				succeeded = true
				break
			}
		} else if err != ErrRunnerOffline {
			klog.Errorf("getGitHubRunner failed. Err: %v\n", err)
			return err
		}

		klog.Infof("Attempt poll %d... sleeping\n", i)
		time.Sleep(time.Duration(defaultSleepBetweenPoll) * time.Second)
	}

	if !succeeded {
		klog.Errorf("createOnlineRunner failed. Delete instance %s\n", instanceID)

		err = deleteEc2Instance(ec2Client, instanceID)
		if err != nil {
			klog.Errorf("deleteEc2Instance failed. Err: %v\n", err)
		}

		return ErrCreateAndConnectRunner
	}

	return nil
}

func createRunner(uniqueID string) error {
	ghClient, err := getGitHubClientWithEnvToken()
	if err != nil {
		klog.Errorf("getGitHubClientWithEnvToken failed. Err: %v\n", err)
		return err
	}
	ec2Client, err := getAwsClientWithEnvToken()
	if err != nil {
		klog.Errorf("getAwsClientWithEnvToken failed. Err: %v\n", err)
		return err
	}

	for i := 0; i < defaultNumOfTimesToRetry; i++ {
		err = createOnlineRunner(ghClient, ec2Client, uniqueID)
		if err == nil {
			klog.Infof("createOnlineRunner succeeded!\n")
			break
		}

		klog.Infof("createOnlineRunner failed... retying. Err: %v\n", err)
	}

	return err
}

func deleteRunner(uniqueID string) error {
	ghClient, err := getGitHubClientWithEnvToken()
	if err != nil {
		klog.Errorf("getGitHubClientWithEnvToken failed. Err: %v\n", err)
		return err
	}
	ec2Client, err := getAwsClientWithEnvToken()
	if err != nil {
		klog.Errorf("getAwsClientWithEnvToken failed. Err: %v\n", err)
		return err
	}

	err = deleteGitHubRunnerByName(ghClient, uniqueID)
	if err != nil {
		// Just a warning because of the new self host ephemeral
		klog.Infof("deleteGitHubRunnerByName failed. Err: %v\n", err)
	}

	err = deleteEc2InstanceByName(ec2Client, uniqueID)
	if err != nil {
		klog.Errorf("deleteEc2InstanceByName failed. Err: %v\n", err)
		return err
	}

	klog.Infof("deleteRunner succeeded\n")
	return nil
}

func handlePing(ping *webhook.PingPayload) {
	// Dump event
	klog.V(6).Infof("---------------------- START DUMP EVENT ----------------------\n\n\n")
	klog.V(6).Infof("%+v\n\n\n", ping)
	klog.V(6).Infof("---------------------- END DUMP EVENT ----------------------\n\n\n")

	klog.Infof("Received ping event %d...\n", ping.HookID)
}

func handlePullRequest(pullRequest *webhook.PullRequestPayload) {
	// Dump event
	klog.V(6).Infof("---------------------- START DUMP EVENT ----------------------\n\n\n")
	klog.V(6).Infof("%+v\n\n\n", pullRequest)
	klog.V(6).Infof("---------------------- END DUMP EVENT ----------------------\n\n\n")

	if strings.EqualFold(pullRequest.Action, "assigned") ||
		strings.EqualFold(pullRequest.Action, "opened") ||
		strings.EqualFold(pullRequest.Action, "synchronize") ||
		strings.EqualFold(pullRequest.Action, "reopened") {
		klog.Infof("PR of interest. ID: %d, Number: %d, Title: %s\n",
			pullRequest.PullRequest.ID, pullRequest.PullRequest.Number, pullRequest.PullRequest.Title)
	}
}

func handleWorkflowJob(workflowJob *webhook.WorkflowJobPayload) {
	// Dump event
	klog.V(6).Infof("---------------------- START DUMP EVENT ----------------------\n\n\n")
	klog.V(6).Infof("%+v\n\n\n", workflowJob)
	klog.V(6).Infof("---------------------- END DUMP EVENT ----------------------\n\n\n")
}

func handleWorkflowRun(workflowRun *webhook.WorkflowRunPayload) {
	// Dump event
	klog.V(6).Infof("---------------------- START DUMP EVENT ----------------------\n\n\n")
	klog.V(6).Infof("%+v\n\n\n", workflowRun)
	klog.V(6).Infof("---------------------- END DUMP EVENT ----------------------\n\n\n")

	workflowName := workflowRun.WorkflowRun.Name

	switch workflowName {
	case workflowBuildStaging:
		doWorkflowRun(workflowRun)

	case workflowCheckMain:
		doWorkflowRun(workflowRun)

	case workflowCheckDockerMgmt:
		doWorkflowRun(workflowRun)

	case workflowCheckDockerSa:
		doWorkflowRun(workflowRun)

	case workflowReleaseFake:
		doWorkflowRun(workflowRun)

	case workflowReleaseNonGa:
		doWorkflowRun(workflowRun)

	case workflowReleaseGa:
		doWorkflowRun(workflowRun)

	default:
		klog.V(6).Infof("Don't need self hosted-runner: %s\n", workflowName)
	}
}

func doWorkflowRun(workflowRun *webhook.WorkflowRunPayload) {
	workflowID := workflowRun.WorkflowRun.ID
	workflowName := workflowRun.WorkflowRun.Name

	uniqueRunnerName := fmt.Sprintf("id-%d-%d", workflowRun.WorkflowRun.ID, workflowRun.WorkflowRun.RunNumber)
	klog.Infof("uniqueRunnerName: %s\n", uniqueRunnerName)

	// if one of these events... create a runner!
	if strings.EqualFold(workflowRun.Action, "requested") {
		klog.Infof("Workflow is requested.  ID: %d, Name: %s\n", workflowID, workflowName)
		err := createRunner(uniqueRunnerName)
		if err != nil {
			klog.Errorf("createRunner failed. Err: %v\n", err)
		}
	} else if strings.EqualFold(workflowRun.Action, "completed") {
		klog.Infof("Workflow is completed.  ID: %d, Name: %s\n", workflowID, workflowName)
		err := deleteRunner(uniqueRunnerName)
		if err != nil {
			klog.Errorf("deleteRunner failed. Err: %v\n", err)
		}
	}
}
