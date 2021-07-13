package e2e_test

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	uexec "k8s.io/utils/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFluentBitE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "fluent-bit Addon Package E2E Test Suite")
}

const (
	deploymentTimeout       = 120 * time.Second
	deploymentCheckInterval = 5 * time.Second
	httpRequestTimeout      = 60 * time.Second
	httpRequestInterval     = 5 * time.Second
)

var (
	// dockerhubProxy is an optional configuration option (provided by using
	// DOCKERHUB_PROXY), that allows you to override docker.io with a proxy to
	// docker.io to avoid any potential issues with rate-limiting.
	dockerhubProxy string

	// packageNamespace is the namespace where the fluent-bit package is
	// installed (i.e this is the namespace tanzu package install is called
	// with)
	packageNamespace string

	// addonNamespace is the namespace where the efluent-bit addon is
	// installed (i.e this is the namespace passed into the fluent-bit addon
	// values.yaml that is provided to the package). This namespace is created
	// by the package installation.
	addonNamespace string
)

var _ = BeforeSuite(func() {
	dockerhubProxy = os.Getenv("DOCKERHUB_PROXY")
	if dockerhubProxy == "" {
		dockerhubProxy = "docker.io"
	}

	suffix := randomSuffix()

	packageNamespace = fmt.Sprintf("e2e-fluent-bit-package-%s", suffix)
	addonNamespace = fmt.Sprintf("e2e-fluent-bit-addon-%s", suffix)

	_, err := kubectl(nil, "create", "namespace", packageNamespace)
	Expect(err).NotTo(HaveOccurred())

})

var _ = AfterSuite(func() {

	_, err := kubectl(nil, "delete", "namespace", packageNamespace)
	Expect(err).NotTo(HaveOccurred())
})

func kubectl(input io.Reader, args ...string) (string, error) {
	return cliRunner("kubectl", input, args...)
}

func tanzu(input io.Reader, args ...string) (string, error) {
	return cliRunner("tanzu", input, args...)
}

func cliRunner(name string, input io.Reader, args ...string) (string, error) {
	fmt.Fprintf(GinkgoWriter, "+ %s %s\n", name, strings.Join(args, " "))
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdin = input
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		rc := -1
		if ee, ok := err.(*exec.ExitError); ok {
			rc = ee.ExitCode()
		}

		fmt.Fprintln(GinkgoWriter, stderr.String())
		return "", uexec.CodeExitError{
			Err:  errors.New(stderr.String()),
			Code: rc,
		}
	}

	fmt.Fprintln(GinkgoWriter, stdout.String())
	return stdout.String(), nil
}

func readFileAndReplaceContents(filename string, findReplaceMap map[string]string) (string, error) {
	byteContents, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	contents := string(byteContents)
	for k, v := range findReplaceMap {
		contents = strings.Replace(contents, k, v, -1)
	}

	return contents, nil
}

func readFileAndReplaceContentsTempFile(filename string, findReplaceMap map[string]string) (string, error) {
	contents, err := readFileAndReplaceContents(filename, findReplaceMap)
	if err != nil {
		return "", err
	}

	file, err := ioutil.TempFile("", fmt.Sprintf("%s-*%s", strings.TrimPrefix(filepath.Base(filename), filepath.Ext(filename)), filepath.Ext(filename)))
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write([]byte(contents))
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func validateDeploymentReady(namespace, name string) {
	EventuallyWithOffset(1, func() (string, error) {
		return kubectl(nil, "-n", namespace, "get", "deployment", name, "-o", "jsonpath={.status.conditions[?(@.type == 'Available')].status}")
	}, deploymentTimeout, deploymentCheckInterval).Should(Equal("True"), fmt.Sprintf("%s/%s deployment was never ready", namespace, name))
}

func validatePodReady(namespace, name string) {
	EventuallyWithOffset(1, func() (string, error) {
		return kubectl(nil, "-n", namespace, "get", "pod", name, "-o", "jsonpath={.status.conditions[?(@.type == 'Ready')].status}")
	}, deploymentTimeout, deploymentCheckInterval).Should(Equal("True"), fmt.Sprintf("%s/%s pod was never ready", namespace, name))
}

func validatePackageReady(namespace, name string) {
	EventuallyWithOffset(1, func() (string, error) {
		return kubectl(nil, "-n", namespace, "get", "installedpackage", name, "-o", "jsonpath={.status.conditions[?(@.type == 'ReconcileSucceeded')].status}")
	}, deploymentTimeout, deploymentCheckInterval).Should(Equal("True"), fmt.Sprintf("%s/%s installedpackage was never ready", namespace, name))
}

func validateDeploymentNotFound(namespace, name string) {
	Eventually(func() error {
		_, err := kubectl(nil, "-n", namespace, "get", "deployment", name)
		return err
	}, deploymentTimeout, deploymentCheckInterval).Should(MatchError(Or(
		ContainSubstring(fmt.Sprintf(`deployments.apps %q not found`, name)),
		ContainSubstring(fmt.Sprintf(`namespaces %q not found`, namespace)),
	)), fmt.Sprintf("%s/%s deployment was never deleted", namespace, name))
}

func validatePodNotFound(namespace, name string) {
	Eventually(func() error {
		_, err := kubectl(nil, "-n", namespace, "get", "pod", name)
		return err
	}, deploymentTimeout, deploymentCheckInterval).Should(MatchError(Or(
		ContainSubstring(fmt.Sprintf(`pods %q not found`, name)),
		ContainSubstring(fmt.Sprintf(`namespaces %q not found`, namespace)),
	)), fmt.Sprintf("%s/%s pod was never deleted", namespace, name))
}

func validatePackageNotFound(namespace, name string) {
	Eventually(func() error {
		_, err := kubectl(nil, "-n", namespace, "get", "installedpackage", name)
		return err
	}, deploymentTimeout, deploymentCheckInterval).Should(MatchError(Or(
		ContainSubstring(fmt.Sprintf(`installedpackages.install.package.carvel.dev %q not found`, name)),
		ContainSubstring(fmt.Sprintf(`namespaces %q not found`, namespace)),
	)), fmt.Sprintf("%s/%s installedpackage was never deleted", namespace, name))
}

func randomSuffix() string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	return fmt.Sprintf("%x", b)
}
