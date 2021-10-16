package errors

import "gorm.io/gorm"

func FromGORM(err error, msg string) *Error {
	e := &Error{
		msg: msg,
		cause: err,
	}

	if err == gorm.ErrRecordNotFound {
		e.code = NotFound
	} else {
		e.code = Internal
	}

	return e
}
