// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"encoding/base64"
	"errors"
	"net/url"
	"strconv"

	"github.com/go-openapi/runtime/middleware"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/pkg/containerruntime"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/docker"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/client"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgconfigproviders"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgconfigupdater"
)

// CheckIfDockerDaemonAvailable gets a status if the docker engine appears to be running. If not, returns error.
func (app *App) CheckIfDockerDaemonAvailable(params docker.CheckIfDockerDaemonAvailableParams) middleware.Responder {
	// Getting our runtime info will either return the docker info or an error if it can't be retrieved
	_, err := containerruntime.GetRuntimeInfo()
	if err != nil {
		return docker.NewCheckIfDockerDaemonAvailableBadRequest().WithPayload(Err(err))
	}
	return docker.NewCheckIfDockerDaemonAvailableOK()
}

// ApplyTKGConfigForDocker applies the TKG configuration for Docker.
func (app *App) ApplyTKGConfigForDocker(params docker.ApplyTKGConfigForDockerParams) middleware.Responder {
	tkgClient, err := app.getTkgClient()
	if err != nil {
		return docker.NewApplyTKGConfigForDockerInternalServerError().WithPayload(Err(err))
	}

	err = app.saveConfig(params.Params, tkgClient)
	if err != nil {
		return docker.NewApplyTKGConfigForDockerInternalServerError().WithPayload(Err(err))
	}

	return docker.NewApplyTKGConfigForDockerOK().WithPayload(&models.ConfigFileInfo{Path: app.getFilePathForSavingConfig()})
}

// CreateDockerManagementCluster creates a new management cluster using the CAPD provider.
func (app *App) CreateDockerManagementCluster(params docker.CreateDockerManagementClusterParams) middleware.Responder {
	tkgClient, err := app.getTkgClient()
	if err != nil {
		return docker.NewCreateDockerManagementClusterInternalServerError().WithPayload(Err(err))
	}

	err = app.saveConfig(params.Params, tkgClient)
	if err != nil {
		return docker.NewCreateDockerManagementClusterInternalServerError().WithPayload(Err(err))
	}

	initOptions := &client.InitRegionOptions{
		InfrastructureProvider: "docker",
		ClusterName:            params.Params.ClusterName,
		Plan:                   "dev",
		Annotations:            params.Params.Annotations,
		Labels:                 params.Params.Labels,
		ClusterConfigFile:      app.getFilePathForSavingConfig(),
		Edition:                "tce",
	}

	if err := tkgClient.ConfigureAndValidateManagementClusterConfiguration(initOptions, false); err != nil {
		return docker.NewCreateDockerManagementClusterInternalServerError().WithPayload(Err(err))
	}

	go app.StartSendingLogsToUI()
	go createManagementCluster(tkgClient, initOptions)

	return docker.NewCreateDockerManagementClusterOK().WithPayload("started creating management cluster")
}

// ExportDockerConfig returns the configuration content based on the passed in configuration values.
func (app *App) ExportDockerConfig(params docker.ExportTKGConfigForDockerParams) middleware.Responder {
	var configString string

	// Initialize our configuration data
	config, err := paramsToDockerConfig(params.Params)
	if err != nil {
		docker.NewExportTKGConfigForDockerInternalServerError().WithPayload(Err(err))
	}

	configString, err = transformConfigToString(config)
	if err != nil {
		docker.NewExportTKGConfigForDockerInternalServerError().WithPayload(Err(err))
	}

	return docker.NewExportTKGConfigForDockerOK().WithPayload(configString)
}

// saveConfig parses and saves a cluster configuration.
func (app *App) saveConfig(params *models.DockerManagementClusterParams, tkgClient *client.TkgClient) error {
	config, err := paramsToDockerConfig(params)
	if err != nil {
		return err
	}

	configReaderWriter := tkgClient.TKGConfigReaderWriter()
	err = tkgconfigupdater.SaveConfig(app.getFilePathForSavingConfig(), configReaderWriter, config)
	if err != nil {
		return err
	}

	return nil
}

// paramsToDockerConfig uses the input parameters to instatiate a DockerConfig.
func paramsToDockerConfig(params *models.DockerManagementClusterParams) (*tkgconfigproviders.DockerConfig, error) {
	var err error
	res := &tkgconfigproviders.DockerConfig{
		ClusterName:            params.ClusterName,
		InfrastructureProvider: "docker",
		ClusterPlan:            "dev", // "dev" is the only plan supported for docker
		ClusterCIDR:            params.Networking.ClusterPodCIDR,
		ServiceCIDR:            params.Networking.ClusterServiceCIDR,
		HTTPProxyEnabled:       falseStr,
	}

	if params.CeipOptIn != nil {
		res.CeipParticipation = strconv.FormatBool(*params.CeipOptIn)
	}

	if params.IdentityManagement != nil {
		res.IdentityManagementType = *params.IdentityManagement.IdmType
		res.OIDCProviderName = params.IdentityManagement.OidcProviderName
		res.OIDCIssuerURL = params.IdentityManagement.OidcProviderURL.String()
		res.OIDCClientID = params.IdentityManagement.OidcClientID
		res.OIDCClientSecret = params.IdentityManagement.OidcClientSecret
		res.OIDCScopes = params.IdentityManagement.OidcScope
		res.OIDCGroupsClaim = params.IdentityManagement.OidcClaimMappings["groups"]
		res.OIDCUsernameClaim = params.IdentityManagement.OidcClaimMappings["username"]
		res.LDAPBindDN = params.IdentityManagement.LdapBindDn
		res.LDAPBindPassword = params.IdentityManagement.LdapBindPassword
		res.LDAPHost = params.IdentityManagement.LdapURL
		res.LDAPUserSearchBaseDN = params.IdentityManagement.LdapUserSearchBaseDn
		res.LDAPUserSearchFilter = params.IdentityManagement.LdapUserSearchFilter
		res.LDAPUserSearchUsername = params.IdentityManagement.LdapUserSearchUsername
		res.LDAPUserSearchNameAttr = params.IdentityManagement.LdapUserSearchNameAttr
		res.LDAPGroupSearchBaseDN = params.IdentityManagement.LdapGroupSearchBaseDn
		res.LDAPGroupSearchFilter = params.IdentityManagement.LdapGroupSearchFilter
		res.LDAPGroupSearchUserAttr = params.IdentityManagement.LdapGroupSearchUserAttr
		res.LDAPGroupSearchGroupAttr = params.IdentityManagement.LdapGroupSearchGroupAttr
		res.LDAPGroupSearchNameAttr = params.IdentityManagement.LdapGroupSearchNameAttr
		res.LDAPRootCAData = base64.StdEncoding.EncodeToString([]byte(params.IdentityManagement.LdapRootCa))
	}

	if params.Networking.HTTPProxyConfiguration != nil && params.Networking.HTTPProxyConfiguration.Enabled {
		res.HTTPProxyEnabled = trueStr
		conf := params.Networking.HTTPProxyConfiguration
		res.ClusterHTTPProxy, err = checkAndGetProxyURL(conf.HTTPProxyUsername, conf.HTTPProxyPassword, conf.HTTPProxyURL)
		if err != nil {
			return res, err
		}
		res.ClusterHTTPSProxy, err = checkAndGetProxyURL(conf.HTTPSProxyUsername, conf.HTTPSProxyPassword, conf.HTTPSProxyURL)
		if err != nil {
			return res, err
		}
		res.ClusterNoProxy = params.Networking.HTTPProxyConfiguration.NoProxy
	}

	if params.MachineHealthCheckEnabled {
		res.MachineHealthCheckEnabled = trueStr
	} else {
		res.MachineHealthCheckEnabled = falseStr
	}

	return res, nil
}

// checkAndGetProxyURL validates and returns the proxy URL
func checkAndGetProxyURL(username, password, proxyURL string) (string, error) {
	httpURL, err := url.Parse(proxyURL)
	if err != nil {
		return "", err
	}

	if httpURL.Scheme == "" {
		return "", errors.New("scheme is missing from the proxy URL")
	}

	if httpURL.Host == "" {
		return "", errors.New("hostname is missing from the proxy URL")
	}

	if username != "" && password != "" {
		httpURL.User = url.UserPassword(username, password)
	} else if username != "" {
		httpURL.User = url.User(username)
	}

	return httpURL.String(), nil
}
