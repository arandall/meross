package main

import (
	"encoding/json"
	"strconv"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

// ResponseTime is the time that the response occurred as per the response sent
// by the server.
// For some reason a seconds past epoch is returned on success and a string is
// returned on errors. eg. 2019-05-19 11:27:20
type ResponseTime time.Time

// UnmarshalJSON unmarshals the time in the response based on the type returned
// by the server.
//
// It is expected to be a seconds past epoch value, if not then it will fall
// back to `TimeFormat` and generate an error should there be a problem
// parsing.
func (t *ResponseTime) UnmarshalJSON(b []byte) error {
	//Assume number
	epoch, err := strconv.Atoi(string(b))
	if err != nil {
		// Fallback to Error time parsing if parse failure.
		var errTime errorTime
		if err := json.Unmarshal(b, &errTime); err != nil {
			return err
		}
		*t = ResponseTime(errTime)
		return nil
	}
	*t = ResponseTime(time.Unix(int64(epoch), 0))
	return nil
}

type errorTime time.Time

func (t *errorTime) UnmarshalText(b []byte) error {
	time, err := time.ParseInLocation(TimeFormat, string(b), serverTimezone)
	if err != nil {
		return err
	}
	*t = errorTime(time)
	return nil
}

// Response is a container structure that is returned by the service.
type Response struct {
	APIStatus int             `json:"apiStatus"`
	SysStatus int             `json:"sysStatus"`
	Info      string          `json:"info"`
	Timestamp ResponseTime    `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}
