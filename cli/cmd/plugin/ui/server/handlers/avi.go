// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/avi"
	aviclient "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/avi"
	mcuimodels "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/web/server/models"
)

// GetAviClouds handles requests to GET avi clouds
func (app *App) GetAviClouds(params avi.GetAviCloudsParams) middleware.Responder {
	if app.aviClient == nil {
		return avi.NewGetAviCloudsInternalServerError().WithPayload(Err(errors.New("avi client is not initialized properly")))
	}

	aviClouds, err := app.aviClient.GetClouds()
	if err != nil {
		return avi.NewGetAviCloudsInternalServerError().WithPayload(Err(errors.Wrap(err, "unable to get avi clouds")))
	}

	return avi.NewGetAviCloudsOK().WithPayload(mcUIAviCloudToAviCloud(aviClouds))
}

// GetAviServiceEngineGroups handles requests to GET avi service engine groups
func (app *App) GetAviServiceEngineGroups(params avi.GetAviServiceEngineGroupsParams) middleware.Responder {
	if app.aviClient == nil {
		return avi.NewGetAviServiceEngineGroupsInternalServerError().WithPayload(Err(errors.New("avi client is not initialized properly")))
	}

	aviServiceEngineGroups, err := app.aviClient.GetServiceEngineGroups()
	if err != nil {
		return avi.NewGetAviServiceEngineGroupsInternalServerError().WithPayload(Err(errors.Wrap(err, "unable to get avi service engine groups")))
	}

	return avi.NewGetAviServiceEngineGroupsOK().WithPayload(mcUIAviServiceEngineGroupToAviServiceEngineGroup(aviServiceEngineGroups))
}

// GetAviVipNetworks handles requests to GET avi VIP networks
func (app *App) GetAviVipNetworks(params avi.GetAviVipNetworksParams) middleware.Responder {
	if app.aviClient == nil {
		return avi.NewGetAviVipNetworksInternalServerError().WithPayload(Err(errors.New("avi client is not initialized properly")))
	}

	aviVipNetworks, err := app.aviClient.GetVipNetworks()
	if err != nil {
		return avi.NewGetAviVipNetworksInternalServerError().WithPayload(Err(errors.Wrap(err, "unable to get avi VIP networks")))
	}

	return avi.NewGetAviVipNetworksOK().WithPayload(mcUIAviVipNetworkToAviVipNetwork(aviVipNetworks))
}

// VerifyAccount validates avi credentials and sets the avi client into web app
func (app *App) VerifyAccount(params avi.VerifyAccountParams) middleware.Responder {
	aviControllerParams := &mcuimodels.AviControllerParams{
		Username: params.Credentials.Username,
		Password: params.Credentials.Password,
		Host:     params.Credentials.Host,
		Tenant:   params.Credentials.Tenant,
		CAData:   params.Credentials.CAData,
	}

	app.aviClient = aviclient.New()
	authed, err := app.aviClient.VerifyAccount(aviControllerParams)
	if err != nil {
		return avi.NewVerifyAccountInternalServerError().WithPayload(Err(err))
	}

	if !authed {
		return avi.NewVerifyAccountInternalServerError().WithPayload(Err(errors.Errorf("unable to authenticate due to incorrect credentials")))
	}

	return avi.NewVerifyAccountCreated()
}

// Need until the aviClient code decouples from the presentation code in the
// management-cluster API logic.
func mcUIAviCloudToAviCloud(instances []*mcuimodels.AviCloud) []*models.AviCloud {
	result := []*models.AviCloud{}

	for _, instance := range instances {
		result = append(result, &models.AviCloud{
			Location: instance.Location,
			Name:     instance.Name,
			UUID:     instance.UUID,
		})
	}

	return result
}

// Need until the aviClient code decouples from the presentation code in the
// management-cluster API logic.
func mcUIAviServiceEngineGroupToAviServiceEngineGroup(instances []*mcuimodels.AviServiceEngineGroup) []*models.AviServiceEngineGroup {
	result := []*models.AviServiceEngineGroup{}

	for _, instance := range instances {
		result = append(result, &models.AviServiceEngineGroup{
			Location: instance.Location,
			Name:     instance.Name,
			UUID:     instance.UUID,
		})
	}

	return result
}

// Need until the aviClient code decouples from the presentation code in the
// management-cluster API logic.
func mcUIAviVipNetworkToAviVipNetwork(instances []*mcuimodels.AviVipNetwork) []*models.AviVipNetwork {
	result := []*models.AviVipNetwork{}

	for _, instance := range instances {
		newInst := &models.AviVipNetwork{
			Cloud: instance.Cloud,
			Name:  instance.Name,
			UUID:  instance.UUID,
		}

		for _, s := range instance.ConfigedSubnets {
			newInst.ConfigedSubnets = append(newInst.ConfigedSubnets, &models.AviSubnet{
				Family: s.Family,
				Subnet: s.Subnet,
			})
		}

		result = append(result, newInst)
	}

	return result
}
