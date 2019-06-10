package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/arandall/meross/testutils/diff"
	"github.com/arandall/meross/testutils/scriptrunner"
)

var updateFlag = flag.Bool("update", false, "update test data")

func newTestServer(code int, r io.Reader) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if _, err := io.Copy(w, r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}))
}

func TestCommander_Run(t *testing.T) {
	if !diff.IsAvailable() {
		t.Skip("diff not available")
	}
	dir := "testdata"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		t.Run(file.Name(), func(t *testing.T) {
			if err := runFile(path.Join(dir, file.Name())); err != nil {
				t.Error(err)
			}
		})
	}
}

func runFile(f string) error {
	var buf bytes.Buffer
	out = &buf
	outErr = &buf
	runFunc := func(inFile io.Reader, args []string) error {
		var server *httptest.Server
		if inFile != nil {
			server = newTestServer(http.StatusOK, inFile)
			defer server.Close()
			args = scriptrunner.ReplaceArg(args, "URL", server.URL)
		}
		return run(args)
	}

	runner, err := scriptrunner.NewScriptRunner(f, runFunc, &buf)
	if err != nil {
		return err
	}
	if err := runner.Run(); err != nil {
		return err
	}

	if err := diff.Run(buf.Bytes(), f); err != nil {
		fmt.Println(err)
		if *updateFlag {
			f, err := os.Create(f)
			if err != nil {
				return err
			}
			if _, err := buf.WriteTo(f); err != nil {
				return err
			}
		}
	}
	return nil
}
