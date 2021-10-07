package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	ec2 "github.com/aws/aws-sdk-go/service/ec2"
)

// Errors
var (
	// ErrEnvVarNotFound Required environment variable not found
	ErrEnvVarNotFound = errors.New("Required environment variable not found")
)

// get github client
func getAwsClientWithEnvToken() (*ec2.EC2, error) {
	// creds from env
	creds := credentials.NewEnvCredentials()

	var awsRegion string
	if v := os.Getenv("AWS_REGION"); v != "" {
		awsRegion = v
	} else {
		awsRegion = "us-west-2"
	}

	// create session object
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: creds,
	}))

	// new EC2 client
	client := ec2.New(sess)

	return client, nil
}

func createEc2Runner(client *ec2.EC2, runnerToken string) (string, error) {
	// setup
	var awsAmi string
	if v := os.Getenv("AWS_AMI_ID"); v != "" {
		awsAmi = v
	} else {
		fmt.Printf("AWS_AMI_ID is empty\n")
		return "", ErrEnvVarNotFound
	}
	var awsSecurityGroup string
	if v := os.Getenv("AWS_SECURITY_GROUP"); v != "" {
		awsSecurityGroup = v
	} else {
		awsSecurityGroup = "sg-0239f52a8acf71c20"
	}
	var awsSubnet string
	if v := os.Getenv("AWS_SUBNET"); v != "" {
		awsSubnet = v
	} else {
		awsSubnet = "subnet-f0837888"
	}

	// Specify the details of the instance that you want to create.
	runResult, err := client.RunInstances(&ec2.RunInstancesInput{
		ImageId:        aws.String(awsAmi),
		InstanceType:   aws.String("t2.2xlarge"),
		MinCount:       aws.Int64(1),
		MaxCount:       aws.Int64(1),
		KeyName:        aws.String("default"),
		SecurityGroups: []*string{aws.String(awsSecurityGroup)},
		SubnetId:       aws.String(awsSubnet),
	})
	if err != nil {
		fmt.Printf("RunInstances failed. Err: %v\n", err)
		return "", err
	}

	instanceID := *runResult.Instances[0].InstanceId
	fmt.Printf("Created instance %s", instanceID)

	// Add tags to the created instance
	_, err = client.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("token"),
				Value: aws.String(runnerToken),
			},
		},
	})
	if err != nil {
		fmt.Printf("CreateTags failed. Err: %v\n", err)
		return "", err
	}

	// wait for instance running
	err = client.WaitUntilInstanceRunning(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	})
	if err != nil {
		fmt.Printf("WaitUntilInstanceRunning failed. Err: %v\n", err)
		return "", err
	}

	fmt.Printf("EC2 instance created successfully\n")
	return instanceID, nil
}

func deleteEc2Instance(client *ec2.EC2, instanceID string) error {
	// Specify the details of the instance that you want to create.
	_, err := client.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	})
	if err != nil {
		fmt.Printf("TerminateInstances failed. Err: %v\n", err)
		return err
	}

	fmt.Printf("EC2 instance deleted successfully\n")
	return nil
}
