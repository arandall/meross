package mdp

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/arandall/meross/testutils/diff"
)

// TestParse checks that the Marshal and Unmarshal functions produce the same result.
// TODO(arandall): test payload marshal/unmarshal.
func TestPacket_EncodingDecoding(t *testing.T) {
	if !diff.IsAvailable() {
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
			if err := diff.Run(json, path); err != nil {
				t.Error(err)
			}
		})
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
