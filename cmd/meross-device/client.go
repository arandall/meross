package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/arandall/meross/mdp"
)

type client struct {
	*http.Client
	key string
	url *url.URL
}

type option func(*client) error

// NewClient returns a new client from a hostname and key that corresponds to
// the credentials of the account linked to the device.
func NewClient(c *http.Client, opts ...option) (*client, error) {
	client := &client{
		Client: c,
	}

	for _, opt := range opts {
		opt(client)
	}
	if client.url == nil {
		return nil, errors.New("url must be specified.")
	}
	return client, nil
}

// WithKey creates an option to configure the client with a key.
func WithKey(key string) option {
	return func(c *client) error {
		c.key = key
		return nil
	}
}

// WithKey creates an option to configure the client with a URL.
func WithUrlString(u string) option {
	url, err := url.Parse(u)
	return func(c *client) error {
		c.url = url
		return err
	}
}

// Do sends a request to the HTTP URL configured.
func (c *client) Do(p *mdp.Packet) (*mdp.Packet, error) {
	p.Sign(c.key)
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.url.String(), bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status: %d - %s", resp.StatusCode, resp.Status)
	}
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return mdp.Parse(d)

}
