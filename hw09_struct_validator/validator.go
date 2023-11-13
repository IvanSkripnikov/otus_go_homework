package main

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrLessThanMin      = errors.New("value is less than min")
	ErrMoreThanMax      = errors.New("value is more than max")
	ErrMismatchedRegexp = errors.New("value is not matching regexp")
	ErrMismatchedLength = errors.New("value is not matching length")
	ErrNotInList        = errors.New("value is not in list")
	ErrUnhandledType    = errors.New("unknown type for value")
	ErrIsNotStructure   = errors.New("value is not structure")
	ErrConvertToNumber  = errors.New("can't convert value to number")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

// пределяем константы названий всех обработчиков.
const (
	funcMin    = "min"
	funcMax    = "max"
	funcLen    = "len"
	funcIn     = "in"
	funcRegexp = "regexp"
)

type CheckRules struct {
	Min           int
	Max           int
	Len           int
	Regexp        string
	InRange       []string
	Types         []string
	ErrValidators []error
}

func (v ValidationErrors) Error() string {
	builder := strings.Builder{}

	for _, validationErr := range v {
		errMessage := fmt.Sprintf("%s: %v; ", validationErr.Field, validationErr.Err)
		builder.WriteString(errMessage)
	}

	return builder.String()
}

func Validate(v interface{}) error {
	itemType := reflect.TypeOf(v)
	if itemType.Kind() != reflect.Struct {
		return ErrIsNotStructure
	}

	itemValue := reflect.ValueOf(v)
	numFields := itemType.NumField()
	validateErrors := make([]ValidationError, 0, numFields)

	// начинаем рассматривать поля
	for i := 0; i < numFields; i++ {
		fieldOfType := itemType.Field(i)
		fieldValue := itemValue.Field(i)
		fieldName := fieldOfType.Name

		// если поле не публичное - пропускаем
		if !fieldOfType.IsExported() {
			continue
		}

		// получаем теги для текщего поля
		tagName := fieldOfType.Tag
		tagValidate, ok := tagName.Lookup("validate")
		if !ok || len(tagName) == 0 {
			continue
		}

		// получаем правила валидации
		checkRules := getCheckRules(tagValidate)
		if len(checkRules.ErrValidators) > 0 {
			validateErrors = append(validateErrors, checkRules.getCheckErrors(fieldName)...)
		}

		switch fieldValue.Kind() { //nolint:exhaustive
		case reflect.Int:
			errsValidate := checkRules.validateInt(int(fieldValue.Int()), fieldName)
			if len(errsValidate) > 0 {
				validateErrors = append(validateErrors, errsValidate...)
			}

		case reflect.String:
			errsValidate := checkRules.validateString(fieldValue.String(), fieldName)
			if len(errsValidate) > 0 {
				validateErrors = append(validateErrors, errsValidate...)
			}

		case reflect.Slice:
			var errsValidate []ValidationError

			switch items := fieldValue.Interface().(type) {
			case []int:
				errsValidate = checkRules.validateSliceInt(items, fieldName)

			case []string:
				errsValidate = checkRules.validateSliceString(items, fieldName)
			}

			if len(errsValidate) > 0 {
				validateErrors = append(validateErrors, errsValidate...)
			}

		default:
			validateErrors = append(validateErrors, ValidationError{
				Field: fieldName,
				Err:   ErrUnhandledType,
			})
		}
	}

	if len(validateErrors) > 0 {
		return ValidationErrors(validateErrors)
	}

	return nil
}

func (checkRules CheckRules) getCheckErrors(fieldName string) []ValidationError {
	validateErrors := make([]ValidationError, 0, len(checkRules.ErrValidators))

	for _, validationErr := range checkRules.ErrValidators {
		validateErrors = append(validateErrors, ValidationError{
			Field: fieldName,
			Err:   validationErr,
		})
	}

	return validateErrors
}

// валидация int.
func (checkRules CheckRules) validateInt(value int, fieldName string) []ValidationError {
	validateErrors := make([]ValidationError, 0, len(checkRules.Types))

	for _, validationType := range checkRules.Types {
		var errMessage error

		switch validationType {
		case funcMin:
			hasValid := value >= checkRules.Min
			if !hasValid {
				errMessage = ErrLessThanMin
			}

		case funcMax:
			hasValid := value <= checkRules.Max
			if !hasValid {
				errMessage = ErrMoreThanMax
			}

		case funcIn:
			hasValid, err := inArray(value, reflect.Int, checkRules.InRange)
			if err != nil {
				validateErrors = append(validateErrors, ValidationError{
					Field: fieldName,
					Err:   err,
				})
			}

			if !hasValid {
				errMessage = ErrNotInList
			}
		}

		if errMessage != nil {
			validateErrors = append(validateErrors, ValidationError{
				Field: fieldName,
				Err:   errMessage,
			})
		}
	}

	return validateErrors
}

// валидация string.
func (checkRules CheckRules) validateString(value string, fieldName string) []ValidationError {
	validateErrors := make([]ValidationError, 0, len(checkRules.Types))

	for _, validationType := range checkRules.Types {
		var errMessage error

		switch validationType {
		case funcLen:
			hasValid := len(value) == checkRules.Len
			if !hasValid {
				errMessage = ErrMismatchedLength
			}

		case funcIn:
			hasValid, err := inArray(value, reflect.String, checkRules.InRange)
			if err != nil {
				validateErrors = append(validateErrors, ValidationError{
					Field: fieldName,
					Err:   err,
				})
			}

			if !hasValid {
				errMessage = ErrNotInList
			}

		case funcRegexp:
			hasValid, err := regexp.MatchString(checkRules.Regexp, value)
			if err != nil {
				validateErrors = append(validateErrors, ValidationError{
					Field: fieldName,
					Err:   err,
				})
			}

			if !hasValid {
				errMessage = ErrMismatchedRegexp
			}
		}

		if errMessage != nil {
			validateErrors = append(validateErrors, ValidationError{
				Field: fieldName,
				Err:   errMessage,
			})
		}
	}

	return validateErrors
}

// валидация []int.
func (checkRules CheckRules) validateSliceInt(items []int, fieldName string) []ValidationError {
	validateErrors := make([]ValidationError, 0, len(items))

	for _, itemValue := range items {
		errorsValidate := checkRules.validateInt(itemValue, fieldName)
		if len(errorsValidate) > 0 {
			validateErrors = append(validateErrors, errorsValidate...)
		}
	}

	return validateErrors
}

// валидация []string.
func (checkRules CheckRules) validateSliceString(items []string, fieldName string) []ValidationError {
	validateErrors := make([]ValidationError, 0, len(items))

	for _, itemValue := range items {
		errsValidate := checkRules.validateString(itemValue, fieldName)
		if len(errsValidate) > 0 {
			validateErrors = append(validateErrors, errsValidate...)
		}
	}

	return validateErrors
}

func getCheckRules(tagValidate string) CheckRules {
	var checkRules CheckRules
	checks := strings.Split(tagValidate, "|")

	for _, check := range checks {
		validationType := strings.Split(check, ":")[0]

		switch validationType {
		case funcMin:
			value, err := convertValidatorToInt(check, funcMin)
			if err != nil {
				checkRules.ErrValidators = append(checkRules.ErrValidators, err)
				continue
			}

			checkRules.Min = value
			checkRules.Types = append(checkRules.Types, funcMin)

		case funcMax:
			value, err := convertValidatorToInt(check, funcMax)
			if err != nil {
				checkRules.ErrValidators = append(checkRules.ErrValidators, err)
				continue
			}

			checkRules.Max = value
			checkRules.Types = append(checkRules.Types, funcMax)

		case funcLen:
			value, err := convertValidatorToInt(check, funcLen)
			if err != nil {
				checkRules.ErrValidators = append(checkRules.ErrValidators, err)
				continue
			}

			checkRules.Len = value
			checkRules.Types = append(checkRules.Types, funcLen)

		case funcIn:
			checkRules.InRange = strings.Split(getValidatorValue(check, funcIn), ",")
			checkRules.Types = append(checkRules.Types, funcIn)

		case funcRegexp:
			checkRules.Regexp = getValidatorValue(check, funcRegexp)
			checkRules.Types = append(checkRules.Types, funcRegexp)
		}
	}

	return checkRules
}

// переконвертировать значение к числу.
func convertValidatorToInt(validator, validatorType string) (int, error) {
	var resErr error
	value, err := strconv.Atoi(getValidatorValue(validator, validatorType))
	if err != nil {
		resErr = ErrConvertToNumber
	}

	return value, resErr
}

func getValidatorValue(validator, replacement string) string {
	return strings.Replace(validator, replacement+":", "", 1)
}

func inArray(value interface{}, kind reflect.Kind, array []string) (bool, error) {
	for _, v := range array {
		switch kind { //nolint:exhaustive
		case reflect.Int:
			itemValue, err := strconv.Atoi(v)
			if err != nil {
				return false, err
			}

			if value == itemValue {
				return true, nil
			}

		case reflect.String:
			if value == v {
				return true, nil
			}
		}
	}

	return false, nil
}
