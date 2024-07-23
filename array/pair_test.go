package array

import (
	"reflect"
	"testing"
)

func TestPairFromArray(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		first    []int
		second   []string
		expected []Pair[int, string]
	}{
		{
			name:     "Empty slices",
			first:    []int{},
			second:   []string{},
			expected: []Pair[int, string]{},
		},
		{
			name:     "Equal length slices",
			first:    []int{1, 2, 3},
			second:   []string{"a", "b", "c"},
			expected: []Pair[int, string]{PairOf(1, "a"), PairOf(2, "b"), PairOf(3, "c")},
		},
		{
			name:     "Different length slices",
			first:    []int{1, 2},
			second:   []string{"a", "b", "c"},
			expected: []Pair[int, string]{PairOf(1, "a"), PairOf(2, "b")},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function under test
			result := PairFromArray(tt.first, tt.second)

			// Check if the result matches the expected value
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, but got %d", len(tt.expected), len(result))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("Expected Pair[%d] to be %v, but got %v", i, tt.expected[i], result[i])
				}
			}
		})
	}
}

func TestPairFirsts(t *testing.T) {
	tests := []struct {
		name string
		data []Pair[int, string]
		want []int
	}{
		{
			name: "Empty slice",
			data: []Pair[int, string]{},
			want: []int{},
		},
		{
			name: "Single element",
			data: []Pair[int, string]{PairOf(1, "one")},
			want: []int{1},
		},
		{
			name: "Multiple elements",
			data: []Pair[int, string]{PairOf(1, "one"), PairOf(2, "two"), PairOf(3, "three")},
			want: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PairFirsts(tt.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PairFirsts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPairSeconds(t *testing.T) {
	tests := []struct {
		name string
		data []Pair[int, string]
		want []string
	}{
		{
			name: "Empty slice",
			data: []Pair[int, string]{},
			want: []string{},
		},
		{
			name: "Single element",
			data: []Pair[int, string]{PairOf(1, "one")},
			want: []string{"one"},
		},
		{
			name: "Multiple elements",
			data: []Pair[int, string]{PairOf(1, "one"), PairOf(2, "two"), PairOf(3, "three")},
			want: []string{"one", "two", "three"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PairSeconds(tt.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PairSeconds() = %v, want %v", got, tt.want)
			}
		})
	}
}
