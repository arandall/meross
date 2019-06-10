package mdp

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

type Method int

const (
	Method_NIL Method = iota
	Method_GET
	Method_GETACK
	Method_SET
	Method_SETACK
	Method_PUSH
	Method_ERROR
)

func (m *Method) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

func ParseMethod(s string) (Method, error) {
	switch s {
	case "GET":
		return Method_GET, nil
	case "GETACK":
		return Method_GETACK, nil
	case "SET":
		return Method_SET, nil
	case "SETACK":
		return Method_SETACK, nil
	case "PUSH":
		return Method_PUSH, nil
	case "ERROR":
		return Method_ERROR, nil
	}
	return Method_NIL, errors.New("invalid method")
}

func (m *Method) UnmarshalText(b []byte) (err error) {
	*m, err = ParseMethod(string(b))
	return err
}

type Unixtime struct {
	time.Time
}

func (t *Unixtime) String() string {
	return strconv.FormatInt(t.Unix(), 10)
}

func (t *Unixtime) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}
func (t *Unixtime) UnmarshalJSON(b []byte) error {
	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	*t = Unixtime{time.Unix(i, 0)}
	return nil
}

// Header is attached to all packets containing data about the request/response. Not all fields are mandatory.
type Header struct {
	// Used to identify MQTT response topic if response required
	From string `json:"from,omitempty"`
	// MessageID is used to identify a request/response pair.
	MessageID      string   `json:"messageId,omitempty"`
	Method         Method   `json:"method,omitempty"`
	Namespace      string   `json:"namespace,omitempty"`
	PayloadVersion int      `json:"payloadVersion,omitempty"`
	Signature      string   `json:"sign,omitempty"`
	Timestamp      Unixtime `json:"timestamp,omitempty"`
	TimestampMS    int64    `json:"timestampMs,omitempty"`
}

// Packet represents the structure of a request/response
type Packet struct {
	Header  Header          `json:"header"`
	Payload json.RawMessage `json:"payload"`
}

func NewPacket(ns string, m Method, p json.RawMessage) *Packet {
	return &Packet{
		Header: Header{
			MessageID:      RandSeq(HEX, 32),
			Method:         m,
			Namespace:      ns,
			PayloadVersion: 1,
			Timestamp:      Unixtime{time.Now()},
		},
		Payload: p,
	}
}

// Sign signs the packet using the key provided and returns a signature
func (p *Packet) Sign(k string) {
	p.Header.Signature = GenerateSignature(p.Header.MessageID, k, p.Header.Timestamp.String())
}

// SignatureValid checks that key was used to sign the packet.
func (p *Packet) SignatureValid(k string) bool {
	return p.Header.Signature == Sign(k, p)
}

func (p *Packet) Unmarshal(v interface{}) error {
	// If Method is an error then return it as an error.
	if p.Header.Method == Method_ERROR {
		var err Error
		json.Unmarshal(p.Payload, &err)
		return err
	}
	return json.Unmarshal(p.Payload, v)
}

// Sign signs the packet using the key provided and returns a signature
func Sign(k string, p *Packet) string {
	return GenerateSignature(p.Header.MessageID, k, p.Header.Timestamp.String())
}

// Parse parses a []byte and returns the a Packet
func Parse(b []byte) (*Packet, error) {
	var p Packet
	err := json.Unmarshal(b, &p)
	return &p, err
}
