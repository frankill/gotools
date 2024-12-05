package agg_test

import (
	"reflect"
	"testing"

	"github.com/frankill/gotools/agg"
)

func TestDistinct(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		slice    [][]int
		expected int
	}{
		{
			name:     "Test with single slice",
			slice:    [][]int{{1, 2, 3, 4, 5, 3}, {4, 5, 6}},
			expected: 6,
		},
		{
			name:     "Test with multiple slices",
			slice:    [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			expected: 9,
		},
		{
			name:     "Test with empty slices",
			slice:    [][]int{{}, {}, {}},
			expected: 0,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := agg.Distinct(tt.slice...)
			if result != tt.expected {
				t.Errorf("Expected %d, but got %d", tt.expected, result)
			}
		})
	}
}

func TestASum(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		slice    [][]int
		expected int
	}{
		{
			name:     "Test with single slice",
			slice:    [][]int{{1, 2, 3, 4, 5, 3}, {4, 5, 6}},
			expected: 33,
		},
		{
			name:     "Test with multiple slices",
			slice:    [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			expected: 45,
		},
		{
			name:     "Test with empty slices",
			slice:    [][]int{{}, {}, {}},
			expected: 0,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := agg.Sum(tt.slice...)
			if result != tt.expected {
				t.Errorf("Expected %d, but got %d", tt.expected, result)
			}
		})
	}
}

func TestAMin(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		slice    [][]int
		expected int
	}{
		{
			name:     "Test with single element slice",
			slice:    [][]int{{1}},
			expected: 1,
		},
		{
			name:     "Test with multiple element slices",
			slice:    [][]int{{9, 5, 2}, {7, 1, 6}, {4, 8, 3, 2, 3, 4, 4}},
			expected: 1,
		},
		{
			name:     "Test with empty slices",
			slice:    [][]int{{}, {6, 3, 7}, {2}},
			expected: 0, // Assuming 0 is the default value for an empty slice
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := agg.Min(test.slice...)
			if result != test.expected {
				t.Errorf("Expected %d, but got %d", test.expected, result)
			}
		})
	}
}

func TestAMax(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		slice    [][]int
		expected int
	}{
		{
			name:     "Test with single slice",
			slice:    [][]int{{1, 2, 3}, {4, 5, 6, 2, 3, 4}, {7, 8, 9}},
			expected: 9,
		},
		{
			name:     "Test with empty slices",
			slice:    [][]int{{-1}, {}, {}},
			expected: 0,
		},
		{
			name:     "Test with negative numbers",
			slice:    [][]int{{-1, -2, 0}, {-4, -5, -6}, {-7, -8, -9}},
			expected: 0,
		},
		// Add more test cases as needed
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := agg.Max(tt.slice...)
			if result != tt.expected {
				t.Errorf("AMax(%v) = %v, expected %v", tt.slice, result, tt.expected)
			}
		})
	}
}

func TestAMaxif(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		fun      func(x int) bool
		slice    []int
		expected int
	}{
		{
			name: "Find maximum even number",
			fun: func(x int) bool {
				return x%2 == 0
			},
			slice:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expected: 10,
		},
		{
			name: "Find maximum odd number",
			fun: func(x int) bool {
				return x%2 != 0
			},
			slice:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expected: 9,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := agg.Maxif(tt.fun, tt.slice)
			if result != tt.expected {
				t.Errorf("Expected %d, but got %d", tt.expected, result)
			}
		})
	}
}

func TestAMinif(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		fun      func(x int) bool
		slice    [][]int
		expected int
	}{
		{
			name: "Find minimum value in a slice of positive integers",
			fun: func(x int) bool {
				return x > 0
			},
			slice:    [][]int{{5, 3, 8, 1, 7}, {4, 5, 6}},
			expected: 1,
		},
		{
			name: "Find minimum value in a slice of negative integers",
			fun: func(x int) bool {
				return x < 0
			},
			slice:    [][]int{{-5, -3, -8, -1, -7}, {-10, -18}},
			expected: -18,
		},
		{
			name: "Find minimum value in an empty slice",
			fun: func(x int) bool {
				return x > 0
			},
			slice:    [][]int{},
			expected: 0,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := agg.Minif(tt.fun, tt.slice...)
			if result != tt.expected {
				t.Errorf("AMinif() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestAargMax(t *testing.T) {
	tests := []struct {
		arg    []int
		val    []int
		expect int
	}{
		{[]int{1, 2, 3}, []int{10, 20, 30}, 3},
		{[]int{5, 2, 7}, []int{5, 2, 7}, 7},
		{[]int{1, 2, 3}, []int{}, 1},
		{[]int{}, []int{10, 20, 30}, 0},
		{[]int{1, 2, 3}, []int{10, 20, 30, 40}, 3},
	}

	for _, tt := range tests {
		result := agg.AargMax(tt.arg, tt.val)
		if result != tt.expect {
			t.Errorf("AargMax(%v, %v) = %v, expect %v", tt.arg, tt.val, result, tt.expect)
		}
	}
}

func TestAargMin(t *testing.T) {
	tests := []struct {
		arg    []int
		val    []int
		expect int
	}{
		{[]int{1, 2, 3}, []int{4, 3, 2}, 3},
		{[]int{5, 2, 8}, []int{1, 2, 3}, 5},
		{[]int{1, 2, 3}, []int{}, 1},
		{[]int{}, []int{1, 2, 3}, 0},
		{[]int{1, 2, 3}, []int{1, 2, 3}, 1},
	}

	for _, test := range tests {
		result := agg.AargMin(test.arg, test.val)
		if !reflect.DeepEqual(result, test.expect) {
			t.Errorf("AargMin(%v, %v) = %v, expect %v", test.arg, test.val, result, test.expect)
		}
	}
}

func TestAccumulateFun(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		fun      func(x, y int) int
		arr      [][]int
		expected int
		init     int
	}{
		{
			name:     "sum of elements",
			fun:      func(x, y int) int { return x + y },
			arr:      [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			expected: 45,
			init:     0,
		},
		{
			name:     "product of elements",
			fun:      func(x, y int) int { return x * y },
			arr:      [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			expected: 362880,
			init:     1,
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := agg.Asccumulate(tt.fun, tt.init, tt.arr...)
			if result != tt.expected {
				t.Errorf("AccumulateFun() = %v, want %v", result, tt.expected)
			}
		})
	}
}
