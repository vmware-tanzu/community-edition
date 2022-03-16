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

	"github.com/fatih/color"
	v1 "k8s.io/api/apps/v1"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/cluster"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/config"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/kapp"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/kubeconfig"
	logger "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/log"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/packages"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/tkr"
)

//nolint
const (
	configFileName        = "config.yaml"
	bootstrapLogName      = "bootstrap.log"
	bomDir                = "bom"
	tkgSysNamespace       = "tkg-system"
	tkgSvcAcctName        = "core-pkgs"
	tkgCoreRepoName       = "tkg-core-repository"
	tkgGlobalPkgNamespace = "tanzu-package-repo-global"
	tceRepoName           = "community-repository"
	tceRepoURL            = "projects.registry.vmware.com/tce/main:0.10.1"
	outputIndent          = 3
	maxProgressLength     = 4
)

// TODO(joshrosso): global logger for the package. This is kind gross, but really convenient.
var log logger.Logger

// value for CNI configuration which represents creating a cluster without a
// CNI.
const cniNoneName = "none"

// Cluster contains information about a cluster.
type Cluster struct {
	Name     string
	Provider string
}

// UnmanagedCluster contains information about an unmanaged Tanzu cluster.
type UnmanagedCluster struct {
	bom                  *tkr.Bom
	kappControllerBundle tkr.ImageReader
	selectedCNIPkg       *CNIPackage
	config               *config.UnmanagedClusterConfig
	clusterDirectory     string
}

type CNIPackage struct {
	fqPkgName  string
	pkgVersion string
}

type Manager interface {
	// Deploy orchestrates all the required steps in order to create an unmanaged Tanzu cluster. This can involve
	// cluster creation, kapp-controller installation, CNI installation, and more. The steps that are taken
	// depend on the configuration passed into Deploy.
	// If something goes wrong during deploy, an error and its corresponding exit code is returned.
	Deploy(scConfig *config.UnmanagedClusterConfig) (int, error)
	// List retrieves all known tanzu clusters are returns a list of them. If it's unable to interact with the
	// underlying cluster provider, it returns an error.
	List() ([]Cluster, error)
	// Delete takes a cluster name and removes the cluster from the underlying cluster provider. If it is unable
	// to communicate with the underlying cluster provider, it returns an error.
	Delete(name string) error
}

// New returns a TanzuMgr for interacting with unmanaged clusters. It is implemented by TanzuUnmanaged.
func New(parentLogger logger.Logger) Manager {
	log = parentLogger
	return &UnmanagedCluster{}
}

// validateConfiguration makes sure the configuration is valid, returning an
// error if there is an issue.
func validateConfiguration(scConfig *config.UnmanagedClusterConfig) error {
	if scConfig.TkrLocation == "" {
		return fmt.Errorf("Tanzu Kubernetes Release (TKR) not specified.") //nolint:revive,stylecheck
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
func (t *UnmanagedCluster) Deploy(scConfig *config.UnmanagedClusterConfig) (int, error) {
	var err error

	// 1. Validate the configuration
	if err := validateConfiguration(scConfig); err != nil {
		return InvalidConfig, err
	}
	t.config = scConfig

	t.clusterDirectory, err = createClusterDirectory(t.config.ClusterName)
	if err != nil {
		return ErrCreatingClusterDirs, err
	}

	// Configure the logger to capture all bootstrap activity
	bootstrapLogsFp := filepath.Join(t.clusterDirectory, "bootstrap.log")
	log.AddLogFile(bootstrapLogsFp)
	log.Event(logger.FolderEmoji, "Created cluster directory")

	// Log a warning if the user has given a ProviderConfiguration
	if len(scConfig.ProviderConfiguration) != 0 {
		log.Style(outputIndent, color.FgYellow).ReplaceLinef("Reading ProviderConfiguration from config file. All other provider specific configs may be ignored.")
	}

	// 2. Download and Read the TKR
	log.Event(logger.WrenchEmoji, "Resolving Tanzu Kubernetes Release (TKR)")
	bomFileName, err := getTkrBom(scConfig.TkrLocation)
	if err != nil {
		return ErrTkrBom, fmt.Errorf("failed getting TKR BOM. Error: %s", err.Error())
	}
	configFp := filepath.Join(t.clusterDirectory, configFileName)
	err = config.RenderConfigToFile(configFp, t.config)
	if err != nil {
		return ErrRenderingConfig, err
	}
	log.Style(outputIndent, color.Faint).Infof("Rendered Config: %s\n", configFp)
	log.Style(outputIndent, color.Faint).Infof("Bootstrap Logs: %s\n", bootstrapLogsFp)

	log.Event(logger.WrenchEmoji, "Processing Tanzu Kubernetes Release")
	t.bom, err = parseTKRBom(bomFileName)
	if err != nil {
		return ErrTkrBomParsing, fmt.Errorf("failed parsing TKR BOM. Error: %s", err.Error())
	}

	// 3. Resolve all required images
	// base image
	log.Event(logger.PictureEmoji, "Selected base image")
	log.Style(outputIndent, color.Faint).Infof("%s\n", t.bom.GetTKRNodeImage())
	scConfig.NodeImage = t.bom.GetTKRNodeImage()

	// core package repository
	log.Event(logger.PackageEmoji, "Selected core package repository")
	log.Style(outputIndent, color.Faint).Infof("%s\n", t.bom.GetTKRCoreRepoBundlePath())
	// core user package repositories
	log.Event(logger.PackageEmoji, "Selected additional package repositories")
	for _, additionalRepo := range t.bom.GetAdditionalRepoBundlesPaths() {
		log.Style(outputIndent, color.Faint).Infof("%s\n", additionalRepo)
	}
	// kapp-controller
	err = resolveKappBundle(t)
	if err != nil {
		return ErrKappBundleResolving, fmt.Errorf("failed resolving kapp-controller bundle. Error: %s", err.Error())
	}
	log.Event(logger.PackageEmoji, "Selected kapp-controller image bundle")
	log.Style(outputIndent, color.Faint).Infof("%s\n", t.kappControllerBundle.GetRegistryURL())

	// 4. Create the cluster
	var clusterToUse *cluster.KubernetesCluster

	if scConfig.ExistingClusterKubeconfig != "" {
		log.Eventf(logger.RocketEmoji, "Using existing cluster\n")
		clusterToUse, err = useExistingCluster(scConfig)
		if err != nil {
			return ErrExistingCluster, fmt.Errorf("failed to use existing cluster, Error: %s", err.Error())
		}
	} else {
		log.Eventf(logger.RocketEmoji, "Creating cluster %s\n", scConfig.ClusterName)
		clusterToUse, err = runClusterCreate(scConfig)
		if err != nil {
			return ErrCreateCluster, fmt.Errorf("failed to create cluster, Error: %s", err.Error())
		}
	}

	kcBytes := clusterToUse.Kubeconfig
	log.Style(outputIndent, color.Faint).Info("To troubleshoot, use:\n")
	log.Style(outputIndent, color.Faint).Infof("kubectl ${COMMAND} --kubeconfig %s\n", scConfig.KubeconfigPath)

	// 5. Install kapp-controller
	kc, err := kapp.New(kcBytes)
	if err != nil {
		return ErrKappInstall, fmt.Errorf("failed to create kapp-controller manager, Error: %s", err.Error())
	}

	log.Event(logger.EnvelopeEmoji, "Installing kapp-controller")
	kappDeployment, err := installKappController(t, kc)
	if err != nil {
		return ErrKappInstall, fmt.Errorf("failed to install kapp-controller, Error: %s", err.Error())
	}
	blockForKappStatus(kappDeployment, kc)

	// 6. Install package repositories
	pkgClient := packages.NewClient(kcBytes)
	log.Event(logger.EnvelopeEmoji, "Installing package repositories")
	createdCoreRepo, err := createPackageRepo(pkgClient, tkgSysNamespace, tkgCoreRepoName, t.bom.GetTKRCoreRepoBundlePath())
	if err != nil {
		return ErrCorePackageRepoInstall, fmt.Errorf("failed to install core package repo. Error: %s", err.Error())
	}
	for _, additionalRepo := range t.bom.GetAdditionalRepoBundlesPaths() {
		_, err = createPackageRepo(pkgClient, tkgGlobalPkgNamespace, tceRepoName, additionalRepo)
		if err != nil {
			return ErrOtherPackageRepoInstall, fmt.Errorf("failed to install adiditonal package repo. Error: %s", err.Error())
		}
	}
	blockForRepoStatus(createdCoreRepo, pkgClient)

	// 7. Install CNI
	// CNI plugins are installed as best effort. If no plugin is resolved in the
	// repository, no CNI is installed, yet the cluster will still run.
	log.Event(logger.GlobeEmoji, "Installing CNI")
	t.selectedCNIPkg, err = resolveCNI(pkgClient, t.config.Cni)

	// No CNI package was resolved to install
	if err != nil {
		log.Style(outputIndent, color.FgYellow).Warnf("No CNI installed: %s.\n", err)
	} else {
		// CNI package resolved, do install
		log.Style(outputIndent, color.Faint).Infof("%s:%s\n", t.selectedCNIPkg.fqPkgName, t.selectedCNIPkg.pkgVersion)
		err = installCNI(pkgClient, t)
		if err != nil {
			return ErrCniInstall, fmt.Errorf("failed to install the CNI package. Error: %s", err.Error())
		}
	}

	// 8. Update kubeconfig and context
	kubeConfigMgr := kubeconfig.NewManager()
	err = mergeKubeconfigAndSetContext(kubeConfigMgr, scConfig.KubeconfigPath, scConfig.ClusterName)
	if err != nil {
		log.Warnf("Failed to merge kubeconfig and set your context. Cluster should still work! Error: %s", err)
	}

	// 8. Return
	log.Event(logger.GreenCheckEmoji, "Cluster created")
	log.Eventf(logger.ControllerEmoji, "kubectl context set to %s\n\n", scConfig.ClusterName)
	// provide user example commands to run
	log.Infof("View available packages:\n")
	log.Style(outputIndent, color.FgGreen).Infof("tanzu package available list\n")
	log.Infof("View running pods:\n")
	log.Style(outputIndent, color.FgGreen).Infof("kubectl get po -A\n")
	log.Infof("Delete this cluster:\n")
	log.Style(outputIndent, color.FgGreen).Infof("tanzu unmanaged delete %s\n", scConfig.ClusterName)
	return Success, nil
}

// List lists the unmanaged clusters.
func (t *UnmanagedCluster) List() ([]Cluster, error) {
	var clusters []Cluster

	configDir, err := config.GetUnmanagedConfigPath()
	if err != nil {
		return nil, err
	}

	dirs, err := os.ReadDir(configDir)
	if err != nil {
		// If the config dir can't be read, it's possible no clusters have been
		// created yet. Or they don't have permissions to something in their own
		// home directory, which could indicate the config was copied in from
		// elsewhere. In either case, if we can't read the unmanaged cluster
		// info, we can just assume there are no clusters for them to see and
		// just return an empty list.
		return clusters, nil
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

		clusters = append(clusters, Cluster{
			Name:     scc.ClusterName,
			Provider: scc.Provider,
		})
	}

	return clusters, nil
}

// Delete deletes an unmanaged cluster.
func (t *UnmanagedCluster) Delete(name string) error {
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

func getUnmanagedBomPath() (path string, err error) {
	tkgUnmanagedConfigDir, err := config.GetUnmanagedConfigPath()
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
	scd, err := config.GetUnmanagedConfigPath()
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
	scd, err := config.GetUnmanagedConfigPath()
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
	scd, err := config.GetUnmanagedConfigPath()
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
	log.Style(outputIndent, color.Faint).Infof("%s\n", registry)
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
			log.Style(outputIndent, color.Faint).Infof("TKR exists at %s\n", filepath.Join(bomPath, file.Name()))
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

func blockForBomImage(b tkr.ImageReader, bomPath, expectedBomName string) error {
	f := filepath.Join(bomPath, expectedBomName)

	// start a go routine to animate the downloading logs while the imgpkg libraries get the bom image
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		log.Style(outputIndent, color.Reset).AnimateProgressWithOptions(
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
	log.Style(outputIndent, color.Faint).ReplaceLinef("Downloaded to: %s", f)

	return nil
}

func parseTKRBom(fileName string) (*tkr.Bom, error) {
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

func resolveKappBundle(t *UnmanagedCluster) error {
	var err error
	t.kappControllerBundle, err = t.bom.GetTKRKappImage()
	if err != nil {
		return err
	}
	return nil
}

func runClusterCreate(scConfig *config.UnmanagedClusterConfig) (*cluster.KubernetesCluster, error) {
	clusterDir, err := resolveClusterDir(scConfig.ClusterName)
	if err != nil {
		return nil, err
	}
	scConfig.KubeconfigPath = filepath.Join(clusterDir, "kube.conf")

	clusterManager := cluster.NewClusterManager(scConfig)

	for _, message := range clusterManager.ProviderNotify() {	
		log.Style(outputIndent, color.Faint).Info(message)
	}
	
	if !scConfig.SkipPreflightChecks {
		if issues := clusterManager.PreflightCheck(); issues != nil {
			return nil, fmt.Errorf("system checks detected issues, please resolve first: %v", issues)
		}
	}

	err = blockForPullingBaseImage(clusterManager, scConfig)
	if err != nil {
		return nil, err
	}

	kc, err := blockForClusterCreate(clusterManager, scConfig)
	if err != nil {
		return nil, err
	}

	return kc, nil
}

func useExistingCluster(scConfig *config.UnmanagedClusterConfig) (*cluster.KubernetesCluster, error) {
	scConfig.KubeconfigPath = scConfig.ExistingClusterKubeconfig
	noopManager := cluster.NewNoopClusterManager()

	kc, err := noopManager.Create(scConfig)
	if err != nil {
		return nil, err
	}

	log.Style(outputIndent, color.FgYellow).Warn("Warning: Components installed using this method will need to be manually removed.\n")

	return kc, nil
}

func blockForPullingBaseImage(cm cluster.Manager, scConfig *config.UnmanagedClusterConfig) error {
	// start a go routine to animate the downloading logs while the docker exec gets the image
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		log.Style(outputIndent, color.Reset).AnimateProgressWithOptions(
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
	log.Style(outputIndent, color.Faint).ReplaceLinef("Base image downloaded")

	return nil
}

func blockForClusterCreate(cm cluster.Manager, scConfig *config.UnmanagedClusterConfig) (*cluster.KubernetesCluster, error) {
	// start a go routine to animate the downloading logs while the docker exec gets the image
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		log.Style(outputIndent, color.Reset).AnimateProgressWithOptions(
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
	log.Style(outputIndent, color.Faint).ReplaceLinef("Cluster created")

	return kc, nil
}

func installKappController(t *UnmanagedCluster, kc kapp.Manager) (*v1.Deployment, error) {
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

	kappControllerCreated, err := kc.Install(kapp.InstallOpts{MergedManifests: kappBytes})
	if err != nil {
		return nil, err
	}

	return kappControllerCreated, nil
}

func blockForKappStatus(kappDeployment *v1.Deployment, kc kapp.Manager) {
	// Create the parent context and fire a go routine to animate the logging progress
	ctx, cancel := context.WithCancel(context.Background())
	status := make(chan string, 1)
	go func(ctx context.Context) {
		log.Style(outputIndent, color.Faint).AnimateProgressWithOptions(
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
			log.Style(outputIndent, color.Faint).ReplaceLinef("kapp-controller status: %s", kappState)
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
		log.Style(outputIndent, color.Reset).AnimateProgressWithOptions(
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
			log.Errorf("failed to check package repository status: %s\n", err.Error())
			return
		}
		if pkgStatus == "Reconcile succeeded" {
			cancel()
			log.Style(outputIndent, color.Faint).ReplaceLinef("Core package repo status: %s", pkgStatus)
			break
		}
		time.Sleep(1 * time.Second)
	}
}

// installCNI installs the CNI package to be satisfied via kapp-controller. If
// the selected CNI package is nil, it returns as a noop.
// TODO(joshrosso): this function is a mess, but waiting on some stuff to
// happen in other packages
func installCNI(pkgClient packages.PackageManager, t *UnmanagedCluster) error {
	if t.selectedCNIPkg == nil {
		return fmt.Errorf("cannot install CNI when value is nil")
	}
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

func mergeKubeconfigAndSetContext(mgr kubeconfig.Manager, kcPath, clusterName string) error {
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

// resolveCNI determines which CNI package to use. It expects to be passed a
// fully qualified package name except for special known CNI values such as
// antrea or calico.
func resolveCNI(mgr packages.PackageManager, cniName string) (*CNIPackage, error) {
	if cniName == cniNoneName {
		return nil, fmt.Errorf("CNI was set to %s", cniName)
	}
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
