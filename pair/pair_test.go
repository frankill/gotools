package pair_test

import (
	"reflect"
	"testing"

	"github.com/frankill/gotools/pair"
)

// TestPairFromMap 测试 PairFromMap 函数的正确性。
func TestPairFromMap(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		input    map[int]string
		expected []pair.Pair[int, string]
	}{
		{
			name: "Basic Case",
			input: map[int]string{
				1: "one",
				2: "two",
				3: "three",
			},
			expected: []pair.Pair[int, string]{
				{First: 1, Second: "one"},
				{First: 2, Second: "two"},
				{First: 3, Second: "three"},
			},
		},
		{
			name:     "Empty Map",
			input:    map[int]string{},
			expected: []pair.Pair[int, string]{}, // 空映射返回空切片
		},
		{
			name: "Single Element",
			input: map[int]string{
				42: "answer",
			},
			expected: []pair.Pair[int, string]{
				{First: 42, Second: "answer"},
			},
		},
		{
			name: "Multiple Elements",
			input: map[int]string{
				1: "one",
				2: "two",
				3: "three",
				4: "four",
				5: "five",
			},
			expected: []pair.Pair[int, string]{
				{First: 1, Second: "one"},
				{First: 2, Second: "two"},
				{First: 3, Second: "three"},
				{First: 4, Second: "four"},
				{First: 5, Second: "five"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 使用 PairFromMap 函数
			result := pair.FromMap(tt.input)

			// 比较结果和期望结果
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("For input %v, expected %v, got %v", tt.input, tt.expected, result)
			}
		})
	}
}
func TestPairFromArray(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		first    []int
		second   []string
		expected []pair.Pair[int, string]
	}{
		{
			name:     "Empty slices",
			first:    []int{},
			second:   []string{},
			expected: []pair.Pair[int, string]{},
		},
		{
			name:     "Equal length slices",
			first:    []int{1, 2, 3},
			second:   []string{"a", "b", "c"},
			expected: []pair.Pair[int, string]{pair.Of(1, "a"), pair.Of(2, "b"), pair.Of(3, "c")},
		},
		{
			name:     "Different length slices",
			first:    []int{1, 2},
			second:   []string{"a", "b", "c"},
			expected: []pair.Pair[int, string]{pair.Of(1, "a"), pair.Of(2, "b")},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function under test
			result := pair.FromArray(tt.first, tt.second)

			// Check if the result matches the expected value
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, but got %d", len(tt.expected), len(result))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("Expected pair.Pair[%d] to be %v, but got %v", i, tt.expected[i], result[i])
				}
			}
		})
	}
}

func TestPairFirsts(t *testing.T) {
	tests := []struct {
		name string
		data []pair.Pair[int, string]
		want []int
	}{
		{
			name: "Empty slice",
			data: []pair.Pair[int, string]{},
			want: []int{},
		},
		{
			name: "Single element",
			data: []pair.Pair[int, string]{pair.Of(1, "one")},
			want: []int{1},
		},
		{
			name: "Multiple elements",
			data: []pair.Pair[int, string]{pair.Of(1, "one"), pair.Of(2, "two"), pair.Of(3, "three")},
			want: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pair.Firsts(tt.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pair.PairFirsts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPairSeconds(t *testing.T) {
	tests := []struct {
		name string
		data []pair.Pair[int, string]
		want []string
	}{
		{
			name: "Empty slice",
			data: []pair.Pair[int, string]{},
			want: []string{},
		},
		{
			name: "Single element",
			data: []pair.Pair[int, string]{pair.Of(1, "one")},
			want: []string{"one"},
		},
		{
			name: "Multiple elements",
			data: []pair.Pair[int, string]{pair.Of(1, "one"), pair.Of(2, "two"), pair.Of(3, "three")},
			want: []string{"one", "two", "three"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pair.Seconds(tt.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pair.PairSeconds() = %v, want %v", got, tt.want)
			}
		})
	}
}
