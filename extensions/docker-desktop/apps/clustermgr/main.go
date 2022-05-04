// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/cluster"
	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/config"
	"github.com/vmware-tanzu/community-edition/extensions/docker-desktop/pkg/kubeconfig"
)

var log = logrus.New()

func init() {
	// The API for setting attributes is a little different than the package level
	// exported logger. See Godoc.
	// log.Out = os.Stdout

	// You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile(filepath.Join(config.GetUserHome(), "cluster.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	// Only log the warning severity or above.
	log.SetLevel(logrus.DebugLevel)
}

//nolint:errcheck
func main() {
	var res config.Response
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) != 1 {
		fmt.Println("This program requires any of the following parameters: create|delete|status|reset|logs|kubeconfig|stats|provision-ingress|provision-certmanager|provision-kubeapps")
		os.Exit(-1)
	}
	c := cluster.New(log)
	switch argsWithoutProg[0] {
	case "create":
		res = c.CreateCluster()
	case "delete":
		res = c.DeleteCluster()
	case "status":
		res = c.ClusterStatus()
	case "reset":
		res = c.Reset() // TODO:
	case "logs":
		res = c.Logs() // TODO:
	case "kubeconfig":
		res = c.GetKubeconfig() // TODO:
	case "stats":
		res = c.Stats() // TODO:
	case "provision-ingress":
		res = c.ProvisionIngress() // TODO:
	case "provision-certmanager":
		res = c.ProvisionCertMan() // TODO:
	case "add-kubeconfig":
		kubeconfig.AddConfig("/home/tanzu/.kube/config", "/opt/kubeconfig/config")
	case "remove-kubeconfig":
		kubeconfig.RemoveNamedConfig(config.DefaultClusterName, "/opt/kubeconfig/config")
	case "test":
		kubeconfig.AddConfig("/Users/jomorales/.kube/config", "/tmp/kubeconfig")
		// TODO: Is this needed?
		// kubeconfig.RemoveNamedConfig(config.DefaultClusterName, "/tmp/kubeconfig")
		// res = c.Test()
	}
	fmt.Println(c.GetJSONResponse(&res))
}
