package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

func DefaultSorting(n string) string {
	arr := strings.Split(n, "\n")
	sort.Strings(arr)
	res := strings.Join(arr, "\n")
	return res
}

func SortByColumn(n string, k int, numeric bool) (string, error) {
	arr := strings.Split(n, "\r\n")
	var columns [][]string
	for _, v := range arr {
		column := strings.Split(v, " ")
		if k >= len(column) {
			return "", fmt.Errorf("column %v doesn't exists at %v string", k, v)
		}
		columns = append(columns, column)
	}
	columns, err := sortByNumeric(columns, k, numeric)
	if err != nil {
		return "", fmt.Errorf("error on sort by numeric: %v", err)
	}

	arr = []string{}
	for _, v := range columns {
		arr = append(arr, strings.Join(v, " "))
	}
	return strings.Join(arr, "\n"), nil
}

func sortByNumeric(r [][]string, k int, numeric bool) ([][]string, error) {
	if numeric == false {
		sort.SliceStable(r, func(i, j int) bool {
			return r[i][k] < r[j][k]
		})
		return r, nil
	}
	var err error
	sort.SliceStable(r, func(i, j int) bool {
		a, errAtoi := strconv.Atoi(r[i][k])
		if errAtoi != nil {
			err = errAtoi
			return false
		}
		b, errAtoi := strconv.Atoi(r[j][k])
		if errAtoi != nil {
			err = errAtoi
			return false
		}
		return a < b
	})
	if err != nil {
		return r, err
	}
	return r, nil
}

func parseCopy(s string) string {
	var m map[string]struct{}

	m = make(map[string]struct{})

	newArr := make([]string, 0)
	arr := strings.Split(s, "\n")

	for _, v := range arr {
		if _, ok := m[v]; ok {
			continue
		}
		newArr = append(newArr, v)
		m[v] = struct{}{}
	}
	return strings.Join(newArr, "\n")
}

func main() {
	var filename string
	var column int
	var byNumber bool
	var skipCopy bool

	flag.IntVar(&column, "k", -1, "column by which will be sorted")
	flag.BoolVar(&byNumber, "n", false, "by number")
	flag.BoolVar(&skipCopy, "u", false, "skip copy of strings")
	flag.Parse()

	filename = flag.Arg(0)

	textBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	body := string(textBytes)

	var result string

	if column == -1 {
		result = DefaultSorting(body)
	} else {
		result, err = SortByColumn(body, column, byNumber)
	}

	if err != nil {
		fmt.Printf("\nerror on sorting by column: %v", err)
		return
	}

	if skipCopy {
		result = parseCopy(result)
	}

	ioutil.WriteFile("sorted_"+filename, []byte(result), 777)
}
