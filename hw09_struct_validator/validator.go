package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
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
	Age    int             `validate:"min:18|max:50"`
	Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	Role   UserRole1       `validate:"in:admin,stuff"`
	Phones []string        `validate:"len:11"`
	meta   json.RawMessage //nolint:unused
}

func main() {
	user := User1{ID: "34242223", Name: "John", Age: 33, Email: "tramak@mail.ru", Role: "admin"}

	t := reflect.TypeOf(user)
	var (
		metaData  []string
		check     []string
		result    map[string]bool
		fieldName string
	)
	for i := 0; i < t.NumField(); i++ {
		fieldName = t.Field(i).Name
		metaData = getMetaData(t.Field(i))
		for j := 0; j < len(metaData); j++ {
			check = strings.Split(metaData[j], ":")
			f := getFuncFromValidator(check[0])
			result[fieldName] = f(reflect.ValueOf(&user).Elem().FieldByName(fieldName), check[1])
		}
		fmt.Println(result)
	}
}

func getMetaData(field reflect.StructField) []string {
	val := field.Tag.Get("validate")
	return strings.Split(val, "|")
}

func getFuncFromValidator(name string) interface{} {
	switch name {
	case "min":
		return func(v int, min int) bool {
			if v < min {
				return false
			}
			return true
		}
	case "max":
		return func(v int, max int) bool {
			if v > max {
				return false
			}
			return true
		}
	case "regexp":
		return func(v string, pattern string) bool {
			matched, err := regexp.MatchString(pattern, v)
			if err != nil {
				return false
			}
			if !matched {
				return false
			}
			return true
		}
	case "len":
		return func(v string, length int) bool {
			if len(v) > length {
				return false
			}
			return true
		}
	default:
		return func(i ...interface{}) bool {
			return i != nil
		}
	}
}
