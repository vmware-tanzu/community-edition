package e2e

import (
	"context"
	"flag"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"os"
	"path"
	"runtime"
	"testing"
	"time"

	packagelib "github.com/vmware-tanzu-private/core/cmd/cli/plugin/package/test/lib"
	"github.com/vmware-tanzu-private/core/pkg/v1/tkg/tkgpackagedatamodel"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	crtClient "sigs.k8s.io/controller-runtime/pkg/client"
	crtConfig "sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var (
	kubeconfigPath          string
	configPath              string
	config                  *Config
	currentDir              string
	packageRepoImgpkgBundle string
	iaas                    string
	clusterType             string
	client                  crtClient.Client
	packagePlugin           packagelib.PackagePlugin
)

const (
	interval = 30 * time.Second
	timeout  = 3 * time.Minute
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Contour test suite")
}

func init() {
	flag.StringVar(&configPath, "e2e.config", "", "path to the e2e config file")
	flag.StringVar(&packageRepoImgpkgBundle, "package-repo-imgpkg-bundle", "", "package repository imgpkg bundle registry url")
	flag.StringVar(&iaas, "iaas", "", "IAAS")
	flag.StringVar(&clusterType, "cluster-type", "", "Cluster Type")

	_, filename, _, _ := runtime.Caller(0)
	currentDir = path.Dir(filename)

	crtClient, err := crtClient.New(crtConfig.GetConfigOrDie(), crtClient.Options{})
	if err != nil {
		fmt.Errorf("Error creating k8s client %v", err)
		os.Exit(1)
	}
	client = crtClient

	packagePlugin = packagelib.NewPackagePlugin(kubeconfigPath, interval, timeout, "", "", 0)
}

var _ = BeforeSuite(func() {
	kubeconfigPath = os.Getenv("KUBECONFIG")

	e2eConfig, err := loadE2EConfig(configPath)
	Expect(err).NotTo(HaveOccurred())
	config = e2eConfig
	Expect(config).NotTo(BeNil())

	// Flag takes precedence than config file
	if iaas != "" {
		config.IAAS = iaas
	}
	if clusterType != "" {
		config.ClusterType = clusterType
	}

	Expect(config.IAAS).NotTo(BeEmpty())
	Expect(config.ClusterType).NotTo(BeEmpty())
	Expect(config.Package).NotTo(BeNil())
	Expect(config.PackageRepository).NotTo(BeNil())

	// Flag takes precedence than config file
	if packageRepoImgpkgBundle != "" {
		config.PackageRepository.ImgpkgBundle = packageRepoImgpkgBundle
	}

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.PackageRepository.Namespace,
		},
	}
	_, err = controllerutil.CreateOrPatch(context.TODO(), client, namespace, nil)
	Expect(err).NotTo(HaveOccurred())

	repoOptions := &tkgpackagedatamodel.RepositoryOptions{
		RepositoryName:     config.PackageRepository.Name,
		RepositoryURL: config.PackageRepository.ImgpkgBundle,
		Namespace:  config.PackageRepository.Namespace,
	}
	packagePluginResult := packagePlugin.AddOrUpdateRepository(repoOptions)
	Expect(packagePluginResult.Error).NotTo(HaveOccurred())
	Expect(packagePlugin.CheckRepositoryAvailable(repoOptions).Error).NotTo(HaveOccurred())

	// TODO: deploy all dependencies instead of statically just deploying cert-manager
	if config.ClusterType == "workload" || config.ClusterType == "standalone" {
		deployCertManager()
	}
})

var _ = AfterSuite(func() {
	if config.ClusterType == "workload" || config.ClusterType == "standalone" {
		deleteCertManager()
	}

	repoOptions := &tkgpackagedatamodel.RepositoryOptions{
		RepositoryName: config.PackageRepository.Name,
		Namespace:config.PackageRepository.Namespace,
	}
	packagePluginResult := packagePlugin.DeleteRepository(repoOptions)
	Expect(packagePluginResult.Error).NotTo(HaveOccurred())
	Expect(packagePlugin.CheckRepositoryDeleted(repoOptions).Error).NotTo(HaveOccurred())

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.PackageRepository.Namespace,
		},
	}
	Expect(client.Delete(context.TODO(), namespace)).NotTo(HaveOccurred())
})

func deployCertManager() {
	packageOptions := &tkgpackagedatamodel.PackageInstalledOptions{
		PkgInstallName: config.Package.Dependencies[0].Name,
		Namespace:        config.PackageRepository.Namespace,
		PackageName:      config.Package.Dependencies[0].RefName,
		Version:          config.Package.Dependencies[0].Version,
		Wait:             true,
	}
	Expect(packagePlugin.CheckPackageAvailable(config.Package.Dependencies[0].RefName, &tkgpackagedatamodel.PackageAvailableOptions{
		Namespace: config.PackageRepository.Namespace,
	}).Error).NotTo(HaveOccurred())
	Expect(packagePlugin.CheckAndInstallPackage(packageOptions).Error).NotTo(HaveOccurred())
	Expect(packagePlugin.CheckPackageInstalled(packageOptions).Error).NotTo(HaveOccurred())
}

func deleteCertManager() {
	packageUninstallOptions := &tkgpackagedatamodel.PackageUninstallOptions{
		PkgInstallName: config.Package.Dependencies[0].Name,
		Namespace: config.PackageRepository.Namespace,
	}
	Expect(packagePlugin.CheckAndUninstallPackage(packageUninstallOptions).Error).NotTo(HaveOccurred())
	Expect(packagePlugin.CheckPackageDeleted(packageUninstallOptions).Error).NotTo(HaveOccurred())
}
