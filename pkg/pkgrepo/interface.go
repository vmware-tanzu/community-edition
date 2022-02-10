package main

type Interface interface {
	// Name is the package repo
	Name() string

	// OutputPath is the directory above a package repo
	OutputPath() string

	// ImgpkgPath returns the location to save imgpkg metadata
	ImgpkgPath() string

	// PkgPath returns the location to save package metadata
	PkgPath()

	// CreateRepo creates a repository on the local filesystem
	CreateRepo(packagesDir string) error

	GenerateFileSystem() error

	ReadPackages() error
}
