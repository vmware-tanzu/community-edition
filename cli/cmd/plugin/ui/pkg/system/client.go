// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package system

import (
	"fmt"
	"path/filepath"

	"github.com/aunum/log"
	"github.com/pkg/errors"

	"github.com/vmware-tanzu/tanzu-framework/apis/config/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/config"
	providergetterclient "github.com/vmware-tanzu/tanzu-framework/pkg/v1/providers/client"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/client"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/clientcreator"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/clusterclient"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/constants"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/providerinterface"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/region"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgconfigupdater"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgctl"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/types"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/utils"
)

// NewTKGClient gets a new client for interacting with the system.
func NewTKGClient(logFile string, logLevel int32) (tkgctl.TKGClient, error) {
	configDir := GetConfigDir()
	xclient, err := tkgctl.New(tkgctl.Options{
		ConfigDir: configDir,
		CustomizerOptions: types.CustomizerOptions{
			RegionManagerFactory: NewFactory(),
		},
		LogOptions: tkgctl.LoggingOptions{Verbosity: logLevel, File: logFile},
	})
	if err != nil {
		return nil, err
	}

	return xclient, err
}

// NewTkgClient gets a new client for interacting with the system. Good luck
// figuring out if you need a TKGClient or a TkgClient.
func NewTkgClient() (*client.TkgClient, error) {
	configDir := GetConfigDir()
	appConfig := types.AppConfig{
		TKGConfigDir:   configDir,
		ProviderGetter: providergetterclient.New(),
		CustomizerOptions: types.CustomizerOptions{
			RegionManagerFactory: region.NewFactory(),
		},
	}

	err := ensureTKGConfigFile(configDir, appConfig.ProviderGetter)
	if err != nil {
		return nil, err
	}

	allClients, err := clientcreator.CreateAllClients(appConfig, nil)
	if err != nil {
		return nil, err
	}

	tkgClient, err := client.New(client.Options{
		ClusterCtlClient:         allClients.ClusterCtlClient,
		ReaderWriterConfigClient: allClients.ConfigClient,
		RegionManager:            allClients.RegionManager,
		TKGConfigDir:             configDir,
		Timeout:                  constants.DefaultOperationTimeout,
		FeaturesClient:           allClients.FeaturesClient,
		TKGConfigProvidersClient: allClients.TKGConfigProvidersClient,
		TKGBomClient:             allClients.TKGBomClient,
		TKGConfigUpdater:         allClients.TKGConfigUpdaterClient,
		TKGPathsClient:           allClients.TKGConfigPathsClient,
		ClusterKubeConfig:        &types.ClusterKubeConfig{},
		ClusterClientFactory:     clusterclient.NewClusterClientFactory(),
		FeatureFlagClient:        allClients.FeatureFlagClient,
	})
	if err != nil {
		return nil, err
	}

	return tkgClient, nil
}

func ensureTKGConfigFile(configDir string, providerGetter providerinterface.ProviderInterface) error {
	var err error

	lock, err := utils.GetFileLockWithTimeOut(filepath.Join(configDir, constants.LocalTanzuFileLock), utils.DefaultLockTimeout)
	if err != nil {
		return errors.Wrap(err, "cannot acquire lock for ensuring local files")
	}

	defer func() {
		if err := lock.Unlock(); err != nil {
			log.Warningf("cannot release lock for ensuring local files, reason: %v", err)
		}
	}()

	_, err = tkgconfigupdater.New(configDir, providerGetter, nil).EnsureTKGConfigFile()
	return err
}

type tanzuRegionManager struct {
}

type tanzuRegionManagerFactory struct {
}

// NewFactory creates a new tanzuRegionManagerFactory which implements
// region.ManagerFactory
func NewFactory() region.ManagerFactory {
	return &tanzuRegionManagerFactory{}
}

// ListRegionContexts will list all region contexts.
func (trm *tanzuRegionManager) ListRegionContexts() ([]region.RegionContext, error) {
	tanzuConfig, err := config.GetClientConfig()
	if err != nil {
		return []region.RegionContext{}, err
	}

	var regionClusters []region.RegionContext
	for _, server := range tanzuConfig.KnownServers {
		if server.Type == v1alpha1.ManagementClusterServerType {
			regionContext := convertServerToRegionContextFull(server,
				server.Name == tanzuConfig.CurrentServer)

			regionClusters = append(regionClusters, regionContext)
		}
	}

	return regionClusters, nil
}

// SaveRegionContext saves the RegionContext.
func (trm *tanzuRegionManager) SaveRegionContext(regionCtxt region.RegionContext) error {
	return config.AddServer(convertRegionContextToServer(regionCtxt), false)
}

// UpstartRegionContext inserts(?) a new region context.
func (trm *tanzuRegionManager) UpsertRegionContext(regionCtxt region.RegionContext) error {
	return config.PutServer(convertRegionContextToServer(regionCtxt), false)
}

// DeleteRegionContext removes a region context.
func (trm *tanzuRegionManager) DeleteRegionContext(clusterName string) error {
	currentServer, err := config.GetCurrentServer()
	if err != nil {
		return err
	}

	if clusterName != "" && clusterName != currentServer.Name {
		return fmt.Errorf("cannot delete cluster %s, it is not the current cluster", clusterName)
	}

	if err := config.RemoveServer(currentServer.Name); err != nil {
		return err
	}

	return nil
}

// SetCurrentContext sets the context being used.
func (trm *tanzuRegionManager) SetCurrentContext(clusterName, contextName string) error {
	return config.SetCurrentServer(clusterName)
}

// GetCurrentContext will get the currently set context.
func (trm *tanzuRegionManager) GetCurrentContext() (region.RegionContext, error) {
	currentServer, err := config.GetCurrentServer()
	if err != nil {
		return region.RegionContext{}, err
	}

	if !currentServer.IsManagementCluster() {
		return region.RegionContext{}, errors.Errorf("The current server is not a management cluster")
	}

	return convertServerToRegionContext(currentServer), nil
}

// CreateManager creates a manager.
func (trmf *tanzuRegionManagerFactory) CreateManager(configFile string) (region.Manager, error) {
	return &tanzuRegionManager{}, nil
}

func convertServerToRegionContext(server *v1alpha1.Server) region.RegionContext {
	return convertServerToRegionContextFull(server, false)
}

func convertServerToRegionContextFull(server *v1alpha1.Server, isCurrentContext bool) region.RegionContext {
	return region.RegionContext{
		ClusterName:      server.Name,
		ContextName:      server.ManagementClusterOpts.Context,
		SourceFilePath:   server.ManagementClusterOpts.Path,
		IsCurrentContext: isCurrentContext,
	}
}

func convertRegionContextToServer(regionContext region.RegionContext) *v1alpha1.Server {
	return &v1alpha1.Server{
		Name: regionContext.ClusterName,
		Type: v1alpha1.ManagementClusterServerType,
		ManagementClusterOpts: &v1alpha1.ManagementClusterServer{
			Path:    regionContext.SourceFilePath,
			Context: regionContext.ContextName,
		},
	}
}
