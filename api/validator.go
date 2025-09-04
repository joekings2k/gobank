package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/joekings2k/gobank/util"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
		if currency,ok := fieldLevel.Field().Interface().(string); ok {
			return util.IsSuppoertedCurrency(currency)
		}
		return false
}