// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/avi"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/aws"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/azure"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/edition"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/ldap"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/provider"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/vsphere"
)

//go:generate swagger generate server --target ../../server --name TanzuUI --spec ../../api/spec.yaml --principal interface{} --exclude-main

func configureFlags(api *operations.TanzuUIAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

//nolint:funlen,gocyclo
func configureAPI(api *operations.TanzuUIAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	if api.AzureApplyTKGConfigForAzureHandler == nil {
		api.AzureApplyTKGConfigForAzureHandler = azure.ApplyTKGConfigForAzureHandlerFunc(func(params azure.ApplyTKGConfigForAzureParams) middleware.Responder {
			return middleware.NotImplemented("operation azure.ApplyTKGConfigForAzure has not yet been implemented")
		})
	}
	if api.VsphereApplyTKGConfigForVsphereHandler == nil {
		api.VsphereApplyTKGConfigForVsphereHandler = vsphere.ApplyTKGConfigForVsphereHandlerFunc(func(params vsphere.ApplyTKGConfigForVsphereParams) middleware.Responder {
			return middleware.NotImplemented("operation vsphere.ApplyTKGConfigForVsphere has not yet been implemented")
		})
	}
	if api.AzureCreateAzureManagementClusterHandler == nil {
		api.AzureCreateAzureManagementClusterHandler = azure.CreateAzureManagementClusterHandlerFunc(func(params azure.CreateAzureManagementClusterParams) middleware.Responder {
			return middleware.NotImplemented("operation azure.CreateAzureManagementCluster has not yet been implemented")
		})
	}
	if api.AzureCreateAzureResourceGroupHandler == nil {
		api.AzureCreateAzureResourceGroupHandler = azure.CreateAzureResourceGroupHandlerFunc(func(params azure.CreateAzureResourceGroupParams) middleware.Responder {
			return middleware.NotImplemented("operation azure.CreateAzureResourceGroup has not yet been implemented")
		})
	}
	if api.AzureCreateAzureVirtualNetworkHandler == nil {
		api.AzureCreateAzureVirtualNetworkHandler = azure.CreateAzureVirtualNetworkHandlerFunc(func(params azure.CreateAzureVirtualNetworkParams) middleware.Responder {
			return middleware.NotImplemented("operation azure.CreateAzureVirtualNetwork has not yet been implemented")
		})
	}
	if api.VsphereCreateVSphereManagementClusterHandler == nil {
		api.VsphereCreateVSphereManagementClusterHandler = vsphere.CreateVSphereManagementClusterHandlerFunc(func(params vsphere.CreateVSphereManagementClusterParams) middleware.Responder {
			return middleware.NotImplemented("operation vsphere.CreateVSphereManagementCluster has not yet been implemented")
		})
	}
	if api.AzureExportTKGConfigForAzureHandler == nil {
		api.AzureExportTKGConfigForAzureHandler = azure.ExportTKGConfigForAzureHandlerFunc(func(params azure.ExportTKGConfigForAzureParams) middleware.Responder {
			return middleware.NotImplemented("operation azure.ExportTKGConfigForAzure has not yet been implemented")
		})
	}
	if api.VsphereExportTKGConfigForVsphereHandler == nil {
		api.VsphereExportTKGConfigForVsphereHandler = vsphere.ExportTKGConfigForVsphereHandlerFunc(func(params vsphere.ExportTKGConfigForVsphereParams) middleware.Responder {
			return middleware.NotImplemented("operation vsphere.ExportTKGConfigForVsphere has not yet been implemented")
		})
	}
	if api.AviGetAviCloudsHandler == nil {
		api.AviGetAviCloudsHandler = avi.GetAviCloudsHandlerFunc(func(params avi.GetAviCloudsParams) middleware.Responder {
			return middleware.NotImplemented("operation avi.GetAviClouds has not yet been implemented")
		})
	}
	if api.AviGetAviServiceEngineGroupsHandler == nil {
		api.AviGetAviServiceEngineGroupsHandler = avi.GetAviServiceEngineGroupsHandlerFunc(func(params avi.GetAviServiceEngineGroupsParams) middleware.Responder {
			return middleware.NotImplemented("operation avi.GetAviServiceEngineGroups has not yet been implemented")
		})
	}
	if api.AviGetAviVipNetworksHandler == nil {
		api.AviGetAviVipNetworksHandler = avi.GetAviVipNetworksHandlerFunc(func(params avi.GetAviVipNetworksParams) middleware.Responder {
			return middleware.NotImplemented("operation avi.GetAviVipNetworks has not yet been implemented")
		})
	}
	if api.AzureGetAzureEndpointHandler == nil {
		api.AzureGetAzureEndpointHandler = azure.GetAzureEndpointHandlerFunc(func(params azure.GetAzureEndpointParams) middleware.Responder {
			return middleware.NotImplemented("operation azure.GetAzureEndpoint has not yet been implemented")
		})
	}
	if api.AzureGetAzureInstanceTypesHandler == nil {
		api.AzureGetAzureInstanceTypesHandler = azure.GetAzureInstanceTypesHandlerFunc(func(params azure.GetAzureInstanceTypesParams) middleware.Responder {
			return middleware.NotImplemented("operation azure.GetAzureInstanceTypes has not yet been implemented")
		})
	}
	if api.AzureGetAzureOSImagesHandler == nil {
		api.AzureGetAzureOSImagesHandler = azure.GetAzureOSImagesHandlerFunc(func(params azure.GetAzureOSImagesParams) middleware.Responder {
			return middleware.NotImplemented("operation azure.GetAzureOSImages has not yet been implemented")
		})
	}
	if api.AzureGetAzureRegionsHandler == nil {
		api.AzureGetAzureRegionsHandler = azure.GetAzureRegionsHandlerFunc(func(params azure.GetAzureRegionsParams) middleware.Responder {
			return middleware.NotImplemented("operation azure.GetAzureRegions has not yet been implemented")
		})
	}
	if api.AzureGetAzureResourceGroupsHandler == nil {
		api.AzureGetAzureResourceGroupsHandler = azure.GetAzureResourceGroupsHandlerFunc(func(params azure.GetAzureResourceGroupsParams) middleware.Responder {
			return middleware.NotImplemented("operation azure.GetAzureResourceGroups has not yet been implemented")
		})
	}
	if api.AzureGetAzureVnetsHandler == nil {
		api.AzureGetAzureVnetsHandler = azure.GetAzureVnetsHandlerFunc(func(params azure.GetAzureVnetsParams) middleware.Responder {
			return middleware.NotImplemented("operation azure.GetAzureVnets has not yet been implemented")
		})
	}
	if api.ProviderGetProviderHandler == nil {
		api.ProviderGetProviderHandler = provider.GetProviderHandlerFunc(func(params provider.GetProviderParams) middleware.Responder {
			return middleware.NotImplemented("operation provider.GetProvider has not yet been implemented")
		})
	}
	if api.EditionGetTanzuEditionHandler == nil {
		api.EditionGetTanzuEditionHandler = edition.GetTanzuEditionHandlerFunc(func(params edition.GetTanzuEditionParams) middleware.Responder {
			return middleware.NotImplemented("operation edition.GetTanzuEdition has not yet been implemented")
		})
	}
	if api.AwsGetVPCsHandler == nil {
		api.AwsGetVPCsHandler = aws.GetVPCsHandlerFunc(func(params aws.GetVPCsParams) middleware.Responder {
			return middleware.NotImplemented("operation aws.GetVPCs has not yet been implemented")
		})
	}
	if api.VsphereGetVSphereComputeResourcesHandler == nil {
		api.VsphereGetVSphereComputeResourcesHandler = vsphere.GetVSphereComputeResourcesHandlerFunc(func(params vsphere.GetVSphereComputeResourcesParams) middleware.Responder {
			return middleware.NotImplemented("operation vsphere.GetVSphereComputeResources has not yet been implemented")
		})
	}
	if api.VsphereGetVSphereDatacentersHandler == nil {
		api.VsphereGetVSphereDatacentersHandler = vsphere.GetVSphereDatacentersHandlerFunc(func(params vsphere.GetVSphereDatacentersParams) middleware.Responder {
			return middleware.NotImplemented("operation vsphere.GetVSphereDatacenters has not yet been implemented")
		})
	}
	if api.VsphereGetVSphereDatastoresHandler == nil {
		api.VsphereGetVSphereDatastoresHandler = vsphere.GetVSphereDatastoresHandlerFunc(func(params vsphere.GetVSphereDatastoresParams) middleware.Responder {
			return middleware.NotImplemented("operation vsphere.GetVSphereDatastores has not yet been implemented")
		})
	}
	if api.VsphereGetVSphereFoldersHandler == nil {
		api.VsphereGetVSphereFoldersHandler = vsphere.GetVSphereFoldersHandlerFunc(func(params vsphere.GetVSphereFoldersParams) middleware.Responder {
			return middleware.NotImplemented("operation vsphere.GetVSphereFolders has not yet been implemented")
		})
	}
	if api.VsphereGetVSphereNetworksHandler == nil {
		api.VsphereGetVSphereNetworksHandler = vsphere.GetVSphereNetworksHandlerFunc(func(params vsphere.GetVSphereNetworksParams) middleware.Responder {
			return middleware.NotImplemented("operation vsphere.GetVSphereNetworks has not yet been implemented")
		})
	}
	if api.VsphereGetVSphereNodeTypesHandler == nil {
		api.VsphereGetVSphereNodeTypesHandler = vsphere.GetVSphereNodeTypesHandlerFunc(func(params vsphere.GetVSphereNodeTypesParams) middleware.Responder {
			return middleware.NotImplemented("operation vsphere.GetVSphereNodeTypes has not yet been implemented")
		})
	}
	if api.VsphereGetVSphereOSImagesHandler == nil {
		api.VsphereGetVSphereOSImagesHandler = vsphere.GetVSphereOSImagesHandlerFunc(func(params vsphere.GetVSphereOSImagesParams) middleware.Responder {
			return middleware.NotImplemented("operation vsphere.GetVSphereOSImages has not yet been implemented")
		})
	}
	if api.VsphereGetVSphereResourcePoolsHandler == nil {
		api.VsphereGetVSphereResourcePoolsHandler = vsphere.GetVSphereResourcePoolsHandlerFunc(func(params vsphere.GetVSphereResourcePoolsParams) middleware.Responder {
			return middleware.NotImplemented("operation vsphere.GetVSphereResourcePools has not yet been implemented")
		})
	}
	if api.VsphereGetVsphereThumbprintHandler == nil {
		api.VsphereGetVsphereThumbprintHandler = vsphere.GetVsphereThumbprintHandlerFunc(func(params vsphere.GetVsphereThumbprintParams) middleware.Responder {
			return middleware.NotImplemented("operation vsphere.GetVsphereThumbprint has not yet been implemented")
		})
	}
	if api.AzureImportTKGConfigForAzureHandler == nil {
		api.AzureImportTKGConfigForAzureHandler = azure.ImportTKGConfigForAzureHandlerFunc(func(params azure.ImportTKGConfigForAzureParams) middleware.Responder {
			return middleware.NotImplemented("operation azure.ImportTKGConfigForAzure has not yet been implemented")
		})
	}
	if api.VsphereImportTKGConfigForVsphereHandler == nil {
		api.VsphereImportTKGConfigForVsphereHandler = vsphere.ImportTKGConfigForVsphereHandlerFunc(func(params vsphere.ImportTKGConfigForVsphereParams) middleware.Responder {
			return middleware.NotImplemented("operation vsphere.ImportTKGConfigForVsphere has not yet been implemented")
		})
	}
	if api.AzureSetAzureEndpointHandler == nil {
		api.AzureSetAzureEndpointHandler = azure.SetAzureEndpointHandlerFunc(func(params azure.SetAzureEndpointParams) middleware.Responder {
			return middleware.NotImplemented("operation azure.SetAzureEndpoint has not yet been implemented")
		})
	}
	if api.VsphereSetVSphereEndpointHandler == nil {
		api.VsphereSetVSphereEndpointHandler = vsphere.SetVSphereEndpointHandlerFunc(func(params vsphere.SetVSphereEndpointParams) middleware.Responder {
			return middleware.NotImplemented("operation vsphere.SetVSphereEndpoint has not yet been implemented")
		})
	}
	if api.AviVerifyAccountHandler == nil {
		api.AviVerifyAccountHandler = avi.VerifyAccountHandlerFunc(func(params avi.VerifyAccountParams) middleware.Responder {
			return middleware.NotImplemented("operation avi.VerifyAccount has not yet been implemented")
		})
	}
	if api.LdapVerifyLdapBindHandler == nil {
		api.LdapVerifyLdapBindHandler = ldap.VerifyLdapBindHandlerFunc(func(params ldap.VerifyLdapBindParams) middleware.Responder {
			return middleware.NotImplemented("operation ldap.VerifyLdapBind has not yet been implemented")
		})
	}
	if api.LdapVerifyLdapCloseConnectionHandler == nil {
		api.LdapVerifyLdapCloseConnectionHandler = ldap.VerifyLdapCloseConnectionHandlerFunc(func(params ldap.VerifyLdapCloseConnectionParams) middleware.Responder {
			return middleware.NotImplemented("operation ldap.VerifyLdapCloseConnection has not yet been implemented")
		})
	}
	if api.LdapVerifyLdapConnectHandler == nil {
		api.LdapVerifyLdapConnectHandler = ldap.VerifyLdapConnectHandlerFunc(func(params ldap.VerifyLdapConnectParams) middleware.Responder {
			return middleware.NotImplemented("operation ldap.VerifyLdapConnect has not yet been implemented")
		})
	}
	if api.LdapVerifyLdapGroupSearchHandler == nil {
		api.LdapVerifyLdapGroupSearchHandler = ldap.VerifyLdapGroupSearchHandlerFunc(func(params ldap.VerifyLdapGroupSearchParams) middleware.Responder {
			return middleware.NotImplemented("operation ldap.VerifyLdapGroupSearch has not yet been implemented")
		})
	}
	if api.LdapVerifyLdapUserSearchHandler == nil {
		api.LdapVerifyLdapUserSearchHandler = ldap.VerifyLdapUserSearchHandlerFunc(func(params ldap.VerifyLdapUserSearchParams) middleware.Responder {
			return middleware.NotImplemented("operation ldap.VerifyLdapUserSearch has not yet been implemented")
		})
	}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
