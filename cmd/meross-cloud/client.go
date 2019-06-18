// Package cloud implements the client API to access the Meross Cloud API.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"
)

type client struct {
	baseURL *url.URL
	client  *http.Client
	acc     Account
}

var serverTimezone *time.Location

func init() {
	setTimezone("Asia/Shanghai")
}

// SetTimezone sets the timezone to use for timestamps returned by the Server.
//
// In the error scenarios I tested it looks like the timezone used is China
// Standard Time so this is the default used. There shouldn't be a need to
// modify this setting but it's early days so thought I'd expose it.
func setTimezone(s string) (err error) {
	serverTimezone, err = time.LoadLocation(s)
	return err
}

// NewClient returns a new client based on the default URL.
func NewClient(c *http.Client, baseURL *url.URL) *client {
	return &client{
		baseURL: baseURL,
		client:  c,
	}
}

// Do performs a request against the client.
func (c *client) Do(in *Request) (*Response, error) {
	url := *c.baseURL
	url.Path = path.Join(url.Path, in.ServicePath)

	req, err := http.NewRequest("POST", url.String(), in.Sign())
	// Add header if credentials have been added to client.
	if c.acc.Token != "" {
		req.Header.Add("Authorization", "Basic "+c.acc.Token)
	}
	if err != nil {
		return nil, err
	}
	// Mimic App UA
	req.Header.Add("User-Agent", "okhttp/3.6.0")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := &Response{}
	if err := json.Unmarshal(buf, result); err != nil {
		return nil, err
	}
	if result.APIStatus != 0 || result.SysStatus != 0 {
		switch result.APIStatus {
		case 1002:
			return nil, errors.New("invalid credentials")
		default:
			return nil, fmt.Errorf("%d-%d: %s", result.SysStatus, result.APIStatus, result.Info)
		}
	}
	return result, nil
}

func (c *client) ApplyCredentials(acc *Account) error {
	if acc.Token == "" {
		return errors.New("token must be defined.")
	}
	c.acc = *acc
	return nil
}
