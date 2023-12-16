package api

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gym_management_system/errors"
	"reflect"
	"strings"
	"time"
	"unicode"
)

type CustomValidator struct {
	validate *validator.Validate
}

func NewCustomValidator() (*CustomValidator, error) {
	v := validator.New()
	if err := v.RegisterValidation("password", PasswordValidation); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("gteCurrentDay", GreaterThanOrEqualCurrentDayValidation); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("gtNow", GreaterThanNowValidation); err != nil {
		return nil, err
	}

	return &CustomValidator{validate: v}, nil
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.validate.Struct(i)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			field := e.Field()
			tag := e.Tag()
			switch tag {
			case "alpha":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must contain only alphabetical characters", field))
			case "email":
				validationErrors = append(validationErrors, fmt.Sprintf("%s is not a valid email format", field))
			case "min":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be at least %s characters long", field, e.Param()))
			case "max":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be at most %s characters long", field, e.Param()))
			case "gt":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be greater than %s", field, e.Param()))
			case "gte":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param()))
			case "gtfield":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be greater than %s", field, e.Param()))
			case "oneof":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be one of %s", field, e.Param()))
			case "gtNow":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be greater than current time", field))
			case "gteCurrentDay":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be greater than or equal to start of current day", field))
			case "password":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must contain number, lower case, upper case and special character", field))
			case "required_with":
				validationErrors = append(validationErrors, fmt.Sprintf("%s is required with %s", field, e.Param()))
			default:
				validationErrors = append(validationErrors, tag)
			}
		}
		return errors.InvalidRequest{Message: "validation error: " + strings.Join(validationErrors, "; ")}
	}
	return nil
}

func IsEmpty(updateRequest *UpdateAccountRequest) bool {
	value := reflect.ValueOf(*updateRequest)

	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i).Interface()
		if fieldValue != nil {
			return false
		}
	}

	return true
}

func PasswordValidation(fl validator.FieldLevel) bool {
	pw := fl.Field().String()
	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, char := range pw {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

func GreaterThanOrEqualCurrentDayValidation(fl validator.FieldLevel) bool {
	startOfDay := time.Now().Truncate(24 * time.Hour)
	inputTime := fl.Field().Interface().(time.Time)
	return !inputTime.Before(startOfDay)
}

func GreaterThanNowValidation(fl validator.FieldLevel) bool {
	date := fl.Field().Interface().(time.Time)
	//if date.After(time.Now()) {
	//	return true
	//}
	//fl.Field().SetString("must be greater than the current time")
	//return false
	return date.After(time.Now())
}
