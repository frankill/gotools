package operation_test

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/frankill/gotools/array"
	"github.com/frankill/gotools/operation"
)

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
			result := operation.Operator(tc.fun, tc.arr...)

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
		result := operation.Cartesian(tc.input...)
		if !reflect.DeepEqual(result, tc.output) {
			t.Errorf("Expected %v, but got %v", tc.output, result)
		}
	}

	operation.ForEach(func(x ...int) {
		fmt.Println(x)
	}, array.Seq(1, 3, 1), array.Seq(2, 10, 1))

}
