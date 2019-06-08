package mdp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

type diffError struct {
	err error
	buf bytes.Buffer
}

func (e *diffError) Error() string {
	return e.buf.String()
}

func diff(d []byte, file string) error {
	cmd := exec.Command("diff", "-u", "--from-file="+file, "-")
	cmd.Stdin = bytes.NewReader(d)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return &diffError{err, out}
	}
	return nil
}

func isDiffAvailable() bool {
	cmd := exec.Command("diff", "--help")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// TestParse checks that the Marshal and Unmarshal functions produce the same result.
func TestPacket_EncodingDecoding(t *testing.T) {
	if !isDiffAvailable() {
		t.Skip("skipped: require diff command")
	}
	root := "testdata"
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		t.Run(path, func(t *testing.T) {
			f, err := os.Open(path)
			if info.IsDir() {
				return
			}
			if err != nil {
				t.Fatalf("opening file: %v", err)
			}
			b, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatalf("reading file: %v", err)
			}
			p, err := Parse(b)
			if err != nil {
				t.Errorf("parse: %v", err)
			}
			if !p.SignatureValid("f00f00f00f00f00f00f00f00f00f0000") {
				t.Error("signature invalid")
			}
			json, err := json.MarshalIndent(p, "", "  ")
			if err != nil {
				t.Errorf("marshal: %v", err)
			}
			if err := diff(json, path); err != nil {
				t.Error(err)
			}
		})
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
