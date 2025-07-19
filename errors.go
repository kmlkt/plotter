package plotter

import (
	"errors"
	"fmt"
)

var (
	errorNoKeys          = errors.New("Need at least 1 key, no keys given")
	errorInvalidKey      = errors.New("Invalid key: %v. Only alphanumeric characters, underscores & dashes are supported")
	errorInvalidFormat   = errors.New("Invalid response format: %v. Supported values: .html (default), .svg, .csv and .json")
	errorInvalidBody     = errors.New("Invalid body: %s, reason: %s. Only numbers with decimal point like 12 or 3.14 are supported")
	errorInvalidKeyCount = errors.New("Only 1 key is supported in POST request, %v given")
	errorKeyExists       = errors.New("Cannot create %s key %s, because it already exists")
	errorKeyNotFound     = errors.New("Table with %s key %s not found")
	errorKeyNoPermission = errors.New("Key %s cannot %s")
	errorInvalidSymbol   = errors.New("%v symbol is %s, expected digit")
	errorStringEmpty     = errors.New("string is empty")
)

func formatError(base error, args ...any) error {
	return FormattedError{base, args}
}

type FormattedError struct {
	Base error
	Args []any
}

func (err FormattedError) Error() string {
	return fmt.Sprintf(err.Base.Error(), err.Args...)
}

func (err FormattedError) Unwrap() error {
	return err.Base
}

func expectError(given error, expected error, notFound error) error {
	switch {
	case given == nil:
		return notFound
	case errors.Is(given, expected):
		return nil
	default:
		return given
	}
}
