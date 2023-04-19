package validation

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/peano88/fizzbuzz-rest/pkg/model"
)

const (
	int1Constraint = "int1 should be a positive integer between 0 (excluding) and 9223372036854775807"
	int2Constraint = "int2 should be a positive integer between 0 (excluding) and 9223372036854775807"
	str1Constraint = "str1 can be any string not including character '-'"
	str2Constraint = "str2 can be any string not including character '-'"
)

// ValidationError is an error created in case of issue with the input parameters
type ValidationError struct {
	err        error
	parameter  string
	constraint string
}

// Error is the error interface implementation
func (ve ValidationError) Error() string {
	return fmt.Sprintf("validation error on parameter: %s", ve.err.Error())
}

// Unwrap allows correct use of errors.As and errors.Is
func (ve ValidationError) Unwrap() error {
	return ve.err
}

// Constraint returns the validation constraint triggering the validation error
func (ve ValidationError) Constraint() string {
	return ve.constraint
}

func mandatoryPositiveInteger(r *http.Request, param string) (int, error) {
	value := r.URL.Query().Get(param)
	if value == "" {
		return 0, errors.New("missing mandatory parameter: " + param)
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil || valueInt <= 0 {
		return 0, errors.New(param + ": " + value + " is not a positive integer")
	}

	return valueInt, nil
}

func mandatoryInteger(r *http.Request, param string) (int, error) {
	value := r.URL.Query().Get(param)
	if value == "" {
		return 0, errors.New("missing mandatory parameter: " + param)
	}

	return strconv.Atoi(value)
}

func optionalInteger(r *http.Request, param string) (int, error, bool) {
	value := r.URL.Query().Get(param)
	if value == "" {
		return 0, nil, false
	}

	intValue, err := strconv.Atoi(value)
	return intValue, err, true
}

func mandatoryString(r *http.Request, param string) (string, error) {
	value := r.URL.Query().Get(param)
	if value == "" {
		return "", errors.New("missing mandatory parameter: " + param)
	}

	if strings.Contains(value, "-") {
		return "", errors.New("parameter " + param + "contains illegal character '-'")
	}

	return value, nil
}

// Validator runs the different validation
type Validator struct {
}

// RunValidations runs the different validations and returns a ValidationError in case of issue
// if the validation is succesfull a modified context.Context is returned. This context is obtained
// by adding a model.FizzBuzzInput in the r.Context()
func (v *Validator) RunValidations(r *http.Request) (context.Context, error) {
	newContext := r.Context()

	int1, err := mandatoryPositiveInteger(r, "int1")
	if err != nil {
		return nil, ValidationError{
			err:        err,
			parameter:  "int1",
			constraint: int1Constraint,
		}
	}

	int2, err := mandatoryPositiveInteger(r, "int2")
	if err != nil {
		return nil, ValidationError{
			err:        err,
			parameter:  "int2",
			constraint: int2Constraint,
		}
	}

	limit, err := mandatoryInteger(r, "limit")
	if err != nil {
		return nil, ValidationError{
			err:       err,
			parameter: "limit",
		}
	}

	start, err, startProvided := optionalInteger(r, "start")
	if err != nil {
		return nil, ValidationError{
			err:       err,
			parameter: "start",
		}
	}

	str1, err := mandatoryString(r, "str1")
	if err != nil {
		return nil, ValidationError{
			err:        err,
			parameter:  "str1",
			constraint: str1Constraint,
		}
	}

	str2, err := mandatoryString(r, "str2")
	if err != nil {
		return nil, ValidationError{
			err:        err,
			parameter:  "str2",
			constraint: str2Constraint,
		}
	}

	input := model.FizzBuzzInput{
		FizzBuzzInputStats: model.FizzBuzzInputStats{
			Int1:  int1,
			Int2:  int2,
			Limit: limit,
			Str1:  str1,
			Str2:  str2,
		},
		Start: 1,
	}

	if startProvided {
		input.Start = start
	}

	newContext = context.WithValue(newContext, model.InputKey, input)

	return newContext, nil
}

// NewFizzBuzzValidator returns a new Validator
func NewFizzBuzzValidator() *Validator {
	return &Validator{}
}
