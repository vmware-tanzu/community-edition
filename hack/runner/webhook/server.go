// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	klog "k8s.io/klog/v2"

	webhook "github.com/go-playground/webhooks/v6/github"
)

const (
	selfhostedRunnerPrefix string = "id-"

	// workflowJobQueued         string = "queued"
	// workflowJobInProgress     string = "in_progress"
	workflowJobCompleted      string = "completed"
	workflowJobSetupRunner    string = "Start self-hosted EC2 runner"
	workflowJobTeardownRunner string = "Stop self-hosted EC2 runner"

	defaultSleepHeadStart       int64 = 30
	defaultSleepBetweenPoll     int64 = 10
	defaultNumOfTimesToPoll     int   = 30
	defaultmustHaveStatusBefore int   = 10
	defaultNumOfTimesToRetry    int   = 3

	defaultGetWorkflowRunTimeout     int   = 3
	defaultGetWorkflowRunRetry       int   = 3
	defaultGetWorkflowRunBetweenPoll int64 = 2
	defaultCleanUpOrphans            int64 = 3600
)

// Errors
var (
	// ErrCreateAndConnectRunner failed to create and connect the runner
	ErrCreateAndConnectRunner = errors.New("failed to create and connect the runner")
)

func createOnlineRunner(g *GitHub, a *Aws, uniqueID string) error {
	klog.Infof("Check to see if %s runner exists...\n", uniqueID)
	runner, err := g.GetGitHubRunner(uniqueID)
	if err == nil {
		klog.Infof("Runner %s exists. Status: %s\n", uniqueID, *runner.Status)
		switch *runner.Status {
		case runnerIdle:
			klog.Infof("Runner %s is idle. Skip creation!\n", uniqueID)
			return nil

		case runnerOnline:
			klog.Infof("Runner %s is online. Skip creation!\n", uniqueID)
			return nil

		case runnerActive:
			klog.Infof("Runner %s is active. Skip creation!\n", uniqueID)
			return nil

		default:
			klog.Infof("Runner %s is %s. Continue with creation!\n", uniqueID, *runner.Status)
		}
	}

	token, err := g.CreateRunnerToken()
	if err != nil {
		klog.Errorf("getGitHubClientWithEnvToken failed. Err: %v\n", err)
		return err
	}

	instanceID, err := a.CreateEc2Runner(uniqueID, token)
	if err != nil {
		klog.Errorf("createEc2Runner failed. Err: %v\n", err)
		return err
	}

	klog.Infof("Giving head start...\n")
	time.Sleep(time.Duration(defaultSleepHeadStart) * time.Second)

	succeeded := false
	for i := 0; i < defaultNumOfTimesToPoll; i++ {
		runner, err := g.GetGitHubRunner(uniqueID)
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

		err = a.DeleteEc2Instance(instanceID)
		if err != nil {
			klog.Errorf("deleteEc2Instance failed. Err: %v\n", err)
		}

		return ErrCreateAndConnectRunner
	}

	klog.Infof("createOnlineRunner succeeded. instance: %s\n", instanceID)
	return nil
}

func createRunner(uniqueID string) error {
	ghClient, err := NewGitHub()
	if err != nil {
		klog.Errorf("NewGitHub failed. Err: %v\n", err)
		return err
	}
	ec2Client, err := NewAws()
	if err != nil {
		klog.Errorf("NewAws failed. Err: %v\n", err)
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
	// insert small delay just like in create. this allows the runner agent time
	// to become not busy
	klog.Infof("Giving settling time...\n")
	time.Sleep(time.Duration(defaultSleepHeadStart) * time.Second)

	ghClient, err := NewGitHub()
	if err != nil {
		klog.Errorf("NewGitHub failed. Err: %v\n", err)
		return err
	}
	ec2Client, err := NewAws()
	if err != nil {
		klog.Errorf("NewAws failed. Err: %v\n", err)
		return err
	}

	err = ghClient.DeleteGitHubRunnerByName(uniqueID)
	if err == ErrRunnerIsBusy {
		klog.Infof("Runner is busy working on other requests. Skipping deletion of runner and ec2 instance.\n")
		return nil
	}
	if err != nil {
		// Just a warning because of the new self host ephemeral
		// Do not error this function out because we need to delete the instance
		klog.Infof("deleteGitHubRunnerByName failed. Err: %v\n", err)
	}

	err = ec2Client.DeleteEc2InstanceByName(uniqueID)
	if err != nil {
		klog.Errorf("deleteEc2InstanceByName failed. Err: %v\n", err)
		return err
	}

	klog.Infof("deleteRunner succeeded\n")
	return nil
}

func getWorkflowRunOnce(uri string) (*webhook.WorkflowRunPayload, error) {
	klog.V(6).Infof("getWorkflowRunOnce(%s)\n", uri)

	client := http.Client{
		Timeout: time.Duration(defaultGetWorkflowRunTimeout) * time.Second,
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, uri, http.NoBody)
	if err != nil {
		klog.Errorf("http.NewRequest failed. Err: %v\n", err)
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		klog.Errorf("client.Do failed. Err: %v\n", err)
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		klog.Errorf("io.ReadAll failed. Err: %v\n", err)
		return nil, err
	}

	workflowRunPayload := webhook.WorkflowRunPayload{}
	err = json.Unmarshal(body, &workflowRunPayload.WorkflowRun)
	if err != nil {
		klog.Errorf("json.Unmarshal failed. Err: %v\n", err)
		return nil, err
	}

	klog.V(6).Infof("getWorkflowRunOnce(%s) succeeded\n", uri)
	return &workflowRunPayload, nil
}

func getWorkflowRun(uri string) (*webhook.WorkflowRunPayload, error) {
	klog.Infof("getWorkflowRun(%s)\n", uri)

	var errRet error
	for i := 0; i < defaultGetWorkflowRunRetry; i++ {
		if i != 0 {
			klog.Infof("Sleeping... Before retrying getWorkflowRunOnce\n")
			time.Sleep(time.Duration(defaultGetWorkflowRunBetweenPoll) * time.Second)
		}

		workflowRunPayload, err := getWorkflowRunOnce(uri)
		if err == nil {
			klog.Infof("getWorkflowRunOnce succeeded!\n")
			return workflowRunPayload, nil
		}

		errRet = err
		klog.Infof("getWorkflowRunOnce failed. Err: %v\n", err)
	}

	klog.Errorf("getWorkflowRunOnce failed. Err: %v\n", errRet)
	return nil, errRet
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

func handleWorkflowJob(workflowJob *webhook.WorkflowJobPayload) error {
	// Dump event
	klog.V(6).Infof("---------------------- START DUMP EVENT ----------------------\n\n\n")
	klog.V(6).Infof("%+v\n\n\n", workflowJob)
	klog.V(6).Infof("---------------------- END DUMP EVENT ----------------------\n\n\n")

	workflowName := workflowJob.WorkflowJob.Name
	klog.Infof("workflowName: %s, Action: %s\n", workflowName, workflowJob.Action)

	switch workflowName {
	case workflowJobSetupRunner:
		klog.Infof("doWorkflowJob using create\n")
		return doWorkflowJob(workflowJob, true)

	case workflowJobTeardownRunner:
		klog.Infof("doWorkflowJob using delete\n")
		return doWorkflowJob(workflowJob, false)

	default:
		klog.V(6).Infof("No create/delete for self hosted-runner: %s\n", workflowName)
		return nil
	}
}

func doWorkflowJob(workflowJob *webhook.WorkflowJobPayload, create bool) error {
	if (create && !strings.EqualFold(workflowJob.Action, workflowJobCompleted)) ||
		(!create && !strings.EqualFold(workflowJob.Action, workflowJobCompleted)) {
		klog.Infof("doWorkflowJob create: %t, status %s. Skipping!\n", create, workflowJob.Action)
		return nil
	}

	// get the WorkflowRun which represents the entire workflow end-to-end
	workflowRun, err := getWorkflowRun(workflowJob.WorkflowJob.RunURL)
	if err != nil {
		klog.Errorf("getWorkflowRun failed. Err: %v\n", err)
		return err
	}

	workflowID := workflowRun.WorkflowRun.ID
	workflowRunNumber := workflowRun.WorkflowRun.RunNumber
	workflowName := workflowRun.WorkflowRun.Name

	uniqueRunnerName := fmt.Sprintf("id-%d-%d", workflowID, workflowRunNumber)
	klog.Infof("uniqueRunnerName: %s\n", uniqueRunnerName)

	klog.Infof("Workflow is requested.  ID: %s, Name: %s\n", uniqueRunnerName, workflowName)
	if create {
		err = createRunner(uniqueRunnerName)
		if err != nil {
			klog.Errorf("createRunner failed. Err: %v\n", err)
			return err
		}
	} else {
		err = deleteRunner(uniqueRunnerName)
		if err != nil {
			klog.Errorf("deleteRunner failed. Err: %v\n", err)
			return err
		}
	}

	klog.Infof("doWorkflowJob succeeded!\n")
	return nil
}

func handleWorkflowRun(workflowRun *webhook.WorkflowRunPayload) {
	// Dump event
	klog.V(6).Infof("---------------------- START DUMP EVENT ----------------------\n\n\n")
	klog.V(6).Infof("%+v\n\n\n", workflowRun)
	klog.V(6).Infof("---------------------- END DUMP EVENT ----------------------\n\n\n")
}

func backgroundClean() {
	// while true loop
	for {
		time.Sleep(time.Duration(defaultCleanUpOrphans) * time.Second)
		klog.V(4).Infof("Do routine clean up check\n")

		ghClient, err := NewGitHub()
		if err != nil {
			klog.Errorf("NewGitHub failed. Err: %v\n", err)
			continue
		}
		ec2Client, err := NewAws()
		if err != nil {
			klog.Errorf("NewAws failed. Err: %v\n", err)
			continue
		}

		runners, err := ghClient.GetGitHubRunners()
		if err != nil {
			klog.Errorf("GetGitHubRunners failed. Err: %v\n", err)
			continue
		}

		if runners.TotalCount == 0 {
			klog.V(4).Infof("runners.TotalCount == 0. Skip!\n")
			continue
		}

		for _, runner := range runners.Runners {
			runnerName := *runner.Name
			runnerStatus := *runner.Status
			klog.Infof("Runner: %s\n", runnerName)

			if !strings.Contains(selfhostedRunnerPrefix, runnerName) {
				klog.Infof("Skipping Runner: %s\n", runnerName)
				continue
			}

			klog.V(4).Infof("Runner %s Status: %s\n", runnerName, runnerStatus)
			switch runnerStatus {
			case runnerIdle:
				err := ghClient.DeleteGitHubRunnerByName(runnerName)
				if err != nil {
					klog.Errorf("DeleteGitHubRunnerByName(%s) failed. Err: %v\n", runnerName, err)
				}
				err = ec2Client.DeleteEc2InstanceByName(runnerName)
				if err != nil {
					klog.Errorf("DeleteEc2InstanceByName(%s) failed. Err: %v\n", runnerName, err)
				}

			default:
				klog.V(4).Infof("Runner %s is %s. Is active!\n", runnerName, runnerStatus)
			}
		}
	}
}
