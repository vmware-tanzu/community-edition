// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package aws

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cloudformation/bootstrap"
	cloudformation "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cloudformation/service"
	awscreds "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/credentials"

	awsclient "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/aws"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/aws/ini"
)

type client struct {
	awsClient   awsclient.Client
	credentials awscreds.AWSCredentials
	session     *session.Session
}

// New creates an AWS client
func New(creds awscreds.AWSCredentials) (Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(creds.Region),
		Credentials: credentials.NewStaticCredentials(
			creds.AccessKeyID,
			creds.SecretAccessKey,
			creds.SessionToken,
		),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create aws session")
	}

	tfClient, err := awsclient.New(creds)
	if err != nil {
		return nil, err
	}

	return &client{awsClient: tfClient, session: sess, credentials: creds}, nil
}

// VerifyAccount will verify AWS account credentials.
func (c *client) VerifyAccount() error {
	if c.awsClient == nil {
		return errors.New("uninitialized aws client")
	}
	return c.awsClient.VerifyAccount()
}

// ListVPCs get a list of all VPCs in the current AWS region.
func (c *client) ListVPCs() ([]*VPC, error) {
	if c.awsClient == nil {
		return nil, errors.New("uninitialized aws client")
	}

	vpcs, err := c.awsClient.ListVPCs()
	if err != nil {
		return nil, err
	}

	result := []*VPC{}
	for _, vpc := range vpcs {
		result = append(result, &VPC{
			CIDR: vpc.Cidr,
			ID:   vpc.ID,
		})
	}

	return result, nil
}

// EncodeCredentials will encode the user's AWS credentials.
func (c *client) EncodeCredentials() (string, error) {
	return c.credentials.RenderBase64EncodedAWSDefaultProfile()
}

// ListAvailabilityZones gets the list of the AZs in the current region.
func (c *client) ListAvailabilityZones() ([]*AvailabilityZone, error) {
	if c.awsClient == nil {
		return nil, errors.New("uninitialized aws client")
	}

	azs, err := c.awsClient.ListAvailabilityZones()
	if err != nil {
		return nil, err
	}

	result := []*AvailabilityZone{}
	for _, az := range azs {
		result = append(result, &AvailabilityZone{
			ID:   az.ID,
			Name: az.Name,
		})
	}

	return result, nil
}

// ListRegionsByUser gets all available regions for the user.
func (c *client) ListRegionsByUser() ([]string, error) {
	if c.awsClient == nil {
		return nil, errors.New("uninitialized aws client")
	}
	return c.awsClient.ListRegionsByUser()
}

// GetSubnetGatewayAssociations gets the subnet gateway associations.
func (c *client) GetSubnetGatewayAssociations(vpcID string) (map[string]bool, error) {
	if c.awsClient == nil {
		return nil, errors.New("uninitialized aws client")
	}
	return c.awsClient.GetSubnetGatewayAssociations(vpcID)
}

// ListSubnets gets the subnets for the given VPC.
func (c *client) ListSubnets(vpcID string) ([]*Subnet, error) {
	if c.awsClient == nil {
		return nil, errors.New("uninitialized aws client")
	}

	subnets, err := c.awsClient.ListSubnets(vpcID)
	if err != nil {
		return nil, err
	}

	result := []*Subnet{}
	for _, subnet := range subnets {
		result = append(result, &Subnet{
			AvailabilityZoneID:   subnet.AvailabilityZoneID,
			AvailabilityZoneName: subnet.AvailabilityZoneName,
			CIDR:                 subnet.Cidr,
			ID:                   subnet.ID,
			IsPublic:             subnet.IsPublic,
			State:                subnet.State,
			VPCID:                subnet.VpcID,
		})
	}

	return result, nil
}

// CreateCloudFormationStack creates a new CloudFormationStack with default settings.
func (c *client) CreateCloudFormationStack() error {
	template, err := c.GenerateBootstrapTemplate(GenerateBootstrapTemplateInput{})
	if err != nil {
		return err
	}
	return c.CreateCloudFormationStackWithTemplate(template)
}

// CreateCloudFormationStackWithTemplate creates a new CloudFormationStack using the provided template.
func (c *client) CreateCloudFormationStackWithTemplate(template *bootstrap.Template) error {
	if c.session == nil {
		return errors.New("uninitialized aws client")
	}

	cfnSvc := cloudformation.NewService(cfn.New(c.session))
	if err := cfnSvc.ReconcileBootstrapStack(template.Spec.StackName, *template.RenderCloudFormation()); err != nil {
		return err
	}

	return cfnSvc.ShowStackResources(template.Spec.StackName)
}

// GenerateBootstrapTemplate generates a wrapped CAPA bootstrapv1 configuration specification that controls
// the generation of CloudFormation stacks
func (c *client) GenerateBootstrapTemplate(i GenerateBootstrapTemplateInput) (*bootstrap.Template, error) {
	if c.awsClient == nil {
		return nil, errors.New("uninitialized aws client")
	}

	template, err := c.awsClient.GenerateBootstrapTemplate(awsclient.GenerateBootstrapTemplateInput{
		BootstrapConfigFile:                   i.BootstrapConfigFile,
		DisableTanzuMissionControlPermissions: i.DisableTanzuMissionControlPermissions,
	})
	if err != nil {
		return nil, err
	}

	return template, nil
}

// ListCloudFormationStacks gets a listing of all existing CloudFormationStacks.
func (c *client) ListCloudFormationStacks() ([]string, error) {
	if c.awsClient == nil {
		return nil, errors.New("uninitialized aws client")
	}

	return c.awsClient.ListCloudFormationStacks()
}

// ListInstanceTypes gets a listing of all instance types.
func (c *client) ListInstanceTypes(optionalAZName string) ([]string, error) {
	if c.awsClient == nil {
		return nil, errors.New("uninitialized aws client")
	}
	return c.awsClient.ListInstanceTypes(optionalAZName)
}

// ListEC2KeyPairs will get a list of defined key pairs for the current region.
func (c *client) ListEC2KeyPairs() ([]*KeyPair, error) {
	if c.session == nil {
		return nil, errors.New("uninitialized aws client")
	}
	svc := ec2.New(c.session)

	results, err := svc.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{})
	if err != nil {
		return nil, errors.Wrap(err, "cannot get the list of vpcs under current account")
	}

	vpcs := []*KeyPair{}
	for _, pair := range results.KeyPairs {
		vpcs = append(vpcs, &KeyPair{
			ID:         *pair.KeyPairId,
			Name:       *pair.KeyName,
			Thumbprint: *pair.KeyFingerprint,
		})
	}

	return vpcs, nil
}

func getCredentialSections(filename string) (ini.Sections, error) {
	if filename == "" {
		filename = credentialsFilename()
	}
	return ini.OpenFile(filename)
}

// ListCredentialProfiles lists the name of all profiles in the credential files
func ListCredentialProfiles(filename string) ([]string, error) {
	config, err := getCredentialSections(filename)
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to load shared credentials file")
	}
	return config.List(), nil
}

func credentialsFilename() string {
	if filename := os.Getenv("AWS_SHARED_CREDENTIALS_FILE"); filename != "" {
		return filename
	}

	return filepath.Join(userHomeDir(), ".aws", "credentials")
}

func userHomeDir() string {
	if runtime.GOOS == "windows" { // Windows
		return os.Getenv("USERPROFILE")
	}

	// *nix
	return os.Getenv("HOME")
}
