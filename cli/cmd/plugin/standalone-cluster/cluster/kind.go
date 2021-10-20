// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cluster

import (
	"fmt"
	"regexp"

	kindCluster "sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/exec"
)

const defaultKindConfig = `kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  #! port forward 80 on the host to 80 on this node
  extraPortMappings:
  - containerPort: 80
    #!hostPort: 80
    #! optional: set the bind address on the host
    #! 0.0.0.0 is the current default
    listenAddress: "127.0.0.1"
    #! optional: set the protocol to one of TCP, UDP, SCTP.
    #! TCP is the default
    protocol: TCP
- role: worker
- role: worker
networking:
  disableDefaultCNI: true`

// KindClusterManager is a ClusterManager implementation for working with
// Kind clusters.
type KindClusterManager struct {
}

// Create will create a new kind cluster or return an error.
func (kcm KindClusterManager) Create(opts *CreateOpts) (*KubernetesCluster, error) {
	kindProvider := kindCluster.NewProvider()
	clusterConfig := kindCluster.CreateWithKubeconfigPath(opts.KubeconfigPath)

	// TODO(stmcginnis): Determine what we need to do for kind configuration
	// generates the kind configuration -- TODO(joshrosso): should not exec ytt; use go lib
	// command := exec.Command("ytt",
	// 	"-f",
	// 	"cli/cmd/plugin/standalone-cluster/hack/kind-config")
	// parsedKindConfig, err := command.Output()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return err
	// }
	parsedKindConfig := []byte(defaultKindConfig)
	kindConfig := kindCluster.CreateWithRawConfig(parsedKindConfig)
	err := kindProvider.Create(opts.Name, clusterConfig, kindConfig)
	if err != nil {
		return nil, err
	}

	kc := &KubernetesCluster{
		Name:       opts.Name,
		Kubeconfig: opts.KubeconfigPath,
	}

	// TODO(stmcginnis): This should be conditional based on the CNI chosen.
	nodes, _ := kindProvider.ListNodes(opts.Name)
	for _, n := range nodes {
		patchForAntrea(n.String())
	}

	return kc, nil
}

// Get retrieves cluster information or return an error indicating a problem.
func (kcm KindClusterManager) Get(clusterName string) (*KubernetesCluster, error) {
	return nil, nil
}

// List gets all kind clusters.
func (kcm KindClusterManager) List() ([]*KubernetesCluster, error) {
	provider := kindCluster.NewProvider()
	clusters, err := provider.List()
	if err != nil {
		return nil, err
	}

	var result []*KubernetesCluster

	// TODO(stmcginnis): Need to figure out a way to filter out only tanzu clusters
	// in case there are other kind clusters present.
	for _, cl := range clusters {
		result = append(result, &KubernetesCluster{
			Name: cl,
		})
	}

	return result, nil
}

// Delete removes a kind cluster.
func (kcm KindClusterManager) Delete(clusterName string) error {
	provider := kindCluster.NewProvider()
	return provider.Delete(clusterName, "")
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
	peerIdx := string(match[1])

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
	out, err = exec.Output(cmd)
	if err != nil {
		return err
	}

	// Finally, enable local routing
	cmd = exec.Command("docker", "exec", nodeName, "sysctl", "-w", "net.ipv4.conf.all.route_localnet=1")
	out, err = exec.Output(cmd)
	if err != nil {
		return err
	}

	return nil
}
