package main

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrInvalidString = errors.New("invalid string")
	zeroRune         = '0'
	backspaceRune    = ' '
)

func isErrorManyDigits(symbolPrev, symbolCurrent rune) bool {
	if unicode.IsDigit(symbolCurrent) && unicode.IsDigit(symbolPrev) {
		return true
	}
	return false
}

func isErrorString(s string) bool {
	for index, symbol := range s {
		if index == 0 && unicode.IsDigit(symbol) {
			return true
		}
		if symbol == backspaceRune {
			return true
		}
		if index > 0 && isErrorManyDigits(rune(s[index-1]), symbol) {
			return true
		}
	}
	return false
}

func Unpack(s string) (string, error) {
	var builder strings.Builder
	lenString := utf8.RuneCountInString(s)
	if lenString == 0 {
		return "", nil
	}

	if isErrorString(s) {
		return "", ErrInvalidString
	}

	stringRunes := []rune(s)
	for index, symbol := range stringRunes {
		if index == lenString-1 && symbol == zeroRune {
			continue
		}

		if unicode.IsLetter(symbol) || unicode.IsSpace(symbol) {
			str := string(symbol)
			if index < lenString-1 && stringRunes[index+1] == zeroRune {
				continue
			}
			if index < lenString-1 && unicode.IsDigit(stringRunes[index+1]) {
				builder.WriteString(strings.Repeat(str, int(stringRunes[index+1]-'0')))
			} else {
				builder.WriteString(str)
			}
		}
	}

	return builder.String(), nil
}
