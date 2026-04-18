package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/arandall/meross/mdp"
)

type genpass struct {
	userID string
	mac    string
}

func (cmd *genpass) FlagSet(fs *flag.FlagSet) {
	fs.StringVar(&cmd.userID, "user-id", "", "user-id")
	fs.StringVar(&cmd.mac, "mac", "", "device MAC address")
}

func (cmd *genpass) Run(c *client) error {
	if cmd.userID == "" {
		return errors.New("-user-id flag is required")
	}
	if cmd.mac == "" {
		return errors.New("-mac flag is required")
	}
	if c.key == "" {
		return errors.New("-key flag is required")
	}
	fmt.Fprintf(out, "Username: %s\n", cmd.mac)
	fmt.Fprintf(out, "Password: %s_%s\n", cmd.userID, mdp.GenerateSignature(cmd.mac, c.key))
	return nil
}
