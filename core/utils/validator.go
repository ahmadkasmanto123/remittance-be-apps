package utils

import (
	"love-remittance-be-apps/lib/model"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) Validate(i interface{}) []model.ErrorData {
	errs := v.validator.Struct(i)
	if errs != nil {
		errorDatas := []model.ErrorData{}
		for _, err := range errs.(validator.ValidationErrors) {
			var errData model.ErrorData
			errData.Field = err.Field()
			errData.Tag = err.Tag()
			errData.Value = err.Value()
			errData.TagValue = err.Param()
			errorDatas = append(errorDatas, errData)
		}
		return errorDatas
	}
	return nil
}
