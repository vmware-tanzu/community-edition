// Package tanzu is responsible for orchestrating the various components that result in a local tanzu cluster creation.
package tanzu

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

type TanzuCluster struct {
	Name string
}

type TanzuLocal struct {
	// kcPath               string
	// config               *LocalClusterConfig
	bom                  *tkr.TKRBom
	kappControllerBundle tkr.TkrImageReader
	selectedCNIPkg       *CNIPackage
}

type CNIPackage struct {
	fqPkgName  string
	pkgVersion string
}

// TODO(joshrosso): global logger for the package. This is kind gross, but really convenient.
var log logger.Logger

type TanzuMgr interface {
	// Deploy orchestrates all the required steps in order to create a local Tanzu cluster. This can involve
	// cluster creation, kapp-controller installation, CNI installation, and more. The steps that are taken
	// depend on the configuration passed into Deploy. If something goes wrong during deploy, an error is
	// returned.
	Deploy(config *LocalClusterConfig) error
	// List retrieves all known tanzu clusters are returns a list of them. If it's unable to interact with the
	// underlying cluster provider, it returns an error.
	List() ([]TanzuCluster, error)
	// Delete takes a cluster name and removes the cluster from the underlying cluster provider. If it is unable
	// to communicate with the underlying cluster provider, it returns an error.
	Delete(name string) error
}

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
	tceRepoUrl            = "projects.registry.vmware.com/tce/main:0.9.1"
)

// New returns a TanzuMgr for interacting with local clusters. It is implemented by TanzuLocal.
func New(logger logger.Logger) TanzuMgr {
	log = logger
	return &TanzuLocal{}
}

// validateConfiguration makes sure the configuration is valid, returning an
// error if there is an issue.
func validateConfiguration(config *LocalClusterConfig) error {
	if config.ClusterName == "" {
		return fmt.Errorf("cluster name is required")
	}

	if config.Provider == "" {
		// Should have been validated earlier, but not an error. We can just
		// default it to kind.
		config.Provider = "kind"
	}

	return nil
}

func (t *TanzuLocal) Deploy(config *LocalClusterConfig) error {
	var err error
	// 1. Validate the configuration
	if err := validateConfiguration(config); err != nil {
		return err
	}

	// 2. Download and Read the TKR
	log.Event("\\U+2692", " Resolving Tanzu Kubernetes Release (TKR)\n")
	bomFileName, err := getTkrBom("projects.registry.vmware.com/tkg/tkr-bom:v1.21.2_vmware.1-tkg.1")
	if err != nil {
		return fmt.Errorf("Failed getting TKR BOM. Error: %s", err.Error())
	}

	log.Event("\\U+2692", " Processing Tanzu Kubernetes Release\n")
	t.bom, err = parseTKRBom(bomFileName)
	if err != nil {
		return fmt.Errorf("Failed parsing TKR BOM. Error: %s", err.Error())
	}

	// 3. Resolve all required images
	// base image
	log.Event("\\U+1F5BC", " Selected base image\n")
	log.Style(2, logger.ColorLightGrey).Infof("%s\n", t.bom.GetTKRNodeImage())
	// core package repository
	log.Event("\\U+1F4E6", "Selected core package repository\n")
	log.Style(2, logger.ColorLightGrey).Infof("%s\n", t.bom.GetTKRCoreRepoBundlePath())
	// core user package repositories
	log.Event("\\U+1F4E6", "Selected additional package repositories\n")
	for _, additionalRepo := range t.bom.GetAdditionalRepoBundlesPaths() {
		log.Style(2, logger.ColorLightGrey).Infof("%s\n", additionalRepo)
	}
	// kapp-controller
	err = resolveKappBundle(t)
	if err != nil {
		return fmt.Errorf("Failed resolving kapp-controller bundle. Error: %s", err.Error())
	}
	log.Event("\\U+1F4E6", "Selected kapp-controller image bundle\n")
	log.Style(2, logger.ColorLightGrey).Infof("%s\n", t.kappControllerBundle.GetRegistryUrl())

	// 4. Create the cluster
	log.Eventf("\\U+1F6F0", " Creating cluster %s\n", config.ClusterName)
	err = runClusterCreate(config)
	if err != nil {
		return fmt.Errorf("Failed to create cluster, Error: %s", err.Error())
	}
	kcPath := config.KubeConfigPath()
	log.Style(2, logger.ColorLightGrey).Info("To troubleshoot, use:\n")
	log.Style(2, logger.ColorLightGrey).Infof("kubectl ${COMMAND} --kubeconfig %s\n", kcPath)

	// 5. Install kapp-controller
	kc := kapp.New(kcPath)
	log.Event("\\U+1F4E7", "Installing kapp-controller\n")
	kappDeployment, err := installKappController(t, kc)
	if err != nil {
		return fmt.Errorf("Failed to install kapp-controller, Error: %s", err.Error())
	}
	blockForKappStatus(kappDeployment, kc)

	// 6. Install package repositories
	pkgClient := packages.NewClient(kcPath)
	log.Event("\\U+1F4E7", "Installing package repositories\n")
	createdCoreRepo, err := createPackageRepo(pkgClient, tkgSysNamespace, tkgCoreRepoName, t.bom.GetTKRCoreRepoBundlePath())
	if err != nil {
		return fmt.Errorf("Failed to install core package repo. Error: %s", err.Error())
	}
	for _, additionalRepo := range t.bom.GetAdditionalRepoBundlesPaths() {
		_, err = createPackageRepo(pkgClient, tkgGlobalPkgNamespace, tkgCoreRepoName, additionalRepo)
		if err != nil {
			return fmt.Errorf("Failed to install adiditonal package repo. Error: %s", err.Error())
		}
	}
	blockForRepoStatus(createdCoreRepo, pkgClient)

	// 7. Install CNI
	log.Event("\\U+1F4E6", "Insatlling CNI\n")
	t.selectedCNIPkg, err = resolveCNI(pkgClient, "antrea")
	if err != nil {
		return fmt.Errorf("Failed to resolve a CNI package. Error: %s", err.Error())
	}
	log.Style(2, logger.ColorLightGrey).Infof("%s:%s\n", t.selectedCNIPkg.fqPkgName, t.selectedCNIPkg.pkgVersion)
	err = installCNI(pkgClient, t)
	if err != nil {
		return fmt.Errorf("Failed to install the CNI package. Error: %s", err.Error())
	}

	// 8. Update kubeconfig and context
	kubeConfigMgr := kubeconfig.NewManager()
	err = mergeKubeconfigAndSetContext(kubeConfigMgr, kcPath, config.ClusterName)
	if err != nil {
		log.Warnf("Failed to merge kubeconfig and set your context. Cluster should still work! Error: %s", err)
	}

	// 8. Return
	log.Event("\\U+2705", "Cluster created\n")
	log.Eventf("\\U+1F3AE", "kubectl context set to %s\n\n", config.ClusterName)
	// provide user example commands to run
	log.Style(0, logger.ColorLightGrey).Infof("View available packages:\n")
	log.Style(2, logger.ColorLightGreen).Infof("tanzu package available list\n")
	log.Style(0, logger.ColorLightGrey).Infof("View running pods:\n")
	log.Style(2, logger.ColorLightGreen).Infof("kubectl get po -A\n")
	log.Style(0, logger.ColorLightGrey).Infof("Delete this cluster:\n")
	log.Style(2, logger.ColorLightGreen).Infof("tanzu local delete %s\n", config.ClusterName)
	return nil
}

func (t *TanzuLocal) List() ([]TanzuCluster, error) {
	panic("implement me")
}

func (t *TanzuLocal) Delete(name string) error {
	panic("implement me")
}

// getTkgConfigDir returns the configuration directory used by tce.
func getTkgConfigDir() (path string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return path, fmt.Errorf("Failed to resolve home dir. Error: %s", err.Error())
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

func getTkrBom(registry string) (string, error) {
	log.Style(2, logger.ColorLightGrey).Infof("%s\n", registry)
	expectedBomName := buildFilesystemSafeBomName(registry)

	bomPath, err := getLocalBomPath()
	if err != nil {
		return "", fmt.Errorf("failed to get tanzu local bom path: %s\n", err)
	}

	_, err = os.Stat(bomPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(bomPath, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to make new tanzu local bom config directories %s\n", err)
		}
	}

	items, err := os.ReadDir(bomPath)
	if err != nil {
		return "", fmt.Errorf("failed to read tanzu local bom directories: %s\n", err)
	}

	// if the expected bom is already in the config directory, don't download it again. return early
	for _, file := range items {
		if file.Name() == expectedBomName {
			log.Style(2, logger.ColorLightGrey).Infof("TKR exists at %s\n", filepath.Join(bomPath, file.Name()))
			return file.Name(), nil
		}
	}

	bomImage, err := tkr.NewTkrImageReader(registry)
	if err != nil {
		return "", fmt.Errorf("failed to create new TkrImageReader: %s\n", err)
	}

	log.Style(2, logger.ColorLightGrey).Infof("Downloading to %s\n", filepath.Join(bomPath, expectedBomName))
	err = bomImage.DownloadImage()
	if err != nil {
		return "", fmt.Errorf("failed to download tkr image: %s\n", err)
	}

	downloadedBomFiles, err := os.ReadDir(bomImage.GetDownloadPath())
	if err != nil {
		return "", fmt.Errorf("failed to read downloaded tkr bom files: %s\n", err)
	}

	// if there is more than 1 file in the downloaded image, fail
	// this is a bit redundent since imgpkg librariers should fail if the image is a bundle with multiple files
	if len(downloadedBomFiles) != 1 {
		return "", fmt.Errorf("more than one file found in TKR bom image. Expected 1 file: %s\n", bomImage.GetDownloadPath())
	}

	downloadedBomFile, err := os.Open(filepath.Join(bomImage.GetDownloadPath(), downloadedBomFiles[0].Name()))
	if err != nil {
		return "", fmt.Errorf("could not open downloaded tkr bom file: %s\n", err)
	}
	defer downloadedBomFile.Close()

	newBomFile, err := os.Create(filepath.Join(bomPath, expectedBomName))
	if err != nil {
		return "", fmt.Errorf("could not create tanzu local bom tkr file: %s\n", err)
	}
	defer newBomFile.Close()

	_, err = io.Copy(newBomFile, downloadedBomFile)
	if err != nil {
		return "", fmt.Errorf("could not copy file contents: %s\n", err)
	}

	return expectedBomName, nil
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

func runClusterCreate(config *LocalClusterConfig) error {
	clusterManager := cluster.NewClusterManager()
	clusterCreateOpts := cluster.CreateOpts{
		Name:           config.ClusterName,
		KubeconfigPath: config.KubeConfigPath(),
		// Config: TBD,
	}
	_, err := clusterManager.Create(&clusterCreateOpts)
	if err != nil {
		return err
	}
	return nil
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
			log.Errorf("Failed to check kapp-controller status: %s", err.Error())
			return
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
}

// TODO(joshrosso) this function is a mess, but waiting on some stuff to happen in other packages
func installCNI(pkgClient packages.PackageManager, t *TanzuLocal) error {

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
		FqPkgName:      t.selectedCNIPkg.fqPkgName,
		Version:        t.selectedCNIPkg.pkgVersion,
		Configuration:  []byte(valueData),
		ServiceAccount: rootSvcAcct.Name,
	}
	_, err = pkgClient.CreatePackageInstall(cniInstallOpts)
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
	if len(pkgs) < 1 {
		return nil, fmt.Errorf("No package was resolved for CNI choice %s", cniName)
	}

	cniPkg := &CNIPackage{}
	for _, pkg := range pkgs {
		if strings.HasPrefix(pkg.Spec.RefName, cniName) {
			cniPkg.fqPkgName = pkg.Spec.RefName
			cniPkg.pkgVersion = pkg.Spec.Version
		}
	}

	return cniPkg, nil
}
