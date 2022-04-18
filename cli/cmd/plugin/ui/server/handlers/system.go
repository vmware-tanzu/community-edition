// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/pkg/containerruntime"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/pkg/system"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/cri"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/edition"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/features"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/provider"
)

// GetContainerRuntimeInfo retrieves information about the docker engine configuration.
func (app *App) GetContainerRuntimeInfo(gcrip cri.GetContainerRuntimeInfoParams) middleware.Responder {
	runtimeInfo, err := containerruntime.GetRuntimeInfo()
	if err != nil {
		return cri.NewGetContainerRuntimeInfoBadRequest().WithPayload(Err(err))
	}

	// convert our internal object to the expected API object
	info := &models.RuntimeInfo{
		Architecture: runtimeInfo.Architecture,
		Containers:   int64(runtimeInfo.Containers),
		CPU:          int64(runtimeInfo.CPU),
		Memory:       runtimeInfo.Memory,
		Name:         runtimeInfo.Name,
		Ostype:       runtimeInfo.OSType,
		Osversion:    runtimeInfo.OSVersion,
	}
	return cri.NewGetContainerRuntimeInfoOK().WithPayload(info)
}

// FeatureFlags gets the set feature flags of the system.
func (app *App) FeatureFlags(params features.GetFeatureFlagsParams) middleware.Responder {
	ff, err := system.FeatureFlags()
	if err != nil {
		return features.NewGetFeatureFlagsBadRequest().WithPayload(Err(err))
	}

	return features.NewGetFeatureFlagsOK().WithPayload(ff)
}

// Edition gets the set edition (TCE, TKG, etc.) of the system.
func (app *App) Edition(params edition.GetTanzuEditionParams) middleware.Responder {
	tanzuEdition, err := system.Edition()
	if err != nil {
		return edition.NewGetTanzuEditionBadRequest().WithPayload(Err(err))
	}

	return edition.NewGetTanzuEditionOK().WithPayload(tanzuEdition)
}

// Providers gets the available infrastructure providers.
func (app *App) Providers(params provider.GetProviderParams) middleware.Responder {
	_, err := system.Providers()
	if err != nil {
		return provider.NewGetProviderBadRequest().WithPayload(Err(err))
	}

	// TODO: Need to change this API to return all available providers, if we can
	// return provider.NewGetProviderOK().WithPayload(tanzuEdition)
	return middleware.NotImplemented("operation provider.GetProvider has not yet been implemented")
}
