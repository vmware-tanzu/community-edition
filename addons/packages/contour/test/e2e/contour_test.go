package e2e

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-tanzu-private/core/pkg/v1/tkg/tkgpackagedatamodel"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	k8sconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
)

func getEnvoyPorts(service *corev1.Service) (int32, int32, int32, int32) {
	var (
		svcHttpPort, svcHttpsPort, nodeHttpPort, nodeHttpsPort int32
	)

	for _, port := range service.Spec.Ports {
		if port.Name == "http" {
			nodeHttpPort = port.NodePort
			svcHttpPort = port.Port
		} else if port.Name == "https" {
			nodeHttpsPort = port.NodePort
			svcHttpsPort = port.Port
		}
	}

	return svcHttpPort, svcHttpsPort, nodeHttpPort, nodeHttpsPort
}

func getEnvoyContainerPorts(ds *v1.DaemonSet) (int32, int32, int32, int32) {
	var (
		httpHostPort, httpsHostPort           int32
		httpContainerPort, httpsContainerPort int32
	)

	for _, container := range ds.Spec.Template.Spec.Containers {
		if container.Name == "envoy" {
			for _, port := range container.Ports {
				if port.Name == "http" {
					httpHostPort = port.HostPort
					httpContainerPort = port.ContainerPort
				} else if port.Name == "https" {
					httpsHostPort = port.HostPort
					httpsContainerPort = port.ContainerPort
				}
			}
		}
	}

	return httpHostPort, httpsHostPort, httpContainerPort, httpsContainerPort
}

func getEnvoyDetails(clientset *kubernetes.Clientset, useHostPort, useHostNetwork, useLocalProxyPort bool) (string, int32, int32) {
	var (
		address             string
		httpPort, httpsPort int32
		err                 error
	)

	if useLocalProxyPort {
		return "127.0.0.1", 65080, 65443
	}

	if useHostPort || useHostNetwork {
		ds, err := clientset.AppsV1().DaemonSets("tanzu-system-ingress").Get(context.TODO(), "envoy", metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())
		Expect(ds).NotTo(BeNil())

		address, err := getAnyNodeIP(clientset)
		Expect(err).NotTo(HaveOccurred())

		httpHostPort, httpsHostPort, httpContainerPort, httpsContainerPort := getEnvoyContainerPorts(ds)

		if useHostNetwork {
			httpPort = httpContainerPort
			httpsPort = httpsContainerPort
		}
		if useHostPort {
			httpPort = httpHostPort
			httpsPort = httpsHostPort
		}
		return address, httpPort, httpsPort
	}

	address = getIngressHost(clientset)
	Expect(address).NotTo(BeEmpty(), "Ingress host cannot be empty")

	service, err := clientset.CoreV1().Services("tanzu-system-ingress").Get(context.TODO(), "envoy", metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())
	Expect(service).NotTo(BeNil())

	httpPort, httpsPort, _, _ = getEnvoyPorts(service)

	return address, httpPort, httpsPort
}

func getAnyNodeIP(clientset *kubernetes.Clientset) (string, error) {
	var (
		address                string
		externalIP, internalIP string
	)

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return address, err
	}

	if len(nodes.Items) > 0 {
		node := nodes.Items[0]
		for _, addr := range node.Status.Addresses {
			if addr.Type == corev1.NodeExternalIP {
				externalIP = addr.Address
			} else if addr.Type == corev1.NodeInternalIP {
				internalIP = addr.Address
			}
		}
	}
	if externalIP != "" {
		address = externalIP
	} else {
		address = internalIP
	}

	return address, nil
}

func getIngressHost(clientset *kubernetes.Clientset) string {
	var (
		address string
	)

	Eventually(func() string {
		service, err := clientset.CoreV1().Services("tanzu-system-ingress").Get(context.TODO(), "envoy", metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())
		Expect(service).NotTo(BeNil())

		switch service.Spec.Type {
		case corev1.ServiceTypeNodePort:
			address, err = getWorkerNodeIP(clientset)
			Expect(address).NotTo(BeEmpty())
			Expect(err).NotTo(HaveOccurred())
			return address
		case corev1.ServiceTypeLoadBalancer:
			if len(service.Status.LoadBalancer.Ingress) > 0 {
				address = service.Status.LoadBalancer.Ingress[0].Hostname
				if address == "" {
					address = service.Status.LoadBalancer.Ingress[0].IP
				}
			}
			return address
		default:
			Fail("Envoy service type should be nodeport or loadbalancer")
			return ""
		}
	},
		900,
	).ShouldNot(BeEmpty(), "ingress address should not be empty")

	return address
}

func getWorkerNodeIP(clientset *kubernetes.Clientset) (string, error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{LabelSelector: "node-role.kubernetes.io/master!="})
	if err != nil {
		return "", err
	}
	if len(nodes.Items) == 0 {
		return "", fmt.Errorf("no worker nodes found")
	}
	var nodeIP string

	nodeAddresses := nodes.Items[0].Status.Addresses
	for _, addr := range nodeAddresses {
		if addr.Type == corev1.NodeExternalIP {
			nodeIP = addr.Address
			break
		}
	}
	if nodeIP == "" {
		return "", fmt.Errorf("no external IP for nodes")
	}

	return nodeIP, nil
}

func runContourE2ETests(clientset *kubernetes.Clientset, runBasic, useHostPort, useHostNetwork, useProxyLocalPort bool, infraProvider string) {
	fmt.Println("Running E2E tests")

	address, httpPort, httpsPort := getEnvoyDetails(clientset, useHostPort, useHostNetwork, useProxyLocalPort)

	contourE2EScript := path.Join(currentDir, fmt.Sprintf("scripts/run-contour-tests.sh"))
	cmd := exec.Command("bash", "-c", fmt.Sprintf("%s %s", contourE2EScript, infraProvider))
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("ADDRESS=%s", address))
	cmd.Env = append(cmd.Env, fmt.Sprintf("HTTP_PORT=%d", httpPort))
	cmd.Env = append(cmd.Env, fmt.Sprintf("HTTPS_PORT=%d", httpsPort))
	if runBasic {
		cmd.Env = append(cmd.Env, fmt.Sprintf("BASIC_E2E=%t", true))
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		fmt.Println(err.Error())
		Fail(err.Error())
	}
	fmt.Println(string(out))
}

var _ = Describe("CONTOUR e2e test", func() {

	var (
		packageName          string
		installedPackageName string
		packageVersion       string
		namespace            string

		clientset *kubernetes.Clientset

		runE2ETests      bool
		runBasicE2ETests bool
	)

	BeforeEach(func() {
		packageName = config.Package.RefName
		installedPackageName = config.Package.Name
		namespace = config.PackageRepository.Namespace
		packageVersion = config.Package.Version

		cfg, err := k8sconfig.GetConfig()
		Expect(err).ToNot(HaveOccurred(), "failed to load kubeconfig")

		clientset, err = kubernetes.NewForConfig(cfg)
		Expect(err).ToNot(HaveOccurred(), "failed to load clientset")
	})

	AfterEach(func() {
		packageUninstallOptions := &tkgpackagedatamodel.PackageUninstallOptions{
			PkgInstallName: installedPackageName,
			Namespace:        namespace,
		}
		Expect(packagePlugin.CheckAndUninstallPackage(packageUninstallOptions).Error).NotTo(HaveOccurred())
		Expect(packagePlugin.CheckPackageDeleted(packageUninstallOptions).Error).NotTo(HaveOccurred())
	})

	JustBeforeEach(func() {
		valuesFile := path.Join(currentDir, fmt.Sprintf("resources/%s.yaml", config.IAAS))
		packageOptions := &tkgpackagedatamodel.PackageInstalledOptions{
			PkgInstallName: installedPackageName,
			PackageName:      packageName,
			Version:          packageVersion,
			Namespace:        namespace,
			Wait:             true,
			ValuesFile:       valuesFile,
		}
		Expect(packagePlugin.CheckPackageAvailable(packageName, &tkgpackagedatamodel.PackageAvailableOptions{
			Namespace: namespace,
		}).Error).NotTo(HaveOccurred())

		packagePluginResult := packagePlugin.CheckAndInstallPackage(packageOptions)
		Expect(packagePluginResult.Error).NotTo(HaveOccurred())
		Expect(packagePlugin.CheckPackageInstalled(packageOptions).Error).NotTo(HaveOccurred())

		if runE2ETests || runBasicE2ETests {
			useProxyLocalPort := false
			if config.IAAS == "kind" || config.IAAS == "capd" {
				useProxyLocalPort = true
			}
			runContourE2ETests(clientset, runBasicE2ETests, false, false, useProxyLocalPort, config.IAAS)
		}
	})

	Context("Contour e2e", func() {
		BeforeEach(func() {
			runE2ETests = true
		})

		It("runs e2e successfully", func() {
		})
	})
})
