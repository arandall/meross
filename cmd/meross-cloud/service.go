package main

import (
	"encoding/json"
	"errors"
)

// Account data is provided as part of the login processes.
// It contains information that is used by devices and used for other
// communications.
type Account struct {
	Token  string `json:"token"`
	Key    string `json:"key"`
	UserID string `json:"userid"`
	Email  string `json:"email"`
}

// User data is used to Login/Create accounts.
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	/* These keys are sent by App but not required */
	Region string `json:"regionCode,omitempty"`
	Vendor string `json:"vendor,omitempty"`
}

// Login authenticates using the provided email/password and returns
// Credentials to be used by future requests.
func (c *client) Login(u *User) (*Account, error) {
	u.Vendor = "meross"
	req, err := NewRequest("v1/Auth/login", "23x17ahWarFH6w29", &u)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	cred := &Account{}
	return cred, json.Unmarshal(res.Data, cred)
}

// Signup creates a new account using the provided email/password and returns
// Credentials to be used by future requests.
func (c *client) Signup(u *User) (*Account, error) {
	u.Vendor = "meross"
	req, err := NewRequest("v1/Auth/reg", "23x17ahWarFH6w29", &u)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	cred := &Account{}
	return cred, json.Unmarshal(res.Data, cred)
}

type AppInfo struct {
	// Unsure what this is for and haven't seen it with any data.
	Extra json.RawMessage `json:"extra"`
	// Model is the make/model of the device running the app.
	Model string `json:"model"`
	// System is "Android" for android, not sure about iOS.
	System string `json:"system"`
	// Looks like a concatenation of a deviceID and the WiFi MAC address.
	UUID string `json:"uuid"`
	// Vendor is hardcoded to "meross"
	Vendor string `json:"vendor"`
	// Version of the OS.
	Version string `json:"version"`
}

// Log logs app information to the server.
//
// I'm not sure if there is any purpose of this call other than for Meross to
// gather information about the versions of people actively using their app.
func (c *client) Log(i *AppInfo) error {
	req, err := NewRequest("v1/log/user", "23x17ahWarFH6w29", i)
	if err != nil {
		return err
	}
	res, err := c.Do(req)
	if res.SysStatus != 0 || res.APIStatus != 0 {
		return errors.New(res.Info)
	}
	return nil
}

type Device struct {
	UUID string `json:"uuid"`
	// 1:online 2:offline
	Status         int               `json:"onlineStatus"`
	Name           string            `json:"devName"`
	IconID         string            `json:"devIconId"`
	BindTime       ResponseTime      `json:"bindTime"`
	Type           string            `json:"deviceType"`
	SubType        string            `json:"subType"`
	Channels       []json.RawMessage `json:"channels"`
	Region         string            `json:"region"`
	FirmwareVer    string            `json:"fmwareVersion"`
	HardwareVer    string            `json:"hdwareVersion"`
	UserIcon       string            `json:"userDevIcon"`
	IconType       int               `json:"iconType"`
	SkillNumber    string            `json:"skillNumber"`
	Domain         string            `json:"domain"`
	ReservedDomain string            `json:"reservedDomain"`
}

// DeviceList returns all devices registered with an account.
func (c *client) DeviceList() ([]*Device, error) {
	req, err := NewRequest("v1/Device/devList", "23x17ahWarFH6w29", json.RawMessage("{}"))
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	devices := []*Device{}
	return devices, json.Unmarshal(res.Data, &devices)
}
