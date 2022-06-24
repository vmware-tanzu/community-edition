// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cluster

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"

	ucconfig "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/config"
	uclogger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/tanzu"
	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/checks"
	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/config"
	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/internal/docker"
	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/kubeconfig"
	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/utils"
)

var log *logrus.Logger

// Cluster provides a way to manage a TCE unmanaged cluster.
type Cluster struct {
}

// TCECluster defines the interface for Cluster operations.
type TCECluster interface {
	CreateCluster() config.Response
	DeleteCluster() config.Response
	ClusterStatus() config.Response
	GetKubeconfig() config.Response
	Logs() config.Response
	GetJSONResponse(res *config.Response) string
}

// New instantiates a new instance of a Cluster.
func New(parentLogger *logrus.Logger) TCECluster {
	log = parentLogger
	return &Cluster{}
}

// CreateCluster creates a new cluster.
//nolint:funlen
func (c *Cluster) CreateCluster() config.Response {
	lock, err := utils.GetFileLockWithTimeOut(utils.GetClusterCreateLockFilename(), utils.DefaultLockTimeout)
	if err != nil {
		log.Errorf("Error creating cluster, cannot init lock, reason: %v", err)
		// Lock is already in place, which means that process is already running, just return
		return config.RunningResponse()
	}

	defer func() {
		if err := lock.Unlock(); err != nil {
			log.Errorf("Error after creating cluster, cannot unlock %q, reason: %v", lock, err)
		}
	}()

	// Get cluster state. If already running, return already running, else
	log.Debug("Checking to see if there's a cluster already running")
	s, _ := checks.GetContainerClusterStatus()
	if s == checks.Running {
		log.Info("Cluster is already running")
		return config.RunningResponse()
	}

	//  If the cluster exists and is not running, delete everything so that it can be safely created
	removeAllConfigFiles(false)
	if s != checks.NotExist {
		log.Info("Deleting cluster container and configuration for a cluster that is not running")
		err := docker.ForceStopAndDeleteCluster()
		if err != nil {
			log.Errorf("Error force deleting the TCE container (%s)", err.Error())
		}
	}

	// Execute preflight checks. If everything ok, proceed, otherwise, return error
	log.Debug("Running preflight checks")
	err = checks.PreflightChecks()
	if err != nil {
		log.Errorf("Preflight check error (%s)", err.Error())
		return config.Response{
			Status:       config.Error,
			Description:  "Cluster can not be created",
			ErrorMessage: err.Error(),
			Error:        true,
		}
	}

	// Generate Config
	configFile := config.GetClusterConfigFileName()
	log.Infof("Processing configuration and storing it at %s", configFile)
	//nolint:gosec
	cmd := exec.Command(config.YttBinary, "-f", config.ClusterInstallTemplateFile, "-f", config.ClusterInstallValuesFile, "--ignore-unknown-comments")
	output, err := cmd.Output()
	if err != nil {
		// Error while running ytt
		log.Errorf("Error processing config (%s)", err.Error())
		return config.Response{
			Status:       config.Error,
			Description:  "Cluster can not be created",
			ErrorMessage: err.Error(),
			Error:        true,
		}
	}

	err = config.WriteConfigFile(output, configFile)
	if err != nil {
		// Error writing config
		log.Errorf("Error writing config (%s)", err.Error())
		return config.Response{
			Status:       config.Error,
			Description:  "Cluster can not be created",
			ErrorMessage: err.Error(),
			Error:        true,
		}
	}

	log.Debugf("Create cluster with config at %s", configFile)
	configArgs := map[string]interface{}{
		ucconfig.ClusterConfigFile:   configFile,
		ucconfig.ClusterName:         config.DefaultClusterName,
		ucconfig.SkipPreflightChecks: true,
	}
	clusterConfig, err := ucconfig.InitializeConfiguration(configArgs)
	if err != nil {
		log.Errorf("Failed to initialize configuration. Error (%s)\n", err.Error())
		return config.Response{
			Status:       config.Error,
			Description:  "Cluster configuration could not be initialized",
			ErrorMessage: err.Error(),
			Error:        true,
		}
	}

	tm := tanzu.New(uclogger.NewLogger(true, 0))
	exitCode, err := tm.Deploy(clusterConfig)
	if err != nil {
		log.Errorf("Error while creating the cluster (%s), code %d", err.Error(), exitCode)
		return config.Response{
			Status:       config.Error,
			Description:  "Cluster can not be created",
			ErrorMessage: err.Error(),
			Error:        true,
		}
	}
	log.Info("Cluster successfully created")

	// Copy kubeconfig to host if everything went ok
	err = kubeconfig.AddConfig(config.DefaultHomeKubeConfig, config.DefaultHostMountedKubeConfig)
	if err != nil {
		// TODO: Maybe print this to the user, although return a running response with this as errorMessage
		log.Errorf("Error while adding kubeconfig to host (%s)", err)
	}
	log.Debug("Kubeconfig copied to host")

	copyConfigFiles()
	log.Info("Tanzu config files copied to host")

	return config.RunningResponse()
}

// DeleteCluster will delete a cluster.
func (c *Cluster) DeleteCluster() config.Response {
	lock, err := utils.GetFileLockWithTimeOut(utils.GetClusterDeleteLockFilename(), utils.DefaultLockTimeout)
	if err != nil {
		log.Info("Cluster already deleting")
		// Lock is already in place, which means that process is already running, just return
		return config.DeletingResponse()
	}

	defer func() {
		if err := lock.Unlock(); err != nil {
			log.Errorf("Error after deleting cluster, cannot unlock %q, reason: %v", lock, err)
		}
	}()

	// Get cluster state. If already running, return already running, else
	status, _ := checks.GetContainerClusterStatus()
	if status == checks.Running {
		log.Debug("There's an existing cluster")
	} else {
		log.Debug("There's no running cluster")
	}

	if status != checks.NotExist {
		log.Info("Deleting cluster")
		tm := tanzu.New(uclogger.NewLogger(true, 0))
		if err := tm.Delete(config.DefaultClusterName); err != nil {
			log.Errorf("Error while deleting the cluster (%s)", err.Error())
			log.Info("Deleting cluster container and configuration for a cluster that is not running")
			err := docker.ForceStopAndDeleteCluster()
			if err != nil {
				log.Errorf("Error force deleting the TCE container (%s)", err)
				return config.Response{
					Status:       config.Error,
					Description:  "Cluster can not be deleted",
					ErrorMessage: err.Error(),
					Error:        true,
				}
			}
		}
		log.Info("Cluster successfully deleted")
	}

	removeAllConfigFiles(true)

	// Copy kubeconfig to host if everything went ok
	err = kubeconfig.RemoveNamedConfig(config.DefaultClusterName, config.DefaultHostMountedKubeConfig)
	if err != nil {
		log.Errorf("Error while removing kubeconfig from host (%s)", err)
	}
	log.Info("Kubeconfig removed from host")

	return config.DeletedResponse()
}

// ClusterStatus queries the current status of a cluster.
func (c *Cluster) ClusterStatus() config.Response {
	creating, err := checks.IsClusterCreating()
	if err != nil {
		log.Debugf("Error while checking cluster status (%s)", err.Error())
		return config.Response{
			Status:       config.Error,
			Description:  "Cluster can not be created",
			ErrorMessage: err.Error(),
			Error:        true,
		}
	}

	if creating {
		log.Info("Status: cluster is being created")
		return config.Response{
			Status:       config.Creating,
			Description:  "Cluster creating",
			ErrorMessage: "",
			Error:        false,
		}
	}

	_, canbecreated, err := checks.IsClusterUpAndRunning()
	if canbecreated {
		log.Debug("Status: TCE cluster does not exist")
		return config.Response{
			Status:       config.NotExists,
			Description:  "Cluster does not exist",
			ErrorMessage: "",
			Error:        false,
		}
	}

	if err != nil {
		// TODO: Check for status Stopped
		log.Errorf("Error while checking cluster status (%s)", err.Error())
		return config.Response{
			Status:       config.Error,
			Description:  "Cluster error",
			ErrorMessage: err.Error(),
			Error:        true,
		}
	}

	deleting, err := checks.IsClusterDeleting()
	if err != nil {
		log.Errorf("Error while checking cluster status (%s)", err.Error())
		return config.Response{
			Status:       config.Error,
			Description:  "Cluster can not be created",
			ErrorMessage: err.Error(),
			Error:        true,
		}
	}

	if deleting {
		log.Debug("Cluster is deleting")
		return config.Response{
			Status:       config.Deleting,
			Description:  "Cluster deleting",
			ErrorMessage: "",
			Error:        false,
		}
	}

	log.Debug("Cluster is running")
	return config.RunningResponse()
}

// copyConfigFiles will make a copy of the TCE configuration files (~/.config/tanzu)
// under the /opt directory.
func copyConfigFiles() {
	// TODO: Better copy the local tanzu files to the volume mount so they are on the host for diagnosis
	//nolint:gosec
	err := exec.Command("cp", "-rf", filepath.Join(config.GetUserHome(), ".config", "tanzu"), "/opt").Run()
	if err != nil {
		log.Errorf("Couldn't copy tanzu config files (%s)", err)
	}
}

// removeAllConfigFiles cleans up all configuration files.
func removeAllConfigFiles(isdelete bool) {
	log.Info("Removing all internal configuration files")

	removalPaths := []string{
		filepath.Join(config.GetUserHome(), ".kube"),
		filepath.Join(config.GetUserHome(), ".tanzu"),
		filepath.Join(config.GetUserHome(), ".config", "tanzu"),
		config.GetClusterConfigFileName(),
	}

	if isdelete {
		removalPaths = append(removalPaths, config.GetLogsFileName())
	}

	for _, path := range removalPaths {
		if err := os.RemoveAll(path); err != nil {
			log.Errorf("Couldn't remove %s (%s)", path, err)
		}
	}
}

// Logs retrieves the cluster log output.
func (c *Cluster) Logs() config.Response {
	up, _, _ := checks.IsClusterUpAndRunning()
	if !up {
		return config.NotExistsResponse()
	}

	content, err := os.ReadFile(config.GetInternalLogsFileName())
	if err != nil {
		log.Errorf("Error reading cluster log file (%s)", err.Error())
		return config.Response{
			Status:       config.Error,
			Description:  "cluster log file can not be read",
			ErrorMessage: err.Error(),
			Error:        true,
		}
	}

	return config.Response{
		Output: string(content),
	}
}

// GetKubeconfig returns the content of the kubeconfig for a cluster.
func (c *Cluster) GetKubeconfig() config.Response {
	up, _, _ := checks.IsClusterUpAndRunning()
	if !up {
		return config.NotExistsResponse()
	}

	// Read a modified kubeconfig
	content, err := kubeconfig.GetConfig(config.GetKubeconfigFileName())
	if err != nil {
		log.Errorf("Error reading kubeconfig file (%s)", err.Error())
		return config.Response{
			Status:       config.Error,
			Description:  "Kubeconfig file can not be read",
			ErrorMessage: err.Error(),
			Error:        true,
		}
	}

	return config.Response{
		Output: string(content),
	}
}

// GetJSONResponse formats a JSON string representation of a command response.
func (c *Cluster) GetJSONResponse(res *config.Response) string {
	byteArray, err := json.MarshalIndent(res, "", "  ")

	if err != nil {
		log.Println(err)
	}

	return (string(byteArray))
}
