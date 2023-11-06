package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"regexp:^\\d+$|len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "24458918609ae919f6592976848b9be123ab",
				Name:   "John Smith",
				Age:    25,
				Email:  "admin@example.ru",
				Role:   "admin",
				Phones: []string{"27401234567"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "4d0acc63962a43f30a2cfff3e4f759839182",
				Name:   "John Doe",
				Age:    52,
				Email:  "test.ru",
				Role:   "guest",
				Phones: []string{"7900123ab67", "790012345"},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Age", Err: ErrMoreThanMax},
				ValidationError{Field: "Email", Err: ErrMismatchedRegexp},
				ValidationError{Field: "Role", Err: ErrNotInList},
				ValidationError{Field: "Phones", Err: ErrMismatchedRegexp},
				ValidationError{Field: "Phones", Err: ErrMismatchedLength},
			},
		},
		{
			in: User{
				ID:     "1f2ffb504440fda52c3fd95ba4252ab1abcd1234",
				Name:   "Jane Doe",
				Age:    16,
				Email:  "janedoe@example.ru",
				Role:   "stuff",
				Phones: []string{"79001234567", "27401234567"},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: ErrMismatchedLength},
				ValidationError{Field: "Age", Err: ErrLessThanMin},
			},
		},
		{
			in: App{
				Version: "1.0.5",
			},
			expectedErr: nil,
		},
		{
			in: App{},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Version", Err: ErrMismatchedLength},
			},
		},
		{
			in: App{
				Version: "1.1.3-rc.2",
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Version", Err: ErrMismatchedLength},
			},
		},
		{
			in: Token{}, expectedErr: nil,
		},
		{
			in: Token{
				Header:    []byte("X-Real-IP"),
				Payload:   []byte("{'id': 34, 'type': 'example'}"),
				Signature: []byte("f6f7ff173a277e3b172e66d98cf4d2ac"),
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
				Body: "Successfully",
			},
			expectedErr: nil,
		},
		{
			in: Response{},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Code", Err: ErrNotInList},
			},
		},
		{
			in: Response{
				Code: 204,
				Body: "Out of range",
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Code", Err: ErrNotInList},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			var (
				expectedErrs ValidationErrors
				errsValidate ValidationErrors
			)

			tt := tt
			t.Parallel()

			// Запускаем валидацию и проверяем ошибки
			resErrors := Validate(tt.in)
			okExpected := errors.As(tt.expectedErr, &expectedErrs)

			if errors.As(resErrors, &errsValidate) && okExpected {
				for i, err := range errsValidate {
					expected := expectedErrs[i]

					require.Equal(t, expected.Field, err.Field, "Field name does not match")
					require.ErrorIs(t, err.Err, expected.Err, "Error text does not match")
					require.ErrorAs(t, err.Err, &expected.Err, "Error data type does not match")
				}
			}
		})
	}
}
