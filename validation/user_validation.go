package validation

import (
	"encoding/json"
	"github.com/alifirfandi/simple-cashier-inventory/exception"
	"github.com/alifirfandi/simple-cashier-inventory/model"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func InsertUserValidation(Request model.UserRequest) (err error) {
	err = validation.ValidateStruct(&Request,
		validation.Field(&Request.Name, validation.Required.Error("NOT_BLANK")),
		validation.Field(&Request.Email, validation.Required.Error("NOT_BLANK")),
		validation.Field(&Request.Password, validation.Required.Error("NOT_BLANK")),
		validation.Field(&Request.Role, validation.Required.Error("NOT_BLANK")),
	)

	if err != nil {
		b, _ := json.Marshal(err)
		return exception.ValidationError{
			Message: string(b),
		}
	}

	return nil
}

func UpdateUserValidation(Request model.UserRequest) (err error) {
	err = validation.ValidateStruct(&Request,
		validation.Field(&Request.Name, validation.Required.Error("NOT_BLANK")),
		validation.Field(&Request.Email, validation.Required.Error("NOT_BLANK")),
		validation.Field(&Request.Role, validation.Required.Error("NOT_BLANK")),
	)

	if err != nil {
		b, _ := json.Marshal(err)
		return exception.ValidationError{
			Message: string(b),
		}
	}

	return nil
}
