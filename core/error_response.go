package core

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type AppError struct {
	Code       int                    `json:"status_code"`
	RootErr    error                  `json:"-"`
	Message    string                 `json:"message"`
	Log        string                 `json:"log"`
	Key        string                 `json:"error_key"`
	Additional map[string]interface{} `json:"additional,omitempty"`
}

func (e *AppError) StatusCode() int {
	return e.Code
}

func (e *AppError) RootError() error {
	if err, ok := e.RootErr.(*AppError); ok {
		return err.RootError()
	}
	return e.RootErr
}

func (e *AppError) Error() string {
	return e.RootError().Error()
}

func (e *AppError) WithAdditional(key string, value interface{}) *AppError {
	if e.Additional == nil {
		e.Additional = make(map[string]interface{})
	}
	e.Additional[key] = value
	return e
}

// Constructors

func NewErrorResponse(root error, msg, log, key string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		RootErr: root,
		Message: msg,
		Log:     log,
		Key:     key,
	}
}

func NewFullErrorResponse(statusCode int, root error, msg, log, key string) *AppError {
	return &AppError{
		Code:    statusCode,
		RootErr: root,
		Message: msg,
		Log:     log,
		Key:     key,
	}
}

func NewCustomError(root error, msg, key string) *AppError {
	if root != nil {
		return NewErrorResponse(root, msg, root.Error(), key)
	}
	return NewErrorResponse(errors.New(msg), msg, msg, key)
}

func SimpleErrorResponse(statusCode int, message, key string) *AppError {
	return &AppError{
		Code:    statusCode,
		RootErr: errors.New(message),
		Message: message,
		Key:     key,
	}
}

// Authentication & Authorization Errors

func ErrUnauthorized(root error, msg, key string) *AppError {
	if msg == "" {
		msg = "Unauthorized access"
	}
	return &AppError{
		Code:    http.StatusUnauthorized,
		RootErr: root,
		Message: msg,
		Key:     key,
	}
}

func ErrInvalidToken(err error) *AppError {
	return NewFullErrorResponse(
		http.StatusUnauthorized,
		err,
		"Invalid or expired token",
		err.Error(),
		"ErrInvalidToken",
	)
}

func ErrNoPermission(err error) *AppError {
	return NewFullErrorResponse(
		http.StatusForbidden,
		err,
		"You don't have permission to perform this action",
		err.Error(),
		"ErrNoPermission",
	)
}

func ErrForbidden(err error, msg string) *AppError {
	if msg == "" {
		msg = "Access forbidden"
	}
	return NewFullErrorResponse(
		http.StatusForbidden,
		err,
		msg,
		err.Error(),
		"ErrForbidden",
	)
}

// Client Errors (4xx)

func ErrInvalidRequest(err error) *AppError {
	return NewErrorResponse(err, "Invalid request", err.Error(), "ErrInvalidRequest")
}

func ErrBadRequest(err error, msg string) *AppError {
	if msg == "" {
		msg = "Bad request"
	}
	return NewErrorResponse(err, msg, err.Error(), "ErrBadRequest")
}

func ErrValidation(err error, msg string) *AppError {
	if msg == "" {
		msg = "Validation failed"
	}
	return NewErrorResponse(err, msg, err.Error(), "ErrValidation")
}

func ErrNotFound(entity string) *AppError {
	return NewFullErrorResponse(
		http.StatusNotFound,
		RecordNotFound,
		fmt.Sprintf("%s not found", entity),
		"record not found",
		"ErrNotFound",
	)
}

func ErrConflict(err error, msg string) *AppError {
	if msg == "" {
		msg = "Resource conflict"
	}
	return NewFullErrorResponse(
		http.StatusConflict,
		err,
		msg,
		err.Error(),
		"ErrConflict",
	)
}

func ErrTooManyRequests(err error) *AppError {
	return NewFullErrorResponse(
		http.StatusTooManyRequests,
		err,
		"Too many requests, please try again later",
		err.Error(),
		"ErrTooManyRequests",
	)
}

func ErrRequestTimeout(err error) *AppError {
	return NewFullErrorResponse(
		http.StatusRequestTimeout,
		err,
		"Request timeout",
		err.Error(),
		"ErrRequestTimeout",
	)
}

func ErrPayloadTooLarge(err error) *AppError {
	return NewFullErrorResponse(
		http.StatusRequestEntityTooLarge,
		err,
		"Payload too large",
		err.Error(),
		"ErrPayloadTooLarge",
	)
}

// Server Errors (5xx)

func ErrInternal(err error) *AppError {
	return NewFullErrorResponse(
		http.StatusInternalServerError,
		err,
		"Something went wrong in the server",
		err.Error(),
		"ErrInternal",
	)
}

func ErrDB(err error) *AppError {
	return NewFullErrorResponse(
		http.StatusInternalServerError,
		err,
		"Database error occurred",
		err.Error(),
		"ErrDatabase",
	)
}

func ErrServiceUnavailable(err error) *AppError {
	return NewFullErrorResponse(
		http.StatusServiceUnavailable,
		err,
		"Service temporarily unavailable",
		err.Error(),
		"ErrServiceUnavailable",
	)
}

func ErrGatewayTimeout(err error) *AppError {
	return NewFullErrorResponse(
		http.StatusGatewayTimeout,
		err,
		"Gateway timeout",
		err.Error(),
		"ErrGatewayTimeout",
	)
}

// Entity-specific CRUD Errors

func ErrCannotCreateEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Cannot create %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotCreate%s", entity),
	)
}

func ErrCannotGetEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Cannot get %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotGet%s", entity),
	)
}

func ErrCannotListEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Cannot list %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotList%s", entity),
	)
}

func ErrCannotUpdateEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Cannot update %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotUpdate%s", entity),
	)
}

func ErrCannotDeleteEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Cannot delete %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotDelete%s", entity),
	)
}

func ErrEntityNotFound(entity string, err error) *AppError {
	return NewFullErrorResponse(
		http.StatusNotFound,
		err,
		fmt.Sprintf("%s not found", strings.ToLower(entity)),
		err.Error(),
		fmt.Sprintf("Err%sNotFound", entity),
	)
}

func ErrEntityExisted(entity string, err error) *AppError {
	return NewFullErrorResponse(
		http.StatusConflict,
		err,
		fmt.Sprintf("%s already exists", strings.ToLower(entity)),
		err.Error(),
		fmt.Sprintf("Err%sAlreadyExists", entity),
	)
}

func ErrEntityDeleted(entity string, err error) *AppError {
	return NewFullErrorResponse(
		http.StatusGone,
		err,
		fmt.Sprintf("%s has been deleted", strings.ToLower(entity)),
		err.Error(),
		fmt.Sprintf("Err%sDeleted", entity),
	)
}

// Business Logic Errors

func ErrInsufficientBalance(err error) *AppError {
	return NewFullErrorResponse(
		http.StatusBadRequest,
		err,
		"Insufficient balance",
		err.Error(),
		"ErrInsufficientBalance",
	)
}

func ErrInvalidCredentials(err error) *AppError {
	return NewFullErrorResponse(
		http.StatusUnauthorized,
		err,
		"Invalid credentials",
		err.Error(),
		"ErrInvalidCredentials",
	)
}

func ErrAccountLocked(err error) *AppError {
	return NewFullErrorResponse(
		http.StatusForbidden,
		err,
		"Account is locked",
		err.Error(),
		"ErrAccountLocked",
	)
}

func ErrExpired(entity string, err error) *AppError {
	return NewFullErrorResponse(
		http.StatusGone,
		err,
		fmt.Sprintf("%s has expired", strings.ToLower(entity)),
		err.Error(),
		fmt.Sprintf("Err%sExpired", entity),
	)
}

// External Service Errors

func ErrExternalService(service string, err error) *AppError {
	return NewFullErrorResponse(
		http.StatusBadGateway,
		err,
		fmt.Sprintf("Error communicating with %s service", service),
		err.Error(),
		"ErrExternalService",
	)
}

func ErrThirdPartyAPI(err error) *AppError {
	return NewFullErrorResponse(
		http.StatusBadGateway,
		err,
		"Third-party API error",
		err.Error(),
		"ErrThirdPartyAPI",
	)
}

// File & Upload Errors

func ErrInvalidFileType(err error) *AppError {
	return NewErrorResponse(
		err,
		"Invalid file type",
		err.Error(),
		"ErrInvalidFileType",
	)
}

func ErrFileTooLarge(err error) *AppError {
	return NewFullErrorResponse(
		http.StatusRequestEntityTooLarge,
		err,
		"File size too large",
		err.Error(),
		"ErrFileTooLarge",
	)
}

func ErrFileUpload(err error) *AppError {
	return NewFullErrorResponse(
		http.StatusInternalServerError,
		err,
		"File upload failed",
		err.Error(),
		"ErrFileUpload",
	)
}

// Predefined Errors
var (
	RecordNotFound       = errors.New("record not found")
	ErrDuplicateKey      = errors.New("duplicate key")
	ErrInvalidInput      = errors.New("invalid input")
	ErrConnectionFailed  = errors.New("connection failed")
	ErrOperationTimeout  = errors.New("operation timeout")
	ErrResourceExhausted = errors.New("resource exhausted")
)
