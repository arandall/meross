package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/arandall/meross/mdp"
)

type configure struct {
	userID     string
	host       string
	port       int
	secondHost string
	secondPort int
	quiet      bool
}

func (cmd *configure) FlagSet(fs *flag.FlagSet) {
	fs.StringVar(&cmd.userID, "user-id", "", "user-id")
	fs.StringVar(&cmd.host, "mqtt-host", "", "MQTT host")
	fs.IntVar(&cmd.port, "mqtt-port", 8883, "MQTT port")
	fs.StringVar(&cmd.secondHost, "second-mqtt-host", "", "Secondary MQTT host (optional)")
	fs.IntVar(&cmd.secondPort, "second-mqtt-port", 8883, "Secondary MQTT port (optional)")
	fs.BoolVar(&cmd.quiet, "quiet", false, "don't show username/password configured on device")
}

func (cmd *configure) Run(c *client) error {
	if cmd.userID == "" {
		return errors.New("-user-id flag is required")
	}
	if cmd.host == "" {
		return errors.New("-mqtt-host flag is required")
	}
	if c.key == "" {
		return errors.New("-key flag is required")
	}
	config := &mdp.MQTTConfig{
		UserID: cmd.userID,
		Key:    c.key,
		Gateway: mdp.Gateway{
			Host: cmd.host,
			Port: cmd.port,
		},
	}
	if cmd.secondHost != "" {
		config.Gateway.SecondHost = cmd.secondHost
		config.Gateway.SecondPort = cmd.secondPort
	}

	mqtt, err := mdp.MQTT(config)
	if err != nil {
		return err
	}
	if _, err := c.Do(mqtt); err != nil {
		return err
	}

	if cmd.quiet {
		return nil
	}
	p, err := c.Do(mdp.GetSystemAll())
	if err != nil {
		return err
	}
	var sys mdp.SystemAll
	if err := p.Unmarshal(&sys); err != nil {
		return err
	}
	fmt.Fprintln(out, "If using auth in your MQTT server use these credentials")
	fmt.Fprintf(out, "Username: %s\n", sys.Info.System.Hardware.MAC)
	fmt.Fprintf(out, "Password: %s_%s\n", config.UserID, mdp.GenerateSignature(sys.Info.System.Hardware.MAC, c.key))
	return nil
}
