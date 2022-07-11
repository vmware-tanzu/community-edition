// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/pkg/content"
	"oras.land/oras-go/pkg/oras"

	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Harbor Addon E2E Test", func() {
	var (
		ociClient *OCIClient
	)

	BeforeEach(func() {
		ociClient = &OCIClient{Host: harborHostname, Username: "admin", Password: harborAdminPassword}
	})

	JustAfterEach(func() {
		if CurrentGinkgoTestDescription().Failed {
			fmt.Fprintf(GinkgoWriter, "\nCollecting diagnostic information just after test failure\n")
			fmt.Fprintf(GinkgoWriter, "\nResources summary:\n")
			_, _ = utils.Kubectl(nil, "-n", packageInstallNamespace, "get", "all,packageinstalls,apps")
			_, _ = utils.Kubectl(nil, "-n", packageComponentsNamespace, "get", "all")

			fmt.Fprintf(GinkgoWriter, "\npackage install status:\n")
			_, _ = utils.Kubectl(nil, "-n", packageInstallNamespace, "get", "app", packageInstallName, "-o", "jsonpath={.status}")

			fmt.Fprintf(GinkgoWriter, "\npackage components status:\n")
			_, _ = utils.Kubectl(nil, "-n", packageComponentsNamespace, "get", "deployment", "harbor", "-o", "jsonpath={.status}")

			fmt.Fprintf(GinkgoWriter, "\nharbor logs:\n")

			_, _ = utils.Kubectl(nil, "-n", packageComponentsNamespace, "logs", "-l", "app=harbor", "--all-containers")
		}
	})

	It("push image to the library project", func() {
		err := ociClient.Push(context.TODO(), "library/hello-world:v1", []byte("Hello pushing!"))

		Expect(err).NotTo(HaveOccurred())
	})

	It("pull image from the library project", func() {
		ref := "library/hello-world:v2"

		err := ociClient.Push(context.TODO(), ref, []byte("Hello pulling!"))
		Expect(err).NotTo(HaveOccurred())

		err = ociClient.Pull(context.TODO(), ref)
		Expect(err).NotTo(HaveOccurred())
	})
})

type OCIClient struct {
	Host     string
	Username string
	Password string
}

func (o *OCIClient) getResolver() remotes.Resolver {
	client := getHTTPClient(true)

	authorizer := docker.NewAuthorizer(client, func(host string) (string, string, error) {
		if host == o.Host {
			return o.Username, o.Password, nil
		}

		return "", "", nil
	})

	return docker.NewResolver(docker.ResolverOptions{
		Hosts: docker.ConfigureDefaultRegistries(
			docker.WithAuthorizer(authorizer),
			docker.WithClient(client),
		),
	})
}

func (o *OCIClient) Push(ctx context.Context, ref string, blobContent []byte) error {
	if !strings.HasPrefix(ref, o.Host) {
		ref = o.Host + "/" + ref
	}

	sum := sha256.Sum256(blobContent)
	filename := fmt.Sprintf("%x", sum)

	store := content.NewMemoryStore()

	desc := store.Add(filename, "tce-x-harbor", blobContent)

	configBytes, _ := json.Marshal(map[string]interface{}{"sha256": filename})
	config := store.Add("config", ocispec.MediaTypeImageConfig, configBytes)

	contents := []ocispec.Descriptor{desc}

	_, err := oras.Push(ctx, o.getResolver(), ref, store, contents, oras.WithConfig(config))

	return err
}

func (o *OCIClient) Pull(ctx context.Context, ref string) error {
	if !strings.HasPrefix(ref, o.Host) {
		ref = o.Host + "/" + ref
	}

	store := content.NewMemoryStore()

	_, _, err := oras.Pull(ctx, o.getResolver(), ref, store)

	return err
}
