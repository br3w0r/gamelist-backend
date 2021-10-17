package errors

import (
	"gorm.io/gorm"
)

func FromGORM(err error, msg string) *Error {
	var code errorCode

	if err == gorm.ErrRecordNotFound {
		code = NotFound
	} else {
		code = Internal
	}

	return New(code, nil, msg)
}

func JSONParseErr(err error) *Error {
	return New(BadInput, err, "failed to parse request to json")
}
