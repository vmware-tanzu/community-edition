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
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/semver"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/tkr"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin"
)

//nolint
const (
	configFileName           = "config.yaml"
	bootstrapLogName         = "bootstrap.log"
	bomDir                   = "bom"
	compatibilityDir         = "compatibility"
	tkgSysNamespace          = "tkg-system"
	tkgSvcAcctName           = "core-pkgs"
	tkgCoreRepoName          = "tkg-core-repository"
	tkgGlobalPkgNamespace    = "tanzu-package-repo-global"
	tceCompatibilityRegistry = "projects.registry.vmware.com/tce/compatibility"
	tceRepoName              = "community-repository"
	outputIndent             = 3
	maxProgressLength        = 4
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
	Status   string
}

// UnmanagedCluster contains information about an unmanaged Tanzu cluster.
type UnmanagedCluster struct {
	bom                  *tkr.Bom
	kappControllerBundle tkr.ImageReader
	selectedCNIPkg       *Package
	profilePkg           *Package
	config               *config.UnmanagedClusterConfig
	clusterDirectory     string
}

type Package struct {
	installName string
	fqPkgName   string
	pkgVersion  string
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
	// Stop takes a cluster name and attempts to stop a running cluster. If it is unable
	// to communicate with the underlying cluster provider, it returns an error.
	Stop(name string) error
	// Start takes a cluster name and attempts to start a stopped cluster. If
	// there are issues starting the cluster or communitcating with the
	// underlying provider, an error is returned.
	Start(name string) error
}

// New returns a TanzuMgr for interacting with unmanaged clusters. It is implemented by TanzuUnmanaged.
func New(parentLogger logger.Logger) Manager {
	log = parentLogger
	return &UnmanagedCluster{}
}

// validateConfiguration makes sure the configuration is valid, returning an
// error if there is an issue.
func validateConfiguration(scConfig *config.UnmanagedClusterConfig) error {
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

	// this var is used to return a non-critical error code
	// (like when installing the CNI fails)
	returnCode := Success

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
	// Use a default log file if config option was not provided by user
	if scConfig.LogFile == "" {
		scConfig.LogFile = filepath.Join(t.clusterDirectory, "bootstrap.log")
	}
	log.AddLogFile(scConfig.LogFile)
	log.Event(logger.FolderEmoji, "Created cluster directory")

	// Log a warning if the user has given a ProviderConfiguration
	if len(scConfig.ProviderConfiguration) != 0 {
		log.Style(outputIndent, color.FgYellow).ReplaceLinef("Reading ProviderConfiguration from config file. Some provider specific flags and configs may be ignored.")
	}

	// 2. Download and Read the compatible TKr

	// Download compatibility file
	log.Event(logger.MagnetEmoji, "Resolving and checking Tanzu Kubernetes release (TKr) compatibility file")
	tkrCompatibility, err := getTkrCompatibility()
	if err != nil {
		return ErrTkrBom, fmt.Errorf("failed downloading and extracting TKr compatibility file. Error: %s", err.Error())
	}

	if scConfig.TkrLocation == "" {
		// read TKr version compatible with version of unmanaged-cluster CLI
		// when the user did _not_ set the TKr via a flag or configuration option
		scConfig.TkrLocation, err = getLatestCompatibleTkr(tkrCompatibility)
		if err != nil {
			return ErrTkrBom, fmt.Errorf("failed parsing the TKr compatibility file. Error: %s", err.Error())
		}
	} else {
		// check if the TKr specified by the user is in the list of compatible TKrs
		// If not, log a warning
		if !isTkrCompatible(tkrCompatibility, scConfig.TkrLocation) {
			log.Style(outputIndent, color.FgYellow).Warnf("Custom TKr %s NOT found in compatibility file. Proceed with caution, the provided TKr may not work with this version of unmanaged-cluster\n", scConfig.TkrLocation)
		} else {
			log.Style(outputIndent, color.Faint).Infof("Custom Tkr %s found in compatibility file\n", scConfig.TkrLocation)
		}
	}

	log.Event(logger.WrenchEmoji, "Resolving TKr")
	bomFileName, err := getTkrBom(scConfig.TkrLocation)
	if err != nil {
		return ErrTkrBom, fmt.Errorf("failed getting TKr BOM. Error: %s", err.Error())
	}
	configFp := filepath.Join(t.clusterDirectory, configFileName)
	err = config.RenderConfigToFile(configFp, t.config)
	if err != nil {
		return ErrRenderingConfig, err
	}
	log.Style(outputIndent, color.Faint).Infof("Rendered Config: %s\n", configFp)
	log.Style(outputIndent, color.Faint).Infof("Bootstrap Logs: %s\n", scConfig.LogFile)

	log.Event(logger.WrenchEmoji, "Processing Tanzu Kubernetes Release")
	t.bom, err = parseTKRBom(bomFileName)
	if err != nil {
		return ErrTkrBomParsing, fmt.Errorf("failed parsing TKr BOM. Error: %s", err.Error())
	}

	// Uses default user package repository found in the TKr if user did not provide one via config/flags
	if len(scConfig.AdditionalPackageRepos) == 0 {
		userRepo := t.bom.GetTKRUserRepoBundlePath()

		if userRepo != "" {
			scConfig.AdditionalPackageRepos = []string{
				userRepo,
			}
		}
	}

	// 3. Resolve all required images
	// base image
	log.Event(logger.PictureEmoji, "Selected base image")
	scConfig.NodeImage = t.bom.GetTKRNodeImage(scConfig.Provider)
	if scConfig.NodeImage == "" {
		return ErrTkrBomParsing, fmt.Errorf("failed parsing TKR BOM. Could not get base node image for provider %s", scConfig.Provider)
	}
	log.Style(outputIndent, color.Faint).Infof("%s\n", scConfig.NodeImage)

	// core package repository
	log.Event(logger.PackageEmoji, "Selected core package repository")
	log.Style(outputIndent, color.Faint).Infof("%s\n", t.bom.GetTKRCoreRepoBundlePath())

	// core user package repositories if they exist
	if len(scConfig.AdditionalPackageRepos) != 0 {
		log.Event(logger.PackageEmoji, "Selected additional package repositories")
		for _, additionalRepo := range scConfig.AdditionalPackageRepos {
			log.Style(outputIndent, color.Faint).Infof("%s\n", additionalRepo)
		}
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
	// Install and wait for core repo first
	pkgClient := packages.NewClient(kcBytes)
	log.Event(logger.EnvelopeEmoji, "Installing package repositories")
	createdCoreRepo, err := createPackageRepo(pkgClient, tkgSysNamespace, tkgCoreRepoName, t.bom.GetTKRCoreRepoBundlePath())
	if err != nil {
		return ErrCorePackageRepoInstall, fmt.Errorf("failed to install core package repo. Error: %s", err.Error())
	}
	blockForRepoStatus(createdCoreRepo, pkgClient)

	// Install the additional package repos
	for _, additionalRepo := range scConfig.AdditionalPackageRepos {
		kappFriendlyRepoName := strings.ReplaceAll(additionalRepo, "/", "-")
		kappFriendlyRepoName = strings.ReplaceAll(kappFriendlyRepoName, ":", "-")
		createdAdditionalRepo, err := createPackageRepo(pkgClient, tkgGlobalPkgNamespace, kappFriendlyRepoName, additionalRepo)
		if err != nil {
			return ErrOtherPackageRepoInstall, fmt.Errorf("failed to install adiditonal package repo. Error: %s", err.Error())
		}
		// Wait for additional package repos to be ready so that we can install a profile latter
		if len(t.config.Profiles) != 0 {
			blockForRepoStatus(createdAdditionalRepo, pkgClient)
		}
	}

	// 7. Install CNI
	// CNI plugins are installed as best effort. If no plugin is resolved in the
	// repository, no CNI is installed, yet the cluster will still run.
	log.Event(logger.GlobeEmoji, "Installing CNI")
	t.selectedCNIPkg, err = resolvePkg(pkgClient, tkgSysNamespace, t.config.Cni, "")
	if err != nil {
		log.Style(outputIndent, color.FgYellow).Warnf("WARNING: failed to select the CNI package. Error: %s", err.Error())
		returnCode = ErrCniInstall
	}

	if t.selectedCNIPkg != nil {
		// CNI package resolved, do install
		log.Style(outputIndent, color.Faint).Infof("%s:%s\n", t.selectedCNIPkg.fqPkgName, t.selectedCNIPkg.pkgVersion)
		err = installCNI(pkgClient, t)
		if err != nil {
			log.Style(outputIndent, color.FgYellow).Warnf("WARNING: failed to install CNI. Error: %s", err.Error())
			returnCode = ErrCniInstall
		}
	}

	// 8. Install profile if specified
	for _, profile := range t.config.Profiles {
		log.Eventf(logger.GlobeEmoji, "Installing Profile %s\n", profile.Name)

		t.profilePkg, err = resolvePkg(pkgClient, tkgGlobalPkgNamespace, profile.Name, profile.Version)
		if err != nil {
			log.Style(outputIndent, color.FgYellow).Warnf("WARNING: failed to install profile %s. Error: %s", profile.Name, err.Error())
			returnCode = ErrProfileInstall
			continue
		}

		log.Style(outputIndent, color.Faint).Infof("Selected package %s\n", t.profilePkg.fqPkgName)

		if profile.Version == "" {
			log.Style(outputIndent, color.FgYellow).Warnf("Installing profile without version specified. Using version %s\n", t.profilePkg.pkgVersion)
		} else {
			log.Style(outputIndent, color.Faint).Infof("Using profile version %s\n", t.profilePkg.pkgVersion)
		}

		if profile.Config == "" {
			log.Style(outputIndent, color.FgYellow).Warnf("Installing profile with no configuration file\n")
		} else {
			log.Style(outputIndent, color.Faint).Infof("Using config %s\n", profile.Config)
		}

		err = installProfile(pkgClient, t, profile.Config)
		if err != nil {
			log.Style(outputIndent, color.FgYellow).Warnf("WARNING: failed to install profile %s. Error: %s", profile.Name, err.Error())
			returnCode = ErrProfileInstall
		}
	}

	// 9. Update kubeconfig and context
	kubeConfigMgr := kubeconfig.NewManager()
	err = mergeKubeconfigAndSetContext(kubeConfigMgr, scConfig.KubeconfigPath)
	if err != nil {
		log.Warnf("Failed to merge kubeconfig and set your context. Cluster should still work! Error: %s", err)
		returnCode = ErrKubeconfigContextSet
	}

	// 10. Return
	log.Event(logger.GreenCheckEmoji, "Cluster created")
	log.Eventf(logger.ControllerEmoji, "kubectl context set to %s\n\n", scConfig.ClusterName)
	// provide user example commands to run
	log.Infof("View available packages:\n")
	log.Style(outputIndent, color.FgGreen).Infof("tanzu package available list\n")
	log.Infof("View running pods:\n")
	log.Style(outputIndent, color.FgGreen).Infof("kubectl get po -A\n")
	log.Infof("Delete this cluster:\n")
	log.Style(outputIndent, color.FgGreen).Infof("tanzu unmanaged delete %s\n", scConfig.ClusterName)
	return returnCode, nil
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
		log.Style(outputIndent, color.FgYellow).Warnf("Warning - could not resolve cluster config directory.\n")
		log.Style(outputIndent, color.FgYellow).Warnf("Cluster NOT deleted.\n")
		log.Style(outputIndent, color.FgYellow).Warnf("Local config files NOT deleted.\n")
		log.Style(outputIndent, color.FgYellow).Warnf("Be sure to manually delete cluster and local config files\n")
		return err
	}

	configPath, err := resolveClusterConfig(name)
	if err != nil {
		log.Style(outputIndent, color.FgYellow).Warnf("Warning - could not resolve cluster config file. Error: %s\n", err.Error())
		log.Style(outputIndent, color.FgYellow).Warnf("Cluster NOT deleted.\n")
		log.Style(outputIndent, color.FgYellow).Warnf("Be sure to manually delete cluster\n")
		deleteErr := os.RemoveAll(t.clusterDirectory)
		if deleteErr != nil {
			log.Style(outputIndent, color.FgRed).Errorf("Failed to remove config %s. Be sure to manually delete files\n", t.clusterDirectory)
			return deleteErr
		}

		log.Style(outputIndent, color.Faint).Infof("Local config files directory deleted: %s\n", t.clusterDirectory)
		return err
	}

	t.config, err = config.RenderFileToConfig(configPath)
	if err != nil {
		log.Style(outputIndent, color.FgYellow).Warnf("Warning - could not create configuration from local config file. Error: %s\n", err.Error())
		log.Style(outputIndent, color.FgYellow).Warnf("Cluster NOT deleted.\n")
		log.Style(outputIndent, color.FgYellow).Warnf("Be sure to manually delete cluster\n")
		deleteErr := os.RemoveAll(t.clusterDirectory)
		if deleteErr != nil {
			log.Style(outputIndent, color.FgRed).Errorf("Failed to remove config %s. Be sure to manually delete files\n", t.clusterDirectory)
			return deleteErr
		}

		log.Style(outputIndent, color.Faint).Infof("Local config files directory deleted: %s\n", t.clusterDirectory)
		return err
	}

	cm := cluster.NewClusterManager(t.config)

	err = cm.Delete(t.config)
	if err != nil {
		log.Style(outputIndent, color.FgYellow).Warnf("Warning - could not delete cluster through provider. Be sure to manually delete cluster. Error: %s\n", err.Error())
	}

	deleteErr := os.RemoveAll(t.clusterDirectory)
	if deleteErr != nil {
		log.Style(outputIndent, color.FgYellow).Warnf("Cluster deleted but failed to remove config %s. Be sure to manually delete files\n", t.clusterDirectory)
		return deleteErr
	}

	log.Style(outputIndent, color.Faint).Infof("Local config files directory deleted: %s\n", t.clusterDirectory)
	return err
}

// Stop tells an unmanaged cluster to no longer continue running.
func (t *UnmanagedCluster) Stop(name string) error {
	configPath, err := resolveClusterConfig(name)
	if err != nil {
		return err
	}
	t.config, err = config.RenderFileToConfig(configPath)
	if err != nil {
		return err
	}

	cm := cluster.NewClusterManager(t.config)

	err = cm.Stop(t.config)
	if err != nil {
		return err
	}

	return nil
}

// Start tells an unmanaged cluster to start a cluster that is not currently running.
func (t *UnmanagedCluster) Start(name string) error {
	configPath, err := resolveClusterConfig(name)
	if err != nil {
		return err
	}
	t.config, err = config.RenderFileToConfig(configPath)
	if err != nil {
		return err
	}

	cm := cluster.NewClusterManager(t.config)

	err = cm.Start(t.config)
	if err != nil {
		return err
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

func getUnmanagedCompatibilityPath() (path string, err error) {
	tkgUnmanagedConfigDir, err := config.GetUnmanagedConfigPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(tkgUnmanagedConfigDir, compatibilityDir), nil
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

func getTkrCompatibility() (*tkr.Compatibility, error) {
	compatibilityFileName, err := getCompatibilityFile()
	if err != nil {
		return nil, err
	}

	compatibilityDirPath, err := getUnmanagedCompatibilityPath()
	if err != nil {
		return nil, err
	}

	// Read compatibility file, returns a compatibility struct
	c, err := tkr.ReadCompatibilityFile(filepath.Join(compatibilityDirPath, compatibilityFileName))
	if err != nil {
		return nil, err
	}

	return c, nil
}

func isTkrCompatible(c *tkr.Compatibility, tkrName string) bool {
	// Inspect CLI version and get most recent compatible version of TKr
	for _, cliVersion := range c.UnmanagedClusterPluginVersions {
		if cliVersion.Version == plugin.Version {
			for _, possibleTkr := range cliVersion.SupportedTkrVersions {
				if possibleTkr.Path == tkrName {
					return true
				}
			}
		}
	}

	return false
}

// Returns the latest TKr compatible image path:tag string
// If none is found for version of CLI, returns an error
func getLatestCompatibleTkr(c *tkr.Compatibility) (string, error) {
	// Inspect CLI version and get most recent compatible version of TKr
	for _, cliVersion := range c.UnmanagedClusterPluginVersions {
		if cliVersion.Version == plugin.Version {
			// We've found a compatible version
			// Check it's filled to prevent a panic. We should never ship a compatibility file with an empty compatibility for a CLI version
			if len(cliVersion.SupportedTkrVersions) == 0 || cliVersion.SupportedTkrVersions[0].Path == "" {
				return "", fmt.Errorf("most recent compatibility image path is invalid. Validate compatibility file")
			}

			return cliVersion.SupportedTkrVersions[0].Path, nil
		}
	}

	return "", fmt.Errorf("could not find compatible CLI version in compatibility file")
}

// Returns the file path of the latest compatibility file
// If the file is _not_ on the system, downloads it
func getCompatibilityFile() (string, error) {
	log.Style(outputIndent, color.Faint).Infof("%s\n", tceCompatibilityRegistry)

	// Get latest versioned tag from the registry for the compatibility file
	tag, err := tkr.GetLatestCompatibilityTag(tceCompatibilityRegistry)
	if err != nil {
		return "", err
	}

	expectedCompatibilityFileName := buildFilesystemSafeBomName(tceCompatibilityRegistry + ":" + tag)
	compatibilityPath, err := getUnmanagedCompatibilityPath()
	if err != nil {
		return "", fmt.Errorf("failed to get tanzu unmanaged compatibility path: %s", err)
	}

	_, err = os.Stat(compatibilityPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(compatibilityPath, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to make new tanzu unmanaged compatibility config directories %s", err)
		}
	}

	items, err := os.ReadDir(compatibilityPath)
	if err != nil {
		return "", fmt.Errorf("failed to read tanzu unmanaged compatibility directories: %s", err)
	}

	// if the expected compatibility file is already in the config directory, don't download it again. return early
	for _, file := range items {
		if file.Name() == expectedCompatibilityFileName {
			log.Style(outputIndent, color.Faint).Infof("Compatibility file exists at %s\n", filepath.Join(compatibilityPath, file.Name()))
			return file.Name(), nil
		}
	}

	registry := tceCompatibilityRegistry + ":" + tag

	compatibilityImage, err := tkr.NewTkrImageReader(registry)
	if err != nil {
		return "", fmt.Errorf("failed to create new TkrImageReader: %s", err)
	}

	err = blockForImageDownload(compatibilityImage, compatibilityPath, expectedCompatibilityFileName)
	if err != nil {
		return "", fmt.Errorf("failed to download compatibility image: %s", err)
	}

	downloadedCompatibilityFiles, err := os.ReadDir(compatibilityImage.GetDownloadPath())
	if err != nil {
		return "", fmt.Errorf("failed to read downloaded compatibility files: %s", err)
	}

	// if there is more than 1 file in the downloaded image, fail
	// this is a bit redundant since imgpkg librariers should fail if the image is a bundle with multiple files
	if len(downloadedCompatibilityFiles) != 1 {
		return "", fmt.Errorf("more than one file found in compatibility image. Expected 1 file: %s", compatibilityImage.GetDownloadPath())
	}

	downloadedCompatibilityFile, err := os.Open(filepath.Join(compatibilityImage.GetDownloadPath(), downloadedCompatibilityFiles[0].Name()))
	if err != nil {
		return "", fmt.Errorf("could not open downloaded compatibility file: %s", err)
	}
	defer downloadedCompatibilityFile.Close()

	newCompatibilityFile, err := os.Create(filepath.Join(compatibilityPath, expectedCompatibilityFileName))
	if err != nil {
		return "", fmt.Errorf("could not create tanzu unmanaged compatibility file: %s", err)
	}
	defer newCompatibilityFile.Close()

	_, err = io.Copy(newCompatibilityFile, downloadedCompatibilityFile)
	if err != nil {
		return "", fmt.Errorf("could not copy compatibility file contents: %s", err)
	}

	return expectedCompatibilityFileName, nil
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
			log.Style(outputIndent, color.Faint).Infof("TKr exists at %s\n", filepath.Join(bomPath, file.Name()))
			return file.Name(), nil
		}
	}

	bomImage, err := tkr.NewTkrImageReader(registry)
	if err != nil {
		return "", fmt.Errorf("failed to create new TkrImageReader: %s", err)
	}

	err = blockForImageDownload(bomImage, bomPath, expectedBomName)
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
		return "", fmt.Errorf("more than one file found in TKr bom image. Expected 1 file: %s", bomImage.GetDownloadPath())
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

func blockForImageDownload(b tkr.ImageReader, path, expectedName string) error {
	f := filepath.Join(path, expectedName)

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

	for _, message := range clusterManager.PreProviderNotify() {
		log.Style(outputIndent, color.Faint).Info(message)
	}

	if !scConfig.SkipPreflightChecks {
		warnings, issues := clusterManager.PreflightCheck()
		if len(issues) > 0 {
			return nil, fmt.Errorf("system checks detected issues, please resolve first: %v", issues)
		}

		for _, warning := range warnings {
			log.Style(outputIndent, color.FgYellow).Warnf("WARNING: %s\n", warning)
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

	for _, l := range cm.PostProviderNotify() {
		log.Style(outputIndent, color.FgYellow).Warnf("%s\n", l)
	}

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
			logger.AnimatorWithMessagef("Package repo status: %s - %s", repo.Name),
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
			log.Style(outputIndent, color.Faint).ReplaceLinef("%s package repo status: %s", repo.Name, pkgStatus)
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

// GetKubeconfigContext returns the current context for a passed in kubeconfig file
// This is a utility function that enables users of the `tanzu` packages
// to utilize an existing cluster with an existing kubeconfig and get it's current context
func ReadClusterContextFromKubeconfig(kcPath string) (string, error) {
	ctx, err := kubeconfig.GetKubeconfigContext(kcPath)
	if err != nil {
		return "", fmt.Errorf("could not get context from kubeconfig found at %s - Error: %s", kcPath, err.Error())
	}

	if ctx == "" {
		return "", fmt.Errorf("no current context set")
	}

	return ctx, nil
}

func mergeKubeconfigAndSetContext(mgr kubeconfig.Manager, kcPath string) error {
	err := mgr.MergeToDefaultConfig(kcPath)
	if err != nil {
		log.Errorf("Failed to merge kubeconfig: %s\n", err.Error())
		return nil
	}

	// Get the current kubeconfig context: this should be the newly created/attached cluster
	kubeContextName, err := ReadClusterContextFromKubeconfig(kcPath)
	if err != nil {
		return err
	}

	err = mgr.SetCurrentContext(kubeContextName)
	if err != nil {
		return err
	}

	return nil
}

// resolvePkg picks the first package in the package repo
// that matches the name and version of the provided profile
// If the user did not specify a profile version, defaults to the first one found which should be the latest
func resolvePkg(mgr packages.PackageManager, namespace, pkgName, pkgVersion string) (*Package, error) {
	keyWordLatest := "latest"

	if pkgName == cniNoneName {
		log.Style(outputIndent, color.FgYellow).Warnf("No CNI installed: CNI was set to %s.\n", cniNoneName)
		return nil, nil
	}

	pkgs, err := mgr.ListPackagesInNamespace(namespace)
	if err != nil {
		return nil, err
	}

	versions := []string{}

	profilePkg := &Package{
		installName: pkgName,
	}

	for _, pkg := range pkgs {
		// Select the package by name directly or by a well known prefix (like calico, antrea, fluent-bit)
		if pkg.Spec.RefName == pkgName || strings.HasPrefix(pkg.Spec.RefName, pkgName) {
			if pkgVersion == "" || pkgVersion == keyWordLatest || pkg.Spec.Version == pkgVersion {
				profilePkg.fqPkgName = pkg.Spec.RefName
				profilePkg.pkgVersion = pkg.Spec.Version

				// Build list of semantic versions for matching package
				versions = append(versions, pkg.Spec.Version)
			}
		}
	}

	if profilePkg.fqPkgName == "" {
		return nil, fmt.Errorf("no package was resolved for name %s with version %s", pkgName, pkgVersion)
	}

	// sort and select latest semver (last in slice) if user did not specify the version to use
	if pkgVersion == "" || pkgVersion == keyWordLatest {
		semver.Sort(versions)
		profilePkg.pkgVersion = versions[len(versions)-1]
	}

	return profilePkg, nil
}

// installProfile installs the profile package to be satisfied via kapp-controller and any provided config file
func installProfile(pkgClient packages.PackageManager, t *UnmanagedCluster, profileConfigPath string) error {
	rootSvcAcct, err := pkgClient.GetRootServiceAccount(tkgSysNamespace, tkgSvcAcctName)
	if err != nil {
		log.Errorf("failed to get root service account: %s\n", err.Error())
		return err
	}

	if rootSvcAcct == nil {
		log.Errorf("the package client root service account is nil and may have not been created successfully")
		return err
	}

	var valueData string

	if profileConfigPath != "" {
		data, err := os.ReadFile(profileConfigPath)
		if err != nil {
			return fmt.Errorf("could not read profile config file. Error: %s", err.Error())
		}

		valueData = string(data)
	}

	profileInstallOpts := packages.PackageInstallOpts{
		Namespace:      tkgSysNamespace,
		InstallName:    t.profilePkg.installName,
		FqPkgName:      t.profilePkg.fqPkgName,
		Version:        t.profilePkg.pkgVersion,
		Configuration:  []byte(valueData),
		ServiceAccount: rootSvcAcct.Name,
	}

	_, err = pkgClient.CreatePackageInstall(&profileInstallOpts)
	if err != nil {
		return fmt.Errorf("could not install profile. Error: %s", err.Error())
	}

	return nil
}
