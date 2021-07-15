// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/vmware-tanzu/community-edition/test/e2e/utils"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/klient/conf"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
)

const (
	STANDALONE = "standalone"
	MANAGED    = "managed"
	DOCKER     = "docker"
	TIMEOUT    = 120
)

func DeployTanzuCluster() error {
	switch ConfigVal.Provider {
	case "docker":
		if ConfigVal.ClusterType == STANDALONE {
			err := createStandaloneCluster()
			if err != nil {
				return err
			}
		}

		if ConfigVal.ClusterType == MANAGED {
			err := createManagedCluster()
			if err != nil {
				return err
			}
		}
	case "aws":
		// TODO
		log.Println("Currently not implemented.")
	case "azure":
		// TODO
		log.Println("Currently not implemented.")
	}

	// ADD tce package repository
	runDeployScript("add-tce-package-repo.sh", "")
	return nil
}

func createStandaloneCluster() error {
	log.Println("Create standalone cluster is in progress.......")

	os.Setenv("CLUSTER_PLAN", ConfigVal.ClusterPlan)
	os.Setenv("CLUSTER_NAME", ConfigVal.GuestClusterName)
	start := time.Now()
	_, err := utils.Tanzu(nil, "standalone-cluster", "create", ConfigVal.GuestClusterName, "-i", ConfigVal.Provider, "-v", "10")
	if err != nil {
		log.Println("Standalone cluster creation failed", err)
		return err
	}
	log.Println("Time taken for standalone cluster provision :", time.Since(start))

	err = CheckClusterHealth(utils.GetClusterContext(ConfigVal.GuestClusterName))
	if err != nil {
		log.Println("Cluster is not healthy", err)
		return err
	}

	log.Println("Standalone docker cluster is up and running......")
	return nil
}

func createManagedCluster() error {
	log.Println("Create management cluster is in progress.......")
	start := time.Now()
	_, err := utils.Tanzu(nil, "management-cluster", "create", "-i", "docker", "--name", ConfigVal.MgmtClusterName, "-v", "10", "--plan", ConfigVal.ClusterPlan, "--ceip-participation=false")
	if err != nil {
		log.Fatal("Management cluster creation failed", err)
	}

	log.Println("Time taken for management cluster provision :", time.Since(start))

	err = CheckClusterHealth(utils.GetClusterContext(ConfigVal.MgmtClusterName))
	if err != nil {
		log.Println("Management cluster is not healthy", err)
		return err
	}
	log.Println("Docker management cluster is up and running......")

	log.Println("Create workload cluster is in progress.......")
	s := time.Now()
	_, err = utils.Tanzu(nil, "cluster", "create", ConfigVal.GuestClusterName, "--plan", ConfigVal.ClusterPlan)
	if err != nil {
		log.Fatal("Workload cluster creation is failed", err)
	}

	log.Println("Time taken for workload cluster provision :", time.Since(s))

	// save workload cluster credentials
	_, err = utils.Tanzu(nil, "cluster", "kubeconfig", "get", ConfigVal.GuestClusterName, "--admin")
	if err != nil {
		log.Println("Error while saving workload cluster credentials", err)
		return err
	}

	err = CheckClusterHealth(utils.GetClusterContext(ConfigVal.GuestClusterName))
	if err != nil {
		log.Println("Workload cluster is not healthy", err)
		return err
	}
	log.Println("Docker worker load cluster is up and running......")

	return nil
}

func CheckClusterHealth(contextName string) error {
	restCfg, err := conf.NewWithContextName(conf.ResolveKubeConfigFile(), contextName)
	if err != nil {
		fmt.Printf("error %s", err)
		return err
	}

	resource, err := resources.New(restCfg)
	if err != nil {
		log.Println("Error while getting resource object", err)
		return err
	}

	pods := &v1.PodList{}
	err = resource.List(context.TODO(), pods)
	if err != nil {
		fmt.Printf("error %s", err)
		return err
	}

	if pods.Items == nil {
		fmt.Printf("Cluster is not healthy or not yet created %s", err)
		return errors.New("cluster is not up")
	}

	_, err = utils.Kubectl(nil, "config", "use-context", contextName)
	if err != nil {
		fmt.Printf("error %s", err)
		return err
	}

	// check all deployments are in ready state
	utils.ValidateAllDeploymentsReady()

	log.Println("Cluster is up and running ...")

	return nil
}

func DeleteCluster() error {
	log.Println("Provider and Cluster type is", ConfigVal.Provider, ConfigVal.ClusterType)
	if ConfigVal.Provider == DOCKER {
		if ConfigVal.ClusterType == STANDALONE {
			log.Println("Executing command delete standard docker cluster", ConfigVal.GuestClusterName)
			_, err := utils.Tanzu(nil, "standalone-cluster", "delete", ConfigVal.GuestClusterName, "--yes")
			if err != nil {
				log.Println("Standalone cluster delete failed", err)
				return err
			}
		}

		if ConfigVal.ClusterType == MANAGED {
			_, err := utils.Tanzu(nil, "cluster", "delete", ConfigVal.GuestClusterName, "--yes")
			if err != nil {
				log.Println("Workload cluster delete failed", err)
				return err
			}

			_, err = utils.Tanzu(nil, "management-cluster", "delete", ConfigVal.MgmtClusterName, "--yes")
			if err != nil {
				log.Println("Management cluster delete failed", err)
				return err
			}
		}
	}

	return nil
}
