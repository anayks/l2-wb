package main

import (
	"io/ioutil"
	"testing"
)

func TestDefaultSort(t *testing.T) {
	textBytes, err := ioutil.ReadFile("test.txt")
	if err != nil {
		return
	}

	rByte, err := ioutil.ReadFile("test_default.txt")
	if err != nil {
		return
	}

	data := string(textBytes)
	result := string(rByte)

	testData := []struct {
		data     string
		expected string
	}{
		{
			data:     data,
			expected: result,
		},
	}

	t.Run("test_default", func(t *testing.T) {
		for id, v := range testData {
			result := DefaultSorting(v.data)
			if result != v.expected {
				t.Errorf("Unexpected result at %v", id)
			}
		}
	})
}

func TestSortByColumnByNumber(t *testing.T) {
	textBytes, err := ioutil.ReadFile("test.txt")
	if err != nil {
		return
	}

	rByte, err := ioutil.ReadFile("test_by_column_int.txt")
	if err != nil {
		return
	}

	data := string(textBytes)
	result := string(rByte)

	testData := []struct {
		data     string
		expected string
	}{
		{
			data:     data,
			expected: result,
		},
	}

	t.Run("test_default", func(t *testing.T) {
		for id, v := range testData {
			result, _ := SortByColumn(v.data, 1, true)
			if result != v.expected {
				t.Errorf("Unexpected result at %v %v \n%v", id, result, v.expected)
			}
		}
	})
}

func TestSortByColumnBySecond(t *testing.T) {
	textBytes, err := ioutil.ReadFile("test.txt")
	if err != nil {
		return
	}

	rByte, err := ioutil.ReadFile("test_by_second.txt")
	if err != nil {
		return
	}

	data := string(textBytes)
	result := string(rByte)

	testData := []struct {
		data     string
		expected string
	}{
		{
			data:     data,
			expected: result,
		},
	}

	t.Run("test_default", func(t *testing.T) {
		for id, v := range testData {
			result, _ := SortByColumn(v.data, 1, false)
			if result != v.expected {
				t.Errorf("Unexpected result at %v %v \n%v", id, result, v.expected)
			}
		}
	})
}
