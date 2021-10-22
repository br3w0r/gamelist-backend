package errors

import (
	"encoding/json"
	"fmt"
	"time"
)

type Error struct {
	code      errorCode
	msg       string
	cause     error
	timestamp int64
}

func New(code errorCode, cause error, msg string) *Error {
	return &Error{
		code:      code,
		cause:     cause,
		msg:       msg,
		timestamp: time.Now().Unix(),
	}
}

func Newf(code errorCode, cause error, format string, args ...interface{}) *Error {
	return New(code, cause, fmt.Sprintf(format, args...))
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) Cause() error {
	return e.cause
}

func (e *Error) Code() errorCode {
	return e.code
}

// JSON parses error to json format.
//
// If safe is true, cause will not be parsed
func (e *Error) JSON(safe bool) []byte {
	var cause string
	if !safe {
		cause = fmt.Sprint(e.cause)
	}

	jsonErr := struct {
		Code      string `json:"code"`
		Message   string `json:"message"`
		Cause     string `json:"cause,omitempty"`
		Timestamp int64  `json:"timestamp"`
	}{
		Code:      e.code.String(),
		Message:   e.msg,
		Cause:     cause,
		Timestamp: e.timestamp,
	}

	res, _ := json.Marshal(&jsonErr)

	return res
}
