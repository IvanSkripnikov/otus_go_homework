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

func isErrorDigitSymbol(index int, sa []rune) bool {
	if index == 0 || unicode.IsDigit(sa[index-1]) {
		return true
	}
	return false
}

func isErrorString(symbol rune, index int, sa []rune) bool {
	isErrorUnknownLetter := !unicode.IsDigit(symbol) && !unicode.IsLetter(symbol) && !unicode.IsSpace(symbol)
	isErrorManyDigits := unicode.IsDigit(symbol) && isErrorDigitSymbol(index, sa)
	if isErrorUnknownLetter || isErrorManyDigits {
		return true
	}
	return false
}

func Unpack(s string) (string, error) {
	var builder strings.Builder
	lenString := utf8.RuneCountInString(s)

	sa := []rune(s)

	for index, symbol := range sa {
		if isErrorString(symbol, index, sa) {
			return "", ErrInvalidString
		}

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
