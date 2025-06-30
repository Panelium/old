package global

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"sync"
)

var (
	initOnce            sync.Once
	validate            *validator.Validate
	usernameRegexString = `^[a-zA-Z0-9][a-zA-Z0-9_.-]{1,30}[a-zA-Z0-9]$` // Username must start and end with an alphanumeric character, can contain alphanumeric characters, underscores, hyphens, and dots in between, and must be between 3 and 32 characters long.
	usernameRegex       = regexp.MustCompile(usernameRegexString)
)

func InitValidator() error {
	var err error

	initOnce.Do(func() {
		validate = validator.New(validator.WithRequiredStructEnabled())

		err = validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
			return usernameRegex.MatchString(fl.Field().String())
		})
		if err != nil {
			return
		}
	})

	if err != nil {
		return err
	}
	return nil
}

func ValidatorInstance() *validator.Validate {
	if validate == nil {
		panic("validator not initialized, call InitValidator() first")
	}
	return validate
}
