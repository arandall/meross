package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/arandall/meross/mdp"
)

type system struct {
	raw bool
}

func (cmd *system) FlagSet(fs *flag.FlagSet) {
	fs.BoolVar(&cmd.raw, "raw", false, "show raw json")
}

func (cmd *system) Run(c *client) error {
	resp, err := c.Do(mdp.GetSystemAll())
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
	var sys mdp.SystemAll
	if err := resp.Unmarshal(&sys); err != nil {
		return err
	}

	hw := sys.Info.System.Hardware
	fmt.Fprintf(out, "UUID: %s\n", hw.UUID)
	fmt.Fprintf(out, "Type: %s-%s (v%s) - %s\n", hw.Type, hw.SubType, hw.Version, hw.ChipType)
	fmt.Fprintf(out, "MAC: %s\n", hw.MAC)
	fw := sys.Info.System.Firmware
	fmt.Fprintf(out, "Firmware: v%s\n", fw.Version)
	fmt.Fprintf(out, "MQTT:\n")
	fmt.Fprintf(out, "\tPrimary: mqtts://%s:%d\n", fw.Server, fw.Port)
	if fw.SecondServer != "" {
		fmt.Fprintf(out, "\tSecondary: mqtts://%s:%d\n\n", fw.SecondServer, fw.SecondPort)
	}
	fmt.Fprintf(out, "Device time: %s\n", sys.Info.System.Time.Timestamp.Format(time.RFC1123))
	fmt.Fprintf(out, "Toggle Status: %t (last changed: %s)\n", sys.Info.Control.Toggle.On, sys.Info.Control.Toggle.LastModified.Format(time.RFC1123))
	return nil
}
