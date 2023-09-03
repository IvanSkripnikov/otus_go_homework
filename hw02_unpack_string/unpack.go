package main

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrInvalidString = errors.New("invalid string")
	zeroRune         = rune("0"[0])
	spaceSymbols     = map[rune]string{
		rune("\n"[0]): `\n`,
		rune("\t"[0]): `\t`,
		rune("\r"[0]): `\r`,
		rune("\v"[0]): `\v`,
		rune("\f"[0]): `\f`,
	}
)

func getStringSymbol(symbol rune) string {
	str := ``
	if unicode.IsLetter(symbol) {
		str = string(symbol)
	} else {
		str = spaceSymbols[symbol]
	}

	return str
}

func isErrorDigitSymbol(index int, sa string) bool {
	if index == 0 || unicode.IsDigit(rune(sa[index-1])) {
		return true
	}
	return false
}

func isErrorUnknownChar(symbol rune) bool {
	return !unicode.IsDigit(symbol) && !unicode.IsLetter(symbol) && isNotPossibleSpace(symbol)
}

func isNotPossibleSpace(symbol rune) bool {
	if unicode.IsSpace(symbol) {
		_, ok := spaceSymbols[symbol]
		if ok == false {
			return true
		}
	}
	return false
}

func isErrorString(s string) bool {
	for index, symbol := range s {
		isErrorManyDigits := unicode.IsDigit(symbol) && isErrorDigitSymbol(index, s)
		if isErrorUnknownChar(symbol) || isErrorManyDigits {
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
	for index, symbol := range sa {
		if index == lenString-1 && symbol == zeroRune {
			continue
		}

		if unicode.IsLetter(symbol) || unicode.IsSpace(symbol) {
			str := getStringSymbol(symbol)
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
