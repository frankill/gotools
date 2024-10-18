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
func TestAConcat(t *testing.T) {
	// Test case 1
	input1 := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	expected1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	result1 := agg.Concat(input1...)
	if len(result1) != len(expected1) || !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test case 1 failed. Expected %v, got %v", expected1, result1)
	}

	// Test case 2
	input2 := [][]string{{"hello", "world"}, {"foo", "bar"}}
	expected2 := []string{"hello", "world", "foo", "bar"}
	result2 := agg.Concat(input2...)
	if len(result2) != len(expected2) || !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Test case 2 failed. Expected %v, got %v", expected2, result2)
	}

	// Test case 3
	input3 := [][]float64{}
	expected3 := []float64{}
	result3 := agg.Concat(input3...)
	if len(result3) != len(expected3) || !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Test case 3 failed. Expected %v, got %v", expected3, result3)
	}

	// Test case 4
	input4 := [][]bool{{true, false}, {true, true}}
	expected4 := []bool{true, false, true, true}
	result4 := agg.Concat(input4...)
	if len(result4) != len(expected4) || !reflect.DeepEqual(result4, expected4) {
		t.Errorf("Test case 4 failed. Expected %v, got %v", expected4, result4)
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
				if x%2 == 0 {
					return true
				}
				return false
			},
			slice:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expected: 10,
		},
		{
			name: "Find maximum odd number",
			fun: func(x int) bool {
				if x%2 != 0 {
					return true
				}
				return false
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
