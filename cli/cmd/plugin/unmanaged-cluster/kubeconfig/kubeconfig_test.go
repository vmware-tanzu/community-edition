// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kubeconfig

import (
	"os"
	"testing"
)

var stubKubeconfig = `---
apiVersion: v1
clusters:
- cluster:
    certificate-authority: /home/my-user/.minikube/ca.crt
    extensions:
    - extension:
        last-update: Mon, 02 May 2022 13:23:44 MDT
        provider: minikube.sigs.k8s.io
        version: v1.25.2
      name: cluster_info
    server: https://192.168.49.2:8443
  name: minikube
contexts:
- context:
    cluster: minikube
    extensions:
    - extension:
        last-update: Mon, 02 May 2022 13:23:44 MDT
        provider: minikube.sigs.k8s.io
        version: v1.25.2
      name: context_info
    namespace: default
    user: minikube
  name: my-context
current-context: my-context
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate: /home/my-user/.minikube/profiles/minikube/client.crt
    client-key: /home/my-user/.minikube/profiles/minikube/client.key`

func TestGetKubeconfigContext(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "kubeconfig-test-")
	if err != nil {
		t.Errorf(err.Error())
	}

	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(stubKubeconfig))
	if err != nil {
		t.Errorf(err.Error())
	}

	context, err := GetKubeconfigContext(tmpFile.Name())
	if err != nil {
		t.Errorf("failed to get context. Error: %s", err.Error())
	}

	if context != "my-context" {
		t.Errorf("did not get correct context. Actual context: %s", context)
	}
}

func TestGetKubeconfigContextEmpty(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "kubeconfig-test-")
	if err != nil {
		t.Errorf(err.Error())
	}

	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte("---"))
	if err != nil {
		t.Errorf(err.Error())
	}

	context, err := GetKubeconfigContext(tmpFile.Name())
	if err != nil {
		t.Errorf("failed to get context. Error: %s", err.Error())
	}

	if context != "" {
		t.Errorf("did not get empty context. Actual context: %s", context)
	}
}

var testJSON = `{
  "json": [
    "rigid",
    "better for data interchange"
  ],
  "yaml": [
    "slim and flexible",
    "better for configuration"
  ],
  "object": {
    "key": "value",
    "array": [
      {
        "null_value": null
      },
      {
        "boolean": true
      },
      {
        "integer": 1
      },
      {
        "alias": "aliases are like variables"
      },
      {
        "alias": "aliases are like variables"
      }
    ]
  },
  "paragraph": "Blank lines denote\nparagraph breaks\n",
  "content": "Or we\ncan auto\nconvert line breaks\nto save space",
  "alias": {
    "bar": "baz"
  },
  "alias_reuse": {
    "bar": "baz"
  }
}`

var testYamls = `alias:
    bar: baz
alias_reuse:
    bar: baz
content: |-
    Or we
    can auto
    convert line breaks
    to save space
json:
    - rigid
    - better for data interchange
object:
    array:
        - null_value: null
        - boolean: true
        - integer: 1
        - alias: aliases are like variables
        - alias: aliases are like variables
    key: value
paragraph: |
    Blank lines denote
    paragraph breaks
yaml:
    - slim and flexible
    - better for configuration
`

func TestJSONToYAML(t *testing.T) {
	yamls, err := jSONToYAML([]byte(testJSON))
	if err != nil {
		t.Errorf("could not convert json to yaml. Error: %s", err.Error())
	}

	if string(yamls) != testYamls {
		t.Errorf("Incorrect yaml returned. Expected ----- :\n%s\nActual ----- :\n%s\n", testYamls, string(yamls))
	}
}
