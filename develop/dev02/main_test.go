package main

import (
	"testing"
)

func TestUnpackString(t *testing.T) {
	testData := []struct {
		data     string
		expected string
	}{
		{
			"qwe\\45\\35test3e",
			"qwe4444433333testtte",
		},
		{
			"a4bc2d5e",
			"aaaabccddddde",
		},
		{
			"abcd",
			"abcd",
		},
		{
			"abcd",
			"abcd",
		},
		{
			"qwe\\4\\5",
			"qwe45",
		},
		{
			"qwe\\45",
			"qwe44444",
		},
		{
			"qwe\\\\5",
			"qwe\\\\\\\\\\",
		},
	}

	t.Run("testing result", func(t *testing.T) {
		for id, v := range testData {
			res, _ := UnpackString(v.data)
			if res != v.expected {
				t.Errorf("error on #%v test: unexpected result %v (%v)", id, v.expected, res)
			}
		}
	})

	testErrorData := []struct {
		data          string
		expectedError bool
	}{
		{
			"45",
			true,
		},
		{
			"\\45",
			false,
		},
		{
			"4\\5",
			true,
		},
		{
			"\\\\4\\5",
			false,
		},
	}

	t.Run("testing errors", func(t *testing.T) {
		for id, v := range testErrorData {
			_, err := UnpackString(v.data)
			if (err == nil) == v.expectedError {
				t.Errorf("error on #%v test: unexpected error", id)
			}
		}
	})
}
