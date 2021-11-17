package validation

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alifirfandi/simple-cashier-inventory/exception"
	"github.com/alifirfandi/simple-cashier-inventory/model"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ProductValidation(Request model.ProductRequest) (err error) {
	err = validation.ValidateStruct(&Request,
		validation.Field(&Request.Name, validation.Required.Error("NOT_BLANK")),
		validation.Field(&Request.Price, validation.Required.Error("NOT_BLANK")),
		validation.Field(&Request.Stock, validation.Required.Error("NOT_BLANK")),
	)

	if err != nil {
		b, _ := json.Marshal(err)
		return exception.ValidationError{
			Message: string(b),
		}
	}

	return nil
}

func SortProductsValidation(sortString string) (err error) {
	err = validation.Validate(strings.ToLower(sortString), validation.In(
		"id_asc",
		"name_asc",
		"name_desc",
		"created_at_asc",
		"created_at_desc",
		"updated_at_desc",
		"price_asc",
		"price_desc",
	).Error("NOT_VALID"))
	if err != nil {
		return exception.ValidationError{
			Message: fmt.Sprintf(`{"sort": "%s"}`, err.Error()),
		}
	}

	return nil
}
