package apiError

import (
	"fmt"
	"net/http"
)

type AErr interface {
	Error() string
	Code() int
}
type APIError struct {
	err  error
	code int
}

func (e *APIError) Error() string {
	return e.err.Error()
}

func (e *APIError) Code() int {
	return e.code
}

func New(err error, code int) *APIError {
	return &APIError{
		err:  err,
		code: code,
	}
}

var (
	ErrInternal     = New(fmt.Errorf("Internal Error"), http.StatusInternalServerError)
	ErrDontHaveUser = New(fmt.Errorf("request dont have user"), http.StatusBadRequest)
)
