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

// Response codes
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

	CodeStatusOK        = 200000
	CodeResourceCreated = 200001
	CodeResourceDeleted = 200002
	CodeUserLoggedOut   = 200003
)

// newResponse is a helper to create a Response.
func newResponse(code uint64, err, msg string) Response {
	return Response{Code: code, Err: err, Msg: msg}
}

// FIELD_VALIDATION_ERROR returns a response for field validation errors.
func FIELD_VALIDATION_ERROR(s ...string) Response {
	data := "Required fields are empty or not valid"
	if len(s) > 0 {
		data = s[0]
	}
	return newResponse(CodeFieldValidationError, data, "")
}

// INVALID_CREDENTIALS returns a response for invalid credentials.
func INVALID_CREDENTIALS(s ...string) Response {
	data := "Invalid credentials"
	if len(s) > 0 {
		data = s[0]
	}
	return newResponse(CodeInvalidCredentials, data, "")
}

// INVALID_SIGNING_METHOD returns a response for invalid signing method.
func INVALID_SIGNING_METHOD(s ...string) Response {
	data := "Invalid signing method"
	if len(s) > 0 {
		data = s[0]
	}
	return newResponse(CodeInvalidSigningMethod, data, "")
}

// INVALID_TOKEN returns a response for invalid token.
func INVALID_TOKEN(s ...string) Response {
	data := "Invalid token"
	if len(s) > 0 {
		data = s[0]
	}
	return newResponse(CodeInvalidToken, data, "")
}

// USER_NOT_FOUND returns a response for user not found.
func USER_NOT_FOUND(s ...string) Response {
	data := "User not found"
	if len(s) > 0 {
		data = s[0]
	}
	return newResponse(CodeUserNotFound, data, "")
}

// TOKEN_EXPIRED returns a response for expired token.
func TOKEN_EXPIRED(s ...string) Response {
	data := "Token expired"
	if len(s) > 0 {
		data = s[0]
	}
	return newResponse(CodeTokenExpired, data, "")
}

// USER_EXISTS returns a response for user already exists.
func USER_EXISTS(s ...string) Response {
	data := "User already exists"
	if len(s) > 0 {
		data = s[0]
	}
	return newResponse(CodeUserExists, data, "")
}

// INTERNAL_SERVICE_ERROR returns a response for internal service errors.
func INTERNAL_SERVICE_ERROR(s ...string) Response {
	data := "Internal service error"
	if len(s) > 0 {
		data = s[0]
	}
	return newResponse(CodeInternalServiceError, data, "")
}

// RESOURCE_NOT_FOUND returns a response for resource not found.
func RESOURCE_NOT_FOUND(s ...string) Response {
	data := "Resource not found"
	if len(s) > 0 {
		data = s[0]
	}
	return newResponse(CodeResourceNotFound, data, "")
}

// CASBIN_UNAUTHORIZED returns a response for authorization failures.
func CASBIN_UNAUTHORIZED(s ...string) Response {
	data := "Access denied - resource authorization failed"
	if len(s) > 0 {
		data = s[0]
	}
	return newResponse(CodeCasbinUnauthorized, data, "")
}

// RESOURCE_CREATION_FAILED returns a response for resource creation failures.
func RESOURCE_CREATION_FAILED(s ...string) Response {
	data := "Resource creation failed"
	if len(s) > 0 {
		data = s[0]
	}
	return newResponse(CodeResourceCreationFailed, data, "")
}

// RESOURCE_EXISTS returns a response for resource already exists.
func RESOURCE_EXISTS(s ...string) Response {
	data := "Resource already exists"
	if len(s) > 0 {
		data = s[0]
	}
	return newResponse(CodeResourceExists, "", data)
}

// STATUS_OK returns a response for successful operations.
func STATUS_OK(s ...string) Response {
	data := "Ok"
	if len(s) > 0 {
		data = s[0]
	}
	return newResponse(CodeStatusOK, "", data)
}

// RESOURCE_CREATED returns a response for successful resource creation.
func RESOURCE_CREATED(s ...string) Response {
	data := "Resource successfully created"
	if len(s) > 0 {
		data = s[0]
	}
	return newResponse(CodeResourceCreated, "", data)
}

// RESOURCE_DELETED returns a response for successful resource deletion.
func RESOURCE_DELETED(s ...string) Response {
	data := "Resource deleted"
	if len(s) > 0 {
		data = s[0]
	}
	return newResponse(CodeResourceDeleted, "", data)
}

// USER_LOGGED_OUT returns a response for successful user logout.
func USER_LOGGED_OUT(s ...string) Response {
	data := "User logged out"
	if len(s) > 0 {
		data = s[0]
	}
	return newResponse(CodeUserLoggedOut, "", data)
}
