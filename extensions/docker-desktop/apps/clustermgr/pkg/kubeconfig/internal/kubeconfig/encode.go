// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kubeconfig

import (
	"bytes"

	yaml "gopkg.in/yaml.v3"
	kubeyaml "sigs.k8s.io/yaml"

	"github.com/vmware-tanzu/community-edition/errors"
)

// Encode encodes the cfg to yaml
func Encode(cfg *Config) ([]byte, error) {
	// NOTE: kubernetes's yaml library doesn't handle inline fields very well
	// so we're not using that to marshal
	encoded, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, errors.NewEncodeFail(err, "failed to encode KUBECONFIG file")
	}

	// normalize with kubernetes's yaml library
	// this is not strictly necessary, but it ensures minimal diffs when
	// modifying kubeconfig files, which is nice to have
	encoded, err = normYaml(encoded)
	if err != nil {
		return nil, errors.NewNormalizationFailed(err, "failed to normalize KUBECONFIG encoding")
	}

	return encoded, nil
}

// normYaml round trips yaml bytes through sigs.k8s.io/yaml to normalize them
// versus other kubernetes ecosystem yaml output
func normYaml(y []byte) ([]byte, error) {
	var unstructured interface{}
	if err := kubeyaml.Unmarshal(y, &unstructured); err != nil {
		return nil, err
	}
	encoded, err := kubeyaml.Marshal(&unstructured)
	if err != nil {
		return nil, err
	}
	// special case: don't write anything when empty
	if bytes.Equal(encoded, []byte("{}\n")) {
		return []byte{}, nil
	}
	return encoded, nil
}
