package errors

import "net/http"

type errorCode uint8

const (
	NotFound errorCode = 1
	BadInput errorCode = 2
	Internal errorCode = 3
	Timeout errorCode = 4
	Unauthorized = 5
	AccessDenied errorCode = 6
)

func (c errorCode) ToHTTP() int {
	switch c {
	case NotFound:
		return http.StatusNotFound
	case BadInput:
		return http.StatusBadRequest
	case Internal:
		return http.StatusInternalServerError
	case Timeout:
		return http.StatusRequestTimeout
	case Unauthorized:
		return http.StatusUnauthorized
	case AccessDenied:
		return http.StatusForbidden
	}

	return http.StatusInternalServerError
}
