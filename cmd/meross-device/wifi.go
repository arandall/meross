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

type configWifi struct {
	bssid string
	password string
}

func (cmd *configWifi) FlagSet(fs *flag.FlagSet) {
	fs.StringVar(&cmd.bssid, "bssid", "", "bssid from scan command")
	fs.StringVar(&cmd.password, "password", "", "password to connect to network")
}

func (cmd *configWifi) Run(c *client) error {
	if cmd.bssid == "" {
		return errors.New("-bssid must be set")
	}
	if cmd.password == "" {
		return errors.New("-password must be set")
	}

	resp, err := c.Do(mdp.WifiScan())
	if err != nil {
		return err
	}
	var wifi mdp.WifiList
	if err := resp.Unmarshal(&wifi); err != nil {
		return err
	}
	for _, network := range wifi.List {
		if network.BSSID != cmd.bssid {
			continue
		}
		network.Password = cmd.password
		cfg, err := mdp.SetWifiConfig(&network)
		if err != nil {
			return err
		}
		_, err = c.Do(cfg)
		return err
	}
	return fmt.Errorf("BSSID %q not found", cmd.bssid)
}