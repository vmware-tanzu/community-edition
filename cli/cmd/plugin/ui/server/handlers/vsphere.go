// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/pkg/system"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/vsphere"
	tfclient "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/client"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/constants"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/log"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgconfigbom"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgconfigproviders"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgconfigupdater"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/vc"
	mcuimodels "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/web/server/models"
)

// ApplyTKGConfigForVsphere applies TKG configuration for vSphere.
func (app *App) ApplyTKGConfigForVsphere(params vsphere.ApplyTKGConfigForVsphereParams) middleware.Responder {
	convertedParams, err := vspMgmtClusterParamsToMCUIManagementParams(params.Params)
	if err != nil {
		return vsphere.NewApplyTKGConfigForVsphereInternalServerError().WithPayload(Err(err))
	}

	filePathForSavingConfig := app.getFilePathForSavingConfig()
	configReaderWriter := app.clientTkg.TKGConfigReaderWriter()
	config, err := tkgconfigproviders.New(system.GetConfigDir(), configReaderWriter).NewVSphereConfig(convertedParams)
	if err != nil {
		return vsphere.NewApplyTKGConfigForVsphereInternalServerError().WithPayload(Err(err))
	}

	err = tkgconfigupdater.SaveConfig(filePathForSavingConfig, configReaderWriter, config)
	if err != nil {
		return vsphere.NewApplyTKGConfigForVsphereInternalServerError().WithPayload(Err(err))
	}

	return vsphere.NewApplyTKGConfigForVsphereOK().WithPayload(&models.ConfigFileInfo{Path: filePathForSavingConfig})
}

// CreateVSphereManagementCluster creates a vSphere management cluster.
func (app *App) CreateVSphereManagementCluster(params vsphere.CreateVSphereManagementClusterParams) middleware.Responder {
	convertedParams, err := vspMgmtClusterParamsToMCUIManagementParams(params.Params)
	if err != nil {
		return vsphere.NewCreateVSphereManagementClusterInternalServerError().WithPayload(Err(err))
	}

	configReaderWriter := app.clientTkg.TKGConfigReaderWriter()
	config, err := tkgconfigproviders.New(system.GetConfigDir(), configReaderWriter).NewVSphereConfig(convertedParams)
	if err != nil {
		return vsphere.NewCreateVSphereManagementClusterInternalServerError().WithPayload(Err(err))
	}

	err = tkgconfigupdater.SaveConfig(app.getFilePathForSavingConfig(), configReaderWriter, config)
	if err != nil {
		return vsphere.NewCreateVSphereManagementClusterInternalServerError().WithPayload(Err(err))
	}

	initOptions := &tfclient.InitRegionOptions{
		InfrastructureProvider:      "vsphere",
		ClusterName:                 convertedParams.ClusterName,
		Plan:                        convertedParams.ControlPlaneFlavor,
		CeipOptIn:                   *convertedParams.CeipOptIn,
		CniType:                     convertedParams.Networking.CniType,
		Annotations:                 convertedParams.Annotations,
		Labels:                      convertedParams.Labels,
		ClusterConfigFile:           app.getFilePathForSavingConfig(),
		VsphereControlPlaneEndpoint: convertedParams.ControlPlaneEndpoint,
		Edition:                     "tce",
	}

	if err := app.clientTkg.ConfigureAndValidateManagementClusterConfiguration(initOptions, false); err != nil {
		return vsphere.NewCreateVSphereManagementClusterInternalServerError().WithPayload(Err(err))
	}
	go app.StartSendingLogsToUI()
	go createManagementCluster(app.clientTkg, initOptions)

	return vsphere.NewCreateVSphereManagementClusterOK().WithPayload("started creating regional cluster")
}

// ExportVSphereConfig creates return payload of config file string from incoming params object.
//nolint:dupl
func (app *App) ExportVSphereConfig(params vsphere.ExportTKGConfigForVsphereParams) middleware.Responder {
	convertedParams, err := vspMgmtClusterParamsToMCUIManagementParams(params.Params)
	if err != nil {
		return vsphere.NewExportTKGConfigForVsphereInternalServerError().WithPayload(Err(err))
	}

	configReaderWriter := app.clientTkg.TKGConfigReaderWriter()
	config, err := tkgconfigproviders.New(system.GetConfigDir(), configReaderWriter).NewVSphereConfig(convertedParams)
	if err != nil {
		return vsphere.NewExportTKGConfigForVsphereInternalServerError().WithPayload(Err(err))
	}

	configString, err := transformConfigToString(config)
	if err != nil {
		return vsphere.NewExportTKGConfigForVsphereInternalServerError().WithPayload(Err(err))
	}

	return vsphere.NewExportTKGConfigForVsphereOK().WithPayload(configString)
}

// GetVsphereComputeResources gets vSphere compute resources.
func (app *App) GetVsphereComputeResources(params vsphere.GetVSphereComputeResourcesParams) middleware.Responder {
	if app.vcClient == nil {
		return vsphere.NewGetVSphereComputeResourcesInternalServerError().WithPayload(Err(errors.New("vSphere client is not initialized properly")))
	}

	results, err := app.vcClient.GetComputeResources(params.HTTPRequest.Context(), params.Dc)
	if err != nil {
		return vsphere.NewGetVSphereComputeResourcesInternalServerError().WithPayload(Err(err))
	}

	return vsphere.NewGetVSphereComputeResourcesOK().WithPayload(vspmcUIObjectToObject(results))
}

// GetVSphereDatacenters returns all the datacenters in vSphere.
func (app *App) GetVSphereDatacenters(params vsphere.GetVSphereDatacentersParams) middleware.Responder {
	if app.vcClient == nil {
		return vsphere.NewGetVSphereDatacentersInternalServerError().WithPayload(Err(errors.New("vSphere client is not initialized properly")))
	}

	datacenters, err := app.vcClient.GetDatacenters(params.HTTPRequest.Context())
	if err != nil {
		return vsphere.NewGetVSphereDatacentersInternalServerError().WithPayload(Err(err))
	}

	return vsphere.NewGetVSphereDatacentersOK().WithPayload(vspmcUIDatacenterToDatacenter(datacenters))
}

// GetVSphereDatastores returns all the datastores in the datacenter.
func (app *App) GetVSphereDatastores(params vsphere.GetVSphereDatastoresParams) middleware.Responder {
	if app.vcClient == nil {
		return vsphere.NewGetVSphereDatastoresInternalServerError().WithPayload(Err(errors.New("vSphere client is not initialized properly")))
	}

	datastores, err := app.vcClient.GetDatastores(params.HTTPRequest.Context(), params.Dc)
	if err != nil {
		return vsphere.NewGetVSphereDatastoresInternalServerError().WithPayload(Err(err))
	}

	return vsphere.NewGetVSphereDatastoresOK().WithPayload(vspmcUIDatastoreToDatastore(datastores))
}

// GetVSphereFolders gets vsphere folders
func (app *App) GetVSphereFolders(params vsphere.GetVSphereFoldersParams) middleware.Responder {
	if app.vcClient == nil {
		return vsphere.NewGetVSphereFoldersInternalServerError().WithPayload(Err(errors.New("vSphere client is not initialized properly")))
	}

	folders, err := app.vcClient.GetFolders(params.HTTPRequest.Context(), params.Dc)
	if err != nil {
		return vsphere.NewGetVSphereFoldersInternalServerError().WithPayload(Err(err))
	}

	return vsphere.NewGetVSphereFoldersOK().WithPayload(vspmcUIFolderToFolder(folders))
}

// GetVSphereNetworks gets all the  networks in the datacenter
func (app *App) GetVSphereNetworks(params vsphere.GetVSphereNetworksParams) middleware.Responder {
	if app.vcClient == nil {
		return vsphere.NewGetVSphereNetworksInternalServerError().WithPayload(Err(errors.New("vSphere client is not initialized properly")))
	}

	networks, err := app.vcClient.GetNetworks(params.HTTPRequest.Context(), params.Dc)
	if err != nil {
		return vsphere.NewGetVSphereNetworksInternalServerError().WithPayload(Err(err))
	}

	return vsphere.NewGetVSphereNetworksOK().WithPayload(vspmcUINetworksToNetworks(networks))
}

// GetVsphereOSImages gets vm templates for deploying kubernetes node
func (app *App) GetVsphereOSImages(params vsphere.GetVSphereOSImagesParams) middleware.Responder {
	if app.vcClient == nil {
		return vsphere.NewGetVSphereOSImagesInternalServerError().WithPayload(Err(errors.New("vSphere client is not initialized properly")))
	}

	vms, err := app.vcClient.GetVirtualMachineImages(params.HTTPRequest.Context(), params.Dc)
	if err != nil {
		return vsphere.NewGetVSphereOSImagesInternalServerError().WithPayload(Err(err))
	}

	configReaderWriter := app.clientTkg.TKGConfigReaderWriter()
	defaultTKRBom, err := tkgconfigbom.New(system.GetConfigDir(), configReaderWriter).GetDefaultTkrBOMConfiguration()
	if err != nil {
		return vsphere.NewGetVSphereOSImagesInternalServerError().WithPayload(Err(errors.Wrap(err, "unable to get the default TanzuKubernetesRelease")))
	}
	matchedTemplates, nonTemplateVms := vc.FindMatchingVirtualMachineTemplate(vms, defaultTKRBom.GetOVAVersions())

	if len(matchedTemplates) == 0 && len(nonTemplateVms) != 0 {
		log.Infof("unable to find any VM Template associated with the TanzuKubernetesRelease %s, but found these VM(s) [%s] that can be used once converted to a VM Template", defaultTKRBom.Release.Version, strings.Join(nonTemplateVms, ","))
	}

	results := []*models.VSphereVirtualMachine{}

	for _, template := range matchedTemplates {
		results = append(results, &models.VSphereVirtualMachine{
			IsTemplate: &template.IsTemplate,
			Name:       template.Name,
			Moid:       template.Moid,
			OsInfo: &models.OSInfo{
				Name:    template.DistroName,
				Version: template.DistroVersion,
				Arch:    template.DistroArch,
			},
		})
	}

	return vsphere.NewGetVSphereOSImagesOK().WithPayload(results)
}

// GetVSphereResourcePools gets all the resource pools in the datacenter
func (app *App) GetVSphereResourcePools(params vsphere.GetVSphereResourcePoolsParams) middleware.Responder {
	if app.vcClient == nil {
		return vsphere.NewGetVSphereResourcePoolsInternalServerError().WithPayload(Err(errors.New("vSphere client is not initialized properly")))
	}

	rps, err := app.vcClient.GetResourcePools(params.HTTPRequest.Context(), params.Dc)
	if err != nil {
		return vsphere.NewGetVSphereResourcePoolsInternalServerError().WithPayload(Err(err))
	}

	return vsphere.NewGetVSphereResourcePoolsOK().WithPayload(vspmcUIResourcePoolToResourcePool(rps))
}

// GetVsphereThumbprint gets the vSphere thumbprint if insecure flag not set
func (app *App) GetVsphereThumbprint(params vsphere.GetVsphereThumbprintParams) middleware.Responder {
	insecure := false

	thumbprint, err := vc.GetVCThumbprint(params.Host)
	if err != nil {
		return vsphere.NewGetVsphereThumbprintInternalServerError().WithPayload(Err(err))
	}

	res := models.VSphereThumbprint{Thumbprint: thumbprint, Insecure: &insecure}

	return vsphere.NewGetVsphereThumbprintOK().WithPayload(&res)
}

// SetVSphereEndpoint validates vsphere credentials and sets the vsphere client into web app
func (app *App) SetVSphereEndpoint(params vsphere.SetVSphereEndpointParams) middleware.Responder {
	host := strings.TrimSpace(params.Credentials.Host)

	if !strings.HasPrefix(host, "http") {
		host = "https://" + host
	}

	vcURL, err := url.Parse(host)
	if err != nil {
		return vsphere.NewSetVSphereEndpointInternalServerError().WithPayload(Err(err))
	}

	vcURL.Path = "/sdk"

	vsphereInsecure := false
	configReaderWriter := app.clientTkg.TKGConfigReaderWriter()
	vsphereInsecureString, err := configReaderWriter.Get(constants.ConfigVariableVsphereInsecure)
	if err == nil {
		vsphereInsecure = (vsphereInsecureString == trueStr)
	}

	if params.Credentials.Insecure != nil && *params.Credentials.Insecure {
		vsphereInsecure = true
	}

	vcClient, err := vc.NewClient(vcURL, params.Credentials.Thumbprint, vsphereInsecure)
	if err != nil {
		return vsphere.NewSetVSphereEndpointInternalServerError().WithPayload(Err(err))
	}

	_, err = vcClient.Login(params.HTTPRequest.Context(), params.Credentials.Username, params.Credentials.Password)
	if err != nil {
		return vsphere.NewSetVSphereEndpointInternalServerError().WithPayload(Err(err))
	}

	app.vcClient = vcClient

	version, build, err := vcClient.GetVSphereVersion()
	if err != nil {
		return vsphere.NewSetVSphereEndpointInternalServerError().WithPayload(Err(err))
	}

	res := models.VsphereInfo{
		Version:    fmt.Sprintf("%s:%s", version, build),
		HasPacific: "no",
	}

	if hasPP, err := vcClient.DetectPacific(params.HTTPRequest.Context()); err == nil && hasPP {
		res.HasPacific = "yes"
	}

	return vsphere.NewSetVSphereEndpointCreated().WithPayload(&res)
}

// Need until the vcClient code decouples from the presentation code in the
// management-cluster API logic.
func vspMgmtClusterParamsToMCUIManagementParams(params *models.VsphereManagementClusterParams) (*mcuimodels.VsphereRegionalClusterParams, error) {
	// Should be same structure, so we can marshal through JSON.
	// Easier this way since there are nested model structs.
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	result := &mcuimodels.VsphereRegionalClusterParams{}
	err = json.Unmarshal(jsonData, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Need until the vcClient code decouples from the presentation code in the
// management-cluster API logic.
func vspmcUIObjectToObject(mobs []*mcuimodels.VSphereManagementObject) []*models.VSphereManagementObject {
	result := []*models.VSphereManagementObject{}

	for _, mob := range mobs {
		result = append(result, &models.VSphereManagementObject{
			Moid:         mob.Moid,
			Name:         mob.Name,
			ParentMoid:   mob.ParentMoid,
			Path:         mob.Path,
			ResourceType: mob.ResourceType,
		})
	}

	return result
}

// Need until the vcClient code decouples from the presentation code in the
// management-cluster API logic.
func vspmcUIDatacenterToDatacenter(dcs []*mcuimodels.VSphereDatacenter) []*models.VSphereDatacenter {
	result := []*models.VSphereDatacenter{}

	for _, dc := range dcs {
		result = append(result, &models.VSphereDatacenter{
			Moid: dc.Moid,
			Name: dc.Name,
		})
	}

	return result
}

// Need until the vcClient code decouples from the presentation code in the
// management-cluster API logic.
func vspmcUIDatastoreToDatastore(dcs []*mcuimodels.VSphereDatastore) []*models.VSphereDatastore {
	result := []*models.VSphereDatastore{}

	for _, dc := range dcs {
		result = append(result, &models.VSphereDatastore{
			Moid: dc.Moid,
			Name: dc.Name,
		})
	}

	return result
}

// Need until the vcClient code decouples from the presentation code in the
// management-cluster API logic.
func vspmcUIFolderToFolder(dcs []*mcuimodels.VSphereFolder) []*models.VSphereFolder {
	result := []*models.VSphereFolder{}

	for _, dc := range dcs {
		result = append(result, &models.VSphereFolder{
			Moid: dc.Moid,
			Name: dc.Name,
		})
	}

	return result
}

// Need until the vcClient code decouples from the presentation code in the
// management-cluster API logic.
func vspmcUINetworksToNetworks(dcs []*mcuimodels.VSphereNetwork) []*models.VSphereNetwork {
	result := []*models.VSphereNetwork{}

	for _, dc := range dcs {
		result = append(result, &models.VSphereNetwork{
			Moid:        dc.Moid,
			Name:        dc.Name,
			DisplayName: dc.DisplayName,
		})
	}

	return result
}

// Need until the vcClient code decouples from the presentation code in the
// management-cluster API logic.
func vspmcUIResourcePoolToResourcePool(dcs []*mcuimodels.VSphereResourcePool) []*models.VSphereResourcePool {
	result := []*models.VSphereResourcePool{}

	for _, dc := range dcs {
		result = append(result, &models.VSphereResourcePool{
			Moid: dc.Moid,
			Name: dc.Name,
		})
	}

	return result
}
