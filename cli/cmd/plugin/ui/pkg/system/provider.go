// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package system

type ProviderInfo struct {
	Name    string
	Version string
}

// Providers gets the available cluster infrastructure providers.
func Providers() (*[]ProviderInfo, error) {
	result := &[]ProviderInfo{}

	// TODO: Need to see if we can get all available providers
	// defaultTKRBom, err := tkgconfigbom.New(GetConfigDir(), nil).GetDefaultTkrBOMConfiguration()
	// if err != nil {
	// 	return result, err
	// }

	return result, nil
}
