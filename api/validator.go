package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/kamalbowselvam/chaintask/util"
)


var validRole validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if role, ok := fieldLevel.Field().Interface().(string); ok {
		// check if role is supported 
		return util.IsSupportedRole(role)
	}

	return false
}