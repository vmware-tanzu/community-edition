// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// package e2e_test implements running the external DNS end to end tests
package e2e_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v1 "k8s.io/api/storage/v1"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestHarborE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Harbor Addon Package E2E Test Suite")
}

type packageDependency struct {
	Name        string // package installed name
	DisplayName string // the display name of the package, used to find the package name
	ValuesFile  string // the custom values file when installing the package
}

const (
	packageWaitCheckInterval = "10s"
	packageWaitTimeout       = "20m"
	provisionerjsonpath      = `jsonpath='{.items[0].provisioner}'`
)

var (
	// packageInstallNamespace is the namespace where the harbor package
	// is installed (i.e this is the namespace tanzu package install is called
	// with). Optionally provided by using PACKAGE_INSTALL_NAMESPACE env var.
	// If PACKAGE_INSTALL_NAMESPACE is not provided the test assumes the available
	// package is in the global namespace and the test will create and install
	// the package into its own test namespace.
	// If PACKAGE_INSTALL_NAMESPACE is provided the test assumes
	// the available package is namespaced and will install the package into
	// the provided namespace and will NOT delete that namespace.
	// Read more on package namespacing here:
	// https://carvel.dev/kapp-controller/docs/latest/package-consumer-concepts/#overview
	packageInstallNamespace string

	// packageComponentsNamespace is the namespace where the harbor
	// package components are installed  (e.g. the  deployment).
	// This is the namespace passed into the harbor values.yaml). This
	// namespace is created by the package installation.
	packageComponentsNamespace string

	// packageInstallName is the app name of the harbor installed
	packageInstallName string

	// packageDependencies is the dependencies for the harbor package
	packageDependencies []*packageDependency

	// harborHostname the hostname of the harbor installed
	harborHostname string

	// harborAdminPassword the password of the harbor installed
	harborAdminPassword string

	// installedPackages record the packages installed by the test
	installedPackages []string

	// storage class name for PVC usage
	storageclass string
	isDefault    bool
)

var _ = BeforeSuite(func() {
	packageInstallNamespace = "test"

	packageInstallName = "harbor"

	packageComponentsNamespace = "goharbor"

	_, err := utils.Kubectl(nil, "create", "ns", packageInstallNamespace)
	Expect(err).NotTo(HaveOccurred())

	//storageclass = getStorageClass()
	storageclass, isDefault = getAvailableStorageClass()
	if isDefault {
		storageclass = ""
	}
	//update storageclass if provisioner is ebs
	updateStorageClass()

	packageDependencies = []*packageDependency{
		{"cert-manager", "cert-manager", ""},
		{"contour", "contour", configContourYamlFile()},
	}

	for _, dependency := range packageDependencies {
		By(fmt.Sprintf("installing %s addon package", dependency.Name))

		packageName := utils.TanzuPackageName(dependency.DisplayName)

		version := findPackageAvailableVersion(packageName, "")
		installPackage(dependency.Name, packageName, version, dependency.ValuesFile)
	}

	By("installing harbor addon package")
	harborHostname = fmt.Sprintf("harbor.%s.nip.io", getReachableContourEnvoyIP())

	harborAdminPassword = generateHarborPassword()

	packageName := utils.TanzuPackageName("harbor")

	version := findPackageAvailableVersion(packageName, "2.6.2")

	valuesFilename, err := utils.ReadFileAndReplaceContentsTempFile(filepath.Join("fixtures", "harbor.yaml"),
		map[string]string{
			"PACKAGE_COMPONENTS_NAMESPACE": packageComponentsNamespace,
			"harbor.yourdomain.com":        harborHostname,
			"Harbor12345":                  harborAdminPassword,
			"STORAGE_CLASS":                storageclass,
		},
	)
	Expect(err).NotTo(HaveOccurred())
	defer os.Remove(valuesFilename)

	installPackage(packageInstallName, packageName, version, valuesFilename)

	By("validating harbor package is ready")
	utils.ValidatePackageInstallReady(packageInstallNamespace, packageInstallName)

	By("validating harbor components are healthy")
	validateHarborHealthy(harborHostname)
})

var _ = AfterSuite(func() {
	for _, installedPackage := range installedPackages {
		By(fmt.Sprintf("cleaning up %s addon package", installedPackage))
		_, err := utils.Tanzu(nil, "package", "installed", "delete", installedPackage,
			"--wait-check-interval", packageWaitCheckInterval,
			"--wait-timeout", packageWaitTimeout,
			"--namespace", packageInstallNamespace, "--yes")
		Expect(err).NotTo(HaveOccurred())
	}

	By("validating the harbor package install no longer exists")
	utils.ValidatePackageInstallNotFound(packageInstallNamespace, packageInstallName)

	By(fmt.Sprintf("cleaning up %s namespace", packageComponentsNamespace))
	utils.Kubectl(nil, "delete", "ns", packageComponentsNamespace) // nolint:errcheck
})

func findPackageAvailableVersion(packageName string, versionSubstr string) string {
	packageVersionJSON, err := utils.Tanzu(nil, "package", "available", "list", packageName, "-o", "json")
	Expect(err).NotTo(HaveOccurred())
	versions := []map[string]string{}

	err = json.Unmarshal([]byte(packageVersionJSON), &versions)
	Expect(err).NotTo(HaveOccurred())
	Expect(len(versions)).To(BeNumerically(">", 0))

	var matchedVersions []string
	for _, v := range versions {
		if versionSubstr == "" || strings.Contains(v["version"], versionSubstr) {
			matchedVersions = append(matchedVersions, v["version"])
		}
	}

	Expect(len(matchedVersions)).To(BeNumerically(">", 0), fmt.Sprintf("version contains %s for package %s not found", versionSubstr, packageName))

	return matchedVersions[len(matchedVersions)-1]
}

/*
apply a "rancher.io/local-path" storageclass for clusters lack of csi driver to dynamically provisioning for PVC usages
clusters included and before tkg1.5(k8s1.22) have in-tree cloud provider plugin
vsphere clusters have own storageclass with csi "csi.vsphere.vmware.com"

only when clusters without csi driver neither have a local-path-storage would have to install it
Temporarily, aws and azure clusters would apply this "rancher.io/local-path" sc instead of the default one

return storageclass name for PVC and boolean isDefaultStorageClass
*/
func getAvailableStorageClass() (string, bool) {
	name, provisioner := getDefaultStorageClass()
	// clusters using kubernetes verions prior to 1.23
	if (getKubernetesVersion() < "1.23") && (provisioner != "") {
		return name, true
	}
	// already has a local-path storageclass or a CSIDriver
	if provisioner == "rancher.io/local-path" || hasCSIDriver(provisioner) {
		return name, true
	}

	// do not have an available storageclass, apply a local-path-stoarge for it
	_, err := utils.Kubectl(nil, "apply", "-f", filepath.Join("fixtures", "local-path-storage.yaml"))
	Expect(err).NotTo(HaveOccurred())
	_, err = utils.Kubectl(nil, "wait", "pod", "-n", "local-path-storage", "-l", "app=local-path-provisioner", "--for", "condition=Ready", "--timeout=300s")
	Expect(err).NotTo(HaveOccurred())

	return "local-path", false
}

func hasCSIDriver(provisioner string) bool {
	csidriver, err := utils.Kubectl(nil, "get", "csidriver", "-o", "json")
	Expect(err).NotTo(HaveOccurred())
	return strings.Contains(csidriver, provisioner)
}

func getKubernetesVersion() string {
	jsonStr, _ := utils.Kubectl(nil, "version", "-o", "json")
	versionmap := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &versionmap)
	Expect(err).NotTo(HaveOccurred())

	for k, v := range versionmap {
		if k == "serverVersion" {
			serverVersionMap := v.(map[string]interface{})
			kubernetesVersion := serverVersionMap["major"].(string) + "." + serverVersionMap["minor"].(string)
			fmt.Println("cluster KubernetesVersion:", kubernetesVersion)
			return kubernetesVersion
		}
	}

	return ""
}

func getDefaultStorageClass() (string, string) {
	jsonPath := `jsonpath={.items}`
	jsonStr, err := utils.Kubectl(nil, "get", "storageclasses", "-o", jsonPath)
	Expect(err).NotTo(HaveOccurred())
	storageclasses := []v1.StorageClass{}
	err = json.Unmarshal([]byte(jsonStr), &storageclasses)
	Expect(err).NotTo(HaveOccurred())

	for _, sc := range storageclasses {
		if sc.GetObjectMeta().GetAnnotations()["storageclass.kubernetes.io/is-default-class"] == "true" {
			return sc.GetObjectMeta().GetName(), sc.Provisioner
		}
	}

	return "", ""
}

func getContourEnvoyLoadBalancerIP() string {
	converts := map[string]func(string) string{
		"ip": func(ip string) string {
			return ip
		},
		"hostname": func(hostname string) string {
			var ip string

			Eventually(func(g Gomega) {
				addr, err := net.LookupIP(hostname)
				g.Expect(err).NotTo(HaveOccurred())
				Expect(len(addr)).To(BeNumerically(">", 0))

				ip = addr[0].String()
			}, time.Second*300, time.Second*5).Should(Succeed(), fmt.Sprintf("failed to lookup ip for %s", hostname))

			return ip
		},
	}

	for key, convert := range converts {
		jsonPath := fmt.Sprintf(`jsonpath='{.status.loadBalancer.ingress[0].%s}'`, key)
		output, err := utils.Kubectl(nil, "-n", "projectcontour", "get", "svc", "envoy", "-o", jsonPath)
		Expect(err).NotTo(HaveOccurred())

		output = strings.ReplaceAll(output, `'`, "")
		if output != "" {
			return convert(output)
		}
	}

	return ""
}

func getContourEnvoyHostIP() string {
	jsonPath := `jsonpath='{.items[0].status.hostIP}'`
	output, err := utils.Kubectl(nil, "-n", "projectcontour", "get", "pod", "-l", "app=envoy", "-o", jsonPath)
	Expect(err).NotTo(HaveOccurred())

	output = strings.ReplaceAll(output, `'`, "")
	Expect(output).NotTo(BeEmpty(), "hostIP of envoy pod not found")

	return output
}

// getReachableContourEnvoyIP returns ip address of the load balancer when it's public,
// otherwise returns the ip address of node host
func getReachableContourEnvoyIP() string {
	if ip := getContourEnvoyLoadBalancerIP(); ip != "" && !isPrivate(net.ParseIP(ip)) {
		return ip
	}

	return getContourEnvoyHostIP()
}

func getHTTPClient(insecure bool) *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure}, // nolint:gosec
	}

	return &http.Client{Transport: transport}
}

func generateHarborPassword() string {
	var (
		upperBytes  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowerBytes  = "abcdefghijklmnopqrstuvwxyz"
		numberBytes = "01234567889"
	)

	validBytes := []string{upperBytes, lowerBytes, numberBytes}

	choice := func(seq string) byte {
		return seq[rand.Intn(len(seq))] // nolint:gosec
	}

	b := []byte{
		choice(upperBytes),
		choice(lowerBytes),
		choice(numberBytes),
	}

	for i := 0; i < 5; i++ {
		bytes := validBytes[rand.Intn(len(validBytes))] // nolint:gosec

		b = append(b, choice(bytes))
	}

	return string(b)
}

func installPackage(name, packageName, version, valuesFilename string) {
	installedPackages = append([]string{name}, installedPackages...)

	args := []string{
		"package", "install", name,
		"--wait-check-interval", packageWaitCheckInterval,
		"--wait-timeout", packageWaitTimeout,
		"--namespace", packageInstallNamespace,
		"--package", packageName,
		"--version", version,
	}

	if valuesFilename != "" {
		args = append(args, "--values-file", valuesFilename)
	}

	_, err := utils.Tanzu(nil, args...)
	Expect(err).NotTo(HaveOccurred())
}

func validateHarborHealthy(hostname string) {
	names := []string{
		"core", "database", "jobservice", "notary",
		"portal", "redis", "registry", "registryctl", "trivy",
	}

	var components []interface{}
	for _, name := range names {
		components = append(components, map[string]string{
			"name":   name,
			"status": "healthy",
		})
	}

	js := map[string]interface{}{
		"status":     "healthy",
		"components": components,
	}

	healthBody, err := json.Marshal(js)
	Expect(err).NotTo(HaveOccurred())

	client := getHTTPClient(true)

	Eventually(func(g Gomega) {
		url := fmt.Sprintf("https://%s/api/v2.0/health", hostname)
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, nil)
		g.Expect(err).NotTo(HaveOccurred())

		resp, err := client.Do(req) // nolint:bodyclose
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(resp).Should(HaveHTTPBody(MatchJSON(healthBody)))
	}, time.Second*600, time.Second*5).Should(Succeed(), "harbor is unhealthy")
}

func isPrivate(ip net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil {
		return ip4[0] == 10 ||
			(ip4[0] == 172 && ip4[1]&0xf0 == 16) ||
			(ip4[0] == 192 && ip4[1] == 168)
	}

	return len(ip) == net.IPv6len && ip[0]&0xfe == 0xfc
}

func configContourYamlFile() string {
	// jsonPath := `jsonpath='{.items[0].provisioner}'` // nolint:goconst
	output, _ := utils.Kubectl(nil, "get", "storageclasses", "-o", provisionerjsonpath)
	//if provider is vc change contour service type to NodePort
	//otherwise use LoadBalancer
	if strings.Contains(output, "csi.vsphere.vmware.com") {
		valuesFilename, err := utils.ReadFileAndReplaceContentsTempFile(filepath.Join("fixtures", "contour.yaml"),
			map[string]string{
				"LoadBalancer": "NodePort",
			},
		)
		Expect(err).NotTo(HaveOccurred())
		return valuesFilename
	}
	return filepath.Join("fixtures", "contour.yaml")
}

func updateStorageClass() {
	// jsonPath := `jsonpath='{.items[0].provisioner}'`
	output, _ := utils.Kubectl(nil, "get", "storageclasses", "-o", provisionerjsonpath)
	if strings.Contains(output, "ebs.csi.aws.com") {
		jsonpath := `jsonpath='{.items[*].metadata.labels.topology\.ebs\.csi\.aws\.com/zone}'`
		zones, _ := utils.Kubectl(nil, "get", "node", "-o", jsonpath)
		zone := strings.Split(strings.Trim(zones, "'"), " ")[0]
		// combine and update storageclass file
		origin, _ := utils.Kubectl(nil, "get", "storageclasses", "-o", "yaml")
		lines := strings.Split(origin, "\n")
		var newlines []string
		added := false
		for _, line := range lines {
			newlines = append(newlines, line)
			if strings.Contains(line, "allowVolumeExpansion") && !added {
				s := "  allowedTopologies:\n  - matchLabelExpressions:\n    - key: failure-domain.beta.kubernetes.io/zone\n      values:\n      - " + zone
				newlines = append(newlines, s)
				added = true
			}
		}
		//write file
		filename := "ebs-storageclass.yaml"
		file, err := os.CreateTemp("", fmt.Sprintf("%s-*%s", strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename)), filepath.Ext(filename)))
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()
		_, err = file.WriteString(strings.Join(newlines, "\n"))
		Expect(err).NotTo(HaveOccurred())
		// apply storage class with new file
		_, err = utils.Kubectl(nil, "delete", "storageclasses", "default")
		Expect(err).NotTo(HaveOccurred())
		_, err = utils.Kubectl(nil, "apply", "-f", file.Name())
		Expect(err).NotTo(HaveOccurred())
	}
}
