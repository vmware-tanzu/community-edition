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
	"sigs.k8s.io/kind/pkg/exec"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/config"
)

const (
	minMemoryBytes     = 2147483648
	minCPUCount        = 1
	kindConfigFileName = "kindconfig.yaml"
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

	// Serlize unstructured data into a kindProviderConfig.
	// Return any error from attempting to read the data
	serializedProviderConfig, err := serializeKindProviderConfig(c.ProviderConfiguration)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize kind config from given ProviderConfiguration. Error was: %s", err)
	}

	// If a user has provided something in the ProviderConfiguration,
	// assume it is a full kind configuration and don't attempt to produce
	// a configuration from the unmanaged-cluster config
	parsedKindConfig := []byte(serializedProviderConfig.rawKindConfig)

	// when the parsed kind config from the provider config is empty,
	// create a kind config from ClusterConfig settings, using the given flags and options
	if len(parsedKindConfig) < 1 {
		parsedKindConfig, err = kindConfigFromClusterConfig(c)
		if err != nil {
			return nil, fmt.Errorf("failed to generate a viable kind config. Error was: %s", err)
		}
	}

	// store our kind config on the filesystem for users to inspect if needed
	err = writeKindConfigFile(parsedKindConfig, c.ClusterName)
	if err != nil {
		return nil, err
	}

	kindConfig := kindcluster.CreateWithRawConfig(parsedKindConfig)
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
		nodes, _ := kindProvider.ListNodes(c.ClusterName)
		for _, n := range nodes {
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

type kindProviderConfig struct {
	rawKindConfig string
}

func serializeKindProviderConfig(pc map[string]interface{}) (kindProviderConfig, error) {
	// Check if key exists. If not, return empty config and continue
	if _, ok := pc["rawKindConfig"]; !ok {
		return kindProviderConfig{}, nil
	}

	// Check if provided data is a string.
	if _, ok := pc["rawKindConfig"].(string); !ok {
		return kindProviderConfig{}, fmt.Errorf("ProviderConfiguration.rawKindConfig wrong type, expected string")
	}

	return kindProviderConfig{
		pc["rawKindConfig"].(string),
	}, nil
}

func kindConfigFromClusterConfig(c *config.UnmanagedClusterConfig) ([]byte, error) {
	// Load the defaults
	kindConfig := &kindconfig.Cluster{}
	kindConfig.Kind = "Cluster"
	kindConfig.APIVersion = "kind.x-k8s.io/v1alpha4"
	kindConfig.Name = c.ClusterName
	nodes, err := setNumberOfNodes(c)
	if err != nil {
		return nil, err
	}
	kindConfig.Nodes = nodes
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

	// Marshal it into the raw bytes we need for creation
	var rawConfig bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&rawConfig)

	err = yamlEncoder.Encode(kindConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Kind configuration. Error: %s", err.Error())
	}
	err = yamlEncoder.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Kind configuration. Error: %s", err.Error())
	}

	return rawConfig.Bytes(), nil
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

	nodes := []kindconfig.Node{}

	for i := 1; i <= cpnc; i++ {
		n := kindconfig.Node{
			Role: kindconfig.ControlPlaneRole,
		}

		nodes = append(nodes, n)
	}

	for i := 1; i <= wnc; i++ {
		n := kindconfig.Node{
			Role: kindconfig.WorkerRole,
		}

		nodes = append(nodes, n)
	}

	return nodes, nil
}

// Get retrieves cluster information or return an error indicating a problem.
// TODO - (jpmcb) We currently do not utilize the Get API on cluster providers
func (kcm *KindClusterManager) Get(clusterName string) (*KubernetesCluster, error) {
	return nil, nil
}

// Delete removes a kind cluster.
func (kcm *KindClusterManager) Delete(c *config.UnmanagedClusterConfig) error {
	provider := kindcluster.NewProvider()
	return provider.Delete(c.ClusterName, "")
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
