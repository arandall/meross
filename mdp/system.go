package mdp

import (
	"strconv"
	"time"
)

const Ability_SystemAll = "Appliance.System.All"

// Time represents a JSON integer in seconds past Epoch.
type Time struct {
	time.Time
}

func (t *Time) String() string {
	return strconv.FormatInt(t.Unix(), 10)
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}
func (t *Time) UnmarshalJSON(b []byte) error {
	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	*t = Time{time.Unix(i, 0)}
	return nil
}

// SystemAll contains all system attributes.
type SystemAll struct {
	Info Info `json:"all"`
}

// Hardware lists details about the device hardware.
type Hardware struct {
	Type     string `json:"type"`
	SubType  string `json:"subType"`
	Version  string `json:"version"`
	ChipType string `json:"chipType"`
	UUID     string `json:"uuid"`
	MAC      string `json:"macAddress"`
}

// Firmware lists details about the devices firmware.
type Firmware struct {
	Version      string `json:"version"`
	CompileTime  string `json:"compileTime"`
	WifiMAC      string `json:"wifiMac"`
	InnerIP      string `json:"innerIp"`
	Server       string `json:"server"`
	Port         int    `json:"port"`
	SecondServer string `json:"secondServer"`
	SecondPort   int    `json:"secondPort"`
	UserID       int    `json:"userId"`
}

// TimeRule contains information about the timezone and daylight savings rules for the device.
type TimeRule []int64

// Time that rule takes effect.
func (tr TimeRule) Time() time.Time {
	return time.Unix(tr[0], 0)
}

// UTCOffset in seconds
func (tr TimeRule) UTCOffset() int64 {
	return tr[1]
}

// DaylightSavings returns if daylight savings is being applied.
func (tr TimeRule) DaylightSavings() bool {
	return tr[2] == 1
}

// SystemTime contains the current system time.
type SystemTime struct {
	Time                Time       `json:"timestamp"`
	Zone                string     `json:"timezone"`
	DaylightSavingRules []TimeRule `json:"timeRule"`
}

// Online indicates that the device is connected to the server.
type Online struct {
	Status int `json:"status"`
}

// System is a container for the devices system components.
type System struct {
	Hardware   Hardware   `json:"hardware"`
	Firmware   Firmware   `json:"firmware"`
	SystemTime SystemTime `json:"time"`
	Online     Online     `json:"online"`
}

// Info is a container for the devices information.
type Info struct {
	System  System  `json:"system"`
	Control Control `json:"control"`
}

func GetSystemAll() *Packet {
	return NewPacket(Ability_SystemAll, Method_GET, []byte(`{}`))
}
