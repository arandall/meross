package mdp

import (
	"encoding/json"
	"errors"
	"time"
)

const Ability_SystemAll = "Appliance.System.All"

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
type TimeRule struct {
	Timestamp             time.Time
	GMTOffsetSeconds      int64
	DaylightSavingsOffset int64
}

func (tr *TimeRule) UnmarshalJSON(text []byte) error {
	aux := []int64{}
	if err := json.Unmarshal(text, &aux); err != nil {
		return err
	}
	if len(aux) != 3 {
		return errors.New("invalid length on timerule")
	}
	*tr = TimeRule{
		time.Unix(aux[0], 0),
		aux[1],
		aux[2],
	}
	return nil
}

// Time shows the current time of the device
type SystemTime struct {
	Timestamp time.Time      `json:"timestamp"`
	Timezone  *time.Location `json:"timezone"`
	TimeRule  []TimeRule     `json:"timeRule"`
}

func (t *SystemTime) UnmarshalJSON(data []byte) error {
	aux := struct {
		Time int64      `json:"timestamp"`
		Loc  string     `json:"timezone"`
		TR   []TimeRule `json:"timeRule"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	loc, err := time.LoadLocation(aux.Loc)
	if err != nil {
		return err
	}

	*t = SystemTime{
		time.Unix(aux.Time, 0),
		loc,
		aux.TR,
	}
	return nil
}

// Connected indicates that the device is connected to the server.
type Connected bool

func (s *Connected) UnmarshalJSON(data []byte) error {
	aux := struct {
		Status int
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*s = aux.Status == 1
	return nil
}

// System is a container for the devices system components.
type System struct {
	Hardware          Hardware   `json:"hardware"`
	Firmware          Firmware   `json:"firmware"`
	Time              SystemTime `json:"time"`
	ConnectedToServer Connected  `json:"online"`
}

// Toggle shows the current state of the switch.
type Toggle struct {
	On           bool
	LastModified time.Time
}

func (t *Toggle) UnmarshalJSON(data []byte) error {
	aux := struct {
		OnOff    int   `json:"onoff"`
		Modified int64 `json:"lmTime"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*t = Toggle{
		aux.OnOff == 1,
		time.Unix(aux.Modified, 0),
	}
	return nil
}

// Control wraps the Toggle information.
type Control struct {
	Toggle Toggle `json:"toggle"`
}

// Info is a container for the devices information.
type Info struct {
	System  System  `json:"system"`
	Control Control `json:"control"`
}

func GetSystemAll() *Packet {
	return NewPacket(Ability_SystemAll, Method_GET, []byte(`{}`))
}
