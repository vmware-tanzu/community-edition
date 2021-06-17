package e2e

import (
	"io/ioutil"
	"path"

	"sigs.k8s.io/yaml"
)

type Config struct {
	IAAS              string             `json:"iaas"`
	ClusterType       string             `json:"clusterType,omitempty"`
	Package           *Package           `json:"package"`
	PackageRepository *PackageRepository `json:"packageRepository"`
}

type Package struct {
	PackageMetadata `json:",inline"`
	Dependencies    []PackageMetadata `json:"dependencies"`
}

type PackageMetadata struct {
	Name    string `json:"name"`
	RefName string `json:"refName"`
	Version string `json:"version"`
}

type PackageRepository struct {
	Name         string `json:"name"`
	ImgpkgBundle string `json:"imgpkgBundle"`
	Namespace    string `json:"namespace"`
}

func loadE2EConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = path.Join(currentDir, "config/e2e_config.yaml")
	}
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	if err := yaml.Unmarshal(configData, config); err != nil {
		return nil, err
	}
	if config.ClusterType == "" {
		config.ClusterType = "workload"
	}
	return config, nil
}
