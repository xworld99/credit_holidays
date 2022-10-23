package models

import (
	"net/http"
)

// HandlerError - specific struct used for error classification in Handler
type HandlerError struct {
	Err  error
	Type int
}

func (e *HandlerError) Error() string {
	return e.Err.Error()
}

// CreateBadRequestError returns new instance of bad request error
func CreateBadRequestError(err error) HandlerError {
	return HandlerError{Type: http.StatusBadRequest, Err: err}
}

// CreateInternalError returns new instance of internal server error
func CreateInternalError(err error) HandlerError {
	return HandlerError{Type: http.StatusInternalServerError, Err: err}
}

// CreateNotFoundError returns new instance of not found error
func CreateNotFoundError(err error) HandlerError {
	return HandlerError{Type: http.StatusNotFound, Err: err}
}
