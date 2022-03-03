package main

// Interface represents an interface for reading Carvel packages and generate a package repository.
// Input could be specified via a file system directory or a manifest file.
type Interface interface {
	// OutputPath is the directory above a package repo
	// All repos will be written into the output path, under a directory with their Name
	OutputPath() string

	// Name is the package repo.
	// This should be `core` or something for a user-defined repo
	Name() string

	// AddPackage adds a package's YAML as bytes to the repository.
	// A byte slice is used in order to avoid requiring particular Go types.
	// Errors returned could be related to filesystem reading
	AddPackage(yamlContents []byte) error

	// AddPackages adds multiple packages to the repository at once.
	// Internally, it should call AddPackage with each YAML byte slice.
	AddPackages(allYaml [][]byte) error

	// PackageNames returns the names of the packages that should be included in this repository.
	// An empty slice means all discovered packages.
	PackageNames() []string

	// ImgpkgPath returns the location to save imgpkg metadata
	// This location should be the canonical .imgpkg/images.yml file, but if not it should be available to query.
	ImgpkgPath() string

	// PkgPath returns the location to save package metadata
	// This location is usually pkgs/packages.yml
	PkgPath() string

	// CreateRepo creates a repository on the local filesystem
	CreateRepo(packagesDir string) error

	// GenerateFileSystem creates a filesystem with the layout that is expected for an imgpkg push
	// Creates something like the following on the filesystem:
	//   OutputPath
	//      - Name
	//          - ImgPkgPath
	//          - PkgPath
	GenerateFileSystem() error

	// CombinePackages will read the contents of package metadata.yml and package.yaml files and combine them.
	// Output should be written to the PkgPath() value
	CombinePackages() error

	// PushBundle will push a bundled repository to an OCI registry.
	// Invokes imgpkg's push functionality to upload the ImgPkgPath and PkgPath to the ociRepo
	PushImages(ociRepo string) error
}
