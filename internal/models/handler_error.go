package controllers

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

func createBadRequestError(err error) HandlerError {
	return HandlerError{Type: http.StatusBadRequest, Err: err}
}

func createInternalError(err error) HandlerError {
	return HandlerError{Type: http.StatusInternalServerError, Err: err}
}

func createNotFoundError(err error) HandlerError {
	return HandlerError{Type: http.StatusNotFound, Err: err}
}
