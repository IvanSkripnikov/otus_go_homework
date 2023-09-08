package main

import (
	"sort"
	"strings"
)

type Word struct {
	Title string
	Count int
}

func Top10(str string) []string {
	return getTopWordsList(str, 10)
}

func getTopWordsList(str string, limit int) []string {
	if len(str) == 0 {
		return make([]string, 0)
	}
	resultMap := getMappedWords(strings.Fields(str))
	resultStruct := getStructuredWords(resultMap)

	if len(resultStruct) < limit {
		limit = len(resultStruct)
	}
	top10Strings := make([]string, limit)
	for i := 0; i < limit; i++ {
		top10Strings[i] = resultStruct[i].Title
	}

	return top10Strings
}

func getStructuredWords(resultMap map[string]int) []Word {
	resultStruct := make([]Word, 0, len(resultMap))
	for key, value := range resultMap {
		w := Word{
			Title: key,
			Count: value,
		}
		resultStruct = append(resultStruct, w)
	}

	sort.Slice(resultStruct, func(i, j int) bool {
		return resultStruct[i].Count > resultStruct[j].Count ||
			(resultStruct[i].Count == resultStruct[j].Count && resultStruct[i].Title < resultStruct[j].Title)
	})

	return resultStruct
}

func getMappedWords(strings []string) map[string]int {
	resultMap := map[string]int{}

	for _, word := range strings {
		resultMap[word]++
	}

	return resultMap
}
