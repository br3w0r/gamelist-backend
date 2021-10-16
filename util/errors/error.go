package errors

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	code  errorCode
	msg   string
	cause error
}

func New(code errorCode, cause error, msg string) error {
	return &Error{
		code: code,
	}
}

func Newf(code errorCode, cause error, format string, args ...interface{}) error {
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

func (e *Error) MarshalJSON() (res []byte, err error) {
	jsonErr := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Cause   string `json:"cause"`
	}{
		Code:    int(e.code),
		Message: e.msg,
		Cause:   fmt.Sprint(e.cause),
	}

	res, err = json.Marshal(&jsonErr)

	return
}
