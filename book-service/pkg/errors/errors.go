package errors

import (
	"net/http"
)

type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *CustomError) Error() string {
	return e.Message
}

func NewError(code int, message string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}

func ErrNotFound(message string) *CustomError {
	return &CustomError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func ErrValidation(message string) *CustomError {
	return &CustomError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func ErrUnauthorized(message string) *CustomError {
	return &CustomError{Code: http.StatusUnauthorized, Message: message}
}

func ErrConflict(message string) *CustomError {
	return &CustomError{Code: http.StatusConflict, Message: message}
}

func ErrInternal(message string) *CustomError {
	return &CustomError{Code: http.StatusInternalServerError, Message: message}
}
