package scriptrunner

import (
	"fmt"
	"strings"
)

func parseCommand(s string) *scriptCommand {
	fields := strings.Fields(s)

	lastElement := fields[len(fields)-1]
	if strings.HasPrefix(lastElement, "<") {
		return &scriptCommand{
			inputFile: lastElement[1:],
			command:   fields[1 : len(fields)-1],
		}
	}

	return &scriptCommand{
		command: fields[1:],
	}
}

type scriptCommand struct {
	inputFile string
	command   []string
	comments  []string
}

func (sc *scriptCommand) appendComment(s string) {
	sc.comments = append(sc.comments, s)
}

// String produces command and comment lines for the command.
func (sc *scriptCommand) String() string {
	out := fmt.Sprintf("# %s", strings.Join(sc.command, " "))
	if sc.inputFile != "" {
		out += " <" + sc.inputFile
	}
	out += "\n"
	for _, comment := range sc.comments {
		out += comment + "\n"
	}
	return out
}
