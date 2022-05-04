// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/cluster"
	tfclient "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/client"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/log"
)

// checkCurrentManagementCluster validates that the request is using the same management cluster as the
// currently set context. This is a limitation for now until we decide if we want to dynamically change
// the users context based on UI operations.
func (app *App) checkCurrentManagementCluster(mgmtClusterName string, tkgClient *tfclient.TkgClient) error {
	currentMC, err := tkgClient.GetCurrentRegionContext()
	if err != nil {
		return err
	}

	if currentMC.ClusterName != mgmtClusterName {
		return fmt.Errorf("requested cluster info for management cluster %q, but current context is using %q", mgmtClusterName, currentMC.ClusterName)
	}

	return nil
}

// CreateCluster creates a new workload cluster.
func (app *App) CreateCluster(params cluster.CreateWorkloadClusterParams) middleware.Responder {
	tkgClient, err := app.getTkgClient()
	if err != nil {
		return cluster.NewCreateWorkloadClusterBadRequest().WithPayload(Err(err))
	}

	err = app.checkCurrentManagementCluster(params.ManagementClusterName, tkgClient)
	if err != nil {
		return cluster.NewCreateWorkloadClusterBadRequest().WithPayload(Err(err))
	}

	createParams := params.Params
	configOpts := tfclient.ClusterConfigOptions{
		ClusterName:     createParams.Name,
		TargetNamespace: createParams.Namespace,
	}
	if createParams.Controlplanecount != 0 {
		configOpts.ControlPlaneMachineCount = &createParams.Controlplanecount
	}
	if createParams.Workernodecount != 0 {
		configOpts.WorkerMachineCount = &createParams.Workernodecount
	}

	nodeSizeOpts := tfclient.NodeSizeOptions{}
	if createParams.Controlplanesize != "" {
		nodeSizeOpts.ControlPlaneSize = createParams.Controlplanesize
	}
	if createParams.Workernodesize != "" {
		nodeSizeOpts.WorkerSize = createParams.Workernodesize
	}

	createOpts := tfclient.CreateClusterOptions{
		ClusterConfigOptions:        configOpts,
		CniType:                     createParams.Cni,
		VsphereControlPlaneEndpoint: createParams.Cpendpoint,
		ClusterType:                 tfclient.WorkloadCluster,
		TKRVersion:                  createParams.Tkrversion,
		NodeSizeOptions:             nodeSizeOpts,
	}

	err = tkgClient.CreateCluster(&createOpts, false)
	if err != nil {
		return cluster.NewCreateWorkloadClusterBadRequest().WithPayload(Err(err))
	}
	return cluster.NewCreateWorkloadClusterOK().WithPayload(nil)
}

// GetCluster gets details of a specific workload cluster.
func (app *App) GetCluster(params cluster.GetWorkloadClusterParams) middleware.Responder {
	tkgClient, err := app.getTkgClient()
	if err != nil {
		return cluster.NewDeleteWorkloadClusterBadRequest().WithPayload(Err(err))
	}
	err = app.checkCurrentManagementCluster(params.ManagementClusterName, tkgClient)
	if err != nil {
		return cluster.NewCreateWorkloadClusterBadRequest().WithPayload(Err(err))
	}

	clusters, err := app.getClusters(params.ClusterName, *params.ClusterNamespace, tkgClient)
	if err != nil {
		return cluster.NewGetWorkloadClusterBadRequest().WithPayload(Err(err))
	}

	if len(clusters) > 0 {
		return cluster.NewGetWorkloadClusterOK().WithPayload(clusters[0])
	}

	return cluster.NewGetWorkloadClusterOK()
}

// GetClusters gets all workload clusters from a given management cluster.
func (app *App) GetClusters(params cluster.GetWorkloadClustersParams) middleware.Responder {
	tkgClient, err := app.getTkgClient()
	if err != nil {
		return cluster.NewDeleteWorkloadClusterBadRequest().WithPayload(Err(err))
	}
	err = app.checkCurrentManagementCluster(params.ManagementClusterName, tkgClient)
	if err != nil {
		return cluster.NewCreateWorkloadClusterBadRequest().WithPayload(Err(err))
	}

	clusters, err := app.getClusters("", *params.ClusterNamespace, tkgClient)
	if err != nil {
		return cluster.NewGetWorkloadClustersInternalServerError().WithPayload(Err(err))
	}
	return cluster.NewGetWorkloadClustersOK().WithPayload(clusters)
}

// DeleteCluster triggers the deletion of a workload cluster.
func (app *App) DeleteCluster(params cluster.DeleteWorkloadClusterParams) middleware.Responder {
	tkgClient, err := app.getTkgClient()
	if err != nil {
		return cluster.NewDeleteWorkloadClusterBadRequest().WithPayload(Err(err))
	}
	err = app.checkCurrentManagementCluster(params.ManagementClusterName, tkgClient)
	if err != nil {
		return cluster.NewCreateWorkloadClusterBadRequest().WithPayload(Err(err))
	}

	clusters, err := app.getClusters(params.ClusterName, *params.ClusterNamespace, tkgClient)
	if err != nil {
		return cluster.NewDeleteWorkloadClusterBadRequest().WithPayload(Err(err))
	}

	if len(clusters) == 0 {
		// Cluster isn't there, so what they wanted is true
		return cluster.NewDeleteWorkloadClusterOK()
	} else if len(clusters) > 1 {
		// There are more than one clusters with a matching name and namespace
		// must not have been provided.
		return cluster.NewDeleteWorkloadClusterBadRequest().WithPayload(Err(fmt.Errorf("more than one cluster named %q, namespace is required", params.ClusterName)))
	}

	deleteOpts := tfclient.DeleteWorkloadClusterOptions{
		ClusterName: clusters[0].Name,
		Namespace:   clusters[0].Namespace,
	}

	err = tkgClient.DeleteWorkloadCluster(deleteOpts)
	if err != nil {
		return cluster.NewDeleteWorkloadClusterInternalServerError().WithPayload(Err(err))
	}

	log.Infof("\nWorkload cluster %q is being deleted!\n", clusters[0].Name)

	return cluster.NewDeleteWorkloadClusterOK()
}

func (app *App) getClusters(name, namespace string, tkgClient tfclient.Client) ([]*models.WorkloadCluster, error) {
	apiClusters := []*models.WorkloadCluster{}

	listOpts := tfclient.ListTKGClustersOptions{}
	if namespace != "" {
		listOpts.Namespace = namespace
	}

	clusters, err := tkgClient.ListTKGClusters(listOpts)
	if err != nil {
		return apiClusters, err
	}

	// TODO: The internal API uses whatever management cluster that is the currently
	// set context. Need to add handling to check if context needs to be changed, or
	// change our API to not be based on management cluster path.
	for _, cluster := range clusters {
		if name != "" && cluster.Name != name {
			continue
		}

		apiCluster := models.WorkloadCluster{
			Cpcount:    cluster.ControlPlaneCount,
			K8sversion: cluster.K8sVersion,
			Name:       cluster.Name,
			Namespace:  cluster.Namespace,
			Status:     cluster.Status,
			Wncount:    cluster.WorkerCount,
			Plan:       cluster.Plan,
		}
		apiClusters = append(apiClusters, &apiCluster)
	}

	return apiClusters, nil
}
