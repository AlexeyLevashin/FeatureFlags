package apperror

import "net/http"

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NotFound(message string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: message}
}

func BadRequest(message string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message}
}

func Forbidden(message string) *AppError {
	return &AppError{Code: http.StatusForbidden, Message: message}
}

func Conflict(message string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: message}
}

func Unauthorized(message string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: message}
}
