// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package tanzu is responsible for orchestrating the various packages that satisfy unmanaged
// operations such as create, configure, list, and delete. This package is meant to be the API
// entrypoint for those calling unmanaged programmatically.
package tanzu

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/config"

	v1 "k8s.io/api/apps/v1"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/cluster"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/kapp"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/kubeconfig"
	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/packages"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/tkr"
)

//nolint
const (
	configDir             = ".config"
	configFileName        = "config.yaml"
	bootstrapLogName      = "bootstrap.log"
	tanzuConfigDir        = "tanzu"
	tkgConfigDir          = "tkg"
	unmanagedConfigDir    = "unmanaged"
	bomDir                = "bom"
	tkgSysNamespace       = "tkg-system"
	tkgSvcAcctName        = "core-pkgs"
	tkgCoreRepoName       = "tkg-core-repository"
	tkgGlobalPkgNamespace = "tanzu-package-repo-global"
	tceRepoName           = "community-repository"
	tceRepoURL            = "projects.registry.vmware.com/tce/main:0.9.1"
	outputIndent          = 3
	maxProgressLength     = 4
)

// TODO(joshrosso): global logger for the package. This is kind gross, but really convenient.
var log logger.Logger

// TanzuCluster contains information about a cluster.
//nolint:golint
type TanzuCluster struct {
	Name     string
	Provider string
}

// TanzuUnmanaged contains information about an unmanaged Tanzu cluster.
//nolint:golint
type TanzuUnmanaged struct {
	bom                  *tkr.TKRBom
	kappControllerBundle tkr.TkrImageReader
	selectedCNIPkg       *CNIPackage
	config               *config.UnmanagedClusterConfig
	clusterDirectory     string
}

type CNIPackage struct {
	fqPkgName  string
	pkgVersion string
}

//nolint:golint
type TanzuMgr interface {
	// Deploy orchestrates all the required steps in order to create an unmanaged Tanzu cluster. This can involve
	// cluster creation, kapp-controller installation, CNI installation, and more. The steps that are taken
	// depend on the configuration passed into Deploy. If something goes wrong during deploy, an error is
	// returned.
	Deploy(scConfig *config.UnmanagedClusterConfig) error
	// List retrieves all known tanzu clusters are returns a list of them. If it's unable to interact with the
	// underlying cluster provider, it returns an error.
	List() ([]TanzuCluster, error)
	// Delete takes a cluster name and removes the cluster from the underlying cluster provider. If it is unable
	// to communicate with the underlying cluster provider, it returns an error.
	Delete(name string) error
}

// New returns a TanzuMgr for interacting with unmanaged clusters. It is implemented by TanzuUnmanaged.
func New(parentLogger logger.Logger) TanzuMgr {
	log = parentLogger
	return &TanzuUnmanaged{}
}

// validateConfiguration makes sure the configuration is valid, returning an
// error if there is an issue.
func validateConfiguration(scConfig *config.UnmanagedClusterConfig) error {
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
func (t *TanzuUnmanaged) Deploy(scConfig *config.UnmanagedClusterConfig) error {
	var err error

	// 1. Validate the configuration
	if err := validateConfiguration(scConfig); err != nil {
		return err
	}
	t.config = scConfig

	t.clusterDirectory, err = createClusterDirectory(t.config.ClusterName)
	if err != nil {
		return err
	}

	// Configure the logger to capture all bootstrap activity
	bootstrapLogsFp := filepath.Join(t.clusterDirectory, "bootstrap.log")
	log.AddLogFile(bootstrapLogsFp)
	log.Event(logger.FolderEmoji, "Created cluster directory\n")

	// 2. Download and Read the TKR
	log.Event(logger.WrenchEmoji, "Resolving Tanzu Kubernetes Release (TKR)\n")
	bomFileName, err := getTkrBom(scConfig.TkrLocation)
	if err != nil {
		return fmt.Errorf("failed getting TKR BOM. Error: %s", err.Error())
	}
	configFp := filepath.Join(t.clusterDirectory, configFileName)
	err = config.RenderConfigToFile(configFp, t.config)
	if err != nil {
		return err
	}
	log.Style(outputIndent, logger.ColorNone).Infof("Rendered Config: %s\n", configFp)
	log.Style(outputIndent, logger.ColorNone).Infof("Bootstrap Logs: %s\n", bootstrapLogsFp)

	log.Event(logger.WrenchEmoji, "Processing Tanzu Kubernetes Release\n")
	t.bom, err = parseTKRBom(bomFileName)
	if err != nil {
		return fmt.Errorf("failed parsing TKR BOM. Error: %s", err.Error())
	}

	// 3. Resolve all required images
	// base image
	log.Event(logger.PictureEmoji, "Selected base image\n")
	log.Style(outputIndent, logger.ColorNone).Infof("%s\n", t.bom.GetTKRNodeImage())
	scConfig.NodeImage = t.bom.GetTKRNodeImage()

	// core package repository
	log.Event(logger.PackageEmoji, "Selected core package repository\n")
	log.Style(outputIndent, logger.ColorNone).Infof("%s\n", t.bom.GetTKRCoreRepoBundlePath())
	// core user package repositories
	log.Event(logger.PackageEmoji, "Selected additional package repositories\n")
	for _, additionalRepo := range t.bom.GetAdditionalRepoBundlesPaths() {
		log.Style(outputIndent, logger.ColorNone).Infof("%s\n", additionalRepo)
	}
	// kapp-controller
	err = resolveKappBundle(t)
	if err != nil {
		return fmt.Errorf("failed resolving kapp-controller bundle. Error: %s", err.Error())
	}
	log.Event(logger.PackageEmoji, "Selected kapp-controller image bundle\n")
	log.Style(outputIndent, logger.ColorNone).Infof("%s\n", t.kappControllerBundle.GetRegistryURL())

	// 4. Create the cluster
	log.Eventf(logger.RocketEmoji, "Creating cluster %s\n", scConfig.ClusterName)
	createdCluster, err := runClusterCreate(scConfig)
	if err != nil {
		return fmt.Errorf("failed to create cluster, Error: %s", err.Error())
	}

	kcBytes := createdCluster.Kubeconfig
	log.Style(outputIndent, logger.ColorNone).Info("To troubleshoot, use:\n")
	log.Style(outputIndent, logger.ColorNone).Infof("kubectl ${COMMAND} --kubeconfig %s\n", scConfig.KubeconfigPath)

	// 5. Install kapp-controller
	kc, err := kapp.New(kcBytes)
	if err != nil {
		return fmt.Errorf("failed to create kapp-controller manager, Error: %s", err.Error())
	}

	log.Event(logger.EnvelopeEmoji, "Installing kapp-controller\n")
	kappDeployment, err := installKappController(t, kc)
	if err != nil {
		return fmt.Errorf("failed to install kapp-controller, Error: %s", err.Error())
	}
	blockForKappStatus(kappDeployment, kc)

	// 6. Install package repositories
	pkgClient := packages.NewClient(kcBytes)
	log.Event(logger.EnvelopeEmoji, "Installing package repositories\n")
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
	log.Event(logger.GlobeEmoji, "Installing CNI\n")
	t.selectedCNIPkg, err = resolveCNI(pkgClient, t.config.Cni)
	if err != nil {
		return fmt.Errorf("failed to resolve a CNI package. Error: %s", err.Error())
	}
	log.Style(outputIndent, logger.ColorNone).Infof("%s:%s\n", t.selectedCNIPkg.fqPkgName, t.selectedCNIPkg.pkgVersion)
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
	log.Event(logger.GreenCheckEmoji, "Cluster created\n")
	log.Eventf(logger.ControllerEmoji, "kubectl context set to %s\n\n", scConfig.ClusterName)
	// provide user example commands to run
	log.Style(0, logger.ColorNone).Infof("View available packages:\n")
	log.Style(outputIndent, logger.ColorLightGreen).Infof("tanzu package available list\n")
	log.Style(0, logger.ColorNone).Infof("View running pods:\n")
	log.Style(outputIndent, logger.ColorLightGreen).Infof("kubectl get po -A\n")
	log.Style(0, logger.ColorNone).Infof("Delete this cluster:\n")
	log.Style(outputIndent, logger.ColorLightGreen).Infof("tanzu unmanaged delete %s\n", scConfig.ClusterName)
	return nil
}

// List lists the unmanaged clusters.
func (t *TanzuUnmanaged) List() ([]TanzuCluster, error) {
	var clusters []TanzuCluster

	configDir, err := getTkgUnmanagedConfigDir()
	if err != nil {
		return nil, err
	}

	dirs, err := os.ReadDir(configDir)
	if err != nil {
		return nil, err
	}

	// 1. enter each directory in the tanzu unmanaged config directory,
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
			Name:     scc.ClusterName,
			Provider: scc.Provider,
		})
	}

	return clusters, nil
}

// Delete deletes an unmanaged cluster.
func (t *TanzuUnmanaged) Delete(name string) error {
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

func getTkgUnmanagedConfigDir() (path string, err error) {
	tkgConfigDir, err := getTkgConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(tkgConfigDir, unmanagedConfigDir), nil
}

func getUnmanagedBomPath() (path string, err error) {
	tkgUnmanagedConfigDir, err := getTkgUnmanagedConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(tkgUnmanagedConfigDir, bomDir), nil
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
	scd, err := getTkgUnmanagedConfigDir()
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
	scd, err := getTkgUnmanagedConfigDir()
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
	scd, err := getTkgUnmanagedConfigDir()
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
	log.Style(outputIndent, logger.ColorNone).Infof("%s\n", registry)
	expectedBomName := buildFilesystemSafeBomName(registry)

	bomPath, err := getUnmanagedBomPath()
	if err != nil {
		return "", fmt.Errorf("failed to get tanzu stanadlone bom path: %s", err)
	}

	_, err = os.Stat(bomPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(bomPath, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to make new tanzu unmanaged bom config directories %s", err)
		}
	}

	items, err := os.ReadDir(bomPath)
	if err != nil {
		return "", fmt.Errorf("failed to read tanzu unmanaged bom directories: %s", err)
	}

	// if the expected bom is already in the config directory, don't download it again. return early
	for _, file := range items {
		if file.Name() == expectedBomName {
			log.Style(outputIndent, logger.ColorNone).Infof("TKR exists at %s\n", filepath.Join(bomPath, file.Name()))
			return file.Name(), nil
		}
	}

	bomImage, err := tkr.NewTkrImageReader(registry)
	if err != nil {
		return "", fmt.Errorf("failed to create new TkrImageReader: %s", err)
	}

	err = blockForBomImage(bomImage, bomPath, expectedBomName)
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
		return "", fmt.Errorf("could not create tanzu unmanaged bom tkr file: %s", err)
	}
	defer newBomFile.Close()

	_, err = io.Copy(newBomFile, downloadedBomFile)
	if err != nil {
		return "", fmt.Errorf("could not copy file contents: %s", err)
	}

	return expectedBomName, nil
}

func blockForBomImage(b tkr.TkrImageReader, bomPath, expectedBomName string) error {
	f := filepath.Join(bomPath, expectedBomName)

	// start a go routine to animate the downloading logs while the imgpkg libraries get the bom image
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		log.Style(outputIndent, logger.ColorNone).AnimateProgressWithOptions(
			logger.AnimatorWithContext(ctx),
			logger.AnimatorWithMaxLen(maxProgressLength),
			logger.AnimatorWithMessagef("Downloading to: %s", f),
		)
	}(ctx)

	// This will block and the go routine will continue to animate the logging
	err := b.DownloadImage()
	if err != nil {
		cancel()
		return err
	}

	// Once downloading is done, cancel the logging animation go routine and log completion
	cancel()
	log.Style(outputIndent, logger.ColorNone).ReplaceLinef("Downloaded to: %s", f)

	return nil
}

func parseTKRBom(fileName string) (*tkr.TKRBom, error) {
	tkgBomPath, err := getUnmanagedBomPath()
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

func resolveKappBundle(t *TanzuUnmanaged) error {
	var err error
	t.kappControllerBundle, err = t.bom.GetTKRKappImage()
	if err != nil {
		return err
	}
	return nil
}

func runClusterCreate(scConfig *config.UnmanagedClusterConfig) (*cluster.KubernetesCluster, error) {
	if scConfig.KubeconfigPath == "" {
		clusterDir, err := resolveClusterDir(scConfig.ClusterName)
		if err != nil {
			return nil, err
		}
		scConfig.KubeconfigPath = filepath.Join(clusterDir, "kube.conf")
	}

	clusterManager := cluster.NewClusterManager(scConfig)

	err := blockForPullingBaseImage(clusterManager, scConfig)
	if err != nil {
		return nil, err
	}

	kc, err := blockForClusterCreate(clusterManager, scConfig)
	if err != nil {
		return nil, err
	}

	return kc, nil
}

func blockForPullingBaseImage(cm cluster.ClusterManager, scConfig *config.UnmanagedClusterConfig) error {
	// start a go routine to animate the downloading logs while the docker exec gets the image
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		log.Style(outputIndent, logger.ColorNone).AnimateProgressWithOptions(
			logger.AnimatorWithContext(ctx),
			logger.AnimatorWithMaxLen(maxProgressLength),
			logger.AnimatorWithMessagef("Pulling base image"),
		)
	}(ctx)

	// This should block
	err := cm.Prepare(scConfig)
	if err != nil {
		cancel()
		return err
	}

	// once we're done, cancel the go routine animation and log final message
	cancel()
	log.Style(outputIndent, logger.ColorNone).ReplaceLinef("Base image downloaded")

	return nil
}

func blockForClusterCreate(cm cluster.ClusterManager, scConfig *config.UnmanagedClusterConfig) (*cluster.KubernetesCluster, error) {
	// start a go routine to animate the downloading logs while the docker exec gets the image
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		log.Style(outputIndent, logger.ColorNone).AnimateProgressWithOptions(
			logger.AnimatorWithContext(ctx),
			logger.AnimatorWithMaxLen(maxProgressLength),
			logger.AnimatorWithMessagef("Creating cluster"),
		)
	}(ctx)

	// This should block
	kc, err := cm.Create(scConfig)
	if err != nil {
		cancel()
		return nil, err
	}

	// Once done, cancel the go routine animations and log final message
	cancel()
	log.Style(outputIndent, logger.ColorNone).ReplaceLinef("Cluster created")

	return kc, nil
}

func installKappController(t *TanzuUnmanaged, kc kapp.KappManager) (*v1.Deployment, error) {
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
	// Create the parent context and fire a go routine to animate the logging progress
	ctx, cancel := context.WithCancel(context.Background())
	status := make(chan string, 1)
	go func(ctx context.Context) {
		log.Style(outputIndent, logger.ColorNone).AnimateProgressWithOptions(
			logger.AnimatorWithContext(ctx),
			logger.AnimatorWithMaxLen(maxProgressLength),
			logger.AnimatorWithMessagef("kapp-controller status: %s"),
			logger.AnimatorWithStatusChan(status),
		)
	}(ctx)

	// Wait for kapp-controller to be running; report status into the status channel
	// TODO: jpmcb - we need to handle timeouts, errors, etc for kapp controller not running
	for {
		kappState := kc.Status(kappDeployment.Namespace, kappDeployment.Name)
		status <- kappState
		if kappState == "Running" {
			cancel()
			log.Style(outputIndent, logger.ColorNone).ReplaceLinef("kapp-controller status: %s", kappState)
			break
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
	// Create the parent context and fire a go routine to animate the logging progress
	ctx, cancel := context.WithCancel(context.Background())
	status := make(chan string, 1)
	go func(ctx context.Context) {
		log.Style(outputIndent, logger.ColorNone).AnimateProgressWithOptions(
			logger.AnimatorWithContext(ctx),
			logger.AnimatorWithMaxLen(maxProgressLength),
			logger.AnimatorWithMessagef("Core package repo status: %s"),
			logger.AnimatorWithStatusChan(status),
		)
	}(ctx)

	// Wait for core packages to be running; report status into the status channel
	// TODO - jpmcb: we need to handle timeouts here
	for {
		pkgStatus, err := pkgClient.GetRepositoryStatus(repo.Namespace, repo.Name)
		status <- pkgStatus
		if err != nil {
			cancel()
			log.Errorf("failed to check package repository status: %s", err.Error())
			return
		}
		if pkgStatus == "Reconcile succeeded" {
			cancel()
			log.Style(outputIndent, logger.ColorNone).ReplaceLinef("Core package repo status: %s", pkgStatus)
			break
		}
		time.Sleep(1 * time.Second)
	}
}

// TODO(joshrosso) this function is a mess, but waiting on some stuff to happen in other packages
func installCNI(pkgClient packages.PackageManager, t *TanzuUnmanaged) error {
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
