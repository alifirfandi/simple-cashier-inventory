package validation

import (
	"encoding/json"
	"github.com/alifirfandi/simple-cashier-inventory/exception"
	"github.com/alifirfandi/simple-cashier-inventory/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func LoginValidation(Request model.AuthRequest) (err error) {
	err = validation.ValidateStruct(&Request,
		validation.Field(&Request.Password, validation.Required.Error("NOT_BLANK")),
		validation.Field(&Request.Email, validation.Required.Error("NOT_BLANK")),
	)

	if err != nil {
		b, _ := json.Marshal(err)
		return exception.ValidationError{
			Message: string(b),
		}
	}

	return nil
}
