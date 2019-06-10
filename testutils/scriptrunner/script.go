/*
Package scriptrunner provides a way of running a series of cli steps to compare the output from a single file.
*/
package scriptrunner

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// RunFunction is called when a command starting with '#' is found.
type RunFunction func(inFile io.Reader, args []string) error

type scriptRunner struct {
	filename string
	commands []*scriptCommand
	runFunc  RunFunction
	w        io.Writer
}

func (sr *scriptRunner) loadCommands(filename string) error {
	sr.filename = filename
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "# ") {
			sr.commands = append(sr.commands, parseCommand(line))
		}
		if strings.HasPrefix(line, "//") {
			sr.commands[len(sr.commands)-1].appendComment(line)
		}
	}
	return nil
}

// NewScriptRunner creates a scriptRunner from a file.
func NewScriptRunner(f string, run RunFunction, w io.Writer) (*scriptRunner, error) {
	sr := &scriptRunner{
		runFunc: run,
		w:       w,
	}
	if err := sr.loadCommands(f); err != nil {
		return nil, err
	}
	return sr, nil
}

// ReplaceArg replaces arguments before running.
//
// This enables dynamic attributes, such as a URL, to be specified at runtime but not change the golden file output.
func ReplaceArg(args []string, find, replace string) []string {
	out := make([]string, len(args))
	for i, s := range args {
		if s == find {
			out[i] = replace
		} else {
			out[i] = s
		}
	}
	return out
}

// Run runs all commands in script and produces the output
func (sr *scriptRunner) Run() error {
	for _, command := range sr.commands {
		var f io.ReadCloser
		if command.inputFile != "" {
			var err error
			if f, err = os.Open(command.inputFile); err != nil {
				return err
			}
			defer f.Close()
		}
		sr.w.Write([]byte(command.String()))
		if err := sr.runFunc(f, command.command); err != nil {
			fmt.Fprintf(sr.w, "ERROR: %v\n", err)
		}
		sr.w.Write([]byte("\n"))
	}
	return nil
}
