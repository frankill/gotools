package array

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

// TestAnd 测试 And 函数
func TestAnd(t *testing.T) {
	type test struct {
		name     string
		arr      []bool
		expected bool
	}

	tests := []test{
		{"empty array", []bool{}, false},
		{"all true", []bool{true, true, true}, true},
		{"one false", []bool{true, false, true}, false},
		{"all false", []bool{false, false, false}, false},
		{"mixed", []bool{true, false, true, true}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := And(tt.arr...)
			if result != tt.expected {
				t.Errorf("And() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestOr(t *testing.T) {
	// Test case with no arguments
	t.Run("no arguments", func(t *testing.T) {
		result := Or()
		if result != false {
			t.Errorf("Or() with no arguments should return false, got %v", result)
		}
	})

	// Test case with all false arguments
	t.Run("all false arguments", func(t *testing.T) {
		result := Or(false, false, false)
		if result != false {
			t.Errorf("Or() with all false arguments should return false, got %v", result)
		}
	})

	// Test case with at least one true argument
	t.Run("at least one true argument", func(t *testing.T) {
		result := Or(false, true, false)
		if result != true {
			t.Errorf("Or() with at least one true argument should return true, got %v", result)
		}
	})
}

func TestNot(t *testing.T) {
	tests := []struct {
		name string
		v    bool
		want bool
	}{
		{"NotTrue", true, false},
		{"NotFalse", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Not(tt.v); got != tt.want {
				t.Errorf("Not() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestOperator(t *testing.T) {
// 	// Define test cases
// 	tests := []struct {
// 		name     string
// 		fun      func(x ...int) int
// 		arr      [][]int
// 		expected []int
// 	}{
// 		{
// 			name:     "Test with single element array",
// 			fun:      func(x ...int) int { return x[0] + 1 },
// 			arr:      [][]int{{1, 2, 3}},
// 			expected: []int{2, 3, 4},
// 		},
// 		{
// 			name:     "Test with multiple element arrays",
// 			fun:      func(x ...int) int { return ArraySum(x) },
// 			arr:      [][]int{{1, 2, 3}, {4, 5, 6}},
// 			expected: []int{5, 7, 9},
// 		},
// 	}

// 	// Run test cases
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Call the function under test
// 			result := Operator(tt.fun, tt.arr...)

//				// Check if the result matches the expected value
//				if !reflect.DeepEqual(result, tt.expected) {
//					t.Errorf("Test %s failed. Expected %v, got %v", tt.name, tt.expected, result)
//				}
//			})
//		}
//	}
func TestOperatorMap(t *testing.T) {
	// 定义测试用例
	testCases := []struct {
		name     string
		fun      func(x ...int) string
		arr      [][]int
		expected []string
	}{
		{
			name:     "Test with empty array",
			fun:      func(x ...int) string { return "" },
			arr:      [][]int{},
			expected: []string{},
		},
		{
			name:     "Test with single element array",
			fun:      func(x ...int) string { return strconv.Itoa(x[0]) },
			arr:      [][]int{{1}},
			expected: []string{"1"},
		},
		{
			name:     "Test with multiple element arrays",
			fun:      func(x ...int) string { return strconv.Itoa(x[0]) + strconv.Itoa(x[1]) },
			arr:      [][]int{{1, 2}, {3, 4}, {5, 6}},
			expected: []string{"13", "24"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 调用待测试函数
			result := Operator(tc.fun, tc.arr...)

			// 检查结果是否与预期相同
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}
func TestCartesian(t *testing.T) {
	// 创建测试用例
	testCases := []struct {
		input  [][]int
		output [][]int
	}{
		{
			input:  [][]int{{1, 2}, {3, 4}},
			output: [][]int{{1, 3}, {2, 3}, {1, 4}, {2, 4}},
		},
		// 添加更多的测试用例
	}

	// 遍历测试用例
	for _, tc := range testCases {
		result := Cartesian(tc.input...)
		if !reflect.DeepEqual(result, tc.output) {
			t.Errorf("Expected %v, but got %v", tc.output, result)
		}
	}

	ForEach(func(x ...int) {
		fmt.Println(x)
	}, ArraySeq(1, 3, 1), ArraySeq(2, 10, 1))

}
