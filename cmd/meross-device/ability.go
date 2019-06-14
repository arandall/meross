package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/arandall/meross/mdp"
)

type ability struct {
	raw bool
}

func (cmd *ability) FlagSet(fs *flag.FlagSet) {
	fs.BoolVar(&cmd.raw, "raw", false, "show raw json")
}

func (cmd *ability) Run(c *client) error {
	resp, err := c.Do(mdp.GetAbilities())
	if err != nil {
		return err
	}

	if cmd.raw {
		d, err := json.MarshalIndent(resp.Payload, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(out, string(d))
		return nil
	}

	var deviceAbility mdp.DeviceAbility
	if err := resp.Unmarshal(&deviceAbility); err != nil {
		return err
	}

	fmt.Fprintf(out, "Payload Version: %d\n", deviceAbility.PayloadVersion)
	for _, ability := range deviceAbility.Abilities {
		fmt.Fprintln(out, ability)
	}
	return nil
}
