package main

import (
	"errors"
	"strings"
	"unicode"
)

var (
	ErrInvalidString = errors.New("invalid string")
	zeroRune         = rune("0"[0])
	nlRune           = rune("\n"[0])
)

func getStringSymbol(symbol byte) (str string) {
	if unicode.IsLetter(rune(symbol)) {
		str = string(symbol)
	} else {
		str = `\n`
	}
	return
}

func isErrorDigitSymbol(index int, sa []byte) bool {
	if index == 0 || unicode.IsDigit(rune(sa[index-1])) {
		return true
	}
	return false
}

func isErrorString(symbol byte, index int, sa []byte) bool {
	symbolRune := rune(symbol)
	isErrorUnknownLetter := !unicode.IsDigit(symbolRune) && !unicode.IsLetter(symbolRune) && symbolRune != nlRune
	isErrorManyDigits := unicode.IsDigit(symbolRune) && isErrorDigitSymbol(index, sa)
	if isErrorUnknownLetter || isErrorManyDigits {
		return true
	}
	return false
}

func Unpack(s string) (string, error) {
	var symbolRune rune
	var builder strings.Builder
	lenString := len(s)

	sa := []byte(s)

	for index, symbol := range sa {
		symbolRune = rune(symbol)

		if isErrorString(symbol, index, sa) {
			return "", ErrInvalidString
		}

		if index == lenString-1 && symbolRune == zeroRune {
			continue
		}

		if unicode.IsLetter(symbolRune) || symbolRune == nlRune {
			str := getStringSymbol(symbol)
			if index < lenString-1 && rune(sa[index+1]) == zeroRune {
				continue
			}
			if index < lenString-1 && unicode.IsDigit(rune(sa[index+1])) {
				builder.WriteString(strings.Repeat(str, int(rune(sa[index+1])-'0')))
			} else {
				builder.WriteString(str)
			}
		}
	}

	return builder.String(), nil
}
