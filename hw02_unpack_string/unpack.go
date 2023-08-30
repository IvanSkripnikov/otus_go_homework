package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(_ string) (string, error) {
	// Place your code here.
	return "", nil
}

func main() {
	s := "a4bc2d5ea0"
	var symbolRune rune
	var builder strings.Builder
	zeroRune := rune("0"[0])
	sa := []byte(s)

	for index, symbol := range sa {
		symbolRune = rune(symbol)
		if !unicode.IsDigit(symbolRune) && !unicode.IsLetter(symbolRune) {
			fmt.Println("shit1")
		}

		if index == len(s)-1 && symbolRune == zeroRune {
			continue
		}

		if unicode.IsDigit(rune(symbol)) {
			if index == 0 {
				fmt.Println("shit2")
			} else {
				if unicode.IsDigit(rune(sa[index-1])) {
					fmt.Println("shit3")
				}
			}
		}

		if unicode.IsLetter(rune(symbol)) {
			if index < len(s)-1 && rune(sa[index+1]) == zeroRune {
				continue
			}
			if unicode.IsDigit(rune(sa[index+1])) {
				builder.WriteString(strings.Repeat(string(symbol), int(rune(sa[index+1])-'0')))
			} else {
				builder.WriteString(string(symbol))
			}
		}
	}
	fmt.Println(builder.String())
	fmt.Println(s)
}
