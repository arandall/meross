package mdp

import (
	"encoding/json"
	"time"
)

// Control wraps the Toggle information.
type Control struct {
	Toggle  Toggle    `json:"toggle"`
	Trigger []Trigger `json:"trigger,omitempty"`
	Timer   []Timer   `json:"timer,omitempty"`
}

// Toggle contains information about the current switch status of a device.
type Toggle struct {
	OnOff        int  `json:"onoff"`
	LastModified Time `json:"lmTime"`
}

// TimerType defines if a Timer or Trigger is a once off or
type TimerType int

const (
	TimerWeekly TimerType = 1
	TimerOnce             = 2
)

// DaysOfWeek is a bitset of the days of the week a Trigger or Timer is for.
type DaysOfWeek uint8

const (
	Sunday DaysOfWeek = 1<<iota + 128
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

// SecondDuration represents a time of day from 0:00 in seconds.
type SecondDuration struct {
	time.Duration
}

func (d *SecondDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(d.Seconds()))
}

func (d *SecondDuration) UnmarshalJSON(b []byte) error {
	duration, err := time.ParseDuration(string(b) + "s")
	if err != nil {
		return err
	}
	*d = SecondDuration{duration}
	return nil
}

// MinuteDuration represents a time of day from 0:00 in minutes.
type MinuteDuration struct {
	time.Duration
}

func (d *MinuteDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(d.Minutes()))
}

func (d *MinuteDuration) UnmarshalJSON(b []byte) error {
	duration, err := time.ParseDuration(string(b) + "m")
	if err != nil {
		return err
	}
	*d = MinuteDuration{duration}
	return nil
}

// Timer is used to schedule a power on/off at a particular time and day of week.
type Timer struct {
	ID         string         `json:"id"`
	Type       TimerType      `json:"type"`
	Enable     int            `json:"enable"`
	Alias      string         `json:"alias"`
	Time       MinuteDuration `json:"time"`
	Week       DaysOfWeek     `json:"week"`
	Duration   int            `json:"duration"` // Have not seen use of this.
	CreateTime Time           `json:"createTime"`
	Extend     Control        `json:"extend"`
}

// Trigger is used to perform a delayed Control action.
// An example of a trigger is, power off 1 hour after power has been turned on.
type Trigger struct {
	ID         string    `json:"id"`
	Type       TimerType `json:"type"`
	Enable     int       `json:"enable"`
	Alias      string    `json:"alias"`
	CreateTime Time      `json:"createTime"`
	Rule       Rule      `json:"rule"`
}

// Rule is the structure of a Trigger rule.
type Rule struct {
	If   Control   `json:"_if_"`
	Then DelayRule `json:"_then_"`
	Do   Control   `json:"_do_"`
}

// DelayRule contains the Delay to apply to a Rule.
type DelayRule struct {
	Delay Delay `json:"delay"`
}

// Delay defines a delay used in a Trigger.
type Delay struct {
	Week     DaysOfWeek     `json:"week"`
	Duration SecondDuration `json:"duration"`
}
