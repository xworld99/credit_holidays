package models

import (
	"net/http"
)

type HandlerError struct {
	Err  error
	Type int
}

func (e *HandlerError) Error() string {
	return e.Err.Error()
}

func CreateBadRequestError(err error) HandlerError {
	return HandlerError{Type: http.StatusBadRequest, Err: err}
}

func CreateInternalError(err error) HandlerError {
	return HandlerError{Type: http.StatusInternalServerError, Err: err}
}

func CreateNotFoundError(err error) HandlerError {
	return HandlerError{Type: http.StatusNotFound, Err: err}
}
