package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"

	"github.com/arandall/meross/mdp"
)

type scanWifi struct {
	raw bool
}

func (cmd *scanWifi) FlagSet(fs *flag.FlagSet) {
	fs.BoolVar(&cmd.raw, "raw", false, "show raw json")
}

func (cmd *scanWifi) Run(c *client) error {
	resp, err := c.Do(mdp.WifiScan())
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

	var wifi mdp.WifiList
	if err := resp.Unmarshal(&wifi); err != nil {
		return err
	}

	if len(wifi.List) < 1 {
		return errors.New("no access points returned")
	}
	for _, ap := range wifi.List {
		fmt.Fprintf(out, "%s (%d%%)\n", ap.SSID, ap.Signal)
		fmt.Fprintf(out, "BSSID: %s\n\n", ap.BSSID)
	}
	return nil
}
