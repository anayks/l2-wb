package main

import (
	"bufio"
	"os"
)

func getFileStrings(filename string) ([]string, error) {
	var strings = []string{}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strings = append(strings, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return strings, err
	}
	return strings, nil
}
