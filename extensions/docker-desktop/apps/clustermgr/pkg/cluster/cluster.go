// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cluster

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/checks"
	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/config"
	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/internal/docker"
	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/kubeconfig"
	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/utils"
)

var log *logrus.Logger

type Cluster struct {
}

type TCECluster interface {
	CreateCluster() config.Response
	DeleteCluster() config.Response
	ClusterStatus() config.Response
	ProvisionIngress() config.Response
	ProvisionCertMan() config.Response
	GetKubeconfig() config.Response
	Logs() config.Response
	Stats() config.Response
	Reset() config.Response
	GetJSONResponse(res *config.Response) string
}

func New(parentLogger *logrus.Logger) TCECluster {
	log = parentLogger
	return &Cluster{}
}

// CreateCluster creates a new cluster.
//nolint:funlen
func (c *Cluster) CreateCluster() config.Response {
	log.Info("Create cluster")
	// LockCreationOrExitIfAlreadyCreating  --> touch $HOME/tce-cluster-create.tag
	lock, err := utils.GetFileLockWithTimeOut(utils.GetClusterCreateLockFilename(), utils.DefaultLockTimeout)
	if err != nil {
		log.Errorf("cannot init lock. reason: %v", err)
		// Lock is already in place, which means that process is already running, just return
		return config.RunningResponse()
	}

	// releaseLocks --> rm -f $HOME/tce-cluster-create.tag $HOME/cluster-config-$$.yaml
	defer func() {
		if err := lock.Unlock(); err != nil {
			log.Errorf("cannot unlock %q, reason: %v", lock, err)
		}
	}()

	// Initial response
	ret := config.RunningResponse()

	// Get cluster state. If already running, return already running, else
	log.Info("Check to see if there's a cluster already running")
	s, _ := checks.GetContainerClusterStatus()
	if s == checks.Running {
		log.Info("Cluster is already running")
		return config.RunningResponse()
	}
	//  If the cluster exists and is not running, delete everything so that it can be safely created
	log.Info("Deleting cluster container and configuration for a cluster that is not running")
	removeAllConfigFiles(false)
	if s != checks.NotExist {
		log.Info("Force delete non running container")
		err := docker.ForceStopAndDeleteCluster()
		if err != nil {
			log.Errorf("Error force deleting the TCE container (%s)", err)
		}
	}

	// Execute preflight checks. If everything ok, proceed, otherwise, return error
	log.Info("Running preflight checks")
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
	log.Infof("Process configuration and store it at %s", config.GetClusterConfigFileName())
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

	err = config.WriteConfigFile(output, config.GetClusterConfigFileName())
	if err != nil {
		// Error writing config
		log.Errorf("error writing config (%s)", err.Error())
		return config.Response{
			Status:       config.Error,
			Description:  "Cluster can not be created",
			ErrorMessage: err.Error(),
			Error:        true,
		}
	}

	log.Infof("Create cluster with config at %s", config.GetClusterConfigFileName())
	// Create Cluster without preflight checks  --> $TCE create "${CLUSTER_NAME}" --skip-preflight -f "$HOME/cluster-config-$$.yaml"
	// TODO: See how we can stream output of the TCE process back or write it to a file
	//nolint:gosec
	cmd = exec.Command(config.UnmanagedClusterBinary, "create", "-v", "0", config.DefaultClusterName, "--skip-preflight", "-f", config.GetClusterConfigFileName())
	if err := cmd.Run(); err != nil {
		log.Errorf("Error while creating the cluster (%s)", err.Error())
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
	log.Info("Kubeconfig copied to host")

	copyConfigFiles()
	log.Info("Tanzu config files copied to host")

	return ret
}

func (c *Cluster) DeleteCluster() config.Response {
	log.Info("Delete Cluster")

	// LockCreationOrExitIfAlreadyCreating  --> touch $HOME/tce-cluster-delete.tag
	log.Info("Check that the cluster is not being deleted")
	lock, err := utils.GetFileLockWithTimeOut(utils.GetClusterDeleteLockFilename(), utils.DefaultLockTimeout)
	if err != nil {
		log.Info("Cluster already deleting")
		// Lock is already in place, which means that process is already running, just return
		return config.DeletingResponse()
	}

	// releaseLocks --> rm -f $HOME/tce-cluster-delete.tag $HOME/cluster-config-$$.yaml
	defer func() {
		if err := lock.Unlock(); err != nil {
			log.Errorf("cannot unlock %q, reason: %v", lock, err)
		}
	}()

	// Get cluster state. If already running, return already running, else
	log.Info("Check to see if there is an existing cluster")
	status, _ := checks.GetContainerClusterStatus()
	if status == checks.Running {
		log.Info("There's an existing cluster")
	} else {
		log.Info("There's no running cluster")
	}

	if status != checks.NotExist {
		log.Info("Deleting cluster")
		//nolint:gosec
		cmd := exec.Command(config.UnmanagedClusterBinary, "delete", "-v", "0", config.DefaultClusterName)
		if err := cmd.Run(); err != nil {
			log.Errorf("Error while deleting the cluster (%s)", err.Error())
			log.Info("Force delete non running container")
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

func (c *Cluster) ClusterStatus() config.Response {
	// Status can be:
	// NotExist
	// Creating
	// Running
	// Deleting
	// Error
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
		log.Info("Cluster creating")
		return config.Response{
			Status:       config.Creating,
			Description:  "Cluster creating",
			ErrorMessage: "",
			Error:        false,
		}
	}
	_, canbecreated, e2 := checks.IsClusterUpAndRunning()
	if canbecreated {
		log.Debug("Cluster does not exist")
		return config.Response{
			Status:       config.NotExists,
			Description:  "Cluster does not exist",
			ErrorMessage: "",
			Error:        false,
		}
	}
	if e2 != nil {
		// TODO: Check for status Stopped
		log.Errorf("Error while checking cluster status (%s)", e2.Error())
		return config.Response{
			Status:       config.Error,
			Description:  "Cluster error",
			ErrorMessage: e2.Error(),
			Error:        true,
		}
	}

	d, e3 := checks.IsClusterDeleting()
	if e3 != nil {
		log.Errorf("Error while checking cluster status (%s)", e3.Error())
		return config.Response{
			Status:       config.Error,
			Description:  "Cluster can not be created",
			ErrorMessage: e3.Error(),
			Error:        true,
		}
	}
	if d {
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

func (c *Cluster) Reset() config.Response {
	// Remove the kubeconfig context
	err := kubeconfig.RemoveNamedConfig(config.DefaultClusterName, config.DefaultHostMountedKubeConfig)
	if err != nil {
		log.Warnf("Error removing cluster kubeconfig (%s)", err.Error())
	}

	// $TCE delete ${CLUSTER_NAME} || true
	//nolint:gosec
	cmd := exec.Command("/backend/tanzu-unmanaged-cluster", "delete", config.DefaultClusterName)
	_ = cmd.Run() // We don't worry about errors

	// docker stop ${CLUSTER_NAME}-control-plane || true
	// docker rm ${CLUSTER_NAME}-control-plane || true
	err = docker.ForceStopAndDeleteCluster()
	if err != nil {
		log.Errorf("Error while stopping cluster (%s)", err.Error())
		return config.Response{
			Status:       config.Error,
			Description:  "Cluster can not be deleted",
			ErrorMessage: err.Error(),
			Error:        true,
		}
	}

	removeAllConfigFiles(true)
	os.RemoveAll(config.GetLogsFileName())

	return config.NotExistsResponse()
}

func copyConfigFiles() {
	// TODO: Copy the local tanzu files to the volume mount so they are on the host for diagnosis
	//nolint:gosec
	err := exec.Command("cp", "-rf", filepath.Join(config.GetUserHome(), ".config", "tanzu"), "/opt").Run()
	if err != nil {
		log.Errorf("Couldn't copy tanzu config files (%s)", err)
	}
}

func removeAllConfigFiles(isdelete bool) {
	log.Info("Removing all internal configuration files")
	err := os.RemoveAll(filepath.Join(config.GetUserHome(), ".kube"))
	if err != nil {
		log.Errorf("Couldn't remove %s (%s)", filepath.Join(config.GetUserHome(), ".kube"), err)
	}
	err = os.RemoveAll(filepath.Join(config.GetUserHome(), ".tanzu"))
	if err != nil {
		log.Errorf("Couldn't remove %s (%s)", filepath.Join(config.GetUserHome(), ".tanzu"), err)
	}
	err = os.RemoveAll(filepath.Join(config.GetUserHome(), ".config", "tanzu"))
	if err != nil {
		log.Errorf("Couldn't remove %s (%s)", filepath.Join(config.GetUserHome(), ".config", "tanzu"), err)
	}
	err = os.RemoveAll(config.GetClusterConfigFileName())
	if err != nil {
		log.Errorf("Couldn't remove %s (%s)", config.GetClusterConfigFileName(), err)
	}
	if isdelete {
		err = os.RemoveAll(config.GetLogsFileName())
		if err != nil {
			log.Errorf("Couldn't remove %s (%s)", config.GetLogsFileName(), err)
		}
	}
}

func (c *Cluster) Logs() config.Response {
	up, _, _ := checks.IsClusterUpAndRunning()
	if !up {
		return config.NotExistsResponse()
	}
	// cat $HOME/.kube/config
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

func (c *Cluster) Stats() config.Response {
	up, _, _ := checks.IsClusterUpAndRunning()
	if !up {
		return config.EmptyStatsResponse()
	}
	stats, _ := docker.GetDockerStats()
	return config.Response{
		Stats: stats,
	}
}

func (c *Cluster) GetKubeconfig() config.Response {
	up, _, _ := checks.IsClusterUpAndRunning()
	if !up {
		return config.NotExistsResponse()
	}
	// cat $HOME/.kube/config
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

func (c *Cluster) ProvisionIngress() config.Response {
	up, _, _ := checks.IsClusterUpAndRunning()
	if !up {
		return config.NotExistsResponse()
	}
	log.Infof("Process ingress configuration and store it at %s", config.GetClusterIngressConfigFileName())
	cmd := exec.Command("ytt", "-f", "/backend/apps/ingress.yaml", "-f", "/backend/apps/apps-values.yaml")
	output, err := cmd.Output()
	if err != nil {
		// Error while running ytt
		log.Errorf("Error procesing ingress config (%s)", err.Error())
		return config.Response{
			Status:       config.Error,
			Description:  "Ingress can not be created",
			ErrorMessage: err.Error(),
			Error:        true,
		}
	}
	err = config.WriteConfigFile(output, config.GetClusterIngressConfigFileName())
	if err != nil {
		log.Warnf("Error ingress config (%s)", err.Error())
	}

	log.Infof("Create ingress with config at %s", config.GetClusterIngressConfigFileName())
	// TODO: See how we can stream output of the TCE process back or write it to a file
	//nolint:gosec
	cmd = exec.Command("kubectl", "apply", "-f", config.GetClusterIngressConfigFileName())
	if err := cmd.Run(); err != nil {
		log.Errorf("Error while creating ingress (%s)", err.Error())
		return config.Response{
			Status:       config.Error,
			Description:  "Ingress can not be created",
			ErrorMessage: err.Error(),
			Error:        true,
		}
	}
	log.Info("Ingress successfully created")

	return config.Response{
		Status:      config.Running,
		Description: "Ingress successfully created",
	}
}

func (c *Cluster) ProvisionCertMan() config.Response {
	up, _, _ := checks.IsClusterUpAndRunning()
	if !up {
		return config.NotExistsResponse()
	}
	log.Infof("Process ingress configuration and store it at %s", config.GetClusterCertmanagerConfigFileName())
	cmd := exec.Command("ytt", "-f", "/backend/apps/certmanager.yaml", "-f", "/backend/apps/apps-values.yaml", "--ignore-unknown-comments")
	output, err := cmd.Output()
	if err != nil {
		// Error while running ytt
		log.Errorf("Error procesing cert-manager config (%s)", err.Error())
		return config.Response{
			Status:       config.Error,
			Description:  "cert-manager can not be created",
			ErrorMessage: err.Error(),
			Error:        true,
		}
	}
	err = config.WriteConfigFile(output, config.GetClusterCertmanagerConfigFileName())
	if err != nil {
		log.Warnf("Error cert-manager config (%s)", err.Error())
	}

	log.Infof("Create cert-manager with config at %s", config.GetClusterCertmanagerConfigFileName())
	// TODO: See how we can stream output of the TCE process back or write it to a file
	//nolint:gosec
	cmd = exec.Command("kubectl", "apply", "-f", config.GetClusterCertmanagerConfigFileName())
	if err := cmd.Run(); err != nil {
		log.Errorf("Error while creating cert-manager (%s)", err.Error())
		return config.Response{
			Status:       config.Error,
			Description:  "cert-manager can not be created",
			ErrorMessage: err.Error(),
			Error:        true,
		}
	}
	log.Info("cert-manager successfully created")

	return config.Response{
		Status:      config.Running,
		Description: "cert-manager successfully created",
	}
}

func (c *Cluster) GetJSONResponse(res *config.Response) string {
	byteArray, err := json.MarshalIndent(res, "", "  ")

	if err != nil {
		log.Println(err)
	}

	return (string(byteArray))
}
