package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	minPathLength = 8
)

type BasicAuth struct {
	Email    string `form:"email" json:"email" validate:"required" example:"username@gmail.com"`
	Password string `form:"password" json:"password" validate:"required" example:"P@ssw0r6"`
}

func (ba BasicAuth) Validate() error {
	return validation.ValidateStruct(&ba,
		validation.Field(&ba.Email, is.Email),
		validation.Field(&ba.Password, validation.Length(minPathLength, 0)),
	)
}

type LoginRequest struct {
	BasicAuth
}

type RegisterRequest struct {
	BasicAuth
	Name string `form:"name" json:"name" validate:"required" example:"John Doe"`
}

func (rr RegisterRequest) Validate() error {
	err := rr.BasicAuth.Validate()
	if err != nil {
		return err
	}

	return validation.ValidateStruct(&rr,
		validation.Field(&rr.Name, validation.Required),
	)
}

type RefreshRequest struct {
	Token string `form:"token" json:"token" validate:"required" example:"refresh_token"`
}

func (rr RefreshRequest) Validate() error {
	return validation.ValidateStruct(&rr,
		validation.Field(&rr.Token, validation.Required),
	)
}

type UpdateRequest struct {
	Name string `form:"name" json:"name" validate:"required" example:"John Doe"`
}

func (rr UpdateRequest) Validate() error {
	return validation.ValidateStruct(&rr,
		validation.Field(&rr.Name, validation.Required),
	)
}
