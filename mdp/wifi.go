package mdp

import (
	"encoding/base64"
	"encoding/json"
)

type WifiNetwork struct {
	SSID       string `json:"ssid"`
	BSSID      string `json:"bssid"`
	Signal     int    `json:"signal"`
	Channel    int    `json:"channel"`
	Encryption int    `json:"encryption"`
	Cipher     int    `json:"cipher"`
}

func (w *WifiNetwork) UnmarshalJSON(data []byte) error {
	type network WifiNetwork

	var aux network
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	ssid, err := base64.StdEncoding.DecodeString(aux.SSID)
	if err != nil {
		return err
	}

	*w = WifiNetwork(aux)
	w.SSID = string(ssid)
	return nil
}

type WifiList struct {
	List []WifiNetwork `json:"wifiList"`
}

const Ability_WifiList = "Appliance.Config.WifiList"

// WifiScan returns a request that can be used to return a []WifiNetwork containing all networks within range of a device.
func WifiScan() *Packet {
	return NewPacket(Ability_WifiList, Method_GET, []byte(`{}`))
}
