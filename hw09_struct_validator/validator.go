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

const (
	typeMin    = "min"
	typeMax    = "max"
	typeLen    = "len"
	typeIn     = "in"
	typeRegexp = "regexp"
)

type ValidationRules struct {
	Min           int
	Max           int
	Len           int
	Regexp        string
	InRange       []string
	Types         []string
	ErrValidators []error
}

type ValidationErrors []ValidationError

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
	resErrors := make([]ValidationError, 0, numFields)

	for i := 0; i < numFields; i++ {
		fieldOfType := itemType.Field(i)
		fieldValue := itemValue.Field(i)
		fieldName := fieldOfType.Name

		// Проверяем является ли поле публичным
		if !fieldOfType.IsExported() {
			continue
		}

		// Получаем теги для текщего поля
		fieldTag := fieldOfType.Tag
		tagValidate, ok := fieldTag.Lookup("validate")
		if !ok || len(fieldTag) == 0 {
			continue
		}

		// Формируем массив правил валидации
		rules := makeValidators(tagValidate)
		if len(rules.ErrValidators) > 0 {
			resErrors = append(resErrors, rules.makeRulesErrors(fieldName)...)
		}

		switch fieldValue.Kind() { //nolint:exhaustive
		case reflect.Int:
			errsValidate := rules.validateInt(int(fieldValue.Int()), fieldName)
			if len(errsValidate) > 0 {
				resErrors = append(resErrors, errsValidate...)
			}

		case reflect.String:
			errsValidate := rules.validateString(fieldValue.String(), fieldName)
			if len(errsValidate) > 0 {
				resErrors = append(resErrors, errsValidate...)
			}

		case reflect.Slice:
			var errsValidate []ValidationError

			switch items := fieldValue.Interface().(type) {
			case []int:
				errsValidate = rules.validateSliceInt(items, fieldName)

			case []string:
				errsValidate = rules.validateSliceString(items, fieldName)
			}

			if len(errsValidate) > 0 {
				resErrors = append(resErrors, errsValidate...)
			}

		default:
			resErrors = append(resErrors, ValidationError{
				Field: fieldName,
				Err:   ErrUnhandledType,
			})
		}
	}

	if len(resErrors) > 0 {
		return ValidationErrors(resErrors)
	}

	return nil
}

// Сформировать ошибки для правил валидации.
func (rules ValidationRules) makeRulesErrors(fieldName string) []ValidationError {
	resErrors := make([]ValidationError, 0, len(rules.ErrValidators))

	for _, validationErr := range rules.ErrValidators {
		resErrors = append(resErrors, ValidationError{
			Field: fieldName,
			Err:   validationErr,
		})
	}

	return resErrors
}

// Валидация целочисленных значений.
func (rules ValidationRules) validateInt(value int, fieldName string) []ValidationError {
	resErrors := make([]ValidationError, 0, len(rules.Types))

	for _, validationType := range rules.Types {
		var errMessage error

		switch validationType {
		case typeMin:
			hasValid := value >= rules.Min
			if !hasValid {
				errMessage = ErrLessThanMin
			}

		case typeMax:
			hasValid := value <= rules.Max
			if !hasValid {
				errMessage = ErrMoreThanMax
			}

		case typeIn:
			hasValid, err := sliceContains(value, reflect.Int, rules.InRange)
			if err != nil {
				resErrors = append(resErrors, ValidationError{
					Field: fieldName,
					Err:   err,
				})
			}

			if !hasValid {
				errMessage = ErrNotInList
			}
		}

		if errMessage != nil {
			resErrors = append(resErrors, ValidationError{
				Field: fieldName,
				Err:   errMessage,
			})
		}
	}

	return resErrors
}

// Валидация строковых значений.
func (rules ValidationRules) validateString(value string, fieldName string) []ValidationError {
	resErrors := make([]ValidationError, 0, len(rules.Types))

	for _, validationType := range rules.Types {
		var errMessage error

		switch validationType {
		case typeLen:
			hasValid := len(value) == rules.Len
			if !hasValid {
				errMessage = ErrMismatchedLength
			}

		case typeIn:
			hasValid, err := sliceContains(value, reflect.String, rules.InRange)
			if err != nil {
				resErrors = append(resErrors, ValidationError{
					Field: fieldName,
					Err:   err,
				})
			}

			if !hasValid {
				errMessage = ErrNotInList
			}

		case typeRegexp:
			hasValid, err := regexp.MatchString(rules.Regexp, value)
			if err != nil {
				resErrors = append(resErrors, ValidationError{
					Field: fieldName,
					Err:   err,
				})
			}

			if !hasValid {
				errMessage = ErrMismatchedRegexp
			}
		}

		if errMessage != nil {
			resErrors = append(resErrors, ValidationError{
				Field: fieldName,
				Err:   errMessage,
			})
		}
	}

	return resErrors
}

// Валидация целочисленных слайсов.
func (rules ValidationRules) validateSliceInt(items []int, fieldName string) []ValidationError {
	resErrors := make([]ValidationError, 0, len(items))

	for _, itemValue := range items {
		errsValidate := rules.validateInt(itemValue, fieldName)
		if len(errsValidate) > 0 {
			resErrors = append(resErrors, errsValidate...)
		}
	}

	return resErrors
}

// Валидация строковых слайсов.
func (rules ValidationRules) validateSliceString(items []string, fieldName string) []ValidationError {
	resErrors := make([]ValidationError, 0, len(items))

	for _, itemValue := range items {
		errsValidate := rules.validateString(itemValue, fieldName)
		if len(errsValidate) > 0 {
			resErrors = append(resErrors, errsValidate...)
		}
	}

	return resErrors
}

// Сформировать правила для валидации.
func makeValidators(tagValidate string) ValidationRules {
	var rules ValidationRules
	validators := strings.Split(tagValidate, "|")

	for _, validator := range validators {
		validationType := strings.Split(validator, ":")[0]

		switch validationType {
		case typeMin:
			value, err := castNumberValidator(validator, typeMin)
			if err != nil {
				rules.ErrValidators = append(rules.ErrValidators, err)
				continue
			}

			rules.Min = value
			rules.Types = append(rules.Types, typeMin)

		case typeMax:
			value, err := castNumberValidator(validator, typeMax)
			if err != nil {
				rules.ErrValidators = append(rules.ErrValidators, err)
				continue
			}

			rules.Max = value
			rules.Types = append(rules.Types, typeMax)

		case typeLen:
			value, err := castNumberValidator(validator, typeLen)
			if err != nil {
				rules.ErrValidators = append(rules.ErrValidators, err)
				continue
			}

			rules.Len = value
			rules.Types = append(rules.Types, typeLen)

		case typeIn:
			rules.InRange = strings.Split(getValidatorValue(validator, typeIn), ",")
			rules.Types = append(rules.Types, typeIn)

		case typeRegexp:
			rules.Regexp = getValidatorValue(validator, typeRegexp)
			rules.Types = append(rules.Types, typeRegexp)
		}
	}

	return rules
}

// Привести значение правила валидатора к числу.
func castNumberValidator(validator, validatorType string) (int, error) {
	var resErr error
	value, err := strconv.Atoi(getValidatorValue(validator, validatorType))
	if err != nil {
		resErr = ErrConvertToNumber
	}

	return value, resErr
}

// Получить значение для валидатора.
func getValidatorValue(validator, replacement string) string {
	return strings.Replace(validator, replacement+":", "", 1)
}

// Содержится ли значение в слайсе.
func sliceContains(value interface{}, kind reflect.Kind, inRanges []string) (bool, error) {
	for _, itemRange := range inRanges {
		switch kind { //nolint:exhaustive
		case reflect.Int:
			itemRangeValue, err := strconv.Atoi(itemRange)
			if err != nil {
				return false, err
			}

			if value == itemRangeValue {
				return true, nil
			}

		case reflect.String:
			if value == itemRange {
				return true, nil
			}
		}
	}

	return false, nil
}
