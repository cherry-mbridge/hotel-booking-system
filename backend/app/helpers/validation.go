package helpers

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidationErrors(model interface{}, errs interface{}) map[string]string {

	errors := make(map[string]string)

	validationErrors := errs.(validator.ValidationErrors)

	t := reflect.TypeOf(model)

	for _, err := range validationErrors {

		field, _ := t.FieldByName(err.Field())

		jsonField :=
			strings.Split(
				field.Tag.Get("json"),
				",",
			)[0]

		label := field.Tag.Get("label")

		if label == "" {
			label = err.Field()
		}

		switch err.Tag() {

		case "required":
			errors[jsonField] =
				label + " is required"

		case "email":
			errors[jsonField] =
				label + " must be a valid email"
		}
	}

	return errors
}
