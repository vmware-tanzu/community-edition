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
	packagePollInterval = "5s"
	packagePollTimeout  = "10m"
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
)

var _ = BeforeSuite(func() {
	packageInstallNamespace = "default"

	packageInstallName = "harbor"

	packageComponentsNamespace = "goharbor"

	packageDependencies = []*packageDependency{
		{"cert-manager", "cert-manager", ""},
		{"contour", "contour", configContourYamlFile()},
	}

	if !hasDefaultStorageClass() {
		packageDependencies = append(packageDependencies,
			&packageDependency{"local-path-storage", "local-path-storage", ""},
		)
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

	version := findPackageAvailableVersion(packageName, "2.3.3")

	valuesFilename, err := utils.ReadFileAndReplaceContentsTempFile(filepath.Join("fixtures", "harbor.yaml"),
		map[string]string{
			"PACKAGE_COMPONENTS_NAMESPACE": packageComponentsNamespace,
			"harbor.yourdomain.com":        harborHostname,
			"Harbor12345":                  harborAdminPassword,
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
			"--poll-interval", packagePollInterval,
			"--poll-timeout", packagePollTimeout,
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

func hasDefaultStorageClass() bool {
	output, err := utils.Kubectl(nil, "get", "storageclasses")
	Expect(err).NotTo(HaveOccurred())

	return strings.Contains(output, "(default)")
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
			}, time.Second*120, time.Second*5).Should(Succeed(), fmt.Sprintf("failed to lookup ip for %s", hostname))

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
		"--poll-interval", packagePollInterval,
		"--poll-timeout", packagePollTimeout,
		"--namespace", packageInstallNamespace,
		"--package-name", packageName,
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
	jsonPath := `jsonpath='{.items[0].provisioner}'`
	output, _ := utils.Kubectl(nil, "get", "storageclasses", "-o", jsonPath)
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
