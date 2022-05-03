// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package app contains overall application settings and common things. This
// duplicates some things from tanzu-framework, but we don't want to create too
// many dependencies on code that wasn't intended as shared resources.
package handlers

import (
	"path/filepath"
	"time"

	awsclient "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/pkg/aws"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/pkg/system"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/avi"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/aws"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/azure"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/cri"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/docker"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/edition"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/features"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/ldap"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/management"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/provider"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/vsphere"
	aviClient "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/avi"
	azureclient "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/azure"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/client"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/constants"
	ldapClient "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/ldap"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/log"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgctl"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/utils"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/vc"
)

const sleepTimeForLogsPropogation = 2 * time.Second

// App application structs consisting init options and clients
type App struct {
	Timeout     time.Duration
	LogLevel    int32
	aviClient   aviClient.Client
	awsClient   awsclient.Client
	azureClient azureclient.Client
	ldapClient  ldapClient.Client
	vcClient    vc.Client
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
//nolint:funlen
func (app *App) ConfigureHandlers(api *operations.TanzuUIAPI) {
	// Handlers for system settings, configuration, and general information
	api.CriGetContainerRuntimeInfoHandler = cri.GetContainerRuntimeInfoHandlerFunc(app.GetContainerRuntimeInfo)
	api.EditionGetTanzuEditionHandler = edition.GetTanzuEditionHandlerFunc(app.Edition)
	api.FeaturesGetFeatureFlagsHandler = features.GetFeatureFlagsHandlerFunc(app.FeatureFlags)
	api.ProviderGetProviderHandler = provider.GetProviderHandlerFunc(app.Providers)

	// Handlers for general management cluster operations
	api.ManagementDeleteMgmtClusterHandler = management.DeleteMgmtClusterHandlerFunc(app.DeleteMgmtCluster)
	api.ManagementGetMgmtClusterHandler = management.GetMgmtClusterHandlerFunc(app.GetMgmtCluster)
	api.ManagementGetMgmtClustersHandler = management.GetMgmtClustersHandlerFunc(app.GetMgmtClusters)
	api.ManagementGetClusterClassHandler = management.GetClusterClassHandlerFunc(app.GetClusterClass)
	api.ManagementGetClusterClassesHandler = management.GetClusterClassesHandlerFunc(app.GetClusterClasses)

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
	api.AwsGetAWSKeyPairsHandler = aws.GetAWSKeyPairsHandlerFunc(app.GetAWSKeyPairs)
	api.AwsGetAWSNodeTypesHandler = aws.GetAWSNodeTypesHandlerFunc(app.GetAWSNodeTypes)
	api.AwsGetAWSOSImagesHandler = aws.GetAWSOSImagesHandlerFunc(app.GetAWSOSImages)
	api.AwsGetAWSRegionsHandler = aws.GetAWSRegionsHandlerFunc(app.GetAWSRegions)
	api.AwsGetAWSSubnetsHandler = aws.GetAWSSubnetsHandlerFunc(app.GetAWSSubnets)
	api.AwsGetVPCsHandler = aws.GetVPCsHandlerFunc(app.GetVPCs)
	api.AwsSetAWSEndpointHandler = aws.SetAWSEndpointHandlerFunc(app.SetAWSEndpoint)

	// Handlers related to the vSphere/CAPV provider
	api.VsphereApplyTKGConfigForVsphereHandler = vsphere.ApplyTKGConfigForVsphereHandlerFunc(app.ApplyTKGConfigForVsphere)
	api.VsphereCreateVSphereManagementClusterHandler = vsphere.CreateVSphereManagementClusterHandlerFunc(app.CreateVSphereManagementCluster)
	api.VsphereExportTKGConfigForVsphereHandler = vsphere.ExportTKGConfigForVsphereHandlerFunc(app.ExportVSphereConfig)
	api.VsphereGetVSphereComputeResourcesHandler = vsphere.GetVSphereComputeResourcesHandlerFunc(app.GetVsphereComputeResources)
	api.VsphereGetVSphereDatacentersHandler = vsphere.GetVSphereDatacentersHandlerFunc(app.GetVSphereDatacenters)
	api.VsphereGetVSphereDatastoresHandler = vsphere.GetVSphereDatastoresHandlerFunc(app.GetVSphereDatastores)
	api.VsphereGetVSphereFoldersHandler = vsphere.GetVSphereFoldersHandlerFunc(app.GetVSphereFolders)
	api.VsphereGetVSphereNetworksHandler = vsphere.GetVSphereNetworksHandlerFunc(app.GetVSphereNetworks)
	api.VsphereGetVSphereOSImagesHandler = vsphere.GetVSphereOSImagesHandlerFunc(app.GetVsphereOSImages)
	api.VsphereGetVSphereResourcePoolsHandler = vsphere.GetVSphereResourcePoolsHandlerFunc(app.GetVSphereResourcePools)
	api.VsphereGetVsphereThumbprintHandler = vsphere.GetVsphereThumbprintHandlerFunc(app.GetVsphereThumbprint)
	api.VsphereSetVSphereEndpointHandler = vsphere.SetVSphereEndpointHandlerFunc(app.SetVSphereEndpoint)

	// Handlers related to the Azure/CAPZ provider
	api.AzureApplyTKGConfigForAzureHandler = azure.ApplyTKGConfigForAzureHandlerFunc(app.ApplyTKGConfigForAzure)
	api.AzureCreateAzureManagementClusterHandler = azure.CreateAzureManagementClusterHandlerFunc(app.CreateAzureManagementCluster)
	api.AzureCreateAzureResourceGroupHandler = azure.CreateAzureResourceGroupHandlerFunc(app.CreateAzureResourceGroup)
	api.AzureCreateAzureVirtualNetworkHandler = azure.CreateAzureVirtualNetworkHandlerFunc(app.CreateAzureVirtualNetwork)
	api.AzureExportTKGConfigForAzureHandler = azure.ExportTKGConfigForAzureHandlerFunc(app.ExportAzureConfig)
	api.AzureGetAzureEndpointHandler = azure.GetAzureEndpointHandlerFunc(app.GetAzureEndpoint)
	api.AzureGetAzureInstanceTypesHandler = azure.GetAzureInstanceTypesHandlerFunc(app.GetAzureInstanceTypes)
	api.AzureGetAzureOSImagesHandler = azure.GetAzureOSImagesHandlerFunc(app.GetAzureOSImages)
	api.AzureGetAzureRegionsHandler = azure.GetAzureRegionsHandlerFunc(app.GetAzureRegions)
	api.AzureGetAzureResourceGroupsHandler = azure.GetAzureResourceGroupsHandlerFunc(app.GetAzureResourceGroups)
	api.AzureGetAzureVnetsHandler = azure.GetAzureVnetsHandlerFunc(app.GetAzureVirtualNetworks)
	api.AzureSetAzureEndpointHandler = azure.SetAzureEndpointHandlerFunc(app.SetAzureEndPoint)

	// Handlers related to AVI
	api.AviGetAviCloudsHandler = avi.GetAviCloudsHandlerFunc(app.GetAviClouds)
	api.AviGetAviServiceEngineGroupsHandler = avi.GetAviServiceEngineGroupsHandlerFunc(app.GetAviServiceEngineGroups)
	api.AviGetAviVipNetworksHandler = avi.GetAviVipNetworksHandlerFunc(app.GetAviVipNetworks)
	api.AviVerifyAccountHandler = avi.VerifyAccountHandlerFunc(app.VerifyAccount)

	// Handlers related to the LDAP and OIDC
	api.LdapVerifyLdapBindHandler = ldap.VerifyLdapBindHandlerFunc(app.VerifyLdapBind)
	api.LdapVerifyLdapCloseConnectionHandler = ldap.VerifyLdapCloseConnectionHandlerFunc(app.VerifyLdapCloseConnection)
	api.LdapVerifyLdapConnectHandler = ldap.VerifyLdapConnectHandlerFunc(app.VerifyLdapConnect)
	api.LdapVerifyLdapGroupSearchHandler = ldap.VerifyLdapGroupSearchHandlerFunc(app.VerifyGroupSearch)
	api.LdapVerifyLdapUserSearchHandler = ldap.VerifyLdapUserSearchHandlerFunc(app.VerifyUserSearch)
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
