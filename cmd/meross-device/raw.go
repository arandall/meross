package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/arandall/meross/mdp"
)

type raw struct {
	method string
	ns     string
}

func (cmd *raw) FlagSet(fs *flag.FlagSet) {
	fs.StringVar(&cmd.method, "method", "", "method for request")
	fs.StringVar(&cmd.ns, "namespace", "", "namespace for request")
}

func (cmd *raw) Run(c *client) error {
	if cmd.ns == "" {
		return errors.New("-namespace required")
	}
	if cmd.method == "" {
		return errors.New("-method required")
	}
	method, err := mdp.ParseMethod(cmd.method)
	if err != nil {
		return err
	}
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return errors.New("payload must be sent via STDIN")
	}
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil
	}

	resp, err := c.Do(mdp.NewPacket(cmd.ns, method, b))
	if err != nil {
		return err
	}

	json, err := json.MarshalIndent(resp.Payload, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(out, string(json))
	return nil
}
