package flownodes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/osmoview/flowdata"
)

var (
	Dir = "./nodes"

	// IgnoreDirs contains a list of directory names that should be skipped
	IgnoreDirs = []string{"vendor", "gopath", "deps"}

	// NodeDrivers contains all available(registered) nodes
	// WARN: maybe required to add sync.Mutex to access to the Nodes map
	nodeDrivers sync.Map // map[string]NodeDriver{}
)

const (
	AccessReadOnly = "read-only" // prevent to change code or a driver
	AccessInternal = "internal"  // internal packages are not exported
	AccessPrivate  = "private"   // private packages cannot be modified or viewed
)

// NodeDriver structure describe how to run and communicate with the node
type NodeDriver struct {
	// node id, should be unique
	ID string `yaml:"id" json:"id"`

	// path to the directory with config file
	Path string `yaml:"-" json:"path,omitempty"`

	// human readable node driver name
	Name string `yaml:"name" json:"name"`

	// INFO: node field is deprecated, use ID instead
	// Node string `yaml:"node" json:"-"`
	Lang string `yaml:"lang,omitempty" json:"lang,omitempty"`

	// access level
	Access string `yaml:"access,omitempty" json:"access,omitempty"`

	IsDataSource bool `yaml:"is_datasource,omitempty" json:"is_datasource,omitempty"`

	Label flowdata.L10n `yaml:"label,omitempty" json:"label,omitempty"`
	Help  flowdata.L10n `yaml:"help,omitempty" json:"help,omitempty"`
	Group string        `yaml:"group,omitempty" json:"group,omitempty"`

	Input  []Port               `yaml:"input,omitempty" json:"input,omitempty"`
	Output []Port               `yaml:"output,omitempty" json:"output,omitempty"`
	Form   []flowdata.FormField `yaml:"form,omitempty" json:"form,omitempty"`

	// INFO: Run section describe how to run this node, and may contains reqired for interpreter(related to 'lang') custom data, like working directory, path to the dependencies directory, etc
	Run struct {
		// path to the source code file, relative to the .yaml file
		File string `yaml:"file" json:"file"`
		// data may contains any optional fields required by different interpreters
		Data `yaml:",inline" json:"data"`
	} `yaml:"run" json:"run"`

	// Errors []string `yaml:"-" json:"errors,omitempty"`
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
	return nd.Path
}

// GetRunFilePath returns path to the source code file or executable file
func (nd NodeDriver) GetRunFilePath() string {
	return filepath.Join(filepath.Dir(nd.GetNodeDriverFilePath()), nd.Run.File)
}

func (nd NodeDriver) GetRunTarget() (string, bool) {
	v, ok := nd.Run.Data["target"]
	if !ok {
		return "", false
	}

	s, ok := v.(string)
	if !ok {
		return "", false
	}

	s = strings.TrimSpace(s)

	return s, s != ""
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

func RangeNodeDrivers(f func(NodeDriver) bool) {
	nodeDrivers.Range(func(k, v any) bool {
		return f(v.(NodeDriver))
	})
}

// GetNodeDriver by id from the loaded nodes
func GetNodeDriver(id string) (NodeDriver, error) {
	v, ok := nodeDrivers.Load(id)
	if !ok {
		return NodeDriver{}, fmt.Errorf("%s %w", id, ErrNodeNotFound)
	}

	return v.(NodeDriver), nil
}

// SetNodeDriver save node driver to the global node drivers storage
func SetNodeDriver(nd NodeDriver) {
	// fix node driver path
	if !strings.HasPrefix(nd.Path, Dir) {
		nd.Path = filepath.Join(filepath.Clean(Dir), nd.Path)
	}

	nodeDrivers.Store(nd.ID, nd)
}

func DelNodeDriver(id string) {
	nodeDrivers.Delete(id)
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

	Dir = dir

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

		relPath := strings.TrimPrefix(fp, path)
		relPath = strings.TrimPrefix(relPath, string(filepath.Separator))

		// useful for symlink directory, trim symlink path and substitute relative
		fp = filepath.Join(dir, relPath)

		// read and check node driver file
		nd, err := ReadNodeDriver(fp)
		if err != nil {
			return fmt.Errorf("failed to load '%s' node: %v", fp, err)
		}

		nd.Path = fp

		SetNodeDriver(nd)

		return nil
	})
}

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

	nd.Path = fp
	if nd.Group == "" {
		if rel, err := filepath.Rel(Dir, filepath.Dir(fp)); err == nil {
			nd.Group = filepath.Base(rel)
		}
	}

	return nd, nd.IsValid()
}
