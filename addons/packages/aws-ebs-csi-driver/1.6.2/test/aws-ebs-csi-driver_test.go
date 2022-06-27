// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package awsebscsidrivertest

import (
	"path/filepath"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"

	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/repo"
	"github.com/vmware-tanzu/community-edition/addons/packages/test/pkg/ytt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AWS EBS CSI Driver Ytt Templates", func() {
	var (
		filePaths []string
		values    string
		output    string
		err       error

		configDir                      = filepath.Join(repo.RootDir(), "addons/packages/aws-ebs-csi-driver/1.6.2/bundle/config")
		fileAWSEBSCSIDriverYaml        = filepath.Join(configDir, "upstream/aws-ebs-csi-driver.yaml")
		fileAWSEBSCSIDriverOverlayYaml = filepath.Join(configDir, "overlays/update-aws-ebs-csi-driver.yaml")
		fileUpdateStrategyOverlayYaml  = filepath.Join(configDir, "overlays/update-strategy-overlay.yaml")
		fileValuesYaml                 = filepath.Join(configDir, "values.yaml")
		fileValuesStar                 = filepath.Join(configDir, "values.star")
		fileSchemaYaml                 = filepath.Join(configDir, "schema.yaml")
	)

	BeforeEach(func() {
		values = ""
	})

	disiredDeploymentToleration := corev1.Toleration{
		Key:      "node-role.kubernetes.io/master",
		Operator: "Exists",
		Effect:   "NoSchedule",
	}

	JustBeforeEach(func() {
		filePaths = []string{fileAWSEBSCSIDriverYaml, fileAWSEBSCSIDriverOverlayYaml, fileUpdateStrategyOverlayYaml, fileValuesYaml, fileValuesStar, fileSchemaYaml}
		output, err = ytt.RenderYTTTemplate(ytt.CommandOptions{}, filePaths, strings.NewReader(values))
	})

	Context("default configuration", func() {
		It("renders a Deployment with node selector", func() {
			Expect(err).NotTo(HaveOccurred())
			deployment := parseDeployment(output)
			tol := findDeploymentTolerationByKey(deployment.Spec.Template.Spec.Tolerations, disiredDeploymentToleration.Key)
			Expect(tol).NotTo(BeNil())
			Expect(tol.Operator).To(Equal(disiredDeploymentToleration.Operator))
			Expect(tol.Effect).To(Equal(disiredDeploymentToleration.Effect))
		})
	})

	Context("configure nodeSelector and updateStrategy", func() {
		BeforeEach(func() {
			values = `#@data/values
---
nodeSelector:
  tanzuKubernetesRelease: 1.22.3
deployment:
  updateStrategy: RollingUpdate
  rollingUpdate:
    maxUnavailable: 0
    maxSurge: 1
daemonset:
  updateStrategy: OnDelete
`
		})
		It("renders the DaemonSet and Deployment with desired nodeSelector and updateStrategy", func() {
			Expect(err).NotTo(HaveOccurred())
			daemonSet := parseDaemonSet(output)
			deployment := parseDeployment(output)
			Expect(deployment.Spec.Template.Spec.NodeSelector).ToNot(BeNil())
			Expect(deployment.Spec.Template.Spec.NodeSelector["tanzuKubernetesRelease"]).To(Equal("1.22.3"))
			Expect(deployment.Spec.Strategy.Type).To(Equal(appsv1.RollingUpdateDeploymentStrategyType))
			Expect(deployment.Spec.Strategy.RollingUpdate).ToNot(BeNil())
			Expect(deployment.Spec.Strategy.RollingUpdate.MaxUnavailable.IntVal).To(Equal(int32(0)))
			Expect(deployment.Spec.Strategy.RollingUpdate.MaxSurge.IntVal).To(Equal(int32(1)))
			Expect(daemonSet.Spec.UpdateStrategy.Type).To(Equal(appsv1.OnDeleteDaemonSetStrategyType))
		})
	})

	Context("configure namespace", func() {
		BeforeEach(func() {
			values = `#@data/values
---
awsEBSCSIDriver:
  namespace: test-namespace
`
		})
		It("renders the DaemonSet and Deployment with desired namespace", func() {
			Expect(err).NotTo(HaveOccurred())
			daemonSet := parseDaemonSet(output)
			deployment := parseDeployment(output)
			Expect(deployment.ObjectMeta.Namespace).To(Equal("test-namespace"))
			Expect(daemonSet.ObjectMeta.Namespace).To(Equal("test-namespace"))
		})

		It("renders the ConfigMap with desired namespace", func() {
			cms := unmarshalConfigMaps(output)
			for _, cm := range cms {
				Expect(cm.ObjectMeta.Namespace).To(Equal("test-namespace"))
			}
		})

		It("renders the ServiceAccount with desired namespace", func() {
			sas := unmarshalServiceAccount(output)
			for _, sa := range sas {
				if sa.Kind == "ServiceAccount" {
					Expect(string(sa.ObjectMeta.Namespace)).To(Equal("test-namespace"))
				}
			}
		})

	})

})

func parseDaemonSet(output string) appsv1.DaemonSet {
	daemonSetDocIndex := 14
	daemonSetDoc := strings.Split(output, "---")[daemonSetDocIndex]
	var daemonSet appsv1.DaemonSet
	err := yaml.Unmarshal([]byte(daemonSetDoc), &daemonSet)
	Expect(err).NotTo(HaveOccurred())
	return daemonSet
}

func parseDeployment(output string) appsv1.Deployment {
	deploymentDocIndex := 12
	deploymentDoc := strings.Split(output, "---")[deploymentDocIndex]
	var deployment appsv1.Deployment
	err := yaml.Unmarshal([]byte(deploymentDoc), &deployment)
	Expect(err).NotTo(HaveOccurred())
	return deployment
}

func findDeploymentTolerationByKey(tolerations []corev1.Toleration, key string) *corev1.Toleration {
	for _, tol := range tolerations {
		if tol.Key == key {
			return &tol
		}
	}
	return nil
}

func unmarshalConfigMaps(output string) []corev1.ConfigMap {
	docs := findDocsWithString(output, "kind: ConfigMap")
	cms := make([]corev1.ConfigMap, len(docs))
	for i, doc := range docs {
		var cm corev1.ConfigMap
		err := yaml.Unmarshal([]byte(doc), &cm)
		Expect(err).NotTo(HaveOccurred())
		cms[i] = cm
	}
	return cms
}

func unmarshalServiceAccount(output string) []corev1.ServiceAccount {
	docs := findDocsWithString(output, "kind: ServiceAccount")
	sas := make([]corev1.ServiceAccount, len(docs))
	for i, doc := range docs {
		var sa corev1.ServiceAccount
		err := yaml.Unmarshal([]byte(doc), &sa)
		Expect(err).NotTo(HaveOccurred())
		sas[i] = sa
	}
	return sas
}

func findDocsWithString(output, selector string) []string {
	var docs []string
	for _, doc := range strings.Split(output, "---") {
		if strings.Contains(doc, selector) {
			docs = append(docs, doc)
		}
	}
	return docs
}
