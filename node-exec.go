package flownodes

import (
	"bytes"
	"os/exec"
	"text/template"
)

func init() {
	RegisterNode("exec", &Exec{})
}

type Exec struct {
	Exec string `yaml:"exec"`

	in     Port
	stdout Port
	stderr Port
}

func (e *Exec) Init(db DB) error {
	e.in = Port{Name: "*", Type: PortMap}
	e.stdout = Port{Name: "stdout", Type: PortString}
	e.stderr = Port{Name: "stderr", Type: PortString}
	return nil
}

func (e *Exec) Execute(c Ctx) (Resp, error) {
	args, err := e.parseArgs(c.ToMap())
	if err != nil {
		return nil, err
	}

	// execute required command and receive stdout and stderr
	cmd := exec.Command("bash", "-c", args)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	// returns stdout, stderr, and an error if this occurred during execution
	return Resp{e.stdout.Name: stdout.String(), e.stderr.Name: stderr.String()}, err
}

func (e *Exec) parseArgs(data map[string]any) (args string, err error) {
	tmpl, err := template.New("").Parse(e.Exec)
	if err != nil {
		return args, err
	}

	var b bytes.Buffer
	if err = tmpl.Execute(&b, data); err != nil {
		return args, err
	}

	// parse exec line to the args slice
	// return shellwords.Parse(b.String())
	return b.String(), nil
}

func (e *Exec) GetPorts() (in, out []Port) {
	return []Port{e.in}, []Port{e.stdout, e.stderr}
}
