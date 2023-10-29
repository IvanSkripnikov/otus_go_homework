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
	ErrLessThanMin      = errors.New("Value is less than min")
	ErrMoreThanMax      = errors.New("Value is more than max")
	ErrMismatchedRegexp = errors.New("Value is not mathing regexp")
	ErrMismatchedLength = errors.New("Value is not mathing length")
	ErrNotInList        = errors.New("Value is not in list")
	ErrUnknownCheck     = errors.New("Unknown check for value")
	ErrIsEmpty          = errors.New("Value is empty")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	// Place your code here.
	return nil
}

type UserRole1 string

type User1 struct {
	ID     string `json:"id" validate:"len:36"`
	Name   string
	Age    int       `validate:"min:18|max:50"`
	Email  string    `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	Role   UserRole1 `validate:"in:admin,stuff"`
	Phones []string  `validate:"len:11"`
}

func main() {
	user := User1{ID: "34242223", Name: "John", Age: 33, Email: "tramak@mail.ru", Role: "admin"}

	t := reflect.TypeOf(user)
	var (
		metaData       []string
		check          []string
		fieldName      string
		validateErrors []error
		value          reflect.Value
	)
	for i := 0; i < t.NumField(); i++ {
		fieldName = t.Field(i).Name
		metaData = getMetaData(t.Field(i))
		if metaData == nil {
			continue
		}
		for j := 0; j < len(metaData); j++ {
			value = reflect.ValueOf(&user).Elem().FieldByName(fieldName)

			if reflect.ValueOf(user).FieldByName(fieldName).IsZero() {
				validateErrors = append(validateErrors, ErrIsEmpty)
				continue
			}

			check = strings.Split(metaData[j], ":")
			obj := getFuncFromValidator(check[0])
			f, _ := obj.(func(string, interface{}) interface{})
			z := getValuesFromReflectValue(value)
			r := f(z[0], check[1])
			if r != nil {
				validateErrors = append(validateErrors, r.(error))
			}
		}
	}
	fmt.Println(validateErrors)
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

			value, _ := strconv.Atoi(v)
			if value > lenValueInt || value < lenValueInt {
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
		var value = make([]string, 1)
		if v.Kind() == reflect.Int {
			value[0] = strconv.Itoa(val.(int))
		} else {
			value[0] = v.String()
		}

		return value
	}
}
