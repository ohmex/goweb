package api

import "github.com/labstack/echo/v4"

// Response represents a standard API response structure.
type Response struct {
	Code uint64 `json:"code"`
	Err  string `json:"error,omitempty"`
	Msg  string `json:"message,omitempty"`
}

// Error implements the error interface for Response.
func (r Response) Error() string {
	return r.Err
}

// WebResponse sends a JSON response with the given status code and data.
func WebResponse(c echo.Context, statusCode int, data interface{}) error {
	return c.JSON(statusCode, data)
}

// Error codes
const (
	CodeFieldValidationError   = 100000
	CodeInvalidCredentials     = 100001
	CodeInvalidSigningMethod   = 100002
	CodeInvalidToken           = 100003
	CodeUserNotFound           = 100004
	CodeTokenExpired           = 100006
	CodeUserExists             = 100007
	CodeInternalServiceError   = 100008
	CodeResourceNotFound       = 100009
	CodeCasbinUnauthorized     = 100011
	CodeResourceCreationFailed = 100012
	CodeResourceExists         = 100013
)

// Status codes
const (
	CodeStatusOK        = 200000
	CodeResourceCreated = 200001
	CodeResourceDeleted = 200002
	CodeUserLoggedOut   = 200003
)

// responseTemplate is a helper to create a Response with flexible error/message placement.
func responseTemplate(code uint64, defaultMsg string, isError bool, s ...string) Response {
	msg := defaultMsg
	if len(s) > 0 {
		msg = s[0]
	}
	if isError {
		return newResponse(code, msg, "")
	}
	return newResponse(code, "", msg)
}

// FIELD_VALIDATION_ERROR returns a response for field validation errors.
func FIELD_VALIDATION_ERROR(s ...string) Response {
	return responseTemplate(CodeFieldValidationError, "Required fields are empty or not valid", true, s...)
}

// INVALID_CREDENTIALS returns a response for invalid credentials.
func INVALID_CREDENTIALS(s ...string) Response {
	return responseTemplate(CodeInvalidCredentials, "Invalid credentials", true, s...)
}

// INVALID_SIGNING_METHOD returns a response for invalid signing method.
func INVALID_SIGNING_METHOD(s ...string) Response {
	return responseTemplate(CodeInvalidSigningMethod, "Invalid signing method", true, s...)
}

// INVALID_TOKEN returns a response for invalid token.
func INVALID_TOKEN(s ...string) Response {
	return responseTemplate(CodeInvalidToken, "Invalid token", true, s...)
}

// USER_NOT_FOUND returns a response for user not found.
func USER_NOT_FOUND(s ...string) Response {
	return responseTemplate(CodeUserNotFound, "User not found", true, s...)
}

// TOKEN_EXPIRED returns a response for expired token.
func TOKEN_EXPIRED(s ...string) Response {
	return responseTemplate(CodeTokenExpired, "Token expired", true, s...)
}

// USER_EXISTS returns a response for user already exists.
func USER_EXISTS(s ...string) Response {
	return responseTemplate(CodeUserExists, "User already exists", true, s...)
}

// INTERNAL_SERVICE_ERROR returns a response for internal service errors.
func INTERNAL_SERVICE_ERROR(s ...string) Response {
	return responseTemplate(CodeInternalServiceError, "Internal service error", true, s...)
}

// RESOURCE_NOT_FOUND returns a response for resource not found.
func RESOURCE_NOT_FOUND(s ...string) Response {
	return responseTemplate(CodeResourceNotFound, "Resource not found", true, s...)
}

// CASBIN_UNAUTHORIZED returns a response for authorization failures.
func CASBIN_UNAUTHORIZED(s ...string) Response {
	return responseTemplate(CodeCasbinUnauthorized, "Access denied - resource authorization failed", true, s...)
}

// RESOURCE_CREATION_FAILED returns a response for resource creation failures.
func RESOURCE_CREATION_FAILED(s ...string) Response {
	return responseTemplate(CodeResourceCreationFailed, "Resource creation failed", true, s...)
}

// RESOURCE_EXISTS returns a response for resource already exists.
func RESOURCE_EXISTS(s ...string) Response {
	return responseTemplate(CodeResourceExists, "Resource already exists", false, s...)
}

// STATUS_OK returns a response for successful operations.
func STATUS_OK(s ...string) Response {
	return responseTemplate(CodeStatusOK, "Ok", false, s...)
}

// RESOURCE_CREATED returns a response for successful resource creation.
func RESOURCE_CREATED(s ...string) Response {
	return responseTemplate(CodeResourceCreated, "Resource successfully created", false, s...)
}

// RESOURCE_DELETED returns a response for successful resource deletion.
func RESOURCE_DELETED(s ...string) Response {
	return responseTemplate(CodeResourceDeleted, "Resource deleted", false, s...)
}

// USER_LOGGED_OUT returns a response for successful user logout.
func USER_LOGGED_OUT(s ...string) Response {
	return responseTemplate(CodeUserLoggedOut, "User logged out", false, s...)
}

// newResponse is a helper to create a Response.
func newResponse(code uint64, err, msg string) Response {
	return Response{Code: code, Err: err, Msg: msg}
}
