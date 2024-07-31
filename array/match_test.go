package array_test

import (
	"reflect"
	"testing"

	"github.com/frankill/gotools/array"
)

func TestMatchOne(t *testing.T) {
	tests := []struct {
		lookup_value []int
		lookup_array []int
		expected     []int
	}{
		{
			lookup_value: []int{3, 5, 2},
			lookup_array: []int{1, 2, 3, 4, 5},
			expected:     []int{2, 4, 1},
		},
		{
			lookup_value: []int{10, 20, 30},
			lookup_array: []int{5, 10, 15, 25, 20},
			expected:     []int{1, 4, -1},
		},
		{
			lookup_value: []int{1, 2, 3},
			lookup_array: []int{3, 2, 1},
			expected:     []int{2, 1, 0},
		},
	}

	for _, test := range tests {
		result := array.MatchOne(test.lookup_value, test.lookup_array)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("MatchOne(%v, %v) = %v, expected %v", test.lookup_value, test.lookup_array, result, test.expected)
		}
	}
}
func TestXlookup(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	lookup_value := []Person{{Name: "Alice", Age: 25}}
	lookup_array := []Person{{Name: "Alice", Age: 25}, {Name: "Bob", Age: 30}, {Name: "Charlie", Age: 35}}
	lookup_result := []int{1, 2, 3}

	expected := []int{1}
	result := array.Xlookup(lookup_value, lookup_array, lookup_result)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
