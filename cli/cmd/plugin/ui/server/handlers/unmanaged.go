// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/go-openapi/runtime/middleware"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
    "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/pkg/containerruntime"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/unmanaged"
)

// unmanagedCluster is the cluster information we get from the CLI output.
type unmanagedCluster struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Status   string `json:"status"`
}

// checkUnmanagedPlugin verifies we can execute the "tanzu unmanaged-cluster" command.
func checkUnmanagedPlugin() error {
	cmd := exec.Command("tanzu", "unmanaged-cluster")
	err := cmd.Run()
	if err != nil {
		return errors.New("tanzu unmanaged-cluster could not be found")
	}
	return nil
}

func checkDockerContainerRunning() error {
    _, err := containerruntime.GetRuntimeInfo()
    return err
}

func runCommand(args ...string) (string, error) {
	cmdArgs := []string{"unmanaged-cluster"}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("tanzu", cmdArgs...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// CreateUnmanagedCluster creates a new unmanaged cluster.
func (app *App) CreateUnmanagedCluster(params unmanaged.CreateUnmanagedClusterParams) middleware.Responder {
    fmt.Println("checking unmanaged plugin...")
	if err := checkUnmanagedPlugin(); err != nil {
	    fmt.Println("ERROR: unmanaged plugin cannot be found")
		return unmanaged.NewDeleteUnmanagedClusterInternalServerError().WithPayload(Err(err))
	}
	fmt.Println("unmanaged plugin found")

    fmt.Println("checking docker container running...")
	if err := checkDockerContainerRunning(); err != nil {
	    fmt.Println("ERROR: docker container cannot be found")
		return unmanaged.NewDeleteUnmanagedClusterInternalServerError().WithPayload(Err(err))
	}
    fmt.Println("docker container running")

	createParams := params.Params
	if createParams.Name == "" {
		return unmanaged.NewCreateUnmanagedClusterBadRequest().WithPayload(Err(fmt.Errorf("cluster name must be provided")))
	}

	args := []string{
		"create",
		createParams.Name,
	}

	if createParams.Provider != "" {
		args = append(args, "--provider", createParams.Provider)
	}

	if createParams.Cni != "" {
		args = append(args, "--cni", createParams.Cni)
	}

	if createParams.Podcidr != "" {
		args = append(args, "--pod-cidr", createParams.Podcidr)
	}

	if createParams.Servicecidr != "" {
		args = append(args, "--service-cidr", createParams.Servicecidr)
	}

	if len(createParams.Portmappings) > 0 {
		for _, pm := range createParams.Portmappings {
			args = append(args, "--port-map", pm)
		}
	}

	if createParams.Controlplanecount > 0 {
		args = append(args, "--control-plane-node-count", strconv.FormatInt(createParams.Controlplanecount, 10))
	}

	if createParams.Workernodecount > 0 {
		args = append(args, "--worker-node-count", strconv.FormatInt(createParams.Workernodecount, 10))
	}

	// TODO: Add streaming of create output
	_, err := runCommand(args...)
	if err != nil {
		return unmanaged.NewDeleteUnmanagedClusterInternalServerError().WithPayload(Err(err))
	}

	cluster, err := app.getUnmanagedCluster(createParams.Name)
	if err != nil {
		e := fmt.Errorf("cluster created but could not be found: %s", err.Error())
		return unmanaged.NewDeleteUnmanagedClusterInternalServerError().WithPayload(Err(e))
	}

	return unmanaged.NewCreateUnmanagedClusterOK().WithPayload(cluster)
}

// GetUnmanagedCluster gets details of a specific unmanaged cluster.
func (app *App) GetUnmanagedCluster(params unmanaged.GetUnmanagedClusterParams) middleware.Responder {
	if err := checkUnmanagedPlugin(); err != nil {
		return unmanaged.NewGetUnmanagedClusterInternalServerError().WithPayload(Err(err))
	}

	cluster, err := app.getUnmanagedCluster(params.ClusterName)
	if err != nil {
		return unmanaged.NewGetUnmanagedClusterInternalServerError().WithPayload(Err(err))
	}

	return unmanaged.NewGetUnmanagedClusterOK().WithPayload(cluster)
}

// getUnmanagedCluster gets details for a specific cluster.
func (app *App) getUnmanagedCluster(clusterName string) (*models.UnmanagedCluster, error) {
	clusters, err := app.getUnmanagedClusters()
	if err != nil {
		return nil, err
	}

	for _, cluster := range clusters {
		if cluster.Name == clusterName {
			return cluster, nil
		}
	}

	return nil, fmt.Errorf("unmanaged cluster %q could not be found", clusterName)
}

// getUnmanagedClusters gets a list of all unmanaged clusters.
func (app *App) getUnmanagedClusters() ([]*models.UnmanagedCluster, error) {
	jsonOutput, err := runCommand("list", "-o", "json")
	if err != nil {
		return nil, err
	}

	var clusters []unmanagedCluster
	err = json.Unmarshal([]byte(jsonOutput), &clusters)
	if err != nil {
		return nil, fmt.Errorf("unable to parse unmanaged cluster information: %s", err.Error())
	}

	results := []*models.UnmanagedCluster{}
	for _, cluster := range clusters {
		results = append(results, &models.UnmanagedCluster{
			Name:     cluster.Name,
			Provider: cluster.Provider,
			Status:   cluster.Status,
		})
	}

	return results, nil
}

// GetUnmanagedClusters gets all unmanaged clusters.
func (app *App) GetUnmanagedClusters(params unmanaged.GetUnmanagedClustersParams) middleware.Responder {
	if err := checkUnmanagedPlugin(); err != nil {
		return unmanaged.NewGetUnmanagedClusterInternalServerError().WithPayload(Err(err))
	}

	clusters, err := app.getUnmanagedClusters()
	if err != nil {
		return unmanaged.NewGetUnmanagedClusterInternalServerError().WithPayload(Err(err))
	}

	return unmanaged.NewGetUnmanagedClustersOK().WithPayload(clusters)
}

// DeleteUnmanagedCluster triggers the deletion of an unmanaged cluster.
func (app *App) DeleteUnmanagedCluster(params unmanaged.DeleteUnmanagedClusterParams) middleware.Responder {
	if err := checkUnmanagedPlugin(); err != nil {
		return unmanaged.NewDeleteUnmanagedClusterInternalServerError().WithPayload(Err(err))
	}

	// TODO: Add streaming of delete output
	// TODO #2: There is some talk of adding an "are you sure" prompt to the
	// tanzu uc rm command. If so, when we update to that release we will need
	// to add a "-y" argument here (if that's what ends up being implemented) to
	// tell it to not prompt for confirmation.
	_, err := runCommand("delete", params.ClusterName)
	if err != nil {
		return unmanaged.NewDeleteUnmanagedClusterInternalServerError().WithPayload(Err(err))
	}

	return unmanaged.NewDeleteUnmanagedClusterOK()
}
