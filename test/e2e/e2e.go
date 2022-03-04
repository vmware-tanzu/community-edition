// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/vmware-tanzu/community-edition/test/e2e/testdata"
	"github.com/vmware-tanzu/community-edition/test/e2e/utils"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/e2e-framework/klient/conf"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
)

type Config struct {
	Kubeconfig             string
	Kubecontext            string
	Provider               string
	ClusterType            string
	Packages               string
	Version                string
	GuestClusterName       string
	ClusterPlan            string
	MgmtClusterName        string
	ClusterInstallRequired bool
	ClusterCleanupRequired bool
	TceVersion             string
}

var (
	addonsCfg        *PackageConfiguration
	ConfigVal        *Config
	MetallbInstalled bool
	VeleroInstalled  bool
)

func init() {
	ConfigVal = New()
}

func New() *Config {
	e2eConfig := new(Config)
	flag.BoolVar(&e2eConfig.ClusterInstallRequired, "create-cluster", false, "Is cluster provision required? Provide true.")
	flag.StringVar(&e2eConfig.Kubeconfig, "kubeconfig", "", "Paths to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&e2eConfig.Kubecontext, "kube-context", "", "Cluster context need to be set. ")
	flag.StringVar(&e2eConfig.Provider, "provider", "", "Provider name in which cluster is running. Values can be docker, aws, vsphere")
	flag.StringVar(&e2eConfig.ClusterType, "cluster-type", "", "Provide cluster type. eg: management")
	flag.StringVar(&e2eConfig.Packages, "packages", "", "Provide package list or 'all'. eg:--packages=all, --packages='antrea, external-dns'")
	flag.StringVar(&e2eConfig.Version, "version", "", "Provide package version. eg: --version='0.11.3,0.8.0'")
	flag.StringVar(&e2eConfig.TceVersion, "tce-version", "", "Provide tce release version to install. If not provided then build it from source code.")
	flag.StringVar(&e2eConfig.GuestClusterName, "guest-cluster-name", "", "Provide a cluster name.")
	flag.StringVar(&e2eConfig.MgmtClusterName, "management-cluster-name", "", "Provide a cluster name.")
	flag.StringVar(&e2eConfig.ClusterPlan, "cluster-plan", "dev", "Provide a cluster plan.")
	flag.BoolVar(&e2eConfig.ClusterCleanupRequired, "cleanup-cluster", false, "Provide true for cluster tear down.")
	return e2eConfig
}

func Initialize() {
	// based on the release version provided installs tce
	err := installTCE()
	if err != nil {
		log.Fatal("error while installing TCE :", err)
	}

	if ConfigVal.ClusterInstallRequired {
		err := installCluster()
		if err != nil {
			log.Println("error while creating cluster, deleting the cluster", err)
			err = DeleteCluster()
			if err != nil {
				log.Fatal("error while deleting cluster", err)
			}

			os.Exit(1)
		}
	}

	addonsCfg, err = readPackageConfig("addons_config.yaml")
	if err != nil {
		log.Println("error while reading config file, deleting cluster", err)
		err = DeleteCluster()
		if err != nil {
			log.Fatal("error while deleting cluster", err)
		}

		os.Exit(1)
	}
}

var (
	Clientset *kubernetes.Clientset
	Res       *resources.Resources
)

// GetKubeConfig will set cluster context passed using flag -kube-context
// and gives you the *rest.Config to access the cluster objects
func GetKubeConfig(clusterContext string) (*rest.Config, error) {
	return conf.NewWithContextName(conf.ResolveKubeConfigFile(), clusterContext)
}

// GetResourceObj sets resource object to communicate with apiserver
func GetResourceObj(clusterContext string) error {
	var err error
	cfg, err := GetKubeConfig(clusterContext)
	if err != nil {
		log.Println("Error while getting kube config", err)
		return err
	}

	Res, err = resources.New(cfg)
	if err != nil {
		log.Println("Error while getting resource object", err)
		return err
	}

	return nil
}

func RunAddonsTests() error {
	// install MetalLB in case of provider Docker
	// as it does not have its own LB support.
	if ConfigVal.Provider == DOCKER {
		err := testdata.InstallMetallb()
		if err != nil {
			log.Println("Error while installing metal-lb", err)
			return err
		}

		MetallbInstalled = true
	}

	if ConfigVal.Packages != "" {
		pkgList := strings.Split(ConfigVal.Packages, ",")
		versionList := strings.Split(ConfigVal.Version, ",")

		if len(pkgList) == 1 {
			if pkgList[0] == "all" {
				// running tests of all pkgs
				err := runAllPackageTest()
				if err != nil {
					return err
				}
				return nil
			}
		}

		for i, pkgName := range pkgList {
			// run the pkg defind by the user
			err := runPackageTest(pkgName, versionList[i])
			if err != nil {
				return err
			}
		}
	} else {
		log.Println("Skipping Package testing....")
	}
	return nil
}

func runAllPackageTest() error {
	// go to addons/package directory and run the tests
	for _, pkgs := range addonsCfg.Packages {
		// run test
		err := runPackageTest(pkgs.Name, pkgs.Version)
		if err != nil {
			return err
		}
	}
	return nil
}

func runPackageTest(pkgName, version string) error {
	// go to addons/package/{packagename} and run the tests
	err := os.Chdir(utils.WorkingDir + "/../../addons/packages/" + pkgName + "/" + version + "/test/")
	if err != nil {
		log.Println("Error while changing directory :", err)
		return err
	}

	// Install velero to run the test
	if pkgName == "velero" {
		runDeployScript("e2e/utils/velero/velero_prefix_setup.sh", "")

		err := testdata.InstallVelero(version)
		if err != nil {
			log.Println("Error while installing Velero", err)
			return err
		}

		runDeployScript("e2e/utils/velero/velero_prefix_cleanup.sh", "")

		// installing AWS CLI
		installAWSCli()

		VeleroInstalled = true
	}

	mydir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	log.Println("Running package testing in ", mydir)

	err = RunCommand("make", "e2e-test")
	if err != nil {
		return err
	}

	return nil
}

func RunCommand(commandName, args string) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	mwriter := io.MultiWriter(os.Stdout)
	cmd := exec.Command(commandName, args)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmd.Stderr = mwriter
	cmd.Stdout = mwriter
	err := cmd.Run()
	if err != nil {
		log.Println(stdout.String(), stderr.String())
		return err
	}

	log.Println(stdout.String(), stderr.String())
	return nil
}
