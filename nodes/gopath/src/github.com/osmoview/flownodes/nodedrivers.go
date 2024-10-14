package flownodes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/osmoview/flowdata"

	"gopkg.in/yaml.v3"
)

var (
	// IgnoreDirs contains a list of directory names that should be skipped
	IgnoreDirs = []string{"vendor", "gopath", "deps"}

	// NodeDrivers contains all available(registered) nodes
	// WARN: maybe required to add sync.Mutex to access to the Nodes map
	NodeDrivers = map[string]NodeDriver{}
)

// NodeDriver structure describe how to run and communicate with the node
type NodeDriver struct {
	filepath string

	// node id, should be unique
	ID string `yaml:"id" json:"id"`

	// path to the directory with config file
	Path string `yaml:"-" json:"path,omitempty"`

	// human readable node driver name
	Name string `yaml:"name" json:"name"`

	// INFO: node field is deprecated, use ID instead
	// Node string `yaml:"node" json:"-"`
	Lang string `yaml:"lang,omitempty" json:"lang,omitempty"`

	IsDataSource bool `yaml:"is_datasource,omitempty" json:"is_datasource,omitempty"`

	Desc  flowdata.L10n `yaml:"desc,omitempty" json:"desc,omitempty"`
	Label flowdata.L10n `yaml:"label,omitempty" json:"label,omitempty"`
	Group string        `yaml:"group,omitempty" json:"group,omitempty"`

	Input  []Port               `yaml:"input,omitempty" json:"input,omitempty"`
	Output []Port               `yaml:"output,omitempty" json:"output,omitempty"`
	Form   []flowdata.FormField `yaml:"form,omitempty" json:"form,omitempty"`

	// INFO: Run section describe how to run this node, and may contains reqired for interpreter(related to 'lang') custom data, like working directory, path to the dependencies directory, etc
	Run struct {
		// path to the source code file, relative to the .yaml file
		File string `yaml:"file" json:"file"`
		// data may contains any optional fields required by different interpreters
		Data `yaml:",inline"`
	} `yaml:"run" json:"run"`

	Errors []string `yaml:"-" json:"errors,omitempty"`
}

// IsValid check if node driver is valid or not
func (nd NodeDriver) IsValid() error {
	if nd.ID == "" {
		return fmt.Errorf("'id' is required and should be unique")
	}
	if nd.Lang == "" {
		return fmt.Errorf("'lang' is required")
	}
	if nd.Run.File == "" {
		return fmt.Errorf("there is no path to the executable file")
	}

	return nil
}

// GetNodeDriverFilePath returns path to the .yaml file
func (nd NodeDriver) GetNodeDriverFilePath() string {
	return nd.filepath
}

// GetRunFilePath returns path to the source code file or executable file
func (nd NodeDriver) GetRunFilePath() string {
	return filepath.Join(nd.Path, nd.Run.File)
}

// Start node with specified data, returns runner interface with 'Run' method
func (nd NodeDriver) Start(data Data, storage Storage, db Database) (Runner, error) {
	i, err := GetInterpreter(nd.Lang)
	if err != nil {
		return nil, err
	}

	return i.NewNode(nd, data, storage, db)
}

//
//
//

// ReadNodeDriver from the .yaml file
func ReadNodeDriver(fp string) (nd NodeDriver, err error) {
	f, err := os.Open(fp)
	if err != nil {
		return nd, err
	}
	defer f.Close()

	if yaml.NewDecoder(f).Decode(&nd); err != nil {
		return nd, err
	}

	nd.filepath = fp
	nd.Path = filepath.Dir(fp)

	return nd, nd.IsValid()
}

// GetNodeDriver by id from the loaded nodes
func GetNodeDriver(node string) (NodeDriver, error) {
	n, ok := NodeDrivers[node]
	if !ok {
		return NodeDriver{}, fmt.Errorf("%s %w", node, ErrNodeNotFound)
	}

	return n, nil
}

// LoadNodes should find and load analytics nodes dynamically connected as modules.
// Such nodes can be written in various programming languages,
// Each "node" should be presented, at least, as a `*.yaml` file containing the necessary information.
func LoadNodeDrivers(dir string) error {
	// find the real path if dir is a symlink
	path, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return err
	}

	return filepath.Walk(path, func(fp string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		for _, pattern := range IgnoreDirs {
			if strings.Contains(fp, pattern) {
				return filepath.SkipDir
			}
		}

		if fi.IsDir() || (!strings.HasSuffix(fp, ".yaml") && !strings.HasSuffix(fp, ".yml")) {
			return nil
		}

		// useful for symlink directory, trim symlink path and substitute relative
		fp = filepath.Join(dir, strings.TrimPrefix(fp, path))

		// read and check node driver file
		drv, err := ReadNodeDriver(fp)
		if err != nil {
			return fmt.Errorf("failed to load '%s' node: %v", fp, err)
		}

		NodeDrivers[drv.ID] = drv

		return nil
	})
}
