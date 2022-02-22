// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"encoding/json"
	"net/http"
)

// SupportedAPIVersions handles requests to get the supported API versions by the server.
func SupportedAPIVersions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	response := []string{"v1"}

	_ = json.NewEncoder(w).Encode(response)
}
