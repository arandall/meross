package main

import (
	"fmt"
	"io"
	"os"
)

var out io.Writer
var outErr io.Writer

func main() {
	out = os.Stdout
	outErr = os.Stderr
	if err := run(os.Args); err != nil {
		fmt.Fprintln(outErr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	c := newCommander()
	if err := c.Run(args...); err != nil {
		return err
	}
	return nil
}
