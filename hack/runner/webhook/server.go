package main

import (
	"fmt"
	"strings"
	"time"

	ec2 "github.com/aws/aws-sdk-go/service/ec2"
	github "github.com/google/go-github/v39/github"
)

const (
	defaultSleepBetweenPoll     int64 = 10
	defaultNumOfTimesToPoll     int   = 30
	defaultmustHaveStatusBefore int   = 10
	defaultNumOfTimesToRetry    int   = 3
)

func createOnlineRunner(ghClient *github.Client, ec2Client *ec2.EC2) error {
	token, err := createRunnerToken(ghClient)
	if err != nil {
		fmt.Printf("getGitHubClientWithEnvToken failed. Err: %v\n", err)
		return err
	}

	instanceID, err := createEc2Runner(ec2Client, token)
	if err != nil {
		fmt.Printf("createEc2Runner failed. Err: %v\n", err)
		return err
	}

	succeeded := false
	for i := 0; i < defaultNumOfTimesToPoll; i++ {
		runner, err := getGitHubRunner(ghClient, instanceID)
		if err != nil {
			fmt.Printf("getGitHubRunner failed. Err: %v\n", err)
			break
		}

		if !strings.EqualFold(*runner.Status, runnerOnline) && i > defaultmustHaveStatusBefore {
			fmt.Printf("The node should have already returned some status... retry\n")
			break
		}

		if strings.EqualFold(*runner.Status, runnerOnline) {
			fmt.Printf("Succeeded!\n")
			succeeded = true
			break
		}

		fmt.Printf("Attempt poll %d... sleeping\n", i)
		time.Sleep(time.Duration(defaultSleepBetweenPoll))
	}

	if !succeeded {
		fmt.Printf("createOnlineRunner failed. Delete instance %s\n", instanceID)
		err = deleteEc2Instance(ec2Client, instanceID)
		if err != nil {
			fmt.Printf("deleteEc2Instance failed. Err: %v\n", err)
		}
	}

	return nil
}

func createRunner() error {

	ghClient, err := getGitHubClientWithEnvToken()
	if err != nil {
		fmt.Printf("getGitHubClientWithEnvToken failed. Err: %v\n", err)
		return err
	}
	ec2Client, err := getAwsClientWithEnvToken()
	if err != nil {
		fmt.Printf("getAwsClientWithEnvToken failed. Err: %v\n", err)
		return err
	}

	for i := 0; i < defaultNumOfTimesToRetry; i++ {
		err = createOnlineRunner(ghClient, ec2Client)
		if err == nil {
			fmt.Printf("createOnlineRunner succeeded!")
			break
		}

		fmt.Printf("createOnlineRunner failed... retying. Err: %v\n", err)
	}

	return err
}

func deleteRunner(instanceID string) error {

	ghClient, err := getGitHubClientWithEnvToken()
	if err != nil {
		fmt.Printf("getGitHubClientWithEnvToken failed. Err: %v\n", err)
		return err
	}
	ec2Client, err := getAwsClientWithEnvToken()
	if err != nil {
		fmt.Printf("getAwsClientWithEnvToken failed. Err: %v\n", err)
		return err
	}

	err = deleteGitHubRunnerByName(ghClient, instanceID)
	if err != nil {
		// Just a warning because of the new self host ephemeral
		fmt.Printf("deleteGitHubRunnerByName failed. Err: %v\n", err)
	}

	err = deleteEc2Instance(ec2Client, instanceID)
	if err != nil {
		fmt.Printf("deleteEc2Instance failed. Err: %v\n", err)
		return err
	}

	fmt.Printf("deleteRunner succeeded\n")
	return nil
}
