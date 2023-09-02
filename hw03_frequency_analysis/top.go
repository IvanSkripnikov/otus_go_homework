package main

import (
	"fmt"
	"sort"
	"strings"
)

func Top10(_ string) []string {
	// Place your code here.
	return nil
}

type Word struct {
	Title string
	Count int
}

func getMappedWords(strings []string) map[string]int {
	resultMap := map[string]int{}

	for _, word := range strings {
		_, ok := resultMap[word]
		if ok == false {
			resultMap[word] = 1
		} else {
			resultMap[word]++
		}
	}

	return resultMap
}

func main() {
	str := `cat and dog, one dog,two cats and one man and one child. Child play with dog and cats. Dogs and cats is a pets`
	resultMap := getMappedWords(strings.Split(str, ` `))

	resultStruct := []Word{}
	for key, value := range resultMap {
		w := Word{
			Title: key,
			Count: value,
		}
		resultStruct = append(resultStruct, w)
	}

	sort.Slice(resultStruct, func(i, j int) bool {
		return resultStruct[i].Count > resultStruct[j].Count
	})

	fmt.Println(resultStruct[:10])
}
