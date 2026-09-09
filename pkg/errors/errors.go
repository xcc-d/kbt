package errors

import (
	"k8s.io/apimachinery/pkg/api/errors"
)

const (
	Success    = 200
	Failed     = 500
	NotFound   = 404
	BadRequest = 400
	Conflict   = 409
)

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewFailed(msg string) *AppError {
	return &AppError{Code: Failed, Message: msg}
}

func NewNotFound(msg string) *AppError {
	return &AppError{Code: NotFound, Message: msg}
}

func NewBadRequest(msg string) *AppError {
	return &AppError{Code: BadRequest, Message: msg}
}

func NewConflict(msg string) *AppError {
	return &AppError{Code: Conflict, Message: msg}
}

func FromError(err error) *AppError {
	if err == nil {
		return nil
	}

	if appError, ok := err.(*AppError); ok {
		return appError
	}

	switch {
	case errors.IsNotFound(err):
		return NewNotFound(err.Error())
	case errors.IsBadRequest(err):
		return NewBadRequest(err.Error())
	case errors.IsConflict(err):
		return NewConflict(err.Error())
	case errors.IsAlreadyExists(err):
		return NewConflict(err.Error())
	default:
		return NewFailed(err.Error())
	}
}
