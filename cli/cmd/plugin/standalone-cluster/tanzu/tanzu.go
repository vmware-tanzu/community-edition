// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package tanzu is responsible for orchestrating the various components that result in
// a local tanzu cluster creation.
package tanzu

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	v1 "k8s.io/api/apps/v1"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/cluster"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/kapp"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/kubeconfig"
	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/packages"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/tkr"
)

//nolint
const (
	configDir             = ".config"
	tanzuConfigDir        = "tanzu"
	tkgConfigDir          = "tkg"
	localConfigDir        = "local"
	bomDir                = "bom"
	tkgSysNamespace       = "tkg-system"
	tkgSvcAcctName        = "core-pkgs"
	tkgCoreRepoName       = "tkg-core-repository"
	tkgGlobalPkgNamespace = "tanzu-package-repo-global"
	tceRepoName           = "community-repository"
	tceRepoURL            = "projects.registry.vmware.com/tce/main:0.9.1"
	outputIndent          = 2
	maxProgressLength     = 4
)

// TODO(joshrosso): global logger for the package. This is kind gross, but really convenient.
var log logger.Logger

// TanzuCluster contains information about a cluster.
//nolint:golint
type TanzuCluster struct {
	Name string
}

// TanzuLocal contains information about a local Tanzu cluster.
//nolint:golint
type TanzuLocal struct {
	name                 string
	kcPath               string
	bom                  *tkr.TKRBom
	kappControllerBundle tkr.TkrImageReader
}

//nolint:golint
type TanzuMgr interface {
	// Deploy orchestrates all the required steps in order to create a local Tanzu cluster. This can involve
	// cluster creation, kapp-controller installation, CNI installation, and more. The steps that are taken
	// depend on the configuration passed into Deploy. If something goes wrong during deploy, an error is
	// returned.
	// TODO(joshrosso): this config will be replaced with the API/struct @stmcginnis comes up with
	Deploy(config string) error
	// List retrieves all known tanzu clusters are returns a list of them. If it's unable to interact with the
	// underlying cluster provider, it returns an error.
	List() ([]TanzuCluster, error)
	// Delete takes a cluster name and removes the cluster from the underlying cluster provider. If it is unable
	// to communicate with the underlying cluster provider, it returns an error.
	Delete(name string) error
}

// New returns a TanzuMgr for interacting with local clusters. It is implemented by TanzuLocal.
func New(name string) TanzuMgr {
	// TODO(joshrosso): bring this in from CMD command
	log = logger.NewLogger(true, 0)
	defaultKubeConfigPath := filepath.Join(os.Getenv("HOME"), configDir, tanzuConfigDir, name+".yaml")
	return &TanzuLocal{name: name, kcPath: defaultKubeConfigPath}
}

// Deploy deploys a new cluster.
//nolint:funlen
func (t *TanzuLocal) Deploy(config string) error {
	var err error
	// 1. Read the Tanzu config (config arg)
	// TODO(joshrosso): if the struct comes in, anything we need to do here?

	// 2. Download and Read the TKR
	log.Event("\\U+2692", " Processing TanzuKubernetesRelease (TKR)\n")
	err = getTkrBom("projects.registry.vmware.com/tkg/tkr-bom:v1.21.2_vmware.1-tkg.1")
	if err != nil {
		return fmt.Errorf("failed getting TKR BOM. Error: %s", err.Error())
	}

	t.bom, err = parseTKRBom("tkr-bom-v1.21.2+vmware.1-tkg.1.yaml")
	if err != nil {
		return fmt.Errorf("failed parsing TKR BOM. Error: %s", err.Error())
	}

	// 3. Resolve all required images
	// base image
	log.Event("\\U+1F5BC", " Selected base image\n")
	log.Style(outputIndent, logger.ColorLightGrey).Infof("%s\n", t.bom.GetTKRNodeImage())
	// core package repository
	log.Event("\\U+1F4E6", "Selected core package repository\n")
	log.Style(outputIndent, logger.ColorLightGrey).Infof("%s\n", t.bom.GetTKRCoreRepoBundlePath())
	// core user package repositories
	log.Event("\\U+1F4E6", "Selected additional package repositories\n")
	for _, additionalRepo := range t.bom.GetAdditionalRepoBundlesPaths() {
		log.Style(outputIndent, logger.ColorLightGrey).Infof("%s\n", additionalRepo)
	}
	// kapp-controller
	err = resolveKappBundle(t)
	if err != nil {
		return fmt.Errorf("failed resolving kapp-controller bundle. Error: %s", err.Error())
	}
	log.Event("\\U+1F4E6", "Selected kapp-controller image bundle\n")
	log.Style(outputIndent, logger.ColorLightGrey).Infof("%s\n", t.kappControllerBundle.GetRegistryURL())

	// 4. Create the cluster
	log.Eventf("\\U+1F6F0", " Creating cluster %s\n", t.name)
	t.kcPath, err = runClusterCreate(t)
	if err != nil {
		return fmt.Errorf("failed to create cluster, Error: %s", err.Error())
	}
	log.Style(outputIndent, logger.ColorLightGrey).Info("To troubleshoot, use:\n")
	log.Style(outputIndent, logger.ColorLightGrey).Infof("kubectl ${COMMAND} --kubeconfig %s\n", t.kcPath)

	// 5. Install kapp-controller
	kc := kapp.New(t.kcPath)
	log.Event("\\U+1F4E7", "Installing kapp-controller\n")
	kappDeployment, err := installKappController(t, kc)
	if err != nil {
		return fmt.Errorf("failed to install kapp-controller, Error: %s", err.Error())
	}
	blockForKappStatus(kappDeployment, kc)

	// 6. Install package repositories
	pkgClient := packages.NewClient(t.kcPath)
	log.Event("\\U+1F4E7", "Installing package repositories\n")
	createdCoreRepo, err := createPackageRepo(pkgClient, tkgSysNamespace, tkgCoreRepoName, t.bom.GetTKRCoreRepoBundlePath())
	if err != nil {
		return fmt.Errorf("failed to install core package repo. Error: %s", err.Error())
	}
	for _, additionalRepo := range t.bom.GetAdditionalRepoBundlesPaths() {
		_, err = createPackageRepo(pkgClient, tkgGlobalPkgNamespace, tkgCoreRepoName, additionalRepo)
		if err != nil {
			return fmt.Errorf("failed to install adiditonal package repo. Error: %s", err.Error())
		}
	}
	blockForRepoStatus(createdCoreRepo, pkgClient)

	// 7. Install CNI
	err = installCNI(pkgClient, t)
	if err != nil {
		return fmt.Errorf("failed to install the CNI package. Error: %s", err.Error())
	}

	// 8. Update kubeconfig and context
	kubeConfigMgr := kubeconfig.NewManager()
	err = mergeKubeconfigAndSetContext(kubeConfigMgr, t)
	if err != nil {
		log.Warnf("Failed to merge kubeconfig and set your context. Cluster should still work! Error: %s", err)
	}

	// 8. Return
	log.Event("\\U+2705", "Cluster created\n")
	log.Eventf("\\U+1F3AE", "kubectl context set to %s\n\n", t.name)
	// provide user example commands to run
	log.Style(0, logger.ColorLightGrey).Infof("View available packages:\n")
	log.Style(outputIndent, logger.ColorLightGreen).Infof("tanzu package available list\n")
	log.Style(0, logger.ColorLightGrey).Infof("View running pods:\n")
	log.Style(outputIndent, logger.ColorLightGreen).Infof("kubectl get po -A\n")
	log.Style(0, logger.ColorLightGrey).Infof("Delete this cluster:\n")
	log.Style(outputIndent, logger.ColorLightGreen).Infof("tanzu local delete %s\n", t.name)
	return nil
}

// List lists the local clusters.
func (t *TanzuLocal) List() ([]TanzuCluster, error) {
	panic("implement me")
}

// Delete deletes a local cluster.
func (t *TanzuLocal) Delete(name string) error {
	panic("implement me")
}

// getTkgConfigDir returns the configuration directory used by tce.
func getTkgConfigDir() (path string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return path, fmt.Errorf("failed to resolve home dir. Error: %s", err.Error())
	}
	path = filepath.Join(home, configDir, tanzuConfigDir, tkgConfigDir)
	return path, nil
}

func getTkgLocalConfigDir() (path string, err error) {
	tkgConfigDir, err := getTkgConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(tkgConfigDir, localConfigDir), nil
}

func getLocalBomPath() (path string, err error) {
	tkgLocalConfigDir, err := getTkgLocalConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(tkgLocalConfigDir, bomDir), nil
}

func getTkrBom(registry string) error {
	bomPath, err := getLocalBomPath()
	if err != nil {
		return err
	}

	_, err = os.Stat(bomPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(bomPath, 0755)
		if err != nil {
			return err
		}
	}

	items, err := os.ReadDir(bomPath)
	if err != nil {
		return err
	}

	// If bom directory isn't empty, assume BOM already downloaded. return early
	if len(items) != 0 {
		return nil
	}

	bomImage, err := tkr.NewTkrImageReader(registry)
	if err != nil {
		return err
	}

	err = bomImage.DownloadImage()
	if err != nil {
		return err
	}

	// Move the file from the temporary directory into the tanzu config dir
	// TODO: these files names should be made configurable (from the LocalConfig ..?)
	err = os.Rename(
		filepath.Join(bomImage.GetDownloadPath(), "tkr-bom-v1.21.2+vmware.1-tkg.1.yaml"),
		filepath.Join(bomPath, "tkr-bom-v1.21.2+vmware.1-tkg.1.yaml"),
	)

	if err != nil {
		return err
	}

	return nil
}

func parseTKRBom(fileName string) (*tkr.TKRBom, error) {
	tkgBomPath, err := getLocalBomPath()
	if err != nil {
		return nil, err
	}

	bomPath := filepath.Join(tkgBomPath, fileName)

	bom, err := tkr.ReadTKRBom(bomPath)
	if err != nil {
		return nil, err
	}
	return bom, nil
}

func resolveKappBundle(t *TanzuLocal) error {
	var err error
	t.kappControllerBundle, err = t.bom.GetTKRKappImage()
	if err != nil {
		return err
	}
	return nil
}

func runClusterCreate(t *TanzuLocal) (string, error) {
	clusterManager := cluster.NewClusterManager()
	clusterCreateOpts := cluster.CreateOpts{
		Name:           t.name,
		KubeconfigPath: t.kcPath,
		// Config: TBD,
	}
	_, err := clusterManager.Create(&clusterCreateOpts)
	if err != nil {
		return "", err
	}
	// TODO(joshrosso): this should always return a path to the kubeconfig
	//                  right now it's a little silly because it's pullin
	//                  already known to the struct
	return t.kcPath, nil
}

func installKappController(t *TanzuLocal, kc kapp.KappManager) (*v1.Deployment, error) {
	err := t.kappControllerBundle.DownloadBundleImage()
	if err != nil {
		return nil, err
	}

	err = t.kappControllerBundle.AddYttYamlValuesBytes([]byte(kapp.DefaultKappValues))
	if err != nil {
		return nil, err
	}
	t.kappControllerBundle.SetRelativeConfigPath("./config")
	kappBytes, err := t.kappControllerBundle.RenderYaml()
	if err != nil {
		return nil, err
	}

	kappControllerCreated, err := kc.Install(kapp.KappInstallOpts{MergedManifests: kappBytes[0]})
	if err != nil {
		return nil, err
	}

	return kappControllerCreated, nil
}

func blockForKappStatus(kappDeployment *v1.Deployment, kc kapp.KappManager) {
	// Wait for kapp-controller be running; report status
	for si := 1; si < 5; si++ {
		kappState := kc.Status(kappDeployment.Namespace, kappDeployment.Name)
		log.Style(outputIndent, logger.ColorLightGrey).Progressf(si, "kapp-controller status: %s", kappState)
		if kappState == "Running" {
			log.Style(outputIndent, logger.ColorLightGrey).Progressf(0, "kapp-controller status: %s", kappState)
			break
		}
		if si == maxProgressLength {
			si = 1
		}
		time.Sleep(1 * time.Second)
	}
}

func createPackageRepo(pkgClient packages.PackageManager, ns, name, url string) (*v1alpha1.PackageRepository, error) {
	createdRepo, err := pkgClient.CreatePackageRepo(ns, name, url)
	if err != nil {
		return nil, err
	}
	return createdRepo, nil
}

func blockForRepoStatus(repo *v1alpha1.PackageRepository, pkgClient packages.PackageManager) {
	for si := 1; si < 5; si++ {
		status, err := pkgClient.GetRepositoryStatus(repo.Namespace, repo.Name)
		if err != nil {
			log.Errorf("failed to check kapp-controller status: %s", err.Error())
			return
		}
		log.Style(outputIndent, logger.ColorLightGrey).Progressf(si, "Core package repo status: %s", status)
		if status == "Reconcile succeeded" {
			log.Style(outputIndent, logger.ColorLightGrey).Progressf(0, "Core package repo status: %s", status)
			break
		}
		if si == maxProgressLength {
			si = 1
		}
		time.Sleep(1 * time.Second)
	}
}

// TODO(joshrosso) this function is a mess, but waiting on some stuff to happen in other packages
func installCNI(pkgClient packages.PackageManager, t *TanzuLocal) error { //nolint:unparam
	// install CNI (TODO(joshrosso): needs to support multiple CNIs
	rootSvcAcct, err := pkgClient.CreateRootServiceAccount(tkgSysNamespace, tkgSvcAcctName)
	if err != nil {
		log.Errorf("failed to create service account: %s\n", err.Error())
		return err
	}

	// TODO(joshrosso): entirely a workaround until we have better plumbing.
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
	_, err = pkgClient.CreatePackageInstall(&cniInstallOpts)
	if err != nil {
		return err
	}

	return nil
}

func mergeKubeconfigAndSetContext(mgr kubeconfig.KubeConfigMgr, t *TanzuLocal) error {
	err := mgr.MergeToDefaultConfig(t.kcPath)
	if err != nil {
		log.Errorf("Failed to merge kubeconfig: %s\n", err.Error())
		return nil
	}
	// TODO(joshrosso): we need to resolve this by introspecting the known kubeconfig
	// 					we cannot assume this syntax will work!
	kubeContextName := fmt.Sprintf("%s-%s", "kind", t.name)
	err = mgr.SetCurrentContext(kubeContextName)
	if err != nil {
		return err
	}

	return nil
}
