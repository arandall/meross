package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type command interface {
	FlagSet(*flag.FlagSet)
	Run(*client) error
}

type commander struct {
	key  string
	url  string
	cmds map[string]command
}

func newCommander() *commander {
	return &commander{
		cmds: map[string]command{
			"scan":   &scanWifi{},
			"system": &system{},
			"config": &configure{},
			"wifi" : &configWifi{},
		},
	}
}

// Run a command specified with arguments.
func (c *commander) Run(args ...string) error {
	if len(args) < 2 {
		c.Usage()
		return errors.New("missing command")
	}

	// Find command
	cmdName := args[1]
	if cmdName == "help" {
		c.Usage()
		return nil
	}
	cmd, ok := c.cmds[cmdName]
	if !ok {
		c.Usage()
		return fmt.Errorf("Command %s not found", cmdName)
	}

	// Configure FlagSet for command
	fs := flag.NewFlagSet(strings.Join(args[0:2], " "), flag.ContinueOnError)
	fs.SetOutput(out)
	c.flagSet(fs)
	cmd.FlagSet(fs)

	if err := fs.Parse(args[2:]); err != nil {
		fs.Usage()
		return err
	}
	if err := c.validate(); err != nil {
		fs.Usage()
		return err
	}
	if err := cmd.Run(c.client()); err != nil {
		return err
	}
	return nil
}

// Usage prints out usage information.
func (c *commander) Usage() {
	fmt.Fprint(out, "Usage of meross-device:\n")
	fmt.Fprint(out, "\tmeross-device cmd [opts]\n\n")
	fmt.Fprint(out, "available commands\n")

	var commands []string
	for cmd := range c.cmds {
		commands = append(commands, cmd)
	}
	sort.Strings(commands)

	// To perform the opertion you want
	for _, cmd := range commands {
		fmt.Fprintf(out, "\t%s\n", cmd)
	}
}

func (c *commander) flagSet(fs *flag.FlagSet) {
	fs.StringVar(&c.key, "key", "", "key used for signing (optional for some commands)")
	fs.StringVar(&c.url, "url", "", "url of device configuration endpoint")
}

func (c *commander) validate() error {
	if c.url == "" {
		return errors.New("-url flag must be provided")
	}
	return nil
}

func (args *commander) client() *client {
	c, err := NewClient(http.DefaultClient, WithUrlString(args.url), WithKey(args.key))
	if err != nil {
		panic(err)
	}
	return c
}
