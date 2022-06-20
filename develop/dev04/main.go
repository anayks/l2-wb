package main

import (
	"fmt"
	"sort"
	"strings"
)

func getMap(arr []string) map[string][]string {
	annMap := map[string]string{}
	resultMap := map[string][]string{}
	existedWord := map[string]struct{}{}

	for _, value := range arr {
		value = strings.ToLower(value)
		if _, ok := existedWord[value]; ok {
			continue
		}
		existedWord[value] = struct{}{}
		ann := getWordSymbols(value)
		if v, ok := annMap[ann]; ok {
			resultMap[v] = append(resultMap[v], value)
			continue
		}
		annMap[ann] = value
		resultMap[value] = make([]string, 0)
		resultMap[value] = append(resultMap[value], value)
	}

	for k := range resultMap {
		sort.Strings(resultMap[k])
	}

	for k := range resultMap {
		if len(resultMap[k]) == 1 {
			delete(resultMap, k)
		}
	}

	return resultMap
}

func getWordSymbols(s string) string {
	runes := []rune(s)
	ints := []int{}
	for _, v := range runes {
		ints = append(ints, int(v))
	}
	sort.Ints(ints)
	runes = []rune{}
	for _, v := range ints {
		runes = append(runes, rune(v))
	}
	return string(runes)
}

func main() {
	values := []string{"TEST", "estt", "test", "rest"}
	result := getMap(values)
	fmt.Printf("%v word symbols", result)
}
