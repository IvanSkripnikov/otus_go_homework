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
	ErrMismatchedRegexp = errors.New("value is not mathing regexp")
	ErrMismatchedLength = errors.New("value is not mathing length")
	ErrNotInList        = errors.New("value is not in list")
	ErrUnknownCheck     = errors.New("unknown check for value")
	ErrIsNotStructure   = errors.New("value is not structure")
)

type ValidationError struct {
	Field string
	Err   error
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
	t := reflect.TypeOf(v)

	if t.Kind() != reflect.Struct {
		return ErrIsNotStructure
	}

	var (
		metaData []string
		check    []string
	)

	itemValue := reflect.ValueOf(v)
	numFields := t.NumField()
	validateErrors := make([]ValidationError, 0, numFields)

	for i := 0; i < numFields; i++ {
		fieldType := t.Field(i)
		value := itemValue.Field(i)
		fieldName := fieldType.Name

		// Проверяем является ли поле публичным
		if !fieldType.IsExported() {
			continue
		}

		metaData = getMetaData(t.Field(i))
		if metaData == nil {
			continue
		}

		for j := 0; j < len(metaData); j++ {
			check = strings.Split(metaData[j], ":")
			obj := getFuncFromValidator(check[0])
			f, _ := obj.(func(string, interface{}) interface{})
			z := getValuesFromReflectValue(value)
			r := f(z[0], check[1])
			if r != nil {
				validateErrors = append(validateErrors, ValidationError{Field: fieldName, Err: r.(error)})
			}
		}
	}

	if len(validateErrors) > 0 {
		return ValidationErrors(validateErrors)
	}

	return nil
}

func getMetaData(field reflect.StructField) []string {
	val := field.Tag.Get("validate")
	if val == "" {
		return nil
	}
	return strings.Split(val, "|")
}

func getFuncFromValidator(name string) interface{} {
	switch name {
	case "min":
		return func(v string, min interface{}) interface{} {
			value, _ := strconv.Atoi(v)

			minValue := min.(string)
			minValueInt, errConvert := strconv.Atoi(minValue)
			if errConvert != nil {
				return errConvert
			}

			if value < minValueInt {
				return ErrLessThanMin
			}
			return nil
		}
	case "max":
		return func(v string, max interface{}) interface{} {
			value, _ := strconv.Atoi(v)

			maxValue := max.(string)
			maxValueInt, errConvert := strconv.Atoi(maxValue)
			if errConvert != nil {
				return errConvert
			}

			if value > maxValueInt {
				return ErrMoreThanMax
			}
			return nil
		}
	case "regexp":
		return func(v string, pattern interface{}) interface{} {
			valPattern := pattern.(string)
			matched, err := regexp.MatchString(valPattern, v)
			if err != nil {
				return err
			}
			if !matched {
				return ErrMismatchedRegexp
			}
			return nil
		}
	case "len":
		return func(v string, length interface{}) interface{} {
			lenValue := length.(string)
			lenValueInt, errConvert := strconv.Atoi(lenValue)

			if errConvert != nil {
				return errConvert
			}
			if len(v) != lenValueInt {
				return ErrMismatchedLength
			}

			return nil
		}
	case "in":
		return func(v string, in interface{}) interface{} {
			valInString := in.(string)
			valInArray := strings.Split(valInString, ",")

			if !inArray(v, valInArray) {
				return ErrNotInList
			}
			return nil
		}
	default:
		return ErrUnknownCheck
	}
}

func inArray(val string, array []string) bool {
	for _, v := range array {
		if val == v {
			return true
		}
	}

	return false
}

func getValuesFromReflectValue(v reflect.Value) []string {
	val := v.Interface()
	switch v.Kind() {
	case reflect.Slice:
		return val.([]string)
	case reflect.Array:
		return val.([]string)
	default:
		value := make([]string, 1)
		if v.Kind() == reflect.Int {
			value[0] = strconv.Itoa(val.(int))
		} else {
			value[0] = v.String()
		}

		return value
	}
}
