package kapp

import (
	"fmt"

	goUi "github.com/cppforlife/go-cli-ui/ui"
	. "github.com/k14s/imgpkg/pkg/imgpkg/cmd"
	"github.com/k14s/ytt/pkg/cmd/template"
	"github.com/k14s/ytt/pkg/cmd/ui"
	"github.com/k14s/ytt/pkg/files"
)

func GetKappImage(bundle string) error {
	confUI := goUi.NewConfUI(goUi.NewNoopLogger())
	defer confUI.Flush()

	po := NewPullOptions(confUI)
	po.BundleFlags = BundleFlags{
		Bundle: bundle,
	}
	po.BundleRecursiveFlags = BundleRecursiveFlags{
		Recursive: true,
	}
	po.OutputPath = "/tmp/kapp-img"

	err := po.Run()
	if err != nil {
		return err
	}

	return nil
}

func GetKappYaml(kappPath string) ([][]byte, error) {
	filesToProcess, err := files.NewSortedFilesFromPaths([]string{kappPath + "/config"}, files.SymlinkAllowOpts{})
	if err != nil {
		return nil, err
	}

	o := template.NewOptions()
	out := o.RunWithFiles(template.Input{Files: filesToProcess}, ui.NewTTY(false))
	if out.Err != nil {
		return nil, fmt.Errorf("Evaluating kapp ytt: %s", out.Err)
	}

	if len(out.Files) == 0 {
		return nil, fmt.Errorf("Expected to find kapp yaml files but saw zero files after ytt processing")
	}

	processedBytes := [][]byte{}
	for _, file := range out.Files {
		processedBytes = append(processedBytes, file.Bytes())
	}

	return processedBytes, nil
}
