package operation_test

import (
	"reflect"
	"strconv"
	"testing"

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
			result := operation.Pmap(tc.fun, tc.arr...)

			// 检查结果是否与预期相同
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}
