package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/morka17/shiny_bank/v1/src/utils"
)


var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return  utils.IsSupportedCurrency(currency)
	}

	return false 
}