package flownodes

import "github.com/osmoview/logger"

var (
	// interpreters is a map contains virtual machines and other solutions
	// to run and execute logic for custom nodes, key is a language.
	interpreters = map[string]Interpreter{}
)

// NodeRunner is an interface for different virtual machines, interpreters or other solutions to run and execute custom nodes
type Interpreter interface {
	NewNode(nd NodeDriver, data Data, storage Storage, db Database) (Runner, error)
	LintNode(nd NodeDriver) error
	GetBlankCode() string
}

// Runner interface required Run function to execute flow ctx with incoming transaction data
type Runner interface {
	// this function is called every time incoming data(`Ctx`) is received,
	// and should returns response(Data:map[string]any),
	// where the key is a port name, and the value contains the data
	Run(c Ctx) (out Data, err error)

	// Call required function by name
	Call(fnc string, c Ctx) (out Data, err error)

	// Check if a function(or variable) exists by the specified name
	Has(fnc string) (bool, error)

	// close should stop interpreter and clear memory
	Close() error

	// returns last log lines from the interpreter
	Logs() []logger.LogLine
}

// RegisterInterpreter add new interpreter for the specified language
func RegisterInterpreter(lang string, interpreter Interpreter) {
	interpreters[lang] = interpreter
}

// GetInterpreter by language
func GetInterpreter(lang string) (Interpreter, error) {
	i, ok := interpreters[lang]
	if !ok {
		return nil, ErrInterpreterNotExist
	}

	return i, nil
}
