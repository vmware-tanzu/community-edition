// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package app contains overall application settings and common things. This
// duplicates some things from tanzu-framework, but we don't want to create too
// many dependencies on code that wasn't intended as shared resources.
package handlers

import (
	"path/filepath"
	"time"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/pkg/system"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/aws"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/cri"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/docker"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/edition"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/features"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/provider"
	awsclient "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/aws"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/client"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/constants"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/log"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgctl"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/utils"
)

const sleepTimeForLogsPropogation = 2 * time.Second

// App application structs consisting init options and clients
type App struct {
	Timeout   time.Duration
	LogLevel  int32
	awsClient awsclient.Client
	// clientTKG         *tkgctl.TKGClient
	// clientTkg         *client.TkgClient
	clusterConfigFile string
}

func (app *App) getTKGClient() (tkgctl.TKGClient, error) {
	return system.NewTKGClient("", app.LogLevel)
}

// Yep, you read that right - there's a TKGClient and a TkgClient. o_O
func (app *App) getTkgClient() (*client.TkgClient, error) {
	return system.NewTkgClient()
}

// ConfigureHandlers configures API handlers func
func (app *App) ConfigureHandlers(api *operations.TanzuUIAPI) {
	// Handlers for system settings, configuration, and general information
	api.CriGetContainerRuntimeInfoHandler = cri.GetContainerRuntimeInfoHandlerFunc(app.GetContainerRuntimeInfo)
	api.EditionGetTanzuEditionHandler = edition.GetTanzuEditionHandlerFunc(app.Edition)
	api.FeaturesGetFeatureFlagsHandler = features.GetFeatureFlagsHandlerFunc(app.FeatureFlags)
	api.ProviderGetProviderHandler = provider.GetProviderHandlerFunc(app.Providers)

	// Handlers related to the docker/CAPD provider
	api.DockerApplyTKGConfigForDockerHandler = docker.ApplyTKGConfigForDockerHandlerFunc(app.ApplyTKGConfigForDocker)
	api.DockerCheckIfDockerDaemonAvailableHandler = docker.CheckIfDockerDaemonAvailableHandlerFunc(app.CheckIfDockerDaemonAvailable)
	api.DockerCreateDockerManagementClusterHandler = docker.CreateDockerManagementClusterHandlerFunc(app.CreateDockerManagementCluster)
	api.DockerExportTKGConfigForDockerHandler = docker.ExportTKGConfigForDockerHandlerFunc(app.ExportDockerConfig)

	// Handlers related to the AWS/CAPA provider
	// api.AwsGetAWSEndpointHandler = aws.GetAWSEndpointHandlerFunc(app.GetAWSEndpoint) // Removed in t-f by #1385
	api.AwsApplyTKGConfigForAWSHandler = aws.ApplyTKGConfigForAWSHandlerFunc(app.ApplyTKGConfigForAWS)
	api.AwsCreateAWSManagementClusterHandler = aws.CreateAWSManagementClusterHandlerFunc(app.CreateAWSManagementCluster)
	api.AwsExportTKGConfigForAWSHandler = aws.ExportTKGConfigForAWSHandlerFunc(app.ExportAWSConfig)
	api.AwsGetAWSCredentialProfilesHandler = aws.GetAWSCredentialProfilesHandlerFunc(app.GetAWSCredentialProfiles)
	api.AwsGetAWSAvailabilityZonesHandler = aws.GetAWSAvailabilityZonesHandlerFunc(app.GetAWSAvailabilityZones)
	api.AwsGetAWSNodeTypesHandler = aws.GetAWSNodeTypesHandlerFunc(app.GetAWSNodeTypes)
	api.AwsGetAWSOSImagesHandler = aws.GetAWSOSImagesHandlerFunc(app.GetAWSOSImages)
	api.AwsGetAWSRegionsHandler = aws.GetAWSRegionsHandlerFunc(app.GetAWSRegions)
	api.AwsGetAWSSubnetsHandler = aws.GetAWSSubnetsHandlerFunc(app.GetAWSSubnets)
	api.AwsGetVPCsHandler = aws.GetVPCsHandlerFunc(app.GetVPCs)
	api.AwsSetAWSEndpointHandler = aws.SetAWSEndpointHandlerFunc(app.SetAWSEndpoint)

	// Handlers related to the vSphere/CAPV provider
	// a.VsphereSetVSphereEndpointHandler = vsphere.SetVSphereEndpointHandlerFunc(app.SetVSphereEndpoint)
	// a.VsphereGetVSphereDatacentersHandler = vsphere.GetVSphereDatacentersHandlerFunc(app.GetVSphereDatacenters)
	// a.VsphereGetVSphereDatastoresHandler = vsphere.GetVSphereDatastoresHandlerFunc(app.GetVSphereDatastores)
	// a.VsphereGetVSphereNetworksHandler = vsphere.GetVSphereNetworksHandlerFunc(app.GetVSphereNetworks)
	// a.VsphereGetVSphereResourcePoolsHandler = vsphere.GetVSphereResourcePoolsHandlerFunc(app.GetVSphereResourcePools)
	// a.VsphereCreateVSphereManagementClusterHandler = vsphere.CreateVSphereManagementClusterHandlerFunc(app.CreateVSphereManagementCluster)
	// a.VsphereGetVSphereOSImagesHandler = vsphere.GetVSphereOSImagesHandlerFunc(app.GetVsphereOSImages)
	// a.VsphereGetVSphereFoldersHandler = vsphere.GetVSphereFoldersHandlerFunc(app.GetVSphereFolders)
	// a.VsphereGetVSphereComputeResourcesHandler = vsphere.GetVSphereComputeResourcesHandlerFunc(app.GetVsphereComputeResources)
	// a.VsphereApplyTKGConfigForVsphereHandler = vsphere.ApplyTKGConfigForVsphereHandlerFunc(app.ApplyTKGConfigForVsphere)
	// a.VsphereGetVsphereThumbprintHandler = vsphere.GetVsphereThumbprintHandlerFunc(app.GetVsphereThumbprint)
	// a.VsphereExportTKGConfigForVsphereHandler = vsphere.ExportTKGConfigForVsphereHandlerFunc(app.ExportVSphereConfig)

	// Handlers related to the Azure/CAPZ provider
	// a.AzureGetAzureEndpointHandler = azure.GetAzureEndpointHandlerFunc(app.GetAzureEndpoint)
	// a.AzureSetAzureEndpointHandler = azure.SetAzureEndpointHandlerFunc(app.SetAzureEndPoint)
	// a.AzureGetAzureResourceGroupsHandler = azure.GetAzureResourceGroupsHandlerFunc(app.GetAzureResourceGroups)
	// a.AzureCreateAzureResourceGroupHandler = azure.CreateAzureResourceGroupHandlerFunc(app.CreateAzureResourceGroup)
	// a.AzureGetAzureVnetsHandler = azure.GetAzureVnetsHandlerFunc(app.GetAzureVirtualNetworks)
	// a.AzureCreateAzureVirtualNetworkHandler = azure.CreateAzureVirtualNetworkHandlerFunc(app.CreateAzureVirtualNetwork)
	// a.AzureGetAzureRegionsHandler = azure.GetAzureRegionsHandlerFunc(app.GetAzureRegions)
	// a.AzureGetAzureInstanceTypesHandler = azure.GetAzureInstanceTypesHandlerFunc(app.GetAzureInstanceTypes)
	// a.AzureApplyTKGConfigForAzureHandler = azure.ApplyTKGConfigForAzureHandlerFunc(app.ApplyTKGConfigForAzure)
	// a.AzureCreateAzureManagementClusterHandler = azure.CreateAzureManagementClusterHandlerFunc(app.CreateAzureManagementCluster)
	// a.AzureGetAzureOSImagesHandler = azure.GetAzureOSImagesHandlerFunc(app.GetAzureOSImages)
	// a.AzureExportTKGConfigForAzureHandler = azure.ExportTKGConfigForAzureHandlerFunc(app.ExportAzureConfig)

	// Handlers related to the AVI
	// a.AviVerifyAccountHandler = avi.VerifyAccountHandlerFunc(app.VerifyAccount)
	// a.AviGetAviCloudsHandler = avi.GetAviCloudsHandlerFunc(app.GetAviClouds)
	// a.AviGetAviServiceEngineGroupsHandler = avi.GetAviServiceEngineGroupsHandlerFunc(app.GetAviServiceEngineGroups)
	// a.AviGetAviVipNetworksHandler = avi.GetAviVipNetworksHandlerFunc(app.GetAviVipNetworks)

	// Handlers related to the LDAP and OIDC
	// a.LdapVerifyLdapConnectHandler = ldap.VerifyLdapConnectHandlerFunc(app.VerifyLdapConnect)
	// a.LdapVerifyLdapBindHandler = ldap.VerifyLdapBindHandlerFunc(app.VerifyLdapBind)
	// a.LdapVerifyLdapUserSearchHandler = ldap.VerifyLdapUserSearchHandlerFunc(app.VerifyUserSearch)
	// a.LdapVerifyLdapGroupSearchHandler = ldap.VerifyLdapGroupSearchHandlerFunc(app.VerifyGroupSearch)
	// a.LdapVerifyLdapCloseConnectionHandler = ldap.VerifyLdapCloseConnectionHandlerFunc(app.VerifyLdapCloseConnection)
}

// Err converts a go error into the API error response.
func Err(err error) *models.Error {
	return &models.Error{Message: err.Error()}
}

func (app *App) getFilePathForSavingConfig() string {
	if app.clusterConfigFile == "" {
		randomFileName := utils.GenerateRandomID(10, true) + ".yaml"
		app.clusterConfigFile = filepath.Join(system.GetConfigDir(), constants.TKGClusterConfigFileDirForUI, randomFileName)
	}
	return app.clusterConfigFile
}

// StartSendingLogsToUI creates logchannel passes it to tkg logger
// retrieves the logs through logChannel and passes it to webSocket
func (app *App) StartSendingLogsToUI() {
	logChannel := make(chan []byte)
	log.SetChannel(logChannel)
	for logMsg := range logChannel {
		SendLog(logMsg)
	}
}
