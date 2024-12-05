package op_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/frankill/gotools/op"
)

func TestOpZip(t *testing.T) {
	// 测试用例：常规zip操作
	regularInput := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	expectedRegularResult := [][]int{{1, 4, 7}, {2, 5, 8}, {3, 6, 9}}
	if result := op.Zip(regularInput...); !reflect.DeepEqual(result, expectedRegularResult) {
		t.Errorf("ArrayZip failed for regular input: expected %v, got %v", expectedRegularResult, result)
	}

	// 边界条件测试用例：空数组
	emptyInput := [][]int{}
	expectedEmptyResult := [][]int{}
	if result := op.Zip(emptyInput...); !reflect.DeepEqual(result, expectedEmptyResult) {
		t.Errorf("ArrayZip did not handle empty input correctly, expected %v, got %v", expectedEmptyResult, result)
	}

	// 边界条件测试用例：数组中只有一个元素
	singleElementInput := [][]int{{42}}
	expectedSingleElementResult := [][]int{{42}}
	if result := op.Zip(singleElementInput...); !reflect.DeepEqual(result, expectedSingleElementResult) {
		t.Errorf("ArrayZip did not handle single element input correctly, expected %v, got %v", expectedSingleElementResult, result)
	}

	// 边界条件测试用例：数组中元素数量不一致
	unevenLengthsInput := [][]int{{1, 2}, {3, 4}, {5, 6}}
	if result := op.Zip(op.Zip(unevenLengthsInput...)...); !reflect.DeepEqual(result, unevenLengthsInput) {
		t.Errorf("ArrayZip did not handle uneven lengths correctly, expected %v, got %v", unevenLengthsInput, result)
	}
}

func TestOpMap(t *testing.T) {
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
			result := op.Map(tc.fun, tc.arr...)

			// 检查结果是否与预期相同
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}
