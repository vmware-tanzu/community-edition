// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"strings"

	klog "k8s.io/klog/v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	ec2 "github.com/aws/aws-sdk-go/service/ec2"
)

// get github client
func getAwsClientWithEnvToken() (*ec2.EC2, error) {
	// creds from env
	creds := credentials.NewEnvCredentials()
	awsRegion := os.Getenv("AWS_REGION")

	// create session object
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: creds,
	}))
	if sess == nil {
		err := ErrClientInvalid
		klog.Errorf("EC2 session is nil. Err: %v\n", err)
		return nil, err
	}

	// new EC2 client
	client := ec2.New(sess)
	if client == nil {
		err := ErrClientInvalid
		klog.Errorf("EC2 client is nil. Err: %v\n", err)
		return nil, err
	}

	return client, nil
}

func createEc2Runner(client *ec2.EC2, uniqueID, runnerToken string) (string, error) {
	klog.V(6).Infof("uniqueID: %s\n", uniqueID)
	klog.V(6).Infof("runnerToken: %s\n", runnerToken)

	// setup
	awsAmi := os.Getenv("AWS_AMI_ID")
	awsSecurityGroup := os.Getenv("AWS_SECURITY_GROUP")
	awsSubnet := os.Getenv("AWS_SUBNET")

	// Specify the details of the instance that you want to create.
	runResult, err := client.RunInstances(&ec2.RunInstancesInput{
		ImageId:          aws.String(awsAmi),
		InstanceType:     aws.String("t2.2xlarge"),
		MinCount:         aws.Int64(1),
		MaxCount:         aws.Int64(1),
		KeyName:          aws.String("default"),
		SecurityGroupIds: []*string{aws.String(awsSecurityGroup)},
		SubnetId:         aws.String(awsSubnet),
	})
	if err != nil {
		klog.Errorf("RunInstances failed. Err: %v\n", err)
		return "", err
	}

	instanceID := *runResult.Instances[0].InstanceId
	klog.Infof("Created instance %s\n", instanceID)

	// Add tags to the created instance
	_, err = client.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(uniqueID),
			},
			{
				Key:   aws.String("Token"),
				Value: aws.String(runnerToken),
			},
		},
	})
	if err != nil {
		klog.Errorf("CreateTags failed. Err: %v\n", err)
		return "", err
	}

	// wait for instance running
	err = client.WaitUntilInstanceRunning(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	})
	if err != nil {
		klog.Errorf("WaitUntilInstanceRunning failed. Err: %v\n", err)
		return "", err
	}

	klog.Infof("EC2 instance created successfully\n")
	return instanceID, nil
}

func deleteEc2InstanceByName(client *ec2.EC2, uniqueID string) error {
	klog.Infof("deleteEc2InstanceByName(%s)\n", uniqueID)

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(strings.Join([]string{"*", uniqueID, "*"}, "")),
				},
			},
		},
	}

	result, err := client.DescribeInstances(params)
	if err != nil {
		klog.Errorf("DescribeInstances failed. Err: %v\n", err.Error())
		return err
	}

	// delete all
	if len(result.Reservations) > 0 {
		klog.Infof("len(result.Reservations) == %d\n", len(result.Reservations))

		for _, reservation := range result.Reservations {
			klog.Infof("result.Reservations().Instances == %d\n", len(reservation.Instances))

			for _, instance := range reservation.Instances {
				klog.Infof("Attempt to Delete InstanceID: %s\n", *instance.InstanceId)

				err := deleteEc2Instance(client, *instance.InstanceId)
				if err != nil {
					klog.Errorf("deleteEc2Instance failed. Err: %v\n", err)
					return err
				}
			}
		}
	}

	klog.Infof("deleteEc2InstanceByName(%s) Succeeded!\n", uniqueID)

	return nil
}

func deleteEc2Instance(client *ec2.EC2, instanceID string) error {
	// Specify the details of the instance that you want to create.
	_, err := client.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	})
	if err != nil {
		klog.Errorf("TerminateInstances failed. Err: %v\n", err)
		return err
	}

	klog.Infof("EC2 instance deleted successfully\n")
	return nil
}
