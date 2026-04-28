package cli

import (
	"io"

	"github.com/tasuku43/cc-bash-guard/internal/app"
)

const (
	exitAllow  = 0
	exitError  = 1
	exitReject = 2
)

type Streams struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type Env = app.Env

func Run(args []string, streams Streams, env Env) int {
	if len(args) == 0 {
		writeUsage(streams.Stdout)
		return exitError
	}

	switch args[0] {
	case "hook":
		return runHook(args[1:], streams, env)
	case "explain":
		return runExplain(args[1:], streams, env)
	case "suggest":
		return runSuggest(args[1:], streams)
	case "doctor":
		return runDoctor(args[1:], streams, env)
	case "verify":
		return runVerify(args[1:], streams, env)
	case "init":
		return runInit(args[1:], streams, env)
	case "version":
		return runVersion(args[1:], streams)
	case "semantic-schema":
		return runSemanticSchema(args[1:], streams)
	case "-h", "--help", "help":
		return writeHelp(streams.Stdout, streams.Stderr, args[1:])
	default:
		writeErr(streams.Stderr, "unknown command: "+args[0])
		return exitError
	}
}
