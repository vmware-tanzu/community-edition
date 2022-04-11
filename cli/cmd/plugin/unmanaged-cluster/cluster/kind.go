// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cluster

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
	kindconfig "sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
	kindcluster "sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cluster/nodes"
	"sigs.k8s.io/kind/pkg/exec"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/config"
)

const (
	minMemoryBytes         = 2147483648
	minCPUCount            = 1
	kindConfigFileName     = "kindconfig.yaml"
	KindTypedataCluster    = "Cluster"
	KindTypedataAPIVersion = "kind.x-k8s.io/v1alpha4"
)

// TODO(stmcginnis): Keeping this here for now for reference, remove once we're
// ready with custom configurations.
// const defaultKindConfig = `kind: Cluster
// apiVersion: kind.x-k8s.io/v1alpha4
// nodes:
// - role: control-plane
//   #! port forward 80 on the host to 80 on this node
//   extraPortMappings:
//   - containerPort: 80
//     #!hostPort: 80
//     #! optional: set the bind address on the host
//     #! 0.0.0.0 is the current default
//     listenAddress: "127.0.0.1"
//     #! optional: set the protocol to one of TCP, UDP, SCTP.
//     #! TCP is the default
//     protocol: TCP
// networking:
//   disableDefaultCNI: true`

// KindClusterManager is a ClusterManager implementation for working with
// Kind clusters.
type KindClusterManager struct {
}

// Create will create a new kind cluster or return an error.
func (kcm *KindClusterManager) Create(c *config.UnmanagedClusterConfig) (*KubernetesCluster, error) {
	var err error

	kindProvider := kindcluster.NewProvider()
	clusterConfig := kindcluster.CreateWithKubeconfigPath(c.KubeconfigPath)

	// Serlize unstructured data into a kindconfig struct
	// Return any error from attempting to read the data
	kcFromProviderConfig, err := serializeKindProviderConfig(c.ProviderConfiguration)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize kind config from given ProviderConfiguration. Error was: %s", err)
	}

	// generate a kindconfig struct from defaults and flag values
	kcFromClusterConfig, err := kindConfigFromClusterConfig(c)
	if err != nil {
		return nil, fmt.Errorf("failed to generate a viable kind config. Error was: %s", err)
	}

	// Merge the kind configs from providerconfig into the clusterconfig
	err = mergeConfigsLeft(kcFromClusterConfig, kcFromProviderConfig)
	if err != nil {
		return nil, err
	}

	// store our kind config on the filesystem for users to inspect if needed
	kcMergedBytes, err := kindClusterToBytes(kcFromClusterConfig)
	if err != nil {
		return nil, err
	}

	err = writeKindConfigFile(kcMergedBytes, c.ClusterName)
	if err != nil {
		return nil, err
	}

	kindConfig := kindcluster.CreateWithV1Alpha4Config(kcFromClusterConfig)
	err = kindProvider.Create(c.ClusterName, clusterConfig, kindConfig)
	if err != nil {
		return nil, fmt.Errorf("kind returned error: %s", err)
	}

	// readkubeconfig in bytes
	kcBytes, err := os.ReadFile(c.KubeconfigPath)
	if err != nil {
		return nil, err
	}

	kc := &KubernetesCluster{
		Name:       c.ClusterName,
		Kubeconfig: kcBytes,
	}

	if strings.Contains(c.Cni, "antrea") {
		kindNodes, _ := kindProvider.ListNodes(c.ClusterName)
		for _, n := range kindNodes {
			if err := patchForAntrea(n.String()); err != nil { //nolint:staticcheck
				// TODO(stmcginnis): We probably don't want to just error out
				// since the cluster has already been created, but we should
				// at least report a warning back to the user that part of the
				// setup failed.
			}
		}
	}

	return kc, nil
}

func serializeKindProviderConfig(pc map[string]interface{}) (*kindconfig.Cluster, error) {
	// Check if key exists. If not, return empty config and continue
	if _, ok := pc["rawKindConfig"]; !ok {
		return &kindconfig.Cluster{}, nil
	}

	// Check if provided data is a string.
	if _, ok := pc["rawKindConfig"].(string); !ok {
		return &kindconfig.Cluster{}, fmt.Errorf("ProviderConfiguration.rawKindConfig wrong type, expected string")
	}

	kc := &kindconfig.Cluster{}
	err := yaml.Unmarshal([]byte(pc["rawKindConfig"].(string)), kc)
	if err != nil {
		return &kindconfig.Cluster{}, fmt.Errorf("ProviderConfiguration.rawKindConfig unable to be unmarsheled. Error: %s", err.Error())
	}

	return kc, nil
}

func kindConfigFromClusterConfig(c *config.UnmanagedClusterConfig) (*kindconfig.Cluster, error) {
	// Load the defaults
	kindConfig := &kindconfig.Cluster{}
	kindConfig.Kind = KindTypedataCluster
	kindConfig.APIVersion = KindTypedataAPIVersion
	kindConfig.Name = c.ClusterName
	kindNodes, err := setNumberOfNodes(c)
	if err != nil {
		return nil, err
	}
	kindConfig.Nodes = kindNodes
	kindconfig.SetDefaultsCluster(kindConfig)

	// Now populate or override with the specified configuration
	kindConfig.Networking.DisableDefaultCNI = true
	if c.PodCidr != "" {
		kindConfig.Networking.PodSubnet = c.PodCidr
	}
	if c.ServiceCidr != "" {
		kindConfig.Networking.ServiceSubnet = c.ServiceCidr
	}

	// Apply the node image to all nodes
	for i := range kindConfig.Nodes {
		kindConfig.Nodes[i].Image = c.NodeImage
	}

	// Do the port mapping for the first node (which should by default be the control plane)
	// If users want a more granular way to apply port mappings, they should use the rawKindConfig
	for _, portToForward := range c.PortsToForward {
		portMapping := kindconfig.PortMapping{}
		if portToForward.ListenAddress != "" {
			portMapping.ListenAddress = portToForward.ListenAddress
		}
		if portToForward.ContainerPort != 0 {
			portMapping.ContainerPort = int32(portToForward.ContainerPort)
		}
		if portToForward.HostPort != 0 {
			portMapping.HostPort = int32(portToForward.HostPort)
		}
		if portToForward.Protocol != "" {
			portMapping.Protocol = kindconfig.PortMappingProtocol(portToForward.Protocol)
		}

		kindConfig.Nodes[0].ExtraPortMappings = append(kindConfig.Nodes[0].ExtraPortMappings, portMapping)
	}

	return kindConfig, nil
}

func kindClusterToBytes(kc *kindconfig.Cluster) ([]byte, error) {
	var rawConfig bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&rawConfig)

	err := yamlEncoder.Encode(kc)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Kind configuration. Error: %s", err.Error())
	}
	err = yamlEncoder.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Kind configuration. Error: %s", err.Error())
	}

	return rawConfig.Bytes(), nil
}

// mergeConfigsLeft merges the second kindconfig into the first kindconfig,
// mutating the first kindconfig.
//
// In the unmanaged-cluster workflow,
// the first argument should be the kindconfig created from CLI flags and arguments.
// This is referred to as the "global" config.
// The second argument should be the kindconfig created via the provider config file.
// This is referred to as the "provider" config.
//
// CLI flags and arguments from the unmanaged cluster config take precedence.
// So options given via `rawKindConfig` take lower order precedence
// and values given via `--flags` and env vars are the highest order.
// Default values take the lowest order precedence and are used when no values are found.
//
// Generally, merges happen from the provider config into global config when:
// - a value for a field is missing OR the default is found for the the global config
// - a value is present in the provider config
// Exceptions are noted via inline comments
//nolint:funlen,gocyclo
func mergeConfigsLeft(l, r *kindconfig.Cluster) error {
	// build the default configuration to compare with the global config to check for default values
	defaultConf, err := kindConfigFromClusterConfig(config.GenerateDefaultConfig())
	if err != nil {
		return err
	}

	// Uses metadata from provider config when the global config value is the default or missing
	if l.Kind == defaultConf.Kind && r.Kind != "" {
		l.Kind = r.Kind
	}

	if l.APIVersion == defaultConf.APIVersion && r.APIVersion != "" {
		l.APIVersion = r.APIVersion
	}

	if l.Name == defaultConf.Name && r.Name != "" {
		l.Name = r.Name
	}

	// Nodes and their fields are merged 1 to 1 via their index
	// using one of the following as a base:
	// - the global config nodes as the base if number of nodes were defined via global config
	// - the provider config nodes as the base if global config nodes are the default
	//   and provider config defined nodes
	// - the default node if neither were configured
	// Extra nodes in either circumstance are ignored
	if len(l.Nodes) > len(defaultConf.Nodes) || len(l.Nodes) == len(r.Nodes) {
		nodes := []kindconfig.Node{}

		for i, n := range l.Nodes {
			role := n.Role
			image := n.Image

			// Take the role from provider config node if defined/available and global config is default
			if n.Role == defaultConf.Nodes[0].Role && len(r.Nodes) > i {
				if r.Nodes[i].Role != "" {
					role = r.Nodes[i].Role
				}
			}

			// Take the image from provider config node if defined/available and global config is default
			if n.Image == defaultConf.Nodes[0].Image && len(r.Nodes) > i {
				if r.Nodes[i].Image != "" {
					image = r.Nodes[i].Image
				}
			}

			// Rebuild port mappings using the current node's extra port mapping as a base
			// and pulling any configured port mappings from the corresponding provider config node
			extraPortMappings := n.ExtraPortMappings
			if len(r.Nodes) > i {
				for _, pm := range r.Nodes[i].ExtraPortMappings {
					extraPortMappings = append(extraPortMappings, kindconfig.PortMapping{
						ContainerPort: pm.ContainerPort,
						HostPort:      pm.HostPort,
						ListenAddress: pm.ListenAddress,
						Protocol:      pm.Protocol,
					})
				}
			}

			nodes = append(nodes, kindconfig.Node{
				Role:              role,
				Image:             image,
				ExtraPortMappings: extraPortMappings,
			})
		}

		l.Nodes = nodes
	} else if len(l.Nodes) == len(defaultConf.Nodes) && len(r.Nodes) > 0 {
		// if the global nodes are the defaults
		// and there were some provider nodes given
		// than use the provider nodes.
		// BUT - still try to merge any options from the global config node

		nodes := []kindconfig.Node{}

		for i, n := range r.Nodes {
			role := n.Role
			image := n.Image

			// Use the global config role if provided or the base, default role if not declared in provider config
			if l.Nodes[0].Role != defaultConf.Nodes[0].Role || n.Role == "" {
				role = l.Nodes[0].Role
			}

			// Use the global config image if provided or the base, default image if not declared in provider config
			if l.Nodes[0].Image != defaultConf.Nodes[0].Image || n.Image == "" {
				image = l.Nodes[0].Image
			}

			// Rebuild port mappings using the current node's extra port mapping as a base
			// and pulling any configured port mappings from the corresponding global config node
			extraPortMappings := n.ExtraPortMappings
			if len(l.Nodes) > i {
				for _, pm := range l.Nodes[i].ExtraPortMappings {
					extraPortMappings = append(extraPortMappings, kindconfig.PortMapping{
						ContainerPort: pm.ContainerPort,
						HostPort:      pm.HostPort,
						ListenAddress: pm.ListenAddress,
						Protocol:      pm.Protocol,
					})
				}
			}

			nodes = append(nodes, kindconfig.Node{
				Role:              role,
				Image:             image,
				ExtraPortMappings: extraPortMappings,
			})
		}

		l.Nodes = nodes
	}

	// Configure network settings
	if l.Networking.IPFamily == defaultConf.Networking.IPFamily && r.Networking.IPFamily != "" {
		l.Networking.IPFamily = r.Networking.IPFamily
	}

	if l.Networking.APIServerPort == defaultConf.Networking.APIServerPort && r.Networking.APIServerPort != 0 {
		l.Networking.APIServerPort = r.Networking.APIServerPort
	}

	if l.Networking.APIServerAddress == defaultConf.Networking.APIServerAddress && r.Networking.APIServerAddress != "" {
		l.Networking.APIServerAddress = r.Networking.APIServerAddress
	}

	if l.Networking.PodSubnet == defaultConf.Networking.PodSubnet && r.Networking.PodSubnet != "" {
		l.Networking.PodSubnet = r.Networking.PodSubnet
	}

	if l.Networking.ServiceSubnet == defaultConf.Networking.ServiceSubnet && r.Networking.ServiceSubnet != "" {
		l.Networking.ServiceSubnet = r.Networking.ServiceSubnet
	}

	if l.Networking.KubeProxyMode == defaultConf.Networking.KubeProxyMode && r.Networking.KubeProxyMode != "" {
		l.Networking.KubeProxyMode = r.Networking.KubeProxyMode
	}

	// The following are not configurable via the global config
	// so we can take the provider config whole sale
	// By default, these values are empty
	l.FeatureGates = r.FeatureGates
	l.RuntimeConfig = r.RuntimeConfig
	l.KubeadmConfigPatches = r.KubeadmConfigPatches
	l.KubeadmConfigPatchesJSON6902 = r.KubeadmConfigPatchesJSON6902
	l.ContainerdConfigPatches = r.ContainerdConfigPatches
	l.ContainerdConfigPatchesJSON6902 = r.ContainerdConfigPatchesJSON6902

	return nil
}

func writeKindConfigFile(configBytes []byte, clusterName string) error {
	configDir, err := config.GetUnmanagedConfigPath()
	if err != nil {
		return err
	}

	configFp := filepath.Join(configDir, clusterName, kindConfigFileName)
	err = os.WriteFile(configFp, configBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write kind config file. Error: %s", err.Error())
	}

	return nil
}

func setNumberOfNodes(c *config.UnmanagedClusterConfig) ([]kindconfig.Node, error) {
	// Get and check control plane count from config
	cpnc, err := strconv.Atoi(c.ControlPlaneNodeCount)
	if err != nil {
		return nil, err
	}

	if cpnc < 1 {
		return nil, fmt.Errorf("cannot have less than 1 control plane node")
	}

	// Get and check worker count from config
	wnc, err := strconv.Atoi(c.WorkerNodeCount)
	if err != nil {
		return nil, err
	}

	if wnc < 0 {
		return nil, fmt.Errorf("cannot have less than 0 worker nodes")
	}

	// TODO (jpmcb): in single node clusters, control-plane nodes act as worker
	// nodes, in which case the taint will be removed.
	//
	// But, in the case where there are no worker nodes but multiple control plane nodes,
	// the taint is _not_ automatically removed so workloads will not be scheduled to the control plane nodes.
	//
	// In the future, the taint should be removed in order to support workloads being scheduled
	// on control-plane nodes without worker nodes.
	// https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#control-plane-node-isolation
	if cpnc > 1 && wnc == 0 {
		return nil, fmt.Errorf("multiple control plane nodes require at least one worker node for workload placement")
	}

	kindNodes := []kindconfig.Node{}

	for i := 1; i <= cpnc; i++ {
		n := kindconfig.Node{
			Role: kindconfig.ControlPlaneRole,
		}

		kindNodes = append(kindNodes, n)
	}

	for i := 1; i <= wnc; i++ {
		n := kindconfig.Node{
			Role: kindconfig.WorkerRole,
		}

		kindNodes = append(kindNodes, n)
	}

	return kindNodes, nil
}

// Get retrieves cluster information. An error is returned if no cluster is
// found or if there is a failure communicating with kind/docker.
func (kcm *KindClusterManager) Get(clusterName string) (*KubernetesCluster, error) {
	var resolvedStatus string

	// use kind APIs to get node name
	provider := kindcluster.NewProvider()
	kindNodes, err := provider.ListNodes(clusterName)
	if err != nil {
		return nil, err
	}

	// no cluster corresponding to name found
	if len(kindNodes) < 1 {
		return nil, fmt.Errorf("cluster %s could not be found by kind", clusterName)
	}

	// get the JSON-representation of the control-plane node and serialize it into a map
	cmdFindStatus := exec.Command("docker",
		"container",
		"ls",
		"-a",
		"--filter",
		// name filter prepends names with ^/ and appends with $ since name filtering is a fuzzy
		// search by default.
		fmt.Sprintf("name=^/%s$", kindNodes[0]),
		"--format",
		"{{json .}}")

	cmdFindStatusOutput, err := exec.Output(cmdFindStatus)
	if err != nil {
		return nil, err
	}

	// serialize into a map. No need to maintain a struct to serialize over time.
	// we are only interested in a single key. If it isn't found, we return status of
	// Unknown.
	var container map[string]interface{}
	err = json.Unmarshal(cmdFindStatusOutput, &container)
	if err != nil {
		return nil, fmt.Errorf("data returned from kind/docker could not be parsed as valid JSON. Error: %s", err)
	}

	if _, ok := container["State"]; !ok {
		return nil, fmt.Errorf("docker returned no status field for container")
	}

	// "running" and "existing" are known-good status from docker. Any other value should be
	// considered unknown.
	switch container["State"] {
	case "running":
		resolvedStatus = StatusRunning
	case "exited":
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
	}

	return kc, nil
}

// Delete removes a kind cluster.
func (kcm *KindClusterManager) Delete(c *config.UnmanagedClusterConfig) error {
	provider := kindcluster.NewProvider()
	return provider.Delete(c.ClusterName, "")
}

// Stop takes a running kind cluster and stops the host.
func (kcm *KindClusterManager) Stop(c *config.UnmanagedClusterConfig) error {
	// verify cluster is in a "Running" state before attempting to stop.
	kc, err := kcm.Get(c.ClusterName)
	if err != nil {
		return fmt.Errorf("cannot stop this cluster. Error occurred retrieving status: %s", err.Error())
	}
	if kc.Status != StatusRunning {
		return fmt.Errorf("cannot stop this cluster. The status must be %s, it was %s", StatusRunning, kc.Status)
	}

	kindNode, err := resolveKindNodesForSingleNodeCluster(c.ClusterName)
	if err != nil {
		return err
	}

	id, err := retrieveContainerIDFromName(kindNode.String())
	if err != nil {
		return err
	}

	// note: originally tried "docker stop" but found clusters could not recover.
	//       perhaps due to an improper signal being sent to them?
	stopCmd := exec.Command("docker", "kill", id)
	_, err = exec.Output(stopCmd)
	if err != nil {
		return err
	}

	return nil
}

// Start attempts to start a kind cluster. It returns an error when:
// 1. The cluster is already running.
// 2. There are issues communicating with kind/docker.
// 3. The cluster fails to start.
func (kcm *KindClusterManager) Start(c *config.UnmanagedClusterConfig) error {
	// verify cluster is in a "Stopped" state before attempting to start.
	kc, err := kcm.Get(c.ClusterName)
	if err != nil {
		return fmt.Errorf("cannot start this cluster. Error occurred retrieving status: %s", err.Error())
	}
	if kc.Status != StatusStopped {
		return fmt.Errorf("cannot start this cluster. The status must be %s, it was %s", StatusStopped, kc.Status)
	}

	kindNode, err := resolveKindNodesForSingleNodeCluster(c.ClusterName)
	if err != nil {
		return err
	}

	id, err := retrieveContainerIDFromName(kindNode.String())
	if err != nil {
		return err
	}

	startCmd := exec.Command("docker", "start", id)
	_, err = exec.Output(startCmd)
	if err != nil {
		return fmt.Errorf("failed to start cluster via kind. Error was: %s", err)
	}

	return nil
}

// Prepare will fetch a container image to the cluster host.
func (kcm *KindClusterManager) Prepare(c *config.UnmanagedClusterConfig) error {
	cmd := exec.Command("docker", "pull", c.NodeImage)
	_, err := exec.Output(cmd)
	if err != nil {
		return err
	}
	return nil
}

// PreflightCheck performs any pre-checks that can find issues up front that
// would cause problems for cluster creation.
func (kcm *KindClusterManager) PreflightCheck() ([]string, []error) {
	// Check presence of docker
	cmd := exec.Command("docker", "ps")
	if err := cmd.Run(); err != nil {
		// In this case we can't check the rest of the settings, so just return
		// the one error.
		return []string{}, []error{fmt.Errorf("docker is not installed or not reachable. Verify it's installed, running, and your user has permissions to interact with it. Error when attempting to run docker ps: %w", err)}
	}

	// Get Docker info
	cmd = exec.Command("docker", "info", "--format", "{{ json . }}")
	output, err := exec.Output(cmd)
	if err != nil {
		return []string{}, []error{fmt.Errorf("unable to get docker info: %w", err)}
	}

	return validateDockerInfo(output)
}

// PreProviderNotify returns the kind provider notification used during cluster bootstrapping
func (kcm *KindClusterManager) PreProviderNotify() []string {
	return []string{
		"Cluster creation using kind!",
		"❤️  Checkout this awesome project at https://kind.sigs.k8s.io",
	}
}

// PostProviderNotify returns the kind provider logs/notifications after bootstrapping
// Noop - nothing to return after bootstrapping
func (kcm *KindClusterManager) PostProviderNotify() []string {
	return []string{}
}

type dockerInfo struct {
	CPUs         int    `json:"NCPU"`
	Memory       int64  `json:"MemTotal"`
	Architecture string `json:"Architecture"`
}

func validateDockerInfo(output []byte) ([]string, []error) {
	info := dockerInfo{}
	if err := json.Unmarshal(output, &info); err != nil {
		// Nothing else we can check, just return this error right away
		return nil, []error{errors.New("unable to parse Docker information")}
	}

	warnings := []string{}
	issues := []error{}

	if !strings.HasSuffix(info.Architecture, "x86_64") {
		// Only amd64 supported right now, arm is experimental. Anything else is not supported.
		if strings.HasSuffix(info.Architecture, "aarch64") {
			warnings = append(warnings, "Arm64 architecture detected. Support is currently experimental. Some packages may not install due to their arm64 image not being available. You can find a list of package that have arm support in the release notes at https://github.com/vmware-tanzu/community-edition/releases/tag/v0.11.0.")
		} else {
			return []string{}, []error{errors.New("only amd64 and arm64 (experimental) architectures are currently supported")}
		}
	}

	if info.CPUs < minCPUCount {
		// Should only hit this if there is an issue getting the docker info
		// correctly, but we can also raise this if we find the need
		issues = append(issues, fmt.Errorf("minimum %d CPU core is required", minCPUCount))
	}

	if info.Memory < minMemoryBytes {
		issues = append(issues, fmt.Errorf("minimum %d GiB of memory is required", (minMemoryBytes/1024/1024/1024)))
	}

	return warnings, issues
}

// resolveKindNodesForSingleNodeCluster resolves the node that makes up a single-node kind cluster.
// This helper-function supports the start/stop functionality of this kind provider. It returns an
// error when:
// a. There is a failure to communicate with kind or docker
// b. It finds no nodes associated with the cluster
// c. It find more than 1 nodes associated with the cluster
//
// c is required as, at the time of writing, kind does not supporting starting/stopping a
// multi-node cluster.
func resolveKindNodesForSingleNodeCluster(clusterName string) (nodes.Node, error) {
	// use kind APIs to get node name
	provider := kindcluster.NewProvider()
	kindNodes, err := provider.ListNodes(clusterName)
	if err != nil {
		return nil, err
	}

	// if kind does not find any nodes associated with the cluster name, fail and return an error.
	if len(kindNodes) < 1 {
		return nil, fmt.Errorf("kind failed to find nodes associated with the cluster %s", clusterName)
	}

	// kind does not support starting/stopping of multi-node clusters. If the cluster contains
	// more than one node, do attempt to stop the cluster and return the error to the user.
	if len(kindNodes) > 1 {
		return nil, fmt.Errorf("cannot stop cluster. Kind does not support stopping and starting multi-node clusters")
	}

	return kindNodes[0], nil
}

// retrieveContainerIDFromName returns a container's ID based on the name provided. It uses docker
// to retrieve the ID. If there are issues communicating with docker, an error is returned.
func retrieveContainerIDFromName(name string) (string, error) {
	// get the json-representation of the control-plane node and serialize it into a map
	cmdFindID := exec.Command("docker",
		"container",
		"ls",
		"-a",
		"--filter",
		// name filter prepends names with ^/ and appends with $ since name filtering is a fuzzy
		// search by default.
		fmt.Sprintf("name=^/%s$", name),
		"--format",
		"{{json .}}")

	findIDOutput, err := exec.Output(cmdFindID)
	if err != nil {
		return "", err
	}

	var container map[string]interface{}
	err = json.Unmarshal(findIDOutput, &container)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve valid JSON from docker when looking up containerId for %s. Error: %s", name, err)
	}

	// If no ID field is present on the container, return an error.
	if _, ok := container["ID"]; !ok {
		return "", fmt.Errorf("found no ID associated with container")
	}

	return container["ID"].(string), nil
}

// patchForAntrea modifies the node network settings to allow local routing.
// this needs to happen for antrea running on kind or else you'll lose network connectivity
// see: https://github.com/antrea-io/antrea/blob/main/hack/kind-fix-networking.sh
func patchForAntrea(nodeName string) error {
	// First need to get the ID of the interface from the cluster node.
	cmd := exec.Command("docker", "exec", nodeName, "ip", "link")
	out, err := exec.Output(cmd)
	if err != nil {
		return err
	}
	re := regexp.MustCompile("eth0@if(.*?):")
	match := re.FindStringSubmatch(string(out))
	peerIdx := match[1]

	// Now that we have the ID, we need to look on the host network to find its name.
	cmd = exec.Command("docker", "run", "--rm", "--net=host", "antrea/ethtool:latest", "ip", "link")
	outLines, err := exec.OutputLines(cmd)
	if err != nil {
		return err
	}
	peerName := ""
	re = regexp.MustCompile(fmt.Sprintf("^%s: (.*?)@.*:", peerIdx))
	for _, line := range outLines {
		match = re.FindStringSubmatch(line)
		if len(match) > 0 {
			peerName = match[1]
			break
		}
	}

	if peerName == "" {
		return fmt.Errorf("unable to find node interface %q on host network", peerIdx)
	}

	// With the name, we can now use ethtool to turn off TX checksumming offload
	cmd = exec.Command("docker", "run", "--rm", "--net=host", "--privileged", "antrea/ethtool:latest", "ethtool", "-K", peerName, "tx", "off")
	_, err = exec.Output(cmd)
	if err != nil {
		return err
	}

	// Finally, enable local routing
	cmd = exec.Command("docker", "exec", nodeName, "sysctl", "-w", "net.ipv4.conf.all.route_localnet=1")
	_, err = exec.Output(cmd)
	if err != nil {
		return err
	}

	return nil
}
