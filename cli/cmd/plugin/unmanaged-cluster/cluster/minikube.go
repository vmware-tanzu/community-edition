// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cluster

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/config"
)

var minikubeBinName = "minikube"

type MinikubeClusterManager struct {
	// Provider configuration options
	driver           string
	containerRuntime string
	rawMinikubeArgs  string

	// Bootstrap logs to bubble up after bootstrapping
	logs []string
}

type MinikubeProfile struct {
	Name   string         `json:"Name"`
	Status string         `json:"Status"`
	Config MinikubeConfig `json:"Config"`
}

type MinikubeProfiles struct {
	Invalid []MinikubeProfile `json:"invalid"`
	Valid   []MinikubeProfile `json:"valid"`
}

type MinikubeConfig struct {
	Driver string `json:"Driver"`
}

// Get retrieves cluster information. An error is returned if no cluster is
// found or if there is a failure communicating with minikube.
func (mkcm *MinikubeClusterManager) Get(clusterName string) (*KubernetesCluster, error) {
	var resolvedStatus string

	mkp, err := mkcm.getProfile(clusterName)
	if err != nil {
		return nil, fmt.Errorf("failed to find a cluster profile with name %s via minikube", clusterName)
	}

	// "Running" and "Stopped" are known-good status from minikube. Any other value should be
	// considered unknown.
	switch mkp.Status {
	case "Running":
		resolvedStatus = StatusRunning
	case "Stopped":
		resolvedStatus = StatusStopped
	default:
		resolvedStatus = StatusUnknown
	}

	kc := &KubernetesCluster{
		Name: clusterName,
		// TODO(joshrosso): We should consider this field in future
		// work. Perhaps when we expose get at the CLI level, we could
		// do a command like `tanzu uc get ${CLUSTER_NAME} --kubeconfig` and return this value.
		Kubeconfig: []byte{},
		Status:     resolvedStatus,
		Driver:     mkp.Config.Driver,
	}

	return kc, nil
}

// Create creates a minikube cluster with a given profile name
//nolint:funlen,gocritic
func (mkcm *MinikubeClusterManager) Create(c *config.UnmanagedClusterConfig) (*KubernetesCluster, error) {
	err := mkcm.setProvierConfigs(c.ProviderConfiguration)
	if err != nil {
		return nil, err
	}

	profile, err := mkcm.getProfile(c.ClusterName)
	if err != nil {
		return nil, fmt.Errorf("could not get minikube profiles. Error: %s", err.Error())
	}

	if profile != nil {
		return nil, fmt.Errorf("minikube profile named: %s already exists with status: %s. This minikube profile should be deleted before proceeding", c.ClusterName, profile.Name)
	}

	// Base start command arguments
	args := []string{
		"start",
		"--driver",
		mkcm.driver,
		"--profile",
		c.ClusterName,
	}

	args = append(
		args,
		"--log_file="+c.LogFile,
	)

	// Set the node image
	args = append(
		args,
		"--base-image="+c.NodeImage,
	)

	// Configure the container container-runtime if provided
	if mkcm.containerRuntime != "" {
		args = append(
			args,
			"--container-runtime="+mkcm.containerRuntime,
		)
	}

	// Configure number of nodes to deploy (including 1st control plane)
	numNodes, err := mkcm.calculatNodes(c)
	if err != nil {
		return nil, fmt.Errorf("could not add worker nodes. Error: %s", err.Error())
	}

	args = append(
		args,
		"--nodes="+numNodes,
	)

	// Configure port mappings
	args = append(
		args,
		mkcm.createPortMappingArgs(c)...,
	)

	// Configure the pod CIDR
	args = append(
		args,
		"--extra-config",
		"kubeadm.pod-network-cidr="+c.PodCidr,
	)

	// Configure the service CIDR
	args = append(
		args,
		"--service-cluster-ip-range",
		c.ServiceCidr,
	)

	// Add any additional user raw minikube args
	if mkcm.rawMinikubeArgs != "" {
		args = append(
			args,
			mkcm.rawMinikubeArgs,
		)
	}

	cmd := exec.Command(minikubeBinName, args...)

	// Set the kubeconfig to land in the configured kubeconfig path
	cmd.Env = append(
		os.Environ(),
		"KUBECONFIG="+c.KubeconfigPath,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error creating minikube cluster. Error: %s. Output: %s", err.Error(), out)
	}

	// No-op: minikube does not support multi-control-plane clusters
	err = mkcm.addControlPlanes(c)
	if err != nil {
		return nil, fmt.Errorf("could not add control plane nodes. Error: %s", err.Error())
	}

	// read in kubeconfig from unmanaged-cluster provider directory
	kcBytes, err := os.ReadFile(c.KubeconfigPath)
	if err != nil {
		return nil, err
	}

	kc := &KubernetesCluster{
		Name:       c.ClusterName,
		Kubeconfig: kcBytes,
	}

	return kc, nil
}

// Delete will delete the minikube cluster given a named profile
func (mkcm *MinikubeClusterManager) Delete(c *config.UnmanagedClusterConfig) error {
	profile, err := mkcm.getProfile(c.ClusterName)
	if err != nil {
		return fmt.Errorf("could not get minikube profiles. Error: %s", err.Error())
	}

	if profile == nil {
		return fmt.Errorf("minikube profile named: %s not found", c.ClusterName)
	}

	args := []string{
		"delete",
		"--profile",
		c.ClusterName,
	}

	cmd := exec.Command(minikubeBinName, args...)

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Prepare gets the minikube provider ready
// Nothing to do here for minikube
func (mkcm *MinikubeClusterManager) Prepare(c *config.UnmanagedClusterConfig) error {
	return nil
}

// PreflightCheck ensures minikube is on the system and runable
func (mkcm *MinikubeClusterManager) PreflightCheck() ([]string, []error) {
	bin, err := exec.LookPath(minikubeBinName)
	if err != nil {
		return []string{}, []error{fmt.Errorf("problem getting minikube binary. Is it in your PATH?")}
	}

	cmd := exec.Command(bin, "version")
	_, err = cmd.Output()
	if err != nil {
		return []string{}, []error{fmt.Errorf("could not run minikube binary. Is it configured?")}
	}

	return nil, nil
}

// PreProviderNotify returns a usage notification about minikube
func (mkcm *MinikubeClusterManager) PreProviderNotify() []string {
	return []string{
		"Cluster creation using Minikube!",
		"Warning: the minikube provider is experimental!",
		"❤️  Checkout this awesome project at https://minikube.sigs.k8s.io",
	}
}

// PostProviderNotify returns the aggregate logs from the minikube bootstrapping
func (mkcm *MinikubeClusterManager) PostProviderNotify() []string {
	return mkcm.logs
}

// Start attempts to start a stopped minikube cluster. An error is returned when:
// 1. The cluster is already running.
// 2. There are issues communicating with minikube.
// 3. The cluster fails to start.
func (mkcm *MinikubeClusterManager) Start(c *config.UnmanagedClusterConfig) error {
	// verify cluster is in a "Stopped" state before attempting to start.
	kc, err := mkcm.Get(c.ClusterName)
	if err != nil {
		return fmt.Errorf("cannot start this cluster. Error occurred retrieving status: %s", err.Error())
	}
	if kc.Status != StatusStopped {
		return fmt.Errorf("cannot start this cluster. The status must be %s, it was %s", StatusStopped, kc.Status)
	}

	args := []string{
		"start",
		"--profile",
		kc.Name,
		// specifying the driver is key as the user may have a global driver config with minikube
		// (via minikube config) and if their global setting is different from the driver used
		// to start this cluser, creation will fail. Since we lookup the cluster details and know
		// the driver, we always specify it when running Start.
		"--driver",
		kc.Driver,
	}
	cmd := exec.Command(minikubeBinName, args...)

	// TODO(joshrosso): starting a cluster can take 1+ minute(s). Could be worth finding a way
	// to propagate details to the client.
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to start cluster via minikube. Error: %s", err)
	}

	return nil
}

// Stop takes a running minikube cluster and stops the host(s).
func (mkcm *MinikubeClusterManager) Stop(c *config.UnmanagedClusterConfig) error {
	// verify cluster is in a "Running" state before attempting to stop.
	kc, err := mkcm.Get(c.ClusterName)
	if err != nil {
		return fmt.Errorf("cannot stop this cluster. Error occurred retrieving status: %s", err.Error())
	}
	if kc.Status != StatusRunning {
		return fmt.Errorf("cannot stop this cluster. The status must be %s, it was %s", StatusRunning, kc.Status)
	}

	args := []string{
		"stop",
		"--profile",
		kc.Name,
	}
	cmd := exec.Command(minikubeBinName, args...)

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to stop cluster via minikube. Error: %s", err)
	}

	return nil
}

// getProfile attempts to retrieve a minikube profile by name
func (mkcm *MinikubeClusterManager) getProfile(name string) (*MinikubeProfile, error) {
	args := []string{
		"profile",
		"list",
		"--output=json",
	}

	cmd := exec.Command(minikubeBinName, args...)

	minikubeJSONOut, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	mkProfiles := &MinikubeProfiles{}

	err = json.Unmarshal(minikubeJSONOut, mkProfiles)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal minikube profile list output: %s", err.Error())
	}

	for _, profile := range mkProfiles.Valid {
		if profile.Name == name {
			return &profile, nil
		}
	}

	return nil, nil
}

// calculatNodes will return a string to be used in a minikube `start` command
// for the `--nodes` flag.
// This represents the 1 control-planes (1 by default until multi-control-plane clusters are supported)
// and the number of worker nodes specified (which may be 0)
func (mkcm *MinikubeClusterManager) calculatNodes(c *config.UnmanagedClusterConfig) (string, error) {
	wnc, err := strconv.Atoi(c.WorkerNodeCount)
	if err != nil {
		return "", err
	}

	// Minikube deploys a control-plane first, and than subsequent worker nodes
	// This includes the 1 control plane and the workers
	wnc++

	return strconv.Itoa(wnc), nil
}

// Adding multiple control planes is not supported in minikube,
// although the flag is there, it will error out with x509 cert errors.
// This will be resolved in a future release
// and is apart of a broader effort to allow for multi-node clusters in minikube.
//
// So therefore, the majority of this function is not utilized but remains for documentation purposes.
//
// Reference:
// - https://github.com/kubernetes/minikube/issues/7461
// - https://github.com/kubernetes/minikube/issues/7538
func (mkcm *MinikubeClusterManager) addControlPlanes(c *config.UnmanagedClusterConfig) error {
	cpnc, err := strconv.Atoi(c.ControlPlaneNodeCount)
	if err != nil {
		return err
	}

	if cpnc > 1 {
		// Always exit early since multiple control planes are not supported and warn user
		mkcm.logs = append(mkcm.logs, "multiple control-plane nodes are not supported in Minikube. Skipping creating additional control plane nodes")
		return nil
	}

	// Add more control planes past 1.
	// Since the default node is a control plane, we need only add more than 1
	for i := 1; i < cpnc; i++ {
		args := []string{
			"node",
			"add",
			"--control-plane",
			"--profile",
			c.ClusterName,
		}

		cmd := exec.Command(minikubeBinName, args...)

		cmd.Env = append(
			os.Environ(),
			"KUBECONFIG="+c.KubeconfigPath,
		)

		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("could not create additional control plane node. Error: %s", err.Error())
		}
	}

	return nil
}

// Builds port mapping string for use in `start` command
func (mkcm *MinikubeClusterManager) createPortMappingArgs(c *config.UnmanagedClusterConfig) []string {
	portArgs := []string{}

	for _, portToForward := range c.PortsToForward {
		portMapping := strings.Builder{}
		if portToForward.ContainerPort != 0 {
			portMapping.WriteString(strconv.Itoa(portToForward.ContainerPort))
		}
		if portToForward.HostPort != 0 {
			portMapping.WriteString(":")
			portMapping.WriteString(strconv.Itoa(portToForward.HostPort))
		}
		if portToForward.Protocol != "" {
			portMapping.WriteString("/")
			portMapping.WriteString(portToForward.Protocol)
		}

		portArgs = append(
			portArgs,
			"--ports",
			portMapping.String(),
		)
	}

	return portArgs
}

// Sets the different supported provider configurations
// in the MinikubeClusterManager struct
func (mkcm *MinikubeClusterManager) setProvierConfigs(pc map[string]interface{}) error {
	if _, ok := pc["driver"]; ok {
		if _, ok := pc["driver"].(string); !ok {
			return fmt.Errorf("ProviderConfiguration.driver wrong type, expected string")
		}

		mkcm.driver = pc["driver"].(string)
	}

	if _, ok := pc["containerRuntime"]; ok {
		if _, ok := pc["containerRuntime"].(string); !ok {
			return fmt.Errorf("ProviderConfiguration.containerRuntime wrong type, expected string")
		}

		mkcm.driver = pc["container-runtime"].(string)
	}

	if _, ok := pc["rawMinikubeArgs"]; ok {
		if _, ok := pc["rawMinikubeArgs"].(string); !ok {
			return fmt.Errorf("ProviderConfiguration.rawMinikubeArgs wrong type, expected string")
		}

		mkcm.driver = pc["containerRuntime"].(string)
	}

	return nil
}
