package mdp

import (
	"fmt"
)

// ErrorDetails give a code and text summary of an error.
type ErrorDetails struct {
	Code    int    `json:"code"`
	Summary string `json:"detail"`
}

// Error encapsulates ErrorDetails.
type Error struct {
	Details ErrorDetails `json:"error"`
}

// Error returns a text representing the error.
func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Details.Code, e.Details.Summary)
}
