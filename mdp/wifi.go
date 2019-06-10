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
	// Used for configuring WIFI network
	Password   string `json:"password,omitempty"`
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
	password, err := base64.StdEncoding.DecodeString(aux.Password)
	if err != nil {
		return err
	}

	*w = WifiNetwork(aux)
	w.SSID = string(ssid)
	w.Password = string(password)
	return nil
}

func (w *WifiNetwork) MarshalJSON() ([]byte, error) {
	type network WifiNetwork
	n := network(*w)
	n.SSID = base64.StdEncoding.EncodeToString([]byte(w.SSID))
	n.Password = base64.StdEncoding.EncodeToString([]byte(w.Password))
	return json.Marshal(n)
}

type WifiList struct {
	List []WifiNetwork `json:"wifiList"`
}

type WifiConfig struct {
	Wifi WifiNetwork `json:"wifi"`
}

const Ability_WifiList = "Appliance.Config.WifiList"
const Ability_ConfigWifi = "Appliance.Config.Wifi"

// WifiScan returns a request that can be used to return a []WifiNetwork containing all networks within range of a device.
func WifiScan() *Packet {
	return NewPacket(Ability_WifiList, Method_GET, []byte(`{}`))
}

func SetWifiConfig(wifi *WifiNetwork) (*Packet, error) {
	b, err := json.Marshal(&WifiConfig{*wifi})
	if err != nil {
		return nil, err
	}

	return NewPacket(Ability_ConfigWifi, Method_SET, b), nil
}