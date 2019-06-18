package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/arandall/meross/mdp"
)

// Request is the container that API request calls must use.
type Request struct {
	// ServicePath is the URL path component from the baseURL
	ServicePath string
	// Params is a JSON object that must match the service being called
	Params []byte
	// Time of the request
	Timestamp time.Time
	// Nonce is random a string used in signing
	Nonce string
	// Key is a PSK used in signing
	Key string
}

// NewRequest creates a new Request
func NewRequest(path, key string, req interface{}) (*Request, error) {
	params, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request unmarshal: %v", err)
	}
	return &Request{
		ServicePath: path,
		Params:      params,
		Timestamp:   time.Now(),
		Nonce:       mdp.RandSeq(mdp.Letters, 16),
		Key:         key,
	}, nil
}

// Sign signs the request returning the URL encoded payload for transmission
// via HTTP.
func (r *Request) Sign() io.Reader {
	vars := url.Values{}
	vars.Add("params", base64.StdEncoding.EncodeToString(r.Params))
	vars.Add("timestamp", strconv.FormatInt(r.Timestamp.Unix(), 10))
	vars.Add("nonce", r.Nonce)
	vars.Add("sign", mdp.GenerateSignature(r.Key, vars.Get("timestamp"), vars.Get("nonce"), vars.Get("params")))

	return strings.NewReader(vars.Encode())
}
