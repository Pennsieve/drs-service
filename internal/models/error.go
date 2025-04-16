package models

// ErrorCode defines error codes used in the API
type ErrorCode int

const (
	// Common errors
	ErrBadRequest       ErrorCode = 400
	ErrUnauthorized     ErrorCode = 401
	ErrForbidden        ErrorCode = 403
	ErrNotFound         ErrorCode = 404
	ErrMethodNotAllowed ErrorCode = 405
	ErrRequestTooLarge  ErrorCode = 413
	ErrInternalServer   ErrorCode = 500
	ErrNotImplemented   ErrorCode = 501

	// Business-specific errors
	ErrObjectNotFound  ErrorCode = 1001
	ErrAccessDenied    ErrorCode = 1002
	ErrInvalidPassport ErrorCode = 1003
	ErrInvalidRequest  ErrorCode = 1004
)

// Error represents an API error response, compliant with the OpenAPI specification
// Used for:
// 1. Standardizing error responses across the API
// 2. Mapping business error codes to appropriate HTTP status codes
// 3. Providing consistent error information to clients
//
// Field descriptions:
// - Msg: Detailed error message (required)
// - StatusCode: HTTP status code (required)
type Error struct {
	Msg        string `json:"msg"`        // Detailed error message
	StatusCode int    `json:"status_code"` // HTTP status code
}

// NewError creates a new error with the specified code and message
// It maps business error codes (>1000) to appropriate HTTP status codes
func NewError(code ErrorCode, msg string) *Error {
	statusCode := int(code)
	if code > 1000 {
		// Map business error codes to HTTP status codes
		statusMap := map[ErrorCode]int{
			ErrObjectNotFound:  404,
			ErrAccessDenied:    403,
			ErrInvalidPassport: 401,
			ErrInvalidRequest:  400,
		}
		if sc, ok := statusMap[code]; ok {
			statusCode = sc
		} else {
			statusCode = 500
		}
	}

	return &Error{
		Msg:        msg,
		StatusCode: statusCode,
	}
}

// NewNotFoundError creates a "not found" error
func NewNotFoundError(msg string) *Error {
	return NewError(ErrNotFound, msg)
}

// NewBadRequestError creates a "bad request" error
func NewBadRequestError(msg string) *Error {
	return NewError(ErrBadRequest, msg)
}

// NewUnauthorizedError creates an "unauthorized" error
func NewUnauthorizedError(msg string) *Error {
	return NewError(ErrUnauthorized, msg)
}

// NewInternalError creates an "internal server error" with optional underlying error details
func NewInternalError(msg string, err error) *Error {
	detail := ""
	if err != nil {
		detail = err.Error()
	}
	
	if detail != "" {
		msg = msg + ": " + detail
	}
	
	return NewError(ErrInternalServer, msg)
}

// NewNotImplementedError creates a "not implemented" error
func NewNotImplementedError(msg string) *Error {
	return NewError(ErrNotImplemented, msg)
}

// Error implements the error interface
func (e *Error) Error() string {
	return e.Msg
}
