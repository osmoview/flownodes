package flownodes

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mattn/go-shellwords"
	"github.com/osmoview/flowdata"
	"github.com/osmoview/settings"
)

func init() {
	RegisterNode("exec", NewExecNode)
	SourceNodes = append(SourceNodes, "exec")
}

type Exec struct {
	Exec string `yaml:"exec" json:"exec"`
	WD   string `yaml:"wd" json:"wd"`

	in     Port
	stdout Port
	stderr Port
}

func NewExecNode() Node {
	return &Exec{}
}

func (e *Exec) Init(db DB) error {
	e.initPorts()
	return nil
}

func (e *Exec) initPorts() {
	e.in = Port{Name: "*", Type: PortMap}
	e.stdout = Port{Name: "stdout", Type: PortString}
	e.stderr = Port{Name: "stderr", Type: PortString}
}

func (e *Exec) Execute(c Ctx) (Resp, error) {
	args, err := e.parseArgs(c.DataMap())
	if err != nil {
		return nil, err
	}

	// name := filepath.Clean(args[0])
	fp := args[0]
	if e.WD != "" && !filepath.IsAbs(fp) {
		dir, _ := os.Getwd()
		fp = filepath.Join(dir, args[0])
	}

	// execute required command and receive stdout and stderr
	// cmd := exec.Command("bash", "-c", args)
	cmd := exec.Command(fp, args[1:]...)
	if e.WD != "" {
		cmd.Dir = e.WD
	}
	cmd.Env = os.Environ()

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	// TODO: kill timeout
	if err != nil {
		err = fmt.Errorf("%w %s", err, stderr.String())
	}

	// returns stdout, stderr, and an error if this occurred during execution
	return Resp{e.stdout.Name: stdout.String(), e.stderr.Name: stderr.String()}, err
}

func (e *Exec) parseArgs(data map[string]any) (args []string, err error) {
	// substite data to the exec line
	line, err := settings.Map(data).SubstituteTo(e.Exec)
	if err != nil {
		return args, err
	}

	// parse exec line to the args slice
	return shellwords.Parse(line)
}

func (e *Exec) GetPorts() (in, out []Port) {
	e.initPorts()
	return []Port{e.in}, []Port{e.stdout, e.stderr}
}

func (e *Exec) GetForm() flowdata.Form {
	return flowdata.Form{{
		Name:  "exec",
		Label: flowdata.L10n{"en": "Execute", "ru": "Запустить"},
		Help:  flowdata.L10n{"en": "some command to execute, ie: date"},
		Level: flowdata.FieldLevelMaster,
	}}
}

func (e *Exec) GetGroup() string {
	return GroupFiles
}
