package diff

import (
	"bytes"
	"os/exec"
)

// Error contains the universal diff returned by the diff command.
type Error struct {
	err error
	buf bytes.Buffer
}

func (e *Error) Error() string {
	return e.buf.String()
}

// Run a universal diff between a []byte and a file.
func Run(d []byte, file string) error {
	cmd := exec.Command("diff", "-u", "--from-file="+file, "-")
	cmd.Stdin = bytes.NewReader(d)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return &Error{err, out}
	}
	return nil
}

// IsAvailable returns if the diff command is available.
func IsAvailable() bool {
	cmd := exec.Command("diff", "--help")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
