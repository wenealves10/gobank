package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/wenealves10/gobank/utils"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return utils.IsValidCurrency(currency)
	}

	return false
}
