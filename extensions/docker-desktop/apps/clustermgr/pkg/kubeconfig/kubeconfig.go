// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package kubeconfig has slightly more generic kubeconfig helpers and
// minimal dependencies on the rest of kind
package kubeconfig

import (
	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/config"
	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/kubeconfig/internal/kubeconfig"
)

func AddConfig(kubeconfigToAddPath, kubeconfigPath string) error {
	origConf, err := kubeconfig.Read(kubeconfigToAddPath)
	if err != nil {
		return err
	}

	// Modify the kubeconfig
	getTCEConfig(origConf)

	return kubeconfig.WriteMerged(origConf, kubeconfigPath)
}

func RemoveConfig(kubeconfigToAddPath, kubeconfigPath string) error {
	origConf, err := kubeconfig.Read(kubeconfigToAddPath)
	if err != nil {
		return err
	}
	return kubeconfig.RemoveKIND(origConf.CurrentContext, kubeconfigPath)
}

func RemoveNamedConfig(context, kubeconfigPath string) error {
	return kubeconfig.RemoveKIND(context, kubeconfigPath)
}

func GetConfig(kubeconfigPath string) ([]byte, error) {
	cfg, err := kubeconfig.Read(kubeconfigPath)
	if err != nil {
		return nil, err
	}
	// Modify the kubeconfig
	getTCEConfig(cfg)
	encoded, err := kubeconfig.Encode(cfg)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

func getTCEConfig(cfg *kubeconfig.Config) {
	// We change all the config Names in the kubeconfig to be the Cluster name
	key := config.DefaultClusterName
	cfg.Clusters[0].Name = key
	cfg.Users[0].Name = key
	cfg.Contexts[0].Name = key
	cfg.Contexts[0].Context.User = key
	cfg.Contexts[0].Context.Cluster = key
	cfg.CurrentContext = key
}

// // Export exports the kubeconfig given the cluster context and a path to write it to
// // This will always be an external kubeconfig
// func Export(p providers.Provider, name, explicitPath string, external bool) error {
// 	// cfg, err := get(p, name, external)
// 	cfg := "kubeconfig"
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	return kubeconfig.WriteMerged(cfg, explicitPath)
// }

// // Remove removes clusterName from the kubeconfig paths detected based on
// // either explicitPath being set or $KUBECONFIG or $HOME/.kube/config, following
// // the rules set by kubectl
// // clusterName must identify a kind cluster.
// func Remove(clusterName, explicitPath string) error {
// 	return kubeconfig.RemoveKIND(clusterName, explicitPath)
// }

// // Get returns the kubeconfig for the cluster
// // external controls if the internal IP address is used or the host endpoint
// func Get(name string, external bool) (string, error) {
// 	cfg, err := get("p", name, external)
// 	if err != nil {
// 		return "", err
// 	}
// 	b, err := kubeconfig.Encode(cfg)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(b), err
// }

// // ContextForCluster returns the context name for a kind cluster based on
// // its name. This key is used for all list entries of kind clusters
// func ContextForCluster(kindClusterName string) string {
// 	return kubeconfig.KINDClusterKey(kindClusterName)
// }

// func get(p providers.Provider, name string, external bool) (*kubeconfig.Config, error) {
// 	// find a control plane node to get the kubeadm config from
// 	n, err := p.ListNodes(name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var buff bytes.Buffer
// 	nodes, err := nodeutils.ControlPlaneNodes(n)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(nodes) < 1 {
// 		return nil, errors.Errorf("could not locate any control plane nodes for cluster named '%s'. "+
// 			"Use the --name option to select a different cluster", name)
// 	}
// 	node := nodes[0]

// 	// grab kubeconfig version from the node
// 	if err := node.Command("cat", "/etc/kubernetes/admin.conf").SetStdout(&buff).Run(); err != nil {
// 		return nil, errors.Wrap(err, "failed to get cluster internal kubeconfig")
// 	}

// 	// if we're doing external we need to override the server endpoint
// 	server := ""
// 	if external {
// 		endpoint, err := p.GetAPIServerEndpoint(name)
// 		if err != nil {
// 			return nil, err
// 		}
// 		server = "https://" + endpoint
// 	}

// 	// actually encode
// 	return kubeconfig.KINDFromRawKubeadm(buff.String(), name, server)
// }
