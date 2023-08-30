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
	s := "d\n5abc"
	var symbolRune rune
	var builder strings.Builder
	zeroRune := rune("0"[0])
	nlRune := rune("\n"[0])

	sa := []byte(s)

	for index, symbol := range sa {
		symbolRune = rune(symbol)
		if !unicode.IsDigit(symbolRune) && !unicode.IsLetter(symbolRune) && symbolRune != nlRune {
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

		if unicode.IsLetter(rune(symbol)) || rune(symbol) == nlRune {
			str := ""
			if unicode.IsLetter(rune(symbol)) {
				str = string(symbol)
			} else {
				str = `\n`
			}
			if index < len(s)-1 && rune(sa[index+1]) == zeroRune {
				continue
			}
			if index < len(s)-1 && unicode.IsDigit(rune(sa[index+1])) {
				builder.WriteString(strings.Repeat(str, int(rune(sa[index+1])-'0')))
			} else {
				builder.WriteString(str)
			}
		}
	}

	fmt.Println(builder.String())
}
