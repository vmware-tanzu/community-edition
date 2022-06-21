// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/pkg/system"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/azure"
	azureclient "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/azure"
	tfclient "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/client"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/constants"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgconfigbom"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgconfigproviders"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgconfigupdater"
	mcuimodels "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/web/server/models"
)

// ApplyTKGConfigForAzure applies TKG configuration for Azure
func (app *App) ApplyTKGConfigForAzure(params azure.ApplyTKGConfigForAzureParams) middleware.Responder {
	if app.azureClient == nil {
		return azure.NewApplyTKGConfigForAzureInternalServerError().WithPayload(Err(errors.New("azure client is not initialized properly")))
	}

	convertedParams, err := azMgmtClusterParamsToMCUIManagementParams(params.Params)
	if err != nil {
		return azure.NewApplyTKGConfigForAzureInternalServerError().WithPayload(Err(err))
	}

	filePathForSavingConfig := app.getFilePathForSavingConfig()
	configReaderWriter := app.clientTkg.TKGConfigReaderWriter()
	config, err := tkgconfigproviders.New(system.GetConfigDir(), configReaderWriter).NewAzureConfig(convertedParams)
	if err != nil {
		return azure.NewApplyTKGConfigForAzureInternalServerError().WithPayload(Err(err))
	}

	err = tkgconfigupdater.SaveConfig(filePathForSavingConfig, configReaderWriter, config)
	if err != nil {
		return azure.NewApplyTKGConfigForAzureInternalServerError().WithPayload(Err(err))
	}

	return azure.NewApplyTKGConfigForAzureOK().WithPayload(&models.ConfigFileInfo{Path: filePathForSavingConfig})
}

// CreateAzureManagementCluster creates an Azure management cluster.
func (app *App) CreateAzureManagementCluster(params azure.CreateAzureManagementClusterParams) middleware.Responder {
	if app.azureClient == nil {
		return azure.NewCreateAzureManagementClusterInternalServerError().WithPayload(Err(errors.New("azure client is not initialized properly")))
	}

	convertedParams, err := azMgmtClusterParamsToMCUIManagementParams(params.Params)
	if err != nil {
		return azure.NewCreateAzureManagementClusterInternalServerError().WithPayload(Err(err))
	}

	filePathForSavingConfig := app.getFilePathForSavingConfig()
	configReaderWriter := app.clientTkg.TKGConfigReaderWriter()
	azureConfig, err := tkgconfigproviders.New(system.GetConfigDir(), configReaderWriter).NewAzureConfig(convertedParams)
	if err != nil {
		return azure.NewCreateAzureManagementClusterInternalServerError().WithPayload(Err(err))
	}

	err = tkgconfigupdater.SaveConfig(filePathForSavingConfig, configReaderWriter, azureConfig)
	if err != nil {
		return azure.NewCreateAzureManagementClusterInternalServerError().WithPayload(Err(err))
	}

	// setting the below configuration to tkgClient to be used during Azure mc creation but not saving them to tkg config
	if params.Params.ResourceGroup == "" {
		return azure.NewCreateAzureManagementClusterInternalServerError().WithPayload(Err(errors.New("azure resource group name cannot be empty")))
	}
	configReaderWriter.Set(constants.ConfigVariableAzureResourceGroup, params.Params.ResourceGroup)

	if params.Params.VnetResourceGroup == "" {
		return azure.NewCreateAzureManagementClusterInternalServerError().WithPayload(Err(errors.New("azure vnet resource group name cannot be empty")))
	}
	configReaderWriter.Set(constants.ConfigVariableAzureVnetResourceGroup, params.Params.VnetResourceGroup)

	if params.Params.VnetName == "" {
		return azure.NewCreateAzureManagementClusterInternalServerError().WithPayload(Err(errors.New("azure vnet name cannot be empty")))
	}
	configReaderWriter.Set(constants.ConfigVariableAzureVnetName, params.Params.VnetName)

	if params.Params.ControlPlaneSubnet == "" {
		return azure.NewCreateAzureManagementClusterInternalServerError().WithPayload(Err(errors.New("azure controlplane subnet name cannot be empty")))
	}
	configReaderWriter.Set(constants.ConfigVariableAzureControlPlaneSubnet, params.Params.ControlPlaneSubnet)

	if params.Params.WorkerNodeSubnet == "" {
		return azure.NewCreateAzureManagementClusterInternalServerError().WithPayload(Err(errors.New("azure node subnet name cannot be empty")))
	}
	configReaderWriter.Set(constants.ConfigVariableAzureWorkerSubnet, params.Params.WorkerNodeSubnet)

	if params.Params.VnetCidr != "" { // create new vnet
		configReaderWriter.Set(constants.ConfigVariableAzureVnetCidr, params.Params.VnetCidr)
		configReaderWriter.Set(constants.ConfigVariableAzureControlPlaneSubnetCidr, params.Params.ControlPlaneSubnetCidr)
		configReaderWriter.Set(constants.ConfigVariableAzureWorkerNodeSubnetCidr, params.Params.WorkerNodeSubnetCidr)
	}

	initOptions := &tfclient.InitRegionOptions{
		InfrastructureProvider: "azure",
		ClusterName:            convertedParams.ClusterName,
		Plan:                   convertedParams.ControlPlaneFlavor,
		CeipOptIn:              *convertedParams.CeipOptIn,
		Annotations:            convertedParams.Annotations,
		Labels:                 convertedParams.Labels,
		ClusterConfigFile:      app.getFilePathForSavingConfig(),
		Edition:                "tce",
	}

	if err := app.clientTkg.ConfigureAndValidateManagementClusterConfiguration(initOptions, false); err != nil {
		return azure.NewCreateAzureManagementClusterInternalServerError().WithPayload(Err(errors.New(err.Message)))
	}
	go app.StartSendingLogsToUI()
	go createManagementCluster(app.clientTkg, initOptions)

	return azure.NewCreateAzureManagementClusterOK().WithPayload("started creating regional cluster")
}

// CreateAzureResourceGroup creates a new Azure resource group
func (app *App) CreateAzureResourceGroup(params azure.CreateAzureResourceGroupParams) middleware.Responder {
	if app.azureClient == nil {
		return azure.NewCreateAzureResourceGroupInternalServerError().WithPayload(Err(errors.New("azure client is not initialized properly")))
	}

	err := app.azureClient.CreateResourceGroup(params.HTTPRequest.Context(), *params.Params.Name, *params.Params.Location)
	if err != nil {
		return azure.NewCreateAzureResourceGroupInternalServerError().WithPayload(Err(err))
	}

	return azure.NewCreateAzureResourceGroupCreated()
}

// CreateAzureVirtualNetwork creates a new Azure Virtual Network
func (app *App) CreateAzureVirtualNetwork(params azure.CreateAzureVirtualNetworkParams) middleware.Responder {
	if app.azureClient == nil {
		return azure.NewCreateAzureVirtualNetworkInternalServerError().WithPayload(Err(errors.New("azure client is not initialized properly")))
	}

	err := app.azureClient.CreateVirtualNetwork(params.HTTPRequest.Context(), params.ResourceGroupName, *params.Params.Name, *params.Params.CidrBlock, *params.Params.Location)
	if err != nil {
		return azure.NewCreateAzureVirtualNetworkInternalServerError().WithPayload(Err(err))
	}

	return azure.NewCreateAzureVirtualNetworkCreated()
}

// ExportAzureConfig returns the config file content as a string from incoming params object.
//nolint:dupl
func (app *App) ExportAzureConfig(params azure.ExportTKGConfigForAzureParams) middleware.Responder {
	convertedParams, err := azMgmtClusterParamsToMCUIManagementParams(params.Params)
	if err != nil {
		return azure.NewExportTKGConfigForAzureInternalServerError().WithPayload(Err(err))
	}

	configReaderWriter := app.clientTkg.TKGConfigReaderWriter()
	config, err := tkgconfigproviders.New(system.GetConfigDir(), configReaderWriter).NewAzureConfig(convertedParams)
	if err != nil {
		return azure.NewExportTKGConfigForAzureInternalServerError().WithPayload(Err(err))
	}

	configString, err := transformConfigToString(config)
	if err != nil {
		return azure.NewExportTKGConfigForAzureInternalServerError().WithPayload(Err(err))
	}

	return azure.NewExportTKGConfigForAzureOK().WithPayload(configString)
}

// GetAzureEndpoint gets Azure account info set in environment variables.
func (app *App) GetAzureEndpoint(params azure.GetAzureEndpointParams) middleware.Responder {
	res := models.AzureAccountParams{
		SubscriptionID: os.Getenv(constants.ConfigVariableAzureSubscriptionID),
		TenantID:       os.Getenv(constants.ConfigVariableAzureTenantID),
		ClientID:       os.Getenv(constants.ConfigVariableAzureClientID),
		ClientSecret:   os.Getenv(constants.ConfigVariableAzureClientSecret),
	}

	return azure.NewGetAzureEndpointOK().WithPayload(&res)
}

// GetAzureInstanceTypes lists the available instance types for a given region.
func (app *App) GetAzureInstanceTypes(params azure.GetAzureInstanceTypesParams) middleware.Responder {
	if app.azureClient == nil {
		return azure.NewGetAzureInstanceTypesInternalServerError().WithPayload(Err(errors.New("azure client is not initialized properly")))
	}

	instanceTypes, err := app.azureClient.GetAzureInstanceTypesForRegion(params.HTTPRequest.Context(), params.Location)
	if err != nil {
		return azure.NewGetAzureInstanceTypesInternalServerError().WithPayload(Err(err))
	}

	return azure.NewGetAzureInstanceTypesOK().WithPayload(azmcUIInstanceTypeToInstanceType(instanceTypes))
}

// GetAzureOSImages gets OS information for Azure.
func (app *App) GetAzureOSImages(params azure.GetAzureOSImagesParams) middleware.Responder {
	configReaderWriter := app.clientTkg.TKGConfigReaderWriter()
	bomConfig, err := tkgconfigbom.New(system.GetConfigDir(), configReaderWriter).GetDefaultTkrBOMConfiguration()
	if err != nil {
		return azure.NewGetAzureOSImagesInternalServerError().WithPayload(Err(err))
	}

	results := []*models.AzureVirtualMachine{}
	for i := range bomConfig.Azure {
		displayName := fmt.Sprintf("%s-%s-%s (%s)", bomConfig.Azure[i].OSInfo.Name, bomConfig.Azure[i].OSInfo.Version, bomConfig.Azure[i].OSInfo.Arch, bomConfig.Azure[i].Version)
		results = append(results, &models.AzureVirtualMachine{
			Name: displayName,
			OsInfo: &models.OSInfo{
				Name:    bomConfig.Azure[i].OSInfo.Name,
				Version: bomConfig.Azure[i].OSInfo.Version,
				Arch:    bomConfig.Azure[i].OSInfo.Arch,
			},
		})
	}

	return azure.NewGetAzureOSImagesOK().WithPayload(results)
}

// GetAzureRegions gets a list of all available regions.
func (app *App) GetAzureRegions(params azure.GetAzureRegionsParams) middleware.Responder {
	if app.azureClient == nil {
		return azure.NewGetAzureRegionsInternalServerError().WithPayload(Err(errors.New("azure client is not initialized properly")))
	}

	regions, err := app.azureClient.GetAzureRegions(params.HTTPRequest.Context())
	if err != nil {
		return azure.NewGetAzureRegionsInternalServerError().WithPayload(Err(err))
	}

	return azure.NewGetAzureRegionsOK().WithPayload(azmcUILocationToLocation(regions))
}

// GetAzureResourceGroups gets the list of all Azure resource groups
func (app *App) GetAzureResourceGroups(params azure.GetAzureResourceGroupsParams) middleware.Responder {
	if app.azureClient == nil {
		return azure.NewGetAzureResourceGroupsInternalServerError().WithPayload(Err(errors.New("azure client is not initialized properly")))
	}

	resourceGroups, err := app.azureClient.ListResourceGroups(params.HTTPRequest.Context(), params.Location)
	if err != nil {
		return azure.NewGetAzureResourceGroupsInternalServerError().WithPayload(Err(err))
	}

	return azure.NewGetAzureResourceGroupsOK().WithPayload(azmcUIResourceGroupToResourceGroup(resourceGroups))
}

// GetAzureVirtualNetworks gets the list of all Azure virtual networks for a resource group.
func (app *App) GetAzureVirtualNetworks(params azure.GetAzureVnetsParams) middleware.Responder {
	if app.azureClient == nil {
		return azure.NewGetAzureVnetsInternalServerError().WithPayload(Err(errors.New("azure client is not initialized properly")))
	}

	vnets, err := app.azureClient.ListVirtualNetworks(params.HTTPRequest.Context(), params.ResourceGroupName, params.Location)
	if err != nil {
		return azure.NewGetAzureVnetsInternalServerError().WithPayload(Err(err))
	}

	return azure.NewGetAzureVnetsOK().WithPayload(azmcUIVirtualNetworkToVirtualNetwork(vnets))
}

// SetAzureEndPoint verifies and sets Azure account.
func (app *App) SetAzureEndPoint(params azure.SetAzureEndpointParams) middleware.Responder {
	creds := azureclient.Credentials{
		SubscriptionID: params.AccountParams.SubscriptionID,
		ClientID:       params.AccountParams.ClientID,
		ClientSecret:   params.AccountParams.ClientSecret,
		TenantID:       params.AccountParams.TenantID,
		AzureCloud:     params.AccountParams.AzureCloud,
	}

	client, err := azureclient.New(&creds)
	if err != nil {
		return azure.NewSetAzureEndpointInternalServerError().WithPayload(Err(err))
	}

	err = client.VerifyAccount(params.HTTPRequest.Context())
	if err != nil {
		return azure.NewSetAzureEndpointInternalServerError().WithPayload(Err(err))
	}

	app.azureClient = client
	return azure.NewSetAzureEndpointCreated()
}

// Need until the azClient code decouples from the presentation code in the
// management-cluster API logic.
func azMgmtClusterParamsToMCUIManagementParams(params *models.AzureManagementClusterParams) (*mcuimodels.AzureRegionalClusterParams, error) {
	// Should be same structure, so we can marshal through JSON.
	// Easier this way since there are nested model structs.
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	result := &mcuimodels.AzureRegionalClusterParams{}
	err = json.Unmarshal(jsonData, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Need until the azClient code decouples from the presentation code in the
// management-cluster API logic.
func azmcUIInstanceTypeToInstanceType(instances []*mcuimodels.AzureInstanceType) []*models.AzureInstanceType {
	result := []*models.AzureInstanceType{}

	for _, instance := range instances {
		result = append(result, &models.AzureInstanceType{
			Family: instance.Family,
			Name:   instance.Name,
			Size:   instance.Size,
			Tier:   instance.Tier,
			Zones:  instance.Zones,
		})
	}

	return result
}

// Need until the azClient code decouples from the presentation code in the
// management-cluster API logic.
func azmcUILocationToLocation(locations []*mcuimodels.AzureLocation) []*models.AzureLocation {
	result := []*models.AzureLocation{}

	for _, location := range locations {
		result = append(result, &models.AzureLocation{
			DisplayName: location.DisplayName,
			Name:        location.Name,
		})
	}

	return result
}

// Need until the azClient code decouples from the presentation code in the
// management-cluster API logic.
func azmcUIResourceGroupToResourceGroup(rgs []*mcuimodels.AzureResourceGroup) []*models.AzureResourceGroup {
	result := []*models.AzureResourceGroup{}

	for _, rg := range rgs {
		result = append(result, &models.AzureResourceGroup{
			ID:       rg.ID,
			Location: rg.Location,
			Name:     rg.Name,
		})
	}

	return result
}

// Need until the azClient code decouples from the presentation code in the
// management-cluster API logic.
func azmcUIVirtualNetworkToVirtualNetwork(vns []*mcuimodels.AzureVirtualNetwork) []*models.AzureVirtualNetwork {
	result := []*models.AzureVirtualNetwork{}

	for _, vn := range vns {
		avn := &models.AzureVirtualNetwork{
			CidrBlock: vn.CidrBlock,
			ID:        vn.ID,
			Location:  vn.Location,
			Name:      vn.Name,
		}

		for _, sn := range vn.Subnets {
			avn.Subnets = append(avn.Subnets, &models.AzureSubnet{
				Cidr: sn.Cidr,
				Name: sn.Name,
			})
		}

		result = append(result, avn)
	}

	return result
}
