package mdp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/arandall/meross/testutils/diff"
)

func unwrap(file string) (*Packet, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("open file: %v", err)
	}
	defer f.Close()
	// Unwrap packet
	b, err := ioutil.ReadAll(f)
	p, err := Parse(b)
	if err != nil {
		return nil, fmt.Errorf("parse: %v", err)
	}
	if !p.SignatureValid("f00f00f00f00f00f00f00f00f00f0000") {
		return nil, fmt.Errorf("signature invalid")
	}
	return p, nil
}

func wrap(packet *Packet, v interface{}) (b []byte, err error) {
	packet.Payload, err = json.MarshalIndent(v, "", "  ")
	if err != nil {
		err = fmt.Errorf("marshal payload: %v", err)
		return
	}
	b, err = json.MarshalIndent(&packet, "", "  ")
	if err != nil {
		err = fmt.Errorf("marshal: %v", err)
	}
	return
}

// TestParse checks that the Marshal and Unmarshal functions produce the same result.
func TestPacket_EncodingDecoding(t *testing.T) {
	if !diff.IsAvailable() {
		t.Skip("skipped: require diff command")
	}
	root := "testdata"

	tt := []struct {
		file string
		t    interface{}
	}{
		{"ERROR-sign-error.json", &Error{}},
		{"GET-Appliance.System.All.json", &SystemAll{}},
		{"GETACK-wifi-scan.json", &WifiList{}},
		{"GETACK-system-all.json", &SystemAll{}},
	}

	for _, tc := range tt {
		t.Run(tc.file, func(t *testing.T) {
			tc.file = path.Join(root, tc.file)
			p, err := unwrap(tc.file)
			if err != nil {
				t.Fatal(err)
			}

			if err := json.Unmarshal(p.Payload, tc.t); err != nil {
				t.Fatalf("unmarshal payload: %v", err)
			}

			// Wrap it back up.
			json, err := wrap(p, tc.t)
			if err != nil {
				t.Fatal(err)
			}

			// Should be the same.
			if err := diff.Run(json, tc.file); err != nil {
				t.Fatal(err)
			}
		})
	}
}
