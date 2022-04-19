// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/credentials"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/pkg/system"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/aws"
	awsclient "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/aws"
	tfclient "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/client"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/log"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgconfigbom"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgconfigproviders"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgconfigupdater"
	mcuimodels "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/web/server/models"
)

const (
	ConfigVariableAWSRegion          = "AWS_REGION"
	ConfigVariableAWSSecretAccessKey = "AWS_SECRET_ACCESS_KEY" //nolint:gosec
	ConfigVariableAWSAccessKeyID     = "AWS_ACCESS_KEY_ID"     //nolint:gosec
	ConfigVariableAWSSessionToken    = "AWS_SESSION_TOKEN"     //nolint:gosec
	ConfigVariableAWSProfile         = "AWS_PROFILE"
	ConfigVariableAWSB64Credentials  = "AWS_B64ENCODED_CREDENTIALS" //nolint:gosec
	ConfigVariableAWSVPCID           = "AWS_VPC_ID"
)

// ApplyTKGConfigForAWS applies TKG configuration for AWS.
func (app *App) ApplyTKGConfigForAWS(params aws.ApplyTKGConfigForAWSParams) middleware.Responder {
	tkgClient, err := app.getTkgClient()
	if err != nil {
		return aws.NewApplyTKGConfigForAWSInternalServerError().WithPayload(Err(err))
	}

	if app.awsClient == nil {
		return aws.NewApplyTKGConfigForAWSInternalServerError().WithPayload(Err(errors.New("aws client is not initialized properly")))
	}

	encodedCreds, err := app.awsClient.EncodeCredentials()
	if err != nil {
		return aws.NewApplyTKGConfigForAWSInternalServerError().WithPayload(Err(err))
	}

	convertedParams, err := mgmtClusterParamsToMCUIRegionalParams(params.Params)
	if err != nil {
		return aws.NewApplyTKGConfigForAWSInternalServerError().WithPayload(Err(err))
	}

	configReaderWriter := tkgClient.TKGConfigReaderWriter()
	awsConfig, err := tkgconfigproviders.New(system.GetConfigDir(), configReaderWriter).NewAWSConfig(convertedParams, encodedCreds)
	if err != nil {
		return aws.NewApplyTKGConfigForAWSInternalServerError().WithPayload(Err(err))
	}

	err = tkgconfigupdater.SaveConfig(app.getFilePathForSavingConfig(), configReaderWriter, awsConfig)
	if err != nil {
		return aws.NewApplyTKGConfigForAWSInternalServerError().WithPayload(Err(err))
	}

	return aws.NewApplyTKGConfigForAWSOK().WithPayload(&models.ConfigFileInfo{Path: app.getFilePathForSavingConfig()})
}

// CreateAWSManagementCluster creates aws management cluster
func (app *App) CreateAWSManagementCluster(params aws.CreateAWSManagementClusterParams) middleware.Responder {
	tkgClient, err := app.getTkgClient()
	if err != nil {
		return aws.NewCreateAWSManagementClusterInternalServerError().WithPayload(Err(err))
	}

	if app.awsClient == nil {
		return aws.NewCreateAWSManagementClusterInternalServerError().WithPayload(Err(errors.New("aws client is not initialized properly")))
	}

	encodedCreds, err := app.awsClient.EncodeCredentials()
	if err != nil {
		return aws.NewCreateAWSManagementClusterInternalServerError().WithPayload(Err(err))
	}

	convertedParams, err := mgmtClusterParamsToMCUIRegionalParams(params.Params)
	if err != nil {
		return aws.NewCreateAWSManagementClusterInternalServerError().WithPayload(Err(err))
	}

	configReaderWriter := tkgClient.TKGConfigReaderWriter()
	awsConfig, err := tkgconfigproviders.New(system.GetConfigDir(), configReaderWriter).NewAWSConfig(convertedParams, encodedCreds)
	if err != nil {
		return aws.NewCreateAWSManagementClusterInternalServerError().WithPayload(Err(err))
	}

	err = tkgconfigupdater.SaveConfig(app.getFilePathForSavingConfig(), configReaderWriter, awsConfig)
	if err != nil {
		return aws.NewCreateAWSManagementClusterInternalServerError().WithPayload(Err(err))
	}

	initOptions := &tfclient.InitRegionOptions{
		InfrastructureProvider: "aws",
		ClusterName:            convertedParams.ClusterName,
		Plan:                   convertedParams.ControlPlaneFlavor,
		CeipOptIn:              *convertedParams.CeipOptIn,
		CniType:                convertedParams.Networking.CniType,
		Annotations:            convertedParams.Annotations,
		Labels:                 convertedParams.Labels,
		ClusterConfigFile:      app.getFilePathForSavingConfig(),
	}
	if err := tkgClient.ConfigureAndValidateManagementClusterConfiguration(initOptions, false); err != nil {
		return aws.NewCreateAWSManagementClusterInternalServerError().WithPayload(Err(err))
	}
	go app.StartSendingLogsToUI()

	go func() {
		if params.Params.CreateCloudFormationStack {
			err = tkgClient.CreateAWSCloudFormationStack()
			if err != nil {
				log.Error(err, "unable to create AWS CloudFormationStack")
				return
			}
		}

		err := tkgClient.InitRegion(initOptions)
		if err != nil {
			log.Error(err, "unable to set up management cluster, ")
		} else {
			log.Infof("\nManagement cluster created!\n\n")
			log.Info("\nYou can now create your first workload cluster by running the following:\n\n")
			log.Info("  tanzu cluster create [name] -f [file]\n\n")
			// wait for the logs to be dispatched to UI before exit
			time.Sleep(sleepTimeForLogsPropogation)
		}
	}()

	return aws.NewCreateAWSManagementClusterOK().WithPayload("started creating regional cluster")
}

// ExportAWSConfig returns the config file content based on passed parameters.
func (app *App) ExportAWSConfig(params aws.ExportTKGConfigForAWSParams) middleware.Responder {
	tkgClient, err := app.getTkgClient()
	if err != nil {
		return aws.NewExportTKGConfigForAWSInternalServerError().WithPayload(Err(err))
	}

	if app.awsClient == nil {
		return aws.NewExportTKGConfigForAWSInternalServerError().WithPayload(Err(errors.New("aws client is not initialized properly")))
	}

	encodedCreds, err := app.awsClient.EncodeCredentials()
	if err != nil {
		return aws.NewExportTKGConfigForAWSInternalServerError().WithPayload(Err(err))
	}

	convertedParams, err := mgmtClusterParamsToMCUIRegionalParams(params.Params)
	if err != nil {
		return aws.NewExportTKGConfigForAWSInternalServerError().WithPayload(Err(err))
	}

	// create the provider object with the configuration data
	configReaderWriter := tkgClient.TKGConfigReaderWriter()
	config, err := tkgconfigproviders.New(system.GetConfigDir(), configReaderWriter).NewAWSConfig(convertedParams, encodedCreds)
	if err != nil {
		return aws.NewExportTKGConfigForAWSInternalServerError().WithPayload(Err(err))
	}

	configString, err := transformConfigToString(config)
	if err != nil {
		return aws.NewExportTKGConfigForAWSInternalServerError().WithPayload(Err(err))
	}

	return aws.NewExportTKGConfigForAWSOK().WithPayload(configString)
}

// GetAWSAvailabilityZones gets availability zones under the user-specified region.
func (app *App) GetAWSAvailabilityZones(params aws.GetAWSAvailabilityZonesParams) middleware.Responder {
	if app.awsClient == nil {
		return aws.NewGetAWSAvailabilityZonesInternalServerError().WithPayload(Err(errors.New("aws client is not initialized properly")))
	}

	azs, err := app.awsClient.ListAvailabilityZones()
	if err != nil {
		return aws.NewGetAWSAvailabilityZonesInternalServerError().WithPayload(Err(err))
	}

	return aws.NewGetAWSAvailabilityZonesOK().WithPayload(mcUIAZtoAZ(azs))
}

// GetAWSCredentialProfiles gets aws credential profiles.
func (app *App) GetAWSCredentialProfiles(params aws.GetAWSCredentialProfilesParams) middleware.Responder {
	res, err := awsclient.ListCredentialProfiles("")
	if err != nil {
		return aws.NewGetAWSCredentialProfilesInternalServerError().WithPayload(Err(err))
	}

	return aws.NewGetAWSCredentialProfilesOK().WithPayload(res)
}

// GetAWSNodeTypes gets aws node types.
func (app *App) GetAWSNodeTypes(params aws.GetAWSNodeTypesParams) middleware.Responder {
	if app.awsClient == nil {
		return aws.NewGetAWSNodeTypesInternalServerError().WithPayload(Err(errors.New("aws client is not initialized properly")))
	}

	az := ""
	if params.Az != nil {
		az = *params.Az
	}

	result, err := app.awsClient.ListInstanceTypes(az)
	if err != nil {
		return aws.NewGetAWSNodeTypesInternalServerError().WithPayload(Err(err))
	}

	return aws.NewGetAWSNodeTypesOK().WithPayload(result)
}

// GetAWSOSImages gets available OS images.
func (app *App) GetAWSOSImages(params aws.GetAWSOSImagesParams) middleware.Responder {
	tkgClient, err := app.getTkgClient()
	if err != nil {
		return aws.NewGetAWSOSImagesInternalServerError().WithPayload(Err(err))
	}

	configReaderWriter := tkgClient.TKGConfigReaderWriter()
	bomConfig, err := tkgconfigbom.New(system.GetConfigDir(), configReaderWriter).GetDefaultTkrBOMConfiguration()
	if err != nil {
		return aws.NewGetAWSOSImagesInternalServerError().WithPayload(Err(err))
	}

	results := []*models.AWSVirtualMachine{}

	amis, exists := bomConfig.AMI[params.Region]
	if !exists {
		return aws.NewGetAWSOSImagesInternalServerError().WithPayload(Err(fmt.Errorf("no AMI found for the provided region %q", params.Region)))
	}

	for _, ami := range amis {
		displayName := fmt.Sprintf("%s-%s-%s (%s)", ami.OSInfo.Name, ami.OSInfo.Version, ami.OSInfo.Arch, ami.ID)
		results = append(results, &models.AWSVirtualMachine{
			Name: displayName,
			OsInfo: &models.OSInfo{
				Name:    ami.OSInfo.Name,
				Version: ami.OSInfo.Version,
				Arch:    ami.OSInfo.Arch,
			},
		})
	}
	return aws.NewGetAWSOSImagesOK().WithPayload(results)
}

// GetAWSRegions returns a list of AWS regions.
func (app *App) GetAWSRegions(params aws.GetAWSRegionsParams) middleware.Responder {
	tkgClient, err := app.getTkgClient()
	if err != nil {
		return aws.NewGetAWSRegionsInternalServerError().WithPayload(Err(err))
	}

	bomConfig, err := tkgconfigbom.New(system.GetConfigDir(), tkgClient.TKGConfigReaderWriter()).GetDefaultTkrBOMConfiguration()
	if err != nil {
		return aws.NewGetAWSRegionsInternalServerError().WithPayload(Err(err))
	}

	regions := []string{}
	for region := range bomConfig.AMI {
		regions = append(regions, region)
	}
	sort.Strings(regions)

	return aws.NewGetAWSRegionsOK().WithPayload(regions)
}

// GetAWSSubnets gets all subnets under given VPC ID.
func (app *App) GetAWSSubnets(params aws.GetAWSSubnetsParams) middleware.Responder {
	if app.awsClient == nil {
		return aws.NewGetAWSSubnetsInternalServerError().WithPayload(Err(errors.New("aws client is not initialized properly")))
	}

	subnets, err := app.awsClient.ListSubnets(params.VpcID)
	if err != nil {
		return aws.NewGetAWSSubnetsInternalServerError().WithPayload(Err(err))
	}

	return aws.NewGetAWSSubnetsOK().WithPayload(mcUISubnetstoSubnets(subnets))
}

// GetVPCs gets all VPCs under current AWS account.
func (app *App) GetVPCs(params aws.GetVPCsParams) middleware.Responder {
	if app.awsClient == nil {
		return aws.NewGetVPCsInternalServerError().WithPayload(Err(errors.New("aws client is not initialized properly")))
	}

	vpcs, err := app.awsClient.ListVPCs()
	if err != nil {
		return aws.NewGetVPCsInternalServerError().WithPayload(Err(err))
	}

	return aws.NewGetVPCsOK().WithPayload(mcUIVPCToVPC(vpcs))
}

// SetAWSEndpoint sets the AWS credentials.
func (app *App) SetAWSEndpoint(params aws.SetAWSEndpointParams) middleware.Responder {
	var err error
	var creds *credentials.AWSCredentials

	if params.AccountParams.AccessKeyID != "" && params.AccountParams.SecretAccessKey != "" {
		creds = &credentials.AWSCredentials{
			Region:          params.AccountParams.Region,
			AccessKeyID:     params.AccountParams.AccessKeyID,
			SecretAccessKey: params.AccountParams.SecretAccessKey,
			SessionToken:    params.AccountParams.SessionToken,
		}
	} else {
		if params.AccountParams.ProfileName != "" {
			os.Setenv(ConfigVariableAWSProfile, params.AccountParams.ProfileName)
		}
		creds, err = credentials.NewAWSCredentialFromDefaultChain(params.AccountParams.Region)
		if err != nil {
			return aws.NewSetAWSEndpointInternalServerError().WithPayload(Err(err))
		}
	}

	client, err := awsclient.New(*creds)
	if err != nil {
		return aws.NewSetAWSEndpointInternalServerError().WithPayload(Err(err))
	}

	err = client.VerifyAccount()
	if err != nil {
		return aws.NewSetAWSEndpointInternalServerError().WithPayload(Err(err))
	}

	app.awsClient = client
	return aws.NewSetAWSEndpointCreated()
}

// Need until the awclient code decouples from the presentation code in the
// management-cluster API logic.
func mcUISubnetstoSubnets(subnets []*mcuimodels.AWSSubnet) []*models.AWSSubnet {
	result := []*models.AWSSubnet{}

	for _, subnet := range subnets {
		result = append(result, &models.AWSSubnet{
			AvailabilityZoneID:   subnet.AvailabilityZoneID,
			AvailabilityZoneName: subnet.AvailabilityZoneName,
			Cidr:                 subnet.Cidr,
			IsPublic:             subnet.IsPublic,
			State:                subnet.State,
			VpcID:                subnet.VpcID,
			ID:                   subnet.ID,
		})
	}

	return result
}

// Need until the awclient code decouples from the presentation code in the
// management-cluster API logic.
func mgmtClusterParamsToMCUIRegionalParams(params *models.AWSManagementClusterParams) (*mcuimodels.AWSRegionalClusterParams, error) {
	// Should be same structure, so we can marshal through JSON.
	// Easier this way since there are nested model structs.
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	result := &mcuimodels.AWSRegionalClusterParams{}
	err = json.Unmarshal(jsonData, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Need until the awclient code decouples from the presentation code in the
// management-cluster API logic.
func mcUIAZtoAZ(azs []*mcuimodels.AWSAvailabilityZone) []*models.AWSAvailabilityZone {
	result := []*models.AWSAvailabilityZone{}

	for _, az := range azs {
		result = append(result, &models.AWSAvailabilityZone{
			ID:   az.ID,
			Name: az.Name,
		})
	}

	return result
}

// Need until the awclient code decouples from the presentation code in the
// management-cluster API logic.
func mcUIVPCToVPC(vpcs []*mcuimodels.Vpc) []*models.Vpc {
	result := []*models.Vpc{}

	for _, vpc := range vpcs {
		result = append(result, &models.Vpc{
			Cidr: vpc.Cidr,
			ID:   vpc.ID,
		})
	}

	return result
}
