package main

import (
	"testing"
)

func TestGetMap(t *testing.T) {
	testData := []struct {
		data     []string
		expected map[string][]string
	}{
		{
			data: []string{"TEST", "estt", "test", "rest"},
			expected: map[string][]string{
				"test": {"estt", "test"},
			},
		},
		{
			data: []string{"estt", "test", "rest", "tser", "rset"},
			expected: map[string][]string{
				"estt": {"estt", "test"},
				"rest": {"rest", "rset", "tser"},
			},
		},
		{
			data: []string{"estt", "test", "rest", "tser", "rset", "roma", "mora", "rma"},
			expected: map[string][]string{
				"estt": {"estt", "test"},
				"rest": {"rest", "rset", "tser"},
				"roma": {"mora", "roma"},
			},
		},
	}

	t.Run("test_map", func(t *testing.T) {
		for k, v := range testData {
			result := getMap(v.data)

			if len(result) != len(v.expected) {
				t.Errorf("error at #%v result: unexpected len", k)
				return
			}

			for k1, v1 := range result {
				res, ok := v.expected[k1]
				if !ok {
					t.Errorf("result not searched at #%v", k)
					return
				}
				if len(v1) != len(res) {
					t.Errorf("result len unexpected at #%v", k)
					return
				}

				for kr, vr := range res {
					if vr != res[kr] {
						t.Errorf("unexpected result item at #%v", k)
						return
					}
				}
			}
		}
	})
}
