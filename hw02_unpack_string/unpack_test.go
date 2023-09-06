package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "aa3b", expected: "aaaab"},
		{input: "d\n5abcd", expected: "d\n\n\n\n\nabcd"},
		{input: "a2b3c4", expected: "aabbbcccc"},
		{input: "aa0a0b", expected: "ab"},
		{input: `a4б4e2`, expected: `aaaaббббee`},
		{input: `本4异2с2`, expected: `本本本本异异сс`},
		{input: "a\t2\n3c4", expected: "a\t\t\n\n\ncccc"},
		{input: "a\v3bcd", expected: "a\v\v\vbcd"},
		{input: "aa🔥0b", expected: "aab"},
		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", "bc2b4a15f", "aa a3b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
