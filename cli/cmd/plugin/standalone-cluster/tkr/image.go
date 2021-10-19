package tkr

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	goUi "github.com/cppforlife/go-cli-ui/ui"
	. "github.com/k14s/imgpkg/pkg/imgpkg/cmd"
	"github.com/k14s/ytt/pkg/cmd/template"
	"github.com/k14s/ytt/pkg/cmd/ui"
	"github.com/k14s/ytt/pkg/files"
)

type TkrImage struct {
	RegistryUrl  string
	DownloadPath string
	ConfigPath   string

	YttValuesFilePath string
	YttRenderedBytes  [][]byte
}

// TkrImageReader enables operations on indivdual image bundles that are referenced from the TKR bom
type TkrImageReader interface {
	// GetRegistryUrl returns the bundle's registry URL
	GetRegistryUrl() string

	// DownloadBundleImage downloads the OCI image bundle using imgpkg libraries.
	// The unpacked bundle's files are stored in a temporary directory
	DownloadBundleImage() error

	// SetRelativeConfigPath sets the _relative_ path for the YTT config bundle in the downloaded OCI image.
	// Example: kapp controller stores it's YTT bundle under "config/" in it's bundle.
	//          So therefore, this function should be called with "config/" as an argument
	SetRelativeConfigPath(string)

	// SetYttValuesFilePath sets the fully realized file path for the intended values yaml for the bundle
	SetYttValuesFilePath(string)

	// RenderYaml renders the OCI bundle using ytt libraries. The returned slice of byte slices contain the rendered yaml
	RenderYaml() ([][]byte, error)
}

// NewTkrImageReader provides a new TkrImageReader through the TkrImage struct
// and automatically populates a temporary directory to download the OCI image
func NewTkrImageReader(imagePath string) (TkrImageReader, error) {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}

	defer os.RemoveAll(tempDir)

	t := &TkrImage{
		RegistryUrl:  imagePath,
		DownloadPath: tempDir,
	}

	return t, nil
}

func (t *TkrImage) GetRegistryUrl() string {
	return t.RegistryUrl
}

func (t *TkrImage) DownloadBundleImage() error {
	confUI := goUi.NewConfUI(goUi.NewNoopLogger())
	defer confUI.Flush()

	po := NewPullOptions(confUI)
	po.BundleFlags = BundleFlags{
		Bundle: t.RegistryUrl,
	}
	po.BundleRecursiveFlags = BundleRecursiveFlags{
		Recursive: true,
	}
	po.OutputPath = t.DownloadPath

	err := po.Run()
	if err != nil {
		return err
	}

	return nil
}

func (t *TkrImage) SetRelativeConfigPath(configPath string) {
	t.ConfigPath = filepath.Join(t.DownloadPath, configPath)
}

func (t *TkrImage) SetYttValuesFilePath(filePath string) {
	t.YttValuesFilePath = filePath
}

func (t *TkrImage) RenderYaml() ([][]byte, error) {
	filesToProcess, err := files.NewSortedFilesFromPaths([]string{t.ConfigPath}, files.SymlinkAllowOpts{})
	if err != nil {
		return nil, err
	}

	o := template.NewOptions()

	if t.YttValuesFilePath != "" {
		o.DataValuesFlags.KVsFromYAML = []string{t.YttValuesFilePath}
	}

	out := o.RunWithFiles(template.Input{Files: filesToProcess}, ui.NewTTY(false))
	if out.Err != nil {
		return nil, fmt.Errorf("Evaluating ytt: %s", out.Err)
	}

	if len(out.Files) == 0 {
		return nil, fmt.Errorf("Expected to find yaml files but saw zero files after ytt processing")
	}

	processedBytes := [][]byte{}
	for _, file := range out.Files {
		processedBytes = append(processedBytes, file.Bytes())
	}

	// This sets the in the image reader itself so they may be referenced elsewhere.
	// These bytes are also returned from this function call
	t.YttRenderedBytes = processedBytes

	return processedBytes, nil
}
