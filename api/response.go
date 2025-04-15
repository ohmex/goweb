package api

import "github.com/labstack/echo/v4"

type Response struct {
	Code uint64 `json:"code"`
	Err  string `json:"error,omitempty"`
	Msg  string `json:"message,omitempty"`
}

func (r Response) Error() string {
	return r.Err
}

func WebResponse(c echo.Context, statusCode int, data interface{}) error {
	return c.JSON(statusCode, data)
}

func FIELD_VALIDATION_ERROR(s ...string) Response {
	data := "Required fields are empty or not valid"
	if s != nil {
		data = s[0]
	}
	return Response{Code: 100000, Err: data, Msg: ""}
}

func INVALID_CREDENTIALS(s ...string) Response {
	data := "Invalid credentials"
	if s != nil {
		data = s[0]
	}
	return Response{Code: 100001, Err: data, Msg: ""}
}

func INVALID_SIGNING_METHOD(s ...string) Response {
	data := "Invalid signing method"
	if s != nil {
		data = s[0]
	}
	return Response{Code: 100002, Err: data, Msg: ""}
}

func INVALID_TOKEN(s ...string) Response {
	data := "Invalid token"
	if s != nil {
		data = s[0]
	}
	return Response{Code: 100003, Err: data, Msg: ""}
}

func USER_NOT_FOUND(s ...string) Response {
	data := "User not found"
	if s != nil {
		data = s[0]
	}
	return Response{Code: 100004, Err: data, Msg: ""}
}

func TOKEN_EXPIRED(s ...string) Response {
	data := "Token expired"
	if s != nil {
		data = s[0]
	}
	return Response{Code: 100006, Err: data, Msg: ""}
}

func USER_EXISTS(s ...string) Response {
	data := "User already exists"
	if s != nil {
		data = s[0]
	}
	return Response{Code: 100007, Err: data, Msg: ""}
}

func INTERNAL_SERVICE_ERROR(s ...string) Response {
	data := "Internal service error"
	if s != nil {
		data = s[0]
	}
	return Response{Code: 100008, Err: data, Msg: ""}
}

func RESOURCE_NOT_FOUND(s ...string) Response {
	data := "Resource not found"
	if s != nil {
		data = s[0]
	}
	return Response{Code: 100009, Err: data, Msg: ""}
}

func CASBIN_UNAUTHORIZED(s ...string) Response {
	data := "Access denied - resource authorization failed"
	if s != nil {
		data = s[0]
	}
	return Response{Code: 100011, Err: data, Msg: ""}
}

func RESOURCE_CREATION_FAILED(s ...string) Response {
	data := "Resource creation failed"
	if s != nil {
		data = s[0]
	}
	return Response{Code: 100012, Err: data, Msg: ""}
}

func RESOURCE_EXISTS(s ...string) Response {
	data := "Resource already exists"
	if s != nil {
		data = s[0]
	}
	return Response{Code: 100013, Err: "", Msg: data}
}

// SUCCESS Responses Start from 200000
func STATUS_OK(s ...string) Response {
	data := "Ok"
	if s != nil {
		data = s[0]
	}
	return Response{Code: 200000, Err: "", Msg: data}
}

func RESOURCE_CREATED(s ...string) Response {
	data := "Resource successfully created"
	if s != nil {
		data = s[0]
	}
	return Response{Code: 200001, Err: "", Msg: data}
}

func RESOURCE_DELETED(s ...string) Response {
	data := "Resource deleted"
	if s != nil {
		data = s[0]
	}
	return Response{Code: 200002, Err: "", Msg: data}
}

func USER_LOGGED_OUT(s ...string) Response {
	data := "User logged out"
	if s != nil {
		data = s[0]
	}
	return Response{Code: 200003, Err: "", Msg: data}
}
