// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package tanzu is responsible for orchestrating the various components that result in
// a standalone tanzu cluster creation.
package tanzu

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster/config"

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
	configFileName        = "config.yaml"
	tanzuConfigDir        = "tanzu"
	tkgConfigDir          = "tkg"
	standaloneConfigDir   = "standalone"
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

// TanzuStandalone contains information about a standalone Tanzu cluster.
//nolint:golint
type TanzuStandalone struct {
	// kcPath               string
	// config               *StandaloneClusterConfig
	bom                  *tkr.TKRBom
	kappControllerBundle tkr.TkrImageReader
	selectedCNIPkg       *CNIPackage
	config               *config.StandaloneClusterConfig
	clusterDirectory     string
}

type CNIPackage struct {
	fqPkgName  string
	pkgVersion string
}

//nolint:golint
type TanzuMgr interface {
	// Deploy orchestrates all the required steps in order to create a standalone Tanzu cluster. This can involve
	// cluster creation, kapp-controller installation, CNI installation, and more. The steps that are taken
	// depend on the configuration passed into Deploy. If something goes wrong during deploy, an error is
	// returned.
	Deploy(scConfig *config.StandaloneClusterConfig) error
	// List retrieves all known tanzu clusters are returns a list of them. If it's unable to interact with the
	// underlying cluster provider, it returns an error.
	List() ([]TanzuCluster, error)
	// Delete takes a cluster name and removes the cluster from the underlying cluster provider. If it is unable
	// to communicate with the underlying cluster provider, it returns an error.
	Delete(name string) error
}

// New returns a TanzuMgr for interacting with standalone clusters. It is implemented by TanzuStandalone.
func New(parentLogger logger.Logger) TanzuMgr {
	log = parentLogger
	return &TanzuStandalone{}
}

// validateConfiguration makes sure the configuration is valid, returning an
// error if there is an issue.
func validateConfiguration(scConfig *config.StandaloneClusterConfig) error {
	if scConfig.TkrLocation == "" {
		return fmt.Errorf("Tanzu Kubernetes Release (TKR) not specified.") //nolint:golint,stylecheck
	}
	if scConfig.ClusterName == "" {
		return fmt.Errorf("cluster name is required")
	}

	if scConfig.Provider == "" {
		// Should have been validated earlier, but not an error. We can just
		// default it to kind.
		scConfig.Provider = cluster.KindClusterManagerProvider
	}

	return nil
}

// Deploy deploys a new cluster.
//nolint:funlen,gocyclo
func (t *TanzuStandalone) Deploy(scConfig *config.StandaloneClusterConfig) error {
	var err error
	// 1. Validate the configuration
	if err := validateConfiguration(scConfig); err != nil {
		return err
	}
	t.config = scConfig

	// 2. Download and Read the TKR
	log.Event("\\U+2692", " Resolving Tanzu Kubernetes Release (TKR)\n")
	bomFileName, err := getTkrBom(scConfig.TkrLocation)
	if err != nil {
		return fmt.Errorf("failed getting TKR BOM. Error: %s", err.Error())
	}

	t.clusterDirectory, err = createClusterDirectory(t.config.ClusterName)
	if err != nil {
		return err
	}
	configFp := filepath.Join(t.clusterDirectory, configFileName)
	err = config.RenderConfigToFile(configFp, t.config)
	if err != nil {
		return err
	}
	log.Event("\\U+1F4C1", "Created cluster directory\n")
	log.Style(outputIndent, logger.ColorLightGrey).Infof("Rendered Config: %s\n", configFp)

	// TODO(joshrosso): this file should be init'd and written to via loggers
	bootstrapLogsFp := filepath.Join(t.clusterDirectory, "bootstrap.log")
	log.Style(outputIndent, logger.ColorLightGrey).Infof("Bootstrap Logs: %s\n", bootstrapLogsFp)

	log.Event("\\U+2692", " Processing Tanzu Kubernetes Release\n")
	t.bom, err = parseTKRBom(bomFileName)
	if err != nil {
		return fmt.Errorf("failed parsing TKR BOM. Error: %s", err.Error())
	}

	// 3. Resolve all required images
	// base image
	log.Event("\\U+1F5BC", " Selected base image\n")
	log.Style(outputIndent, logger.ColorLightGrey).Infof("%s\n", t.bom.GetTKRNodeImage())
	scConfig.NodeImage = t.bom.GetTKRNodeImage()

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
	log.Eventf("\\U+1F6F0", " Creating cluster %s\n", scConfig.ClusterName)
	createdCluster, err := runClusterCreate(scConfig)
	if err != nil {
		return fmt.Errorf("failed to create cluster, Error: %s", err.Error())
	}

	kcBytes := createdCluster.Kubeconfig
	log.Style(outputIndent, logger.ColorLightGrey).Info("To troubleshoot, use:\n")
	log.Style(outputIndent, logger.ColorLightGrey).Infof("kubectl ${COMMAND} --kubeconfig %s\n", scConfig.KubeconfigPath)

	// 5. Install kapp-controller
	kc, err := kapp.New(kcBytes)
	if err != nil {
		return fmt.Errorf("failed to create kapp-controller manager, Error: %s", err.Error())
	}

	log.Event("\\U+1F4E7", "Installing kapp-controller\n")
	kappDeployment, err := installKappController(t, kc)
	if err != nil {
		return fmt.Errorf("failed to install kapp-controller, Error: %s", err.Error())
	}
	blockForKappStatus(kappDeployment, kc)

	// 6. Install package repositories
	pkgClient := packages.NewClient(kcBytes)
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
	log.Event("\\U+1F4E6", "Installing CNI\n")
	t.selectedCNIPkg, err = resolveCNI(pkgClient, t.config.Cni)
	if err != nil {
		return fmt.Errorf("failed to resolve a CNI package. Error: %s", err.Error())
	}
	log.Style(outputIndent, logger.ColorLightGrey).Infof("%s:%s\n", t.selectedCNIPkg.fqPkgName, t.selectedCNIPkg.pkgVersion)
	err = installCNI(pkgClient, t)
	if err != nil {
		return fmt.Errorf("failed to install the CNI package. Error: %s", err.Error())
	}

	// 8. Update kubeconfig and context
	kubeConfigMgr := kubeconfig.NewManager()
	err = mergeKubeconfigAndSetContext(kubeConfigMgr, scConfig.KubeconfigPath, scConfig.ClusterName)
	if err != nil {
		log.Warnf("Failed to merge kubeconfig and set your context. Cluster should still work! Error: %s", err)
	}

	// 8. Return
	log.Event("\\U+2705", "Cluster created\n")
	log.Eventf("\\U+1F3AE", "kubectl context set to %s\n\n", scConfig.ClusterName)
	// provide user example commands to run
	log.Style(0, logger.ColorLightGrey).Infof("View available packages:\n")
	log.Style(outputIndent, logger.ColorLightGreen).Infof("tanzu package available list\n")
	log.Style(0, logger.ColorLightGrey).Infof("View running pods:\n")
	log.Style(outputIndent, logger.ColorLightGreen).Infof("kubectl get po -A\n")
	log.Style(0, logger.ColorLightGrey).Infof("Delete this cluster:\n")
	log.Style(outputIndent, logger.ColorLightGreen).Infof("tanzu standalone delete %s\n", scConfig.ClusterName)
	return nil
}

// List lists the standalone clusters.
func (t *TanzuStandalone) List() ([]TanzuCluster, error) {
	var clusters []TanzuCluster

	configDir, err := getTkgStandaloneConfigDir()
	if err != nil {
		return nil, err
	}

	dirs, err := os.ReadDir(configDir)
	if err != nil {
		return nil, err
	}

	// 1. enter each directory in the tanzu standalone config directory,
	// 2. assess if there is a config.yaml file which was generated during the `create` command
	// 3. render the config file to a config struct and add the named cluster to the tanzu clusters
	for _, dir := range dirs {
		info, err := os.Stat(filepath.Join(configDir, dir.Name()))
		if !info.IsDir() || err != nil {
			continue
		}

		configFilePath := filepath.Join(configDir, dir.Name(), configFileName)
		_, err = os.Stat(configFilePath)
		if os.IsNotExist(err) {
			continue
		}

		scc, err := config.RenderFileToConfig(configFilePath)
		if err != nil {
			return nil, err
		}

		clusters = append(clusters, TanzuCluster{
			Name: scc.ClusterName,
		})
	}

	return clusters, nil
}

// Delete deletes a standalone cluster.
func (t *TanzuStandalone) Delete(name string) error {
	var err error
	t.clusterDirectory, err = resolveClusterDir(name)
	if err != nil {
		return err
	}
	configPath, err := resolveClusterConfig(name)
	if err != nil {
		return err
	}
	t.config, err = config.RenderFileToConfig(configPath)
	if err != nil {
		return err
	}

	cm := cluster.NewClusterManager(t.config)

	err = cm.Delete(t.config)
	if err != nil {
		return err
	}

	err = os.RemoveAll(t.clusterDirectory)
	if err != nil {
		log.Warnf("Cluster deleted but failed to remove config %s. Be sure to manually delete.", t.clusterDirectory)
	}

	return nil
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

func getTkgStandaloneConfigDir() (path string, err error) {
	tkgConfigDir, err := getTkgConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(tkgConfigDir, standaloneConfigDir), nil
}

func getStandaloneBomPath() (path string, err error) {
	tkgStandaloneConfigDir, err := getTkgStandaloneConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(tkgStandaloneConfigDir, bomDir), nil
}

func buildFilesystemSafeBomName(bomFileName string) (path string) {
	var sb strings.Builder
	for _, char := range bomFileName {
		if char == '/' || char == ':' {
			sb.WriteRune('_')
		}

		if char == '.' || char == '_' || char == '-' {
			sb.WriteRune(char)
		}

		if char >= 'a' && char <= 'z' {
			sb.WriteRune(char)
		}

		if char >= 'A' && char <= 'Z' {
			sb.WriteRune(char)
		}

		if char >= '0' && char <= '9' {
			sb.WriteRune(char)
		}
	}

	return sb.String()
}

func resolveClusterDir(clusterName string) (string, error) {
	scd, err := getTkgStandaloneConfigDir()
	if err != nil {
		return "", err
	}

	// determine if directory pre-exists
	fp := filepath.Join(scd, clusterName)
	_, err = os.ReadDir(fp)

	if os.IsNotExist(err) {
		return "", fmt.Errorf("failed to locate the cluster's config file at %s", fp)
	}
	return fp, nil
}

func resolveClusterConfig(clusterName string) (string, error) {
	scd, err := getTkgStandaloneConfigDir()
	if err != nil {
		return "", err
	}

	// determine if directory pre-exists
	fp := filepath.Join(scd, clusterName)
	files, err := os.ReadDir(fp)

	if os.IsNotExist(err) {
		return "", fmt.Errorf("failed to locate the cluster's config file at %s", fp)
	}

	var resolvedConfigFile string
	for _, file := range files {
		if file.Name() == configFileName {
			resolvedConfigFile = file.Name()
		}
	}

	expectFp := filepath.Join(fp, configFileName)
	if resolvedConfigFile == "" {
		return "", fmt.Errorf("failed to locate a config file at %s", expectFp)
	}

	return filepath.Join(fp, resolvedConfigFile), nil
}

func createClusterDirectory(clusterName string) (string, error) {
	scd, err := getTkgStandaloneConfigDir()
	if err != nil {
		return "", err
	}
	// determine if directory pre-exists
	fp := filepath.Join(scd, clusterName)
	_, err = os.ReadDir(fp)

	// if it does not exist, which is expected, create it
	if !os.IsNotExist(err) {
		return "", fmt.Errorf("directory %s already exists, this cluster must be deleted before proceeding", fp)
	}

	err = os.MkdirAll(fp, 0755)
	if err != nil {
		return "", err
	}

	return fp, nil
}

func getTkrBom(registry string) (string, error) {
	log.Style(outputIndent, logger.ColorLightGrey).Infof("%s\n", registry)
	expectedBomName := buildFilesystemSafeBomName(registry)

	bomPath, err := getStandaloneBomPath()
	if err != nil {
		return "", fmt.Errorf("failed to get tanzu stanadlone bom path: %s", err)
	}

	_, err = os.Stat(bomPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(bomPath, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to make new tanzu standalone bom config directories %s", err)
		}
	}

	items, err := os.ReadDir(bomPath)
	if err != nil {
		return "", fmt.Errorf("failed to read tanzu standalone bom directories: %s", err)
	}

	// if the expected bom is already in the config directory, don't download it again. return early
	for _, file := range items {
		if file.Name() == expectedBomName {
			log.Style(outputIndent, logger.ColorLightGrey).Infof("TKR exists at %s\n", filepath.Join(bomPath, file.Name()))
			return file.Name(), nil
		}
	}

	bomImage, err := tkr.NewTkrImageReader(registry)
	if err != nil {
		return "", fmt.Errorf("failed to create new TkrImageReader: %s", err)
	}

	log.Style(outputIndent, logger.ColorLightGrey).Infof("Downloading to %s\n", filepath.Join(bomPath, expectedBomName))
	err = bomImage.DownloadImage()
	if err != nil {
		return "", fmt.Errorf("failed to download tkr image: %s", err)
	}

	downloadedBomFiles, err := os.ReadDir(bomImage.GetDownloadPath())
	if err != nil {
		return "", fmt.Errorf("failed to read downloaded tkr bom files: %s", err)
	}

	// if there is more than 1 file in the downloaded image, fail
	// this is a bit redundant since imgpkg librariers should fail if the image is a bundle with multiple files
	if len(downloadedBomFiles) != 1 {
		return "", fmt.Errorf("more than one file found in TKR bom image. Expected 1 file: %s", bomImage.GetDownloadPath())
	}

	downloadedBomFile, err := os.Open(filepath.Join(bomImage.GetDownloadPath(), downloadedBomFiles[0].Name()))
	if err != nil {
		return "", fmt.Errorf("could not open downloaded tkr bom file: %s", err)
	}
	defer downloadedBomFile.Close()

	newBomFile, err := os.Create(filepath.Join(bomPath, expectedBomName))
	if err != nil {
		return "", fmt.Errorf("could not create tanzu standalone bom tkr file: %s", err)
	}
	defer newBomFile.Close()

	_, err = io.Copy(newBomFile, downloadedBomFile)
	if err != nil {
		return "", fmt.Errorf("could not copy file contents: %s", err)
	}

	return expectedBomName, nil
}

func parseTKRBom(fileName string) (*tkr.TKRBom, error) {
	tkgBomPath, err := getStandaloneBomPath()
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

func resolveKappBundle(t *TanzuStandalone) error {
	var err error
	t.kappControllerBundle, err = t.bom.GetTKRKappImage()
	if err != nil {
		return err
	}
	return nil
}

func runClusterCreate(scConfig *config.StandaloneClusterConfig) (*cluster.KubernetesCluster, error) {
	if scConfig.KubeconfigPath == "" {
		clusterDir, err := resolveClusterDir(scConfig.ClusterName)
		if err != nil {
			return nil, err
		}
		scConfig.KubeconfigPath = filepath.Join(clusterDir, "kube.conf")
	}

	clusterManager := cluster.NewClusterManager(scConfig)
	kc, err := clusterManager.Create(scConfig)
	if err != nil {
		return nil, err
	}
	return kc, nil
}

func installKappController(t *TanzuStandalone, kc kapp.KappManager) (*v1.Deployment, error) {
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
func installCNI(pkgClient packages.PackageManager, t *TanzuStandalone) error {
	// install CNI (TODO(joshrosso): needs to support multiple CNIs
	rootSvcAcct, err := pkgClient.CreateRootServiceAccount(tkgSysNamespace, tkgSvcAcctName)
	if err != nil {
		log.Errorf("failed to create service account: %s\n", err.Error())
		return err
	}
	var valueData string

	if strings.Contains(t.config.Cni, "antrea") {
		// TODO(joshrosso): entirely a workaround until we have better plumbing.
		valueData = `---
infraProvider: docker
`
	}

	cniInstallOpts := packages.PackageInstallOpts{
		Namespace:      tkgSysNamespace,
		InstallName:    "cni",
		FqPkgName:      t.selectedCNIPkg.fqPkgName,
		Version:        t.selectedCNIPkg.pkgVersion,
		Configuration:  []byte(valueData),
		ServiceAccount: rootSvcAcct.Name,
	}
	_, err = pkgClient.CreatePackageInstall(&cniInstallOpts)
	if err != nil {
		return err
	}

	return nil
}

func mergeKubeconfigAndSetContext(mgr kubeconfig.KubeConfigMgr, kcPath, clusterName string) error {
	err := mgr.MergeToDefaultConfig(kcPath)
	if err != nil {
		log.Errorf("Failed to merge kubeconfig: %s\n", err.Error())
		return nil
	}
	// TODO(joshrosso): we need to resolve this by introspecting the known kubeconfig
	// 					we cannot assume this syntax will work!
	kubeContextName := fmt.Sprintf("%s-%s", "kind", clusterName)
	err = mgr.SetCurrentContext(kubeContextName)
	if err != nil {
		return err
	}

	return nil
}

// resolveCNI determines which CNI package to use. It expects to be passed a fully qualified package name
// except for special known CNI values such as antrea or calico.
func resolveCNI(mgr packages.PackageManager, cniName string) (*CNIPackage, error) {
	pkgs, err := mgr.ListPackagesInNamespace(tkgSysNamespace)
	if err != nil {
		return nil, err
	}

	cniPkg := &CNIPackage{}
	for _, pkg := range pkgs {
		if strings.HasPrefix(pkg.Spec.RefName, cniName) {
			cniPkg.fqPkgName = pkg.Spec.RefName
			cniPkg.pkgVersion = pkg.Spec.Version
		}
	}
	if err != nil {
		return nil, err
	}
	if cniPkg.fqPkgName == "" {
		return nil, fmt.Errorf("no package was resolved for CNI choice %s", cniName)
	}

	return cniPkg, nil
}
