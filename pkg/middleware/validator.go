package middleware

import (
	"database/sql/driver"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"reflect"
	"regexp"
)

// Precompiled regular expressions for performance
var (
	noSpecialCharPattern      = regexp.MustCompile(`[<>/"'/{}/\[\]/]`)
	noSpecialCharSlashPattern = regexp.MustCompile(`[<>"'{}[\]]`)
)

// CustomValidator holds the validator instance.
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator creates a new instance of CustomValidator.
func NewCustomValidator() *CustomValidator {
	v := validator.New()
	registerValidations(v)
	return &CustomValidator{validator: v}
}

func registerValidations(v *validator.Validate) {
	_ = v.RegisterValidation("no_special_char", validateNoSpecialChar)
	_ = v.RegisterValidation("no_special_char_slash", validateNoSpecialCharSlash)

	v.RegisterCustomTypeFunc(validateValuer, pgtype.Text{})
}

// Validate validates the given struct using the validator instance.
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// validateValuer implements validator.CustomTypeFunc.
func validateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
	}
	return nil
}

// validateNoSpecialChar disallows specific special characters.
func validateNoSpecialChar(fl validator.FieldLevel) bool {
	return !noSpecialCharPattern.MatchString(fl.Field().String())
}

// validateNoSpecialCharSlash disallows specific special characters except slashes.
func validateNoSpecialCharSlash(fl validator.FieldLevel) bool {
	return !noSpecialCharSlashPattern.MatchString(fl.Field().String())
}
