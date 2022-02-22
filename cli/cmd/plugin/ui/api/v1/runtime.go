// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1

import (
	"encoding/json"
	"net/http"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/objects"
)

// GetContainerRuntime gets information about the container runtime. An empty
// response indicates there is no runtime available (not installed or not running).
func GetContainerRuntime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	response := objects.Runtime{
		Name:       "test",
		OSType:     "linux",
		OSVersion:  "20.04",
		CPU:        4,
		Memory:     6235168768,
		Containers: 1,
	}

	_ = json.NewEncoder(w).Encode(response)
}
