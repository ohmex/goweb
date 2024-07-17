package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RoleRequest struct {
	Name string `form:"name" json:"name" validate:"required" example:"Operator"`
}

func (rr RoleRequest) Validate() error {
	return validation.ValidateStruct(&rr,
		validation.Field(&rr.Name, validation.Required),
	)
}
