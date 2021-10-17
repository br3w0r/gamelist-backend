package errors

import (
	"errors"

	"gorm.io/gorm"
)

func FromGORM(tx *gorm.DB, msg string) *Error {
	var code errorCode

	if (tx.Error == nil && tx.RowsAffected == 0) || errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		code = NotFound
	} else {
		code = Internal
	}

	return New(code, nil, msg)
}

func JSONParseErr(err error) *Error {
	return New(BadInput, err, "failed to parse request to json")
}
