package mdp

import (
	"encoding/json"
)

const Ability_ConfigKey = "Appliance.Config.Key"

type KeyConfig struct {
	Key MQTTConfig `json:"key"`
}

type MQTTConfig struct {
	Gateway Gateway `json:"gateway"`
	Key     string  `json:"key"`
	UserID  string  `json:"userId"`
}

type Gateway struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	SecondHost string `json:"secondHost,omitempty"`
	SecondPort int    `json:"secondPort,omitempty"`
}

// MQTT creates a packet to configure MQTT servers for a device.
func MQTT(config *MQTTConfig) (*Packet, error) {
	json, err := json.Marshal(KeyConfig{*config})
	if err != nil {
		return nil, err
	}
	return NewPacket(Ability_ConfigKey, Method_SET, json), nil
}
