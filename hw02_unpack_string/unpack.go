package main

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrInvalidString = errors.New("invalid string")
	zeroRune         = rune(`0`[0])
	slashRune        = rune(`\`[0])
)

func getStringSymbol(symbol rune, withSlash bool) string {
	str := string(symbol)
	if withSlash {
		str = `\` + str
	}

	return str
}

func isErrorUnknownChar(symbol rune) bool {
	return !unicode.IsDigit(symbol) && !unicode.IsLetter(symbol) && symbol != slashRune
}

func isErrorManyDigits(symbol rune, index int, s string) bool {
	if unicode.IsDigit(symbol) {
		if index == 0 || unicode.IsDigit(rune(s[index-1])) {
			return true
		}
	}
	return false
}

func isErrorString(s string) bool {
	for index, symbol := range s {
		if isErrorUnknownChar(symbol) || isErrorManyDigits(symbol, index, s) {
			return true
		}
	}
	return false
}

func Unpack(s string) (string, error) {
	var builder strings.Builder
	lenString := utf8.RuneCountInString(s)

	if isErrorString(s) {
		return "", ErrInvalidString
	}

	sa := []rune(s)
	var flagSlash bool
	for index, symbol := range sa {
		if index == lenString-1 && symbol == zeroRune {
			continue
		}

		if symbol == slashRune {
			flagSlash = true
			continue
		}

		if unicode.IsLetter(symbol) || unicode.IsControl(symbol) {
			str := getStringSymbol(symbol, flagSlash)
			flagSlash = false

			if index < lenString-1 && sa[index+1] == zeroRune {
				continue
			}
			if index < lenString-1 && unicode.IsDigit(sa[index+1]) {
				builder.WriteString(strings.Repeat(str, int(sa[index+1]-'0')))
			} else {
				builder.WriteString(str)
			}
		}
	}

	return builder.String(), nil
}
