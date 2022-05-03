// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/management"
	tfclient "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/client"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/log"
)

// DeleteMgmtCluster triggers the deletion of a management cluster.
func (app *App) DeleteMgmtCluster(params management.DeleteMgmtClusterParams) middleware.Responder {
	tkgClient, err := app.getTkgClient()
	if err != nil {
		return management.NewDeleteMgmtClusterBadRequest().WithPayload(Err(err))
	}

	deleteOptions := tfclient.DeleteRegionOptions{
		ClusterName: params.ManagementClusterName,
		Force:       true,
	}

	err = tkgClient.DeleteRegion(deleteOptions)
	if err != nil {
		return management.NewDeleteMgmtClusterInternalServerError().WithPayload(Err(err))
	}

	log.Infof("\nManagement cluster deleted!\n")

	return management.NewDeleteMgmtClusterOK()
}

// GetMgmtClusters gets all management clusters.
func (app *App) GetMgmtClusters(params management.GetMgmtClustersParams) middleware.Responder {
	clusters, err := app.getMgmtClusters("")
	if err != nil {
		return management.NewGetMgmtClustersInternalServerError().WithPayload(Err(err))
	}
	return management.NewGetMgmtClustersOK().WithPayload(clusters)
}

// GetMgmtCluster gets details of a specific management cluster.
func (app *App) GetMgmtCluster(params management.GetMgmtClusterParams) middleware.Responder {
	clusters, err := app.getMgmtClusters(params.ManagementClusterName)
	if err != nil {
		return management.NewGetMgmtClusterBadRequest().WithPayload(Err(err))
	}

	if len(clusters) > 0 {
		return management.NewGetMgmtClusterOK().WithPayload(clusters[0])
	}

	return management.NewGetMgmtClusterOK()
}

func (app *App) getMgmtClusters(name string) ([]*models.ManagementCluster, error) {
	apiClusters := []*models.ManagementCluster{}
	tkgClient, err := app.getTkgClient()
	if err != nil {
		return apiClusters, err
	}

	clusters, err := tkgClient.GetRegionContexts(name)
	if err != nil {
		return apiClusters, err
	}

	// TODO: We need to reconcile what we want to be able to show to end users
	// with what is actually available through the internal API.
	for _, cluster := range clusters {
		apiCluster := models.ManagementCluster{
			Context: cluster.ContextName,
			Name:    cluster.ClusterName,
		}
		apiClusters = append(apiClusters, &apiCluster)
	}

	return apiClusters, nil
}

// GetClusterClasses gets all cluster classes available on a given management cluster.
func (app *App) GetClusterClasses(params management.GetClusterClassesParams) middleware.Responder {
	// TODO: We are pinned to v0.11.x of tanzu-framework where ClusterClass is not available yet.
	// This is a placeholder until we are able to upgrade to a version that does have support.
	return management.NewGetClusterClassesOK()
}

// GetClusterClass gets details of a cluster class on a given management cluster.
func (app *App) GetClusterClass(params management.GetClusterClassParams) middleware.Responder {
	// TODO: We are pinned to v0.11.x of tanzu-framework where ClusterClass is not available yet.
	// This is a placeholder until we are able to upgrade to a version that does have support.
	return management.NewGetClusterClassesOK()
}
