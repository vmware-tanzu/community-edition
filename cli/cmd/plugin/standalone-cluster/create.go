// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/tanzu"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/cluster"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/kapp"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/kubeconfig"
	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/packages"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/tkr"
)

const (
	configDir             = ".config"
	tanzuConfigDir        = "tanzu"
	tkgConfigDir          = "tkg"
	tkgSysNamespace       = "tkg-system"
	tkgSvcAcctName        = "core-pkgs"
	tkgCoreRepoName       = "tkg-core-repository"
	tkgGlobalPkgNamespace = "tanzu-package-repo-global"
	tceRepoName           = "community-repository"
	tceRepoUrl            = "projects.registry.vmware.com/tce/main:0.9.1"
)

type initStandaloneOptions struct {
	clusterConfigFile      string
	infrastructureProvider string
	ui                     bool
	bind                   string
	browser                string
}

type createOptions struct {
	tty bool
}

// CreateCmd creates a standalone workload cluster.
var CreateCmd = &cobra.Command{
	Use:   "create <cluster name> -f <configuration location>",
	Short: "create a local tanzu cluster",
	RunE:  create,
}

var iso = initStandaloneOptions{}
var createOpts = createOptions{}

func init() {
	CreateCmd.Flags().StringVarP(&iso.clusterConfigFile, "config", "f", "", "Configuration file for local cluster creation")
	CreateCmd.Flags().StringVarP(&iso.clusterConfigFile, "kind-config", "k", "", "Kind configuration file; fully overwrites Tanzu defaults")
	CreateCmd.Flags().StringVarP(&iso.clusterConfigFile, "port-forward", "p", "", "Port to forward from host to container")
	CreateCmd.Flags().BoolVar(&createOpts.tty, "tty", true, "Specify whether terminal is tty;\\nSet to false to disable styled ouput; default: true")
}

func create(cmd *cobra.Command, args []string) error {
	var clusterName string

	// validate a cluster name was passed when not using the kickstart UI
	if len(args) < 1 && !iso.ui {
		return Error(nil, "no cluster name specified")
	} else if len(args) == 1 {
		clusterName = args[0]
	}
	log := logger.NewLogger(createOpts.tty, 0)

	tm := tanzu.New(clusterName)
	err := tm.Deploy("")
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	return nil
}

// TODO(joshrosso): a lot of this functionality should be moved into pkg/* so that it's importable and not tied
// to the cobra command creation.
func create2(cmd *cobra.Command, args []string) error {
	var clusterName string

	// validate a cluster name was passed when not using the kickstart UI
	if len(args) < 1 && !iso.ui {
		return Error(nil, "no cluster name specified")
	} else if len(args) == 1 {
		clusterName = args[0]
	}
	log := logger.NewLogger(createOpts.tty, 0)

	// Resolve the BOM location
	// TODO(joshrosso): BOM is currently static. Need to be resolvable/configurable
	tkgConfigPath, err := getTkgConfigDir()
	if err != nil {
		log.Errorf("Unable to resolve TKG config path. Error: %s", err.Error())
	}
	bomPath := filepath.Join(tkgConfigPath, "bom", "tkr-bom-v1.21.2+vmware.1-tkg.1.yaml")

	log.Event("\\U+2692", " Processing TanzuKubernetesRelease (TKR)\n")
	bom, err := tkr.ReadTKRBom(bomPath)
	if err != nil {
		return err
	}

	// Resolve the base image to use
	kindNodeImage := bom.GetTKRNodeImage()
	log.Event("\\U+1F5BC", " Selected base image\n")
	log.Style(2, logger.ColorLightGrey).Info(kindNodeImage + "\n")

	// Resolve the kapp-controller image to use
	kappControllerBundle, err := bom.GetTKRKappImage()
	if err != nil {
		return err
	}

	log.Event("\\U+1F4E6", "Selected kapp-controller image bundle\n")
	log.Style(2, logger.ColorLightGrey).Info(kappControllerBundle.GetRegistryUrl() + "\n")

	// Resolve the core-package image bundle to use
	corePackageBundle := bom.GetTKRCoreRepoBundlePath()
	log.Event("\\U+1F4E6", "Selected core package repository\n")
	log.Style(2, logger.ColorLightGrey).Info(corePackageBundle + "\n")
	// Resolve user package bundle repo(s) to use
	userPackageBundles := []string{tceRepoUrl}
	log.Event("\\U+1F4E6", "Selected additional package repositories\n")
	for _, upb := range userPackageBundles {
		log.Style(2, logger.ColorLightGrey).Info(upb + "\n")
	}

	// Create the cluster
	kubeConfigPath := filepath.Join(os.Getenv("HOME"), configDir, tanzuConfigDir, clusterName+".yaml")
	log.Eventf("\\U+1F6F0", " Creating cluster %s\n", clusterName)
	log.Style(2, logger.ColorLightGrey).Info("To troubleshoot, use:\n")
	log.Style(2, logger.ColorLightGrey).Infof("kubectl ${COMMAND} --kubeconfig %s\n", kubeConfigPath)
	clusterManager := cluster.NewClusterManager()
	clusterCreateOpts := cluster.CreateOpts{
		Name:           clusterName,
		KubeconfigPath: kubeConfigPath,
		// Config: TBD,
	}
	_, err = clusterManager.Create(&clusterCreateOpts)
	if err != nil {
		log.Errorf("Failed to create cluster: %s", err.Error())
		return nil
	}

	log.Event("\\U+1F4E7", "Installing kapp-controller\n")

	err = kappControllerBundle.DownloadBundleImage()
	if err != nil {
		return err
	}

	kappControllerBundle.SetRelativeConfigPath("config/")
	kappValues, err := ioutil.ReadFile("cli/cmd/plugin/standalone-cluster/hack/kapp-values.yaml")
	if err != nil {
		return err
	}

	kappControllerBundle.AddYttYamlValuesBytes(kappValues)
	kappBytes, err := kappControllerBundle.RenderYaml()
	if err != nil {
		return err
	}

	kc := kapp.New(kubeConfigPath)
	kappControllerCreated, err := kc.Install(kapp.KappInstallOpts{MergedManifests: kappBytes[0]})
	if err != nil {
		return err
	}

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// Wait for kapp-controller be running; report status
	for si := 1; si < 5; si++ {
		kappState := kc.Status(kappControllerCreated.Namespace, kappControllerCreated.Name)
		log.Style(2, logger.ColorLightGrey).Progressf(si, "kapp-controller status: %s", kappState)
		if kappState == "Running" {
			log.Style(2, logger.ColorLightGrey).Progressf(0, "kapp-controller status: %s", kappState)
			break
		}
		if si == 4 {
			si = 1
		}
		time.Sleep(1 * time.Second)
	}

	// Create package client used for package installs
	pkgClient := packages.NewClient(kubeConfigPath)

	// Install package repositories
	log.Event("\\U+1F4E7", "Installing package repositories\n")
	// install core package repository
	_, err = pkgClient.CreatePackageRepo(tkgSysNamespace, tkgCoreRepoName, corePackageBundle)
	if err != nil {
		log.Errorf("failed to add core package repo: %s\n", err.Error())
	}
	// install user package repos
	for _, upb := range userPackageBundles {
		_, err = pkgClient.CreatePackageRepo(tkgGlobalPkgNamespace, tceRepoName, upb)
		if err != nil {
			log.Errorf("failed to add core package repo: %s\n", err.Error())
		}
	}

	// Wait for core package repository to reconcile; report status
	for si := 1; si < 5; si++ {
		status, err := pkgClient.GetRepositoryStatus(tkgSysNamespace, tkgCoreRepoName)
		if err != nil {
			log.Errorf("Failed to check kapp-controller status: %s", err.Error())
			return nil
		}
		log.Style(2, logger.ColorLightGrey).Progressf(si, "Core package repo status: %s", status)
		if status == "Reconcile succeeded" {
			log.Style(2, logger.ColorLightGrey).Progressf(0, "Core package repo status: %s", status)
			break
		}
		if si == 4 {
			si = 1
		}
		time.Sleep(1 * time.Second)
	}

	// install CNI (TODO(joshrosso): needs to support multiple CNIs
	rootSvcAcct, err := pkgClient.CreateRootServiceAccount(tkgSysNamespace, tkgSvcAcctName)
	if err != nil {
		log.Errorf("failed to create service account: %s\n", err.Error())
		return nil
	}

	// run the antrea patch for kind-specific deployments
	nodes := ListNodes(clusterName)
	for _, node := range nodes {
		err := cluster.PatchForAntrea(node)
		if err != nil {
			log.Errorf("Failed to patch node!!! %s\n", err.Error())
		}
	}

	// antrea data
	valueData := `---
infraProvider: docker
`
	cniInstallOpts := packages.PackageInstallOpts{
		Namespace:      tkgSysNamespace,
		InstallName:    "cni",
		FqPkgName:      "antrea.tanzu.vmware.com",
		Version:        "0.13.3+vmware.1-tkg.1",
		Configuration:  []byte(valueData),
		ServiceAccount: rootSvcAcct.Name,
	}
	_, err = pkgClient.CreatePackageInstall(cniInstallOpts)
	if err != nil {
		log.Errorf("failed to add package install for CNI: %s\n", err.Error())
	}

	log.Event("\\U+2705", "Cluster created\n")

	// Merge working kubeconfig into main
	kubeConfigMgr := kubeconfig.NewManager()
	err = kubeConfigMgr.MergeToDefaultConfig(kubeConfigPath)
	if err != nil {
		log.Errorf("Failed to merge kubeconfig: %s\n", err.Error())
		return nil
	}
	kubeContextName := fmt.Sprintf("%s-%s", "kind", clusterName)
	err = kubeConfigMgr.SetCurrentContext(kubeContextName)
	if err != nil {
		log.Errorf("Failed to set default contxt: %s\n", err.Error())
		return nil
	}
	log.Eventf("\\U+1F3AE", "kubectl context set to %s\n\n", clusterName)

	// provide user example commands to run
	log.Style(0, logger.ColorLightGrey).Infof("View available packages:\n")
	log.Style(2, logger.ColorLightGreen).Infof("tanzu package available list\n")
	log.Style(0, logger.ColorLightGrey).Infof("View running pods:\n")
	log.Style(2, logger.ColorLightGreen).Infof("kubectl get po -A\n")
	log.Style(0, logger.ColorLightGrey).Infof("Delete this cluster:\n")
	log.Style(2, logger.ColorLightGreen).Infof("tanzu local delete %s\n", clusterName)

	return nil
}

func createAntreaConfig(kubeConfig string) (*v1.Secret, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		fmt.Printf("Unable to create client config to contact cluster: %s", err.Error())
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// for docker you need this or else assumptions are made like availability of OVS
	valueData := `---
infraProvider: docker
`

	values := make(map[string]string)
	values["values.yml"] = valueData

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "antrea-values",
			Namespace: tkgSysNamespace,
		},
		StringData: values,
	}

	createdSecret, err := clientset.CoreV1().Secrets(tkgSysNamespace).Create(secret)
	if err != nil {
		fmt.Printf("Failed to create secret: %s", err.Error())
		return nil, err
	}

	return createdSecret, err
}

// getTkgConfigDir returns the configuration directory used by tce.
func getTkgConfigDir() (path string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return path, fmt.Errorf("Failed to resolve home dir. Error: %s", err.Error())
	}
	path = filepath.Join(home, configDir, tanzuConfigDir, tkgConfigDir)
	return
}
