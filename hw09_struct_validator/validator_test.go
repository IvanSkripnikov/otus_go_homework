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
		Phones []string        `validate:"len:11"`
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
				ID:     "18858ab860111919f6592976848b9be12301",
				Name:   "Neuvillet",
				Age:    30,
				Email:  "neuvillet@fontaine.ru",
				Role:   "admin",
				Phones: []string{"27999999999"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "44adf98870bc0943453535345545657567345554332",
				Name:   "Zhongli",
				Age:    6500,
				Email:  "zhongli@teyvat.ru",
				Role:   "master",
				Phones: []string{"2100123ooo7", "890012334"},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: ErrMismatchedLength},
				ValidationError{Field: "Age", Err: ErrMoreThanMax},
				ValidationError{Field: "Role", Err: ErrNotInList},
				ValidationError{Field: "Phones", Err: ErrMismatchedLength},
			},
		},
		{
			in: User{
				ID:     "1c0acc63962a43f30a2cf2f3e5f759839223",
				Name:   "Nahida",
				Age:    8,
				Email:  "nahida@mail.sum.teyv",
				Role:   "stuff",
				Phones: []string{"81077734567", "21401234000"},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Age", Err: ErrLessThanMin},
				ValidationError{Field: "Email", Err: ErrMismatchedRegexp},
			},
		},
		{
			in: App{
				Version: "4.0.1",
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
				Version: "1.16.0.2",
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
				Header:    []byte("X-Forwarded-For"),
				Payload:   []byte("{'id': 1, 'type': 'example'}"),
				Signature: []byte("7777ff173a27222b172e66d98cf4d2a1"),
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
				Body: "Success",
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

			// валидируем переданные структуры
			results := Validate(tt.in)
			okExpected := errors.As(tt.expectedErr, &expectedErrs)

			if errors.As(results, &errsValidate) && okExpected {
				for i, err := range errsValidate {
					expected := expectedErrs[i]

					require.Equal(t, expected.Field, err.Field, "Field Name does not match")
					require.ErrorIs(t, err.Err, expected.Err, "Error Text does not match")
				}
			}
		})
	}
}

func TestNotStructureError(t *testing.T) {
	t.Run("case error ErrIsNotStructure", func(t *testing.T) {
		err := Validate("just string")
		require.Truef(t, errors.Is(err, ErrIsNotStructure), "actual error %q", err)
	})
}
