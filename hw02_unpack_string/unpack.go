package main

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func getStringSymbol(symbol byte) (str string) {
	if unicode.IsLetter(rune(symbol)) {
		str = string(symbol)
	} else {
		str = `\n`
	}
	return
}

func Unpack(s string) (string, error) {
	var symbolRune rune
	var builder strings.Builder
	zeroRune := rune("0"[0])
	nlRune := rune("\n"[0])
	lenString := len(s)

	sa := []byte(s)

	for index, symbol := range sa {
		symbolRune = rune(symbol)
		if !unicode.IsDigit(symbolRune) && !unicode.IsLetter(symbolRune) && symbolRune != nlRune {
			return "", ErrInvalidString
		}

		if index == lenString-1 && symbolRune == zeroRune {
			continue
		}

		if unicode.IsDigit(rune(symbol)) {
			if index == 0 {
				return "", ErrInvalidString
			} else if unicode.IsDigit(rune(sa[index-1])) {
				return "", ErrInvalidString
			}
		}

		if unicode.IsLetter(rune(symbol)) || rune(symbol) == nlRune {
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
