package array_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/frankill/gotools/array"
)

func TestArrayFromAny(t *testing.T) {
	// 测试不同类型的输入
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{
			name:  "Test int slice",
			input: []int{1, 2, 3},
			want:  []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := array.From(tt.input...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayFromAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test case for the ArrayShif function
func TestArrayShif(t *testing.T) {
	input1 := []int{1, 2, 3, 4, 5, 6}
	output1 := array.Shif(input1, 2)
	expected1 := []int{0, 0, 1, 2, 3, 4}
	if !reflect.DeepEqual(output1, expected1) {
		t.Errorf("Expected %v but got %v", expected1, output1)
	}

	input2 := []int{1, 2, 3, 4, 5, 6}
	output2 := array.Shif(input2, -2)
	expected2 := []int{3, 4, 5, 6, 0, 0}
	if !reflect.DeepEqual(output2, expected2) {
		t.Errorf("Expected %v but got %v", expected2, output2)
	}

	// Add more test cases as needed
}

// Test case for the ArrayRotate function
func TestArrayRotate(t *testing.T) {
	input1 := []int{1, 2, 3, 4, 5, 6}
	output1 := array.Rotate(input1, 2)
	expected1 := []int{3, 4, 5, 6, 1, 2}
	if !reflect.DeepEqual(output1, expected1) {
		t.Errorf("Expected %v but got %v", expected1, output1)
	}

	input2 := []int{1, 2, 3, 4, 5, 6}
	output2 := array.Rotate(input2, -2)
	expected2 := []int{5, 6, 1, 2, 3, 4}
	if !reflect.DeepEqual(output2, expected2) {
		t.Errorf("Expected %v but got %v", expected2, output2)
	}

	// Add more test cases as needed
}

// Test case for the ArrayProduct function
func TestArrayProduct(t *testing.T) {
	input1 := []int{1, 2, 3, 4, 5, 6}
	output1 := array.Product(input1)
	expected1 := float64(720)
	if output1 != expected1 {
		t.Errorf("Expected %v but got %v", expected1, output1)
	}

	input2 := []float32{2.5, 3.5, 1.5}
	output2 := array.Product(input2)
	expected2 := float64(13.125)
	if output2 != expected2 {
		t.Errorf("Expected %v but got %v", expected2, output2)
	}

	input3 := []int{}
	output3 := array.Product(input3)
	expected3 := float64(1.0)
	if output3 != expected3 {
		t.Errorf("Expected %v but got %v", expected3, output3)
	}

	// Add more test cases as needed
}

// Test case for the ArrayAvg function
func TestArrayAvg(t *testing.T) {
	input1 := []int{1, 2, 3, 4, 5}
	output1 := array.Mean(input1)
	expected1 := int(3)
	if output1 != expected1 {
		t.Errorf("Expected %v but got %v", expected1, output1)
	}

	input2 := []float32{1.5, 2.5, 3.5, 4.5, 5.5}
	output2 := array.Mean(input2)
	expected2 := float32(3.5)
	if output2 != expected2 {
		t.Errorf("Expected %v but got %v", expected2, output2)
	}

	// Add more test cases as needed
}

// Test case for the ArraySum function
func TestArraySum(t *testing.T) {
	input1 := []int{1, 2, 3, 4, 5}
	output1 := array.Sum(input1)
	expected1 := int(15)
	if output1 != expected1 {
		t.Errorf("Expected %v but got %v", expected1, output1)
	}

	input2 := []float32{1.5, 2.5, 3.5, 4.5, 5.5}
	output2 := array.Sum(input2)
	expected2 := float32(17.5)
	if output2 != expected2 {
		t.Errorf("Expected %v but got %v", expected2, output2)
	}

	// Add more test cases as needed
}

// Test case for the ArrayMax function
func TestArrayMax(t *testing.T) {
	input1 := []int{1, 5, 3, 7, 2, 6}
	output1 := array.Max(input1)
	expected1 := int(7)
	if output1 != expected1 {
		t.Errorf("Expected %v but got %v", expected1, output1)
	}

	input2 := []float32{1.5, 2.5, 3.5, 4.5, 5.5}
	output2 := array.Max(input2)
	expected2 := float32(5.5)
	if output2 != expected2 {
		t.Errorf("Expected %v but got %v", expected2, output2)
	}

	// Add more test cases as needed
}

// Test case for the ArrayMin function
func TestArrayMin(t *testing.T) {
	input1 := []int{1, 5, 3, 7, 2, 6}
	output1 := array.Min(input1)
	expected1 := int(1)
	if output1 != expected1 {
		t.Errorf("Expected %v but got %v", expected1, output1)
	}

	input2 := []float32{1.5, 2.5, 3.5, 4.5, 5.5}
	output2 := array.Min(input2)
	expected2 := float32(1.5)
	if output2 != expected2 {
		t.Errorf("Expected %v but got %v", expected2, output2)
	}

	// Add more test cases as needed
}

func TestArrayFindLast(t *testing.T) {
	// Test case 1: Empty array
	arr1 := []int{}
	result1 := array.FindLast(func(x int) bool { return false }, arr1)
	expected1 := -1
	if result1 != expected1 {
		t.Errorf("Test case 1 failed. Expected %d, got %d", expected1, result1)
	}

	// Test case 2: Single array
	f := func(x int) bool {
		return x%2 == 0
	}
	result2 := array.FindLast(f, []int{1, 2, 3, 4, 5})
	expected2 := 3
	if result2 != expected2 {
		t.Errorf("Test case 2 failed. Expected %d, got %d", expected2, result2)
	}

	// Test case 3: Multiple arrays
	arr3 := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	f1 := func(x []int) bool {
		return array.Sum(x) > 15
	}
	result3 := array.FindLast(f1, arr3)
	expected3 := 2
	if result3 != expected3 {
		t.Errorf("Test case 3 failed. Expected %d, got %d", expected3, result3)
	}
}

func TestArrayFindFirst(t *testing.T) {
	// Test case 1: Empty array
	arr1 := [][]int{}
	result1 := array.FindFirst(func(x []int) bool { return false }, arr1)
	expected1 := -1
	if result1 != expected1 {
		t.Errorf("Test case 1 failed. Expected %d, got %d", expected1, result1)
	}

	// Test case 2: Single array
	arr2 := []int{1, 2, 3, 4, 5}

	f := func(x int) bool {
		return x%2 == 0
	}
	result2 := array.FindFirst(f, arr2)
	expected2 := 1
	if result2 != expected2 {
		t.Errorf("Test case 2 failed. Expected %d, got %d", expected2, result2)
	}

	// Test case 3: Multiple arrays
	arr3 := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	f1 := func(x []int) bool {
		return array.Sum(x) > 12
	}
	result3 := array.FindFirst(f1, arr3)
	expected3 := 1
	if result3 != expected3 {
		t.Errorf("Test case 3 failed. Expected %d, got %d", expected3, result3)
	}
}

func TestArrayLast(t *testing.T) {
	// Test case 1: Empty array
	arr1 := []int{}
	result1 := array.Last(func(x int) bool { return false }, arr1)
	expected1 := 0
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test case 1 failed. Expected %v, got %v", expected1, result1)
	}

	// Test case 2: Single array
	arr2 := []int{1, 2, 3, 4, 5}
	f := func(x int) bool {
		return x%2 == 0
	}
	result2 := array.Last(f, arr2)
	expected2 := 4
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Test case 2 failed. Expected %v, got %v", expected2, result2)
	}

	// Test case 3: Multiple arrays
	arr3 := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	f1 := func(x []int) bool {
		return array.Sum(x) > 12
	}
	result3 := array.Last(f1, arr3)
	expected3 := []int{7, 8, 9}
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Test case 3 failed. Expected %v, got %v", expected3, result3)
	}
}

func TestArrayFirst(t *testing.T) {
	// Test case 1: Empty array
	arr1 := []int{}
	result1 := array.First(func(x int) bool { return false }, arr1)
	expected1 := 0
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test case 1 failed. Expected %v, got %v", expected1, result1)
	}

	// Test case 2: Single array
	arr2 := []int{1, 2, 3, 4, 5}
	f := func(x int) bool {
		return x%2 == 0
	}
	result2 := array.First(f, arr2)
	expected2 := 2
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Test case 2 failed. Expected %v, got %v", expected2, result2)
	}

	// Test case 3: Multiple arrays
	arr3 := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	f1 := func(x []int) bool {
		return array.Sum(x) > 12
	}
	result3 := array.First(f1, arr3)
	expected3 := []int{4, 5, 6}
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Test case 3 failed. Expected %v, got %v", expected3, result3)
	}
}

func TestArrayAll(t *testing.T) {
	// Test case 1: Empty array
	arr1 := []int{}
	result1 := array.All(func(x int) bool { return false }, arr1)
	expected1 := false
	if result1 != expected1 {
		t.Errorf("Test case 1 failed. Expected %v, got %v", strconv.FormatBool(expected1), strconv.FormatBool(result1))
	}

	// Test case 2: Single array with all elements greater than 0
	arr2 := []int{1, 2, 3, 4, 5}
	result2 := array.All(func(x int) bool {
		return x%2 == 0
	}, arr2)
	expected2 := false
	if result2 != expected2 {
		t.Errorf("Test case 2 failed. Expected %v, got %v", strconv.FormatBool(expected2), strconv.FormatBool(result2))
	}

	// Test case 3: Multiple arrays with all elements greater than 0
	arr3 := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	result3 := array.All(func(x []int) bool {
		return array.Sum(x) > 11
	}, arr3)
	expected3 := false
	if result3 != expected3 {
		t.Errorf("Test case 3 failed. Expected %v, got %v", strconv.FormatBool(expected3), strconv.FormatBool(result3))
	}

}

func TestArrayAny(t *testing.T) {

	// Test case 2: Single array with all elements greater than 0
	arr2 := []int{1, 2, 3, 4, 5}
	result2 := array.Any(func(x int) bool {
		return x%2 == 0
	}, arr2)
	expected2 := true
	if result2 != expected2 {
		t.Errorf("Test case 2 failed. Expected %v, got %v", strconv.FormatBool(expected2), strconv.FormatBool(result2))
	}

	// Test case 3: Multiple arrays with all elements greater than 0
	arr3 := [][]int{{1, 2, 3}, {4, 5, 6}, {4, 8, 9}}
	result3 := array.All(func(x []int) bool {
		return array.Sum(x) > 11
	}, arr3)
	expected3 := false
	if result3 != expected3 {
		t.Errorf("Test case 3 failed. Expected %v, got %v", strconv.FormatBool(expected3), strconv.FormatBool(result3))
	}

}

func TestArrayReverseSplit(t *testing.T) {
	// Test case 1: Empty array
	arr1 := []int{1}
	result1 := array.ReverseSplit(func(x int) bool { return true }, arr1)
	expected1 := [][]int{{1}}
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test case 1 failed. Expected %v, got %v", expected1, result1)
	}

	// Test case 2: Single array with no split points
	arr2 := []int{}
	f := func(x int) bool {
		return x > 2
	}
	result2 := array.ReverseSplit(f, arr2)
	expected2 := [][]int{}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Test case 2 failed. Expected %v, got %v", expected2, result2)
	}

	// Test case 4: Multiple arrays with no split points
	arr4 := []int{1, 2, 3}

	f1 := func(x int) bool {
		return x > 1
	}

	result4 := array.ReverseSplit(f1, arr4)
	expected4 := [][]int{{1, 2}, {3}}
	if !reflect.DeepEqual(result4, expected4) {
		t.Errorf("Test case 4 failed. Expected %v, got %v", expected4, result4)
	}

}

func TestArraySplit(t *testing.T) {
	// Test case 1: Empty input
	arr := []int{}
	expected := [][]int{}
	result := array.Split(func(x int) bool { return true }, arr)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Test case 1 failed. Expected: %v, Got: %v", expected, result)
	}

	// Test case 2: Single-element input
	arr = []int{1}
	expected = [][]int{{1}}
	result = array.Split(func(x int) bool { return true }, arr)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Test case 2 failed. Expected: %v, Got: %v", expected, result)
	}

	arr = []int{1, 2}
	expected = [][]int{{1}, {2}}
	result = array.Split(func(x int) bool { return true }, arr)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Test case 2 failed. Expected: %v, Got: %v", expected, result)
	}
	// Test case 3: Multiple elements, no split
	arr = []int{1, 2, 3}
	expected = [][]int{{1}, {2, 3}}
	result = array.Split(func(x int) bool { return x%2 == 0 }, arr)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Test case 3 failed. Expected: %v, Got: %v", expected, result)
	}

}

func TestArrayReverseFill(t *testing.T) {
	// Test case 1: Empty input
	arr := []int{}
	expected := []int{}
	result := array.ReverseFill(func(x int) bool { return false }, arr)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Test case 1 failed. Expected: %v, Got: %v", expected, result)
	}

	// Test case 2: Single-element input
	arr = []int{1}
	expected = []int{1}
	result = array.ReverseFill(func(x int) bool { return false }, arr)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Test case 2 failed. Expected: %v, Got: %v", expected, result)
	}

	// Test case 3: Multiple elements, no reverse fill
	arr = []int{1, 2, 3}
	expected = []int{2, 2, 3}
	result = array.ReverseFill(func(x int) bool { return x%2 == 0 }, arr)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Test case 3 failed. Expected: %v, Got: %v", expected, result)
	}

}

func alwaysTrue(x int) bool {
	return true
}

// 示例函数，用于测试逻辑，这里简单定义一个总是返回false的函数
func alwaysFalse(x int) bool {
	return false
}

func equalSlices[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestArrayFill(t *testing.T) {
	// 边界条件：空切片
	emptyResult := array.Fill(alwaysTrue, []int{})
	if len(emptyResult) != 0 {
		t.Errorf("Expected empty slice for empty input, got %v", emptyResult)
	}

	// 单元素切片
	singleElementResult := array.Fill(alwaysTrue, []int{1})
	expectedSingle := []int{1}
	if !equalSlices(singleElementResult, expectedSingle) {
		t.Errorf("Expected %v for single element input, got %v", expectedSingle, singleElementResult)
	}

	// 多维切片一般情况
	multiDimensionalInput := []int{1, 2, 3}
	multiDimensionalResult := array.Fill(alwaysTrue, multiDimensionalInput)
	expectedMulti := []int{1, 2, 3}
	if !equalSlices(multiDimensionalResult, expectedMulti) {
		t.Errorf("Expected %v for multi-dimensional input with alwaysTrue, got %v", expectedMulti, multiDimensionalResult)
	}

	multiDimensionalInput = []int{1, 2, 3}
	multiDimensionalResult = array.Fill(func(x int) bool { return x == 3 }, multiDimensionalInput)
	expectedMulti = []int{1, 1, 3}
	if !equalSlices(multiDimensionalResult, expectedMulti) {
		t.Errorf("Expected %v for multi-dimensional input with alwaysTrue, got %v", expectedMulti, multiDimensionalResult)
	}

	// 测试函数逻辑应用
	logicTestInput := []int{1, 2, 3} // 假设falseVal类型为int，但值表示逻辑上的"假"
	logicTestResult := array.Fill(alwaysFalse, logicTestInput)
	expectedLogic := []int{1, 1, 1} // 根据alwaysFalse逻辑，第二个元素应被替换为前一个元素的值
	if !equalSlices(logicTestResult, expectedLogic) {
		t.Errorf("Expected %v for applying custom logic, got %v", expectedLogic, logicTestResult)
	}
}

func TestArrayFilterBoundary(t *testing.T) {
	// 边界条件测试用例：空数组
	emptyArrayFilter := func(x int) bool {
		return true
	}
	emptyInput := []int{}
	expectedEmptyResult := []int{}
	if result := array.Filter(emptyArrayFilter, emptyInput); !reflect.DeepEqual(result, expectedEmptyResult) {
		t.Errorf("ArrayFilter did not handle empty input correctly, expected %v, got %v", expectedEmptyResult, result)
	}

	// 边界条件测试用例：数组中只有一个元素
	singleElementFilter := func(x int) bool {
		return true
	}
	singleElementInput := []int{42}
	expectedSingleElementResult := []int{42}
	if result := array.Filter(singleElementFilter, singleElementInput); !reflect.DeepEqual(result, expectedSingleElementResult) {
		t.Errorf("ArrayFilter did not handle single element input correctly, expected %v, got %v", expectedSingleElementResult, result)
	}

	// 边界条件测试用例：
	allFalseFilter := func(x int) bool {
		return x%2 == 0
	}
	allFalseInput := []int{1, 2, 3}
	expectedAllFalseResult := []int{2}
	if result := array.Filter(allFalseFilter, allFalseInput); !reflect.DeepEqual(result, expectedAllFalseResult) {
		t.Errorf("ArrayFilter did not handle all false condition correctly, expected %v, got %v", expectedAllFalseResult, result)
	}
}

func TestArrayMap(t *testing.T) {
	// 测试用例：常规映射
	regularMap := func(x int) int {
		return x * x
	}
	regularInput := []int{1, 2}
	expectedRegularResult := []int{1, 4}
	if result := array.Map(regularMap, regularInput); !reflect.DeepEqual(result, expectedRegularResult) {
		t.Errorf("ArrayMap failed for regular input: expected %v, got %v", expectedRegularResult, result)
	}

	// 边界条件测试用例：空数组
	emptyArrayMap := func(x int) int {
		return 0
	}
	emptyInput := []int{}
	expectedEmptyResult := []int{}
	if result := array.Map(emptyArrayMap, emptyInput); !reflect.DeepEqual(result, expectedEmptyResult) {
		t.Errorf("ArrayMap did not handle empty input correctly, expected %v, got %v", expectedEmptyResult, result)
	}

	// 边界条件测试用例：数组中只有一个元素
	singleElementMap := func(x int) int {
		return x
	}
	singleElementInput := []int{42}
	expectedSingleElementResult := []int{42}
	if result := array.Map(singleElementMap, singleElementInput); !reflect.DeepEqual(result, expectedSingleElementResult) {
		t.Errorf("ArrayMap did not handle single element input correctly, expected %v, got %v", expectedSingleElementResult, result)
	}

	// 边界条件测试用例：数组中元素数量不一致
	unevenLengthsMap := func(x int) int {
		return x //
	}
	unevenLengthsInput := []int{1, 2}
	expectedUnevenLengthsResult := []int{1, 2} // 只映射第一个元素
	if result := array.Map(unevenLengthsMap, unevenLengthsInput); !reflect.DeepEqual(result, expectedUnevenLengthsResult) {
		t.Errorf("ArrayMap did not handle uneven lengths correctly, expected %v, got %v", expectedUnevenLengthsResult, result)
	}
}

func TestArrayCompact(t *testing.T) {
	// 测试用例：常规去重
	regularInput := []int{1, 2, 2, 3, 3, 3, 4}
	expectedRegularResult := []int{1, 2, 3, 4}
	if result := array.Compact(regularInput); !reflect.DeepEqual(result, expectedRegularResult) {
		t.Errorf("ArrayCompact failed for regular input: expected %v, got %v", expectedRegularResult, result)
	}

	// 边界条件测试用例：空数组
	emptyInput := []int{}
	expectedEmptyResult := []int{}
	if result := array.Compact(emptyInput); !reflect.DeepEqual(result, expectedEmptyResult) {
		t.Errorf("ArrayCompact did not handle empty input correctly, expected %v, got %v", expectedEmptyResult, result)
	}

	// 边界条件测试用例：所有元素相同
	identicalElementsInput := []int{1, 1, 1, 1}
	expectedIdenticalElementsResult := []int{1}
	if result := array.Compact(identicalElementsInput); !reflect.DeepEqual(result, expectedIdenticalElementsResult) {
		t.Errorf("ArrayCompact did not handle identical elements correctly, expected %v, got %v", expectedIdenticalElementsResult, result)
	}

	// 边界条件测试用例：数组中只有一个元素
	singleElementInput := []int{42}
	expectedSingleElementResult := []int{42}
	if result := array.Compact(singleElementInput); !reflect.DeepEqual(result, expectedSingleElementResult) {
		t.Errorf("ArrayCompact did not handle single element input correctly, expected %v, got %v", expectedSingleElementResult, result)
	}
}
func TestArrayReverse(t *testing.T) {
	// 测试用例：常规反转
	regularInput := []int{1, 2, 3, 4, 5}
	expectedRegularResult := []int{5, 4, 3, 2, 1}
	if result := array.Reverse(regularInput); !reflect.DeepEqual(result, expectedRegularResult) {
		t.Errorf("ArrayReverse failed for regular input: expected %v, got %v", expectedRegularResult, result)
	}

	// 边界条件测试用例：空数组
	emptyInput := []int{}
	expectedEmptyResult := []int{}
	if result := array.Reverse(emptyInput); !reflect.DeepEqual(result, expectedEmptyResult) {
		t.Errorf("ArrayReverse did not handle empty input correctly, expected %v, got %v", expectedEmptyResult, result)
	}

	// 边界条件测试用例：只有一个元素的数组
	singleElementInput := []int{42}
	expectedSingleElementResult := []int{42}
	if result := array.Reverse(singleElementInput); !reflect.DeepEqual(result, expectedSingleElementResult) {
		t.Errorf("ArrayReverse did not handle single element input correctly, expected %v, got %v", expectedSingleElementResult, result)
	}

	// 边界条件测试用例：两个元素的数组
	twoElementsInput := []int{1, 2, 3, 4}
	expectedTwoElementsResult := []int{4, 3, 2, 1}
	if result := array.Reverse(twoElementsInput); !reflect.DeepEqual(result, expectedTwoElementsResult) {
		t.Errorf("ArrayReverse failed for two elements input: expected %v, got %v", expectedTwoElementsResult, result)
	}
}

func TestArrayEnumerateDense(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		expected []int
	}{
		{
			name:     "Test with duplicate elements",
			arr:      []int{1, 2, 3, 2, 1, 4, 5, 3},
			expected: []int{0, 1, 2, 1, 0, 5, 6, 2},
		},
		{
			name:     "Test with empty slice",
			arr:      []int{},
			expected: []int{},
		},
		{
			name:     "Test with all elements the same",
			arr:      []int{1, 1, 1, 1},
			expected: []int{0, 0, 0, 0},
		},
		{
			name:     "Test with negative numbers",
			arr:      []int{-1, -2, -3, -2, -1},
			expected: []int{0, 1, 2, 1, 0},
		},
		{
			name:     "Test with zero",
			arr:      []int{0, 1, 0, 2, 0},
			expected: []int{0, 1, 0, 3, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := array.EnumerateDense(tt.arr); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ArrayEnumerateDense() = %v, want %v", got, tt.expected)
			}
		})
	}
}
func TestArraySort(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		expected []int
	}{
		{"Test with multiple elements", []int{5, 3, 8, 6, 2}, []int{2, 3, 5, 6, 8}},
		{"Test with single element", []int{1}, []int{1}},
		{"Test with empty slice", []int{}, []int{}},
		{"Test with all elements the same", []int{1, 1, 1, 1}, []int{1, 1, 1, 1}},
		{"Test with negative numbers", []int{-1, -3, -2, -5, -4}, []int{-5, -4, -3, -2, -1}},
		{"Test with zero and negative numbers", []int{0, -1, -2, 0, 1}, []int{-2, -1, 0, 0, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortFunc := func(x, y int) bool {
				return x < y
			}
			if result := array.SortFun(sortFunc, tt.arr); !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ArraySort() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestArrayDifference(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		expected []int
	}{
		{"Test with positive numbers", []int{5, 2, 9, 1}, []int{5, -3, 7, -8}},
		{"Test with a single element", []int{42}, []int{42}},
		{"Test with negative numbers", []int{-5, -8, -1, -10}, []int{-5, -3, 7, -9}},
		{"Test with mixed positive and negative numbers", []int{-3, 7, -1, 4}, []int{-3, 10, -8, 5}},
		{"Test with zero", []int{0, 5, 0, -5}, []int{0, 5, -5, -5}},
		{"Test with large numbers", []int{100, 200, 300, 400}, []int{100, 100, 100, 100}},
		{"Test with decreasing order", []int{10, 5, 0, -5}, []int{10, -5, -5, -5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := array.Difference(tt.arr)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ArrayDifference() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestArrayCount(t *testing.T) {
	tests := []struct {
		name     string
		fun      func(int) bool
		arrs     []int
		expected int
	}{
		{"Test with all matching pairs", func(x int) bool { return x%2 == 0 }, []int{1, 2, 2, 3}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := array.CountIf(tt.fun, tt.arrs)
			if result != tt.expected {
				t.Errorf("ArrayCount() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestArrayHasAll(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		elems    []int
		expected bool
	}{
		{"1", []int{1, 2, 3}, []int{1, 2, 3}, true},
		{"2", []int{1, 2, 3}, []int{1, 2, 4}, false},
		{"3", []int{}, []int{}, false},
		{"4", []int{}, []int{1}, false},
		{"5", []int{1, 2, 3}, []int{}, false},
		{"6", []int{1, 2, 3, 4}, []int{1, 2, 3}, true},
		{"7", []int{-1, -2, -3}, []int{-1, -2}, true},
		{"8", []int{0, 1, 2}, []int{0}, true},
		{"9", []int{1, 2, 2, 3}, []int{2}, false},
		{"10", []int{1, 2, 3}, []int{2, 2}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := array.HasAll(tt.arr, tt.elems...)
			if result != tt.expected {
				t.Errorf("%s, ArrayHasAll() = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestArrayHasAny(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		elems    []int
		expected bool
	}{
		{"Test with any element present", []int{1, 2, 3}, []int{3, 4, 5}, true},
		{"Test with no element present", []int{1, 2, 3}, []int{4, 5, 6}, false},
		{"Test with empty arr", []int{}, []int{1, 2, 3}, false},
		{"Test with empty elems", []int{1, 2, 3}, []int{}, false},
		{"Test with arr containing one of elems", []int{1, 2, 3}, []int{2}, true},
		{"Test with arr and elems containing duplicates", []int{1, 2, 2, 3}, []int{2, 2}, true},
		{"Test with negative numbers", []int{-1, -2, -3}, []int{-2, -4, -5}, true},
		{"Test with zero", []int{0, 1, 2}, []int{0}, true},
		{"Test with all elements in arr also in elems", []int{1, 2, 3}, []int{1, 2, 3}, true},
		{"Test with large numbers", []int{100, 200, 300}, []int{400, 500, 300}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := array.HasAny(tt.arr, tt.elems...)
			if result != tt.expected {
				t.Errorf("ArrayHasAny() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestArrayHasSequence(t *testing.T) {
	// Test case 1: arr2 is a subsequence of arr1
	arr1 := []int{1, 2, 3, 4, 5, 6}
	arr2 := []int{3, 4, 5}
	expected := true
	expectedIndex := 4
	result, index := array.HasSequence(arr1, arr2)
	if result != expected || index != expectedIndex {
		t.Errorf("Test case 1 failed. Expected (%v, %v), got (%v, %v)", expected, expectedIndex, result, index)
	}

	// Test case 2: arr2 is not a subsequence of arr1
	arr1 = []int{1, 2, 3, 4, 5, 6}
	arr2 = []int{7, 8, 9}
	expected = false
	expectedIndex = 0
	result, index = array.HasSequence(arr1, arr2)
	if result != expected || index != expectedIndex {
		t.Errorf("Test case 2 failed. Expected (%v, %v), got (%v, %v)", expected, expectedIndex, result, index)
	}

	// Test case 3: arr2 is an empty sequence
	arr1 = []int{1, 2, 3, 4, 5, 6}
	arr2 = []int{}
	expected = false
	expectedIndex = 0
	result, index = array.HasSequence(arr1, arr2)
	if result != expected || index != expectedIndex {
		t.Errorf("Test case 3 failed. Expected (%v, %v), got (%v, %v)", expected, expectedIndex, result, index)
	}

	// Add more test cases as needed
}

func TestArraySequenceCount(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		arr1     []int
		arr2     []int
		expected int
	}{
		{
			name:     "Sequence found once",
			arr1:     []int{1, 2, 3, 4, 3, 2, 1},
			arr2:     []int{3, 2, 1},
			expected: 1,
		},
		{
			name:     "Sequence found multiple times",
			arr1:     []int{1, 2, 3, 4, 3, 2, 1, 3, 2, 1},
			arr2:     []int{3, 2, 1},
			expected: 2,
		},
		{
			name:     "Sequence not found",
			arr1:     []int{1, 2, 3, 4, 5},
			arr2:     []int{6, 7, 8},
			expected: 0,
		},
		// 可以根据需要添加更多测试用例
	}

	// 遍历测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 调用待测试的函数
			count := array.ArrSequenceCount(tt.arr1, tt.arr2)

			// 验证结果是否符合预期
			if count != tt.expected {
				t.Errorf("ArraySequenceCount(%v, %v) = %d; expected %d", tt.arr1, tt.arr2, count, tt.expected)
			}
		})
	}
}

func TestArrayExtractByIndex(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		index    []int
		expected []int
	}{
		{
			name:     "Test with valid index",
			arr:      []int{5, 2, 7, 3, 9},
			index:    []int{3, 1, 4, 2, 0},
			expected: []int{3, 2, 9, 7, 5},
		},
		{
			name:     "Test with index containing -1",
			arr:      []int{5, 2, 7, 3, 9},
			index:    []int{3, -1, 4, 2, 0},
			expected: []int{3, 0, 9, 7, 5},
		},
		{
			name:     "Test with different length index",
			arr:      []int{5, 2, 7, 3, 9},
			index:    []int{3, 1, 4},
			expected: []int{3, 2, 9},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := array.Choose(test.index, test.arr)
			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("Test %s failed. Expected %v, got %v", test.name, test.expected, actual)
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
			output: [][]int{{1, 3}, {1, 4}, {2, 3}, {2, 4}},
		},
		// 添加更多的测试用例
	}

	// 遍历测试用例
	for _, tc := range testCases {
		result := array.Cartesian(tc.input...)
		if !reflect.DeepEqual(result, tc.output) {
			t.Errorf("Expected %v, but got %v", tc.output, result)
		}
	}

}
func TestArrayReduce(t *testing.T) {
	// 测试空输入
	fun := func(x int, y int) int { return x + y }
	arr := []int{}
	result := array.Reduce(fun, 0, arr)
	expected := 0
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}

	// 测试非空输入
	fun = func(x int, y int) int { return x + y }
	arr = []int{1, 2, 3, 4, 5}
	result = array.Reduce(fun, 0, arr)
	expected = 15
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
func TestArrayReduceR(t *testing.T) {
	// 测试空输入
	fun := func(x int, y int) int { return x + y }
	arr := []int{}
	result := array.ReduceR(fun, 0, arr)
	expected := 0
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}

	// 测试非空输入
	fun = func(x int, y int) int { return x + y }
	arr = []int{1, 2, 3, 4, 5}
	result = array.ReduceR(fun, 0, arr)
	expected = 15
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}

	// 测试非交换运算
	fun = func(x int, y int) int { return x - y }
	arr = []int{1, 2, 3, 4, 5}
	result = array.ReduceR(fun, 0, arr)
	expected = -15
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestArrayUniqueCount(t *testing.T) {
	// 测试空输入
	arr := []int{}
	result := array.UniqueCount(arr)
	expected := 0
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}

	// 测试单个元素
	arr = []int{1}
	result = array.UniqueCount(arr)
	expected = 1
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}

	// 测试多个元素
	arr = []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	result = array.UniqueCount(arr)
	expected = 4
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}

	// 测试所有元素相同
	arr = []int{1, 1, 1, 1, 1}
	result = array.UniqueCount(arr)
	expected = 1
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
func TestArrayDistinctCount(t *testing.T) {
	// 测试空输入
	arr := []int{}
	result := array.DistinctCount(arr)
	expected := 0
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}

	// 测试单个元素
	arr = []int{1}
	result = array.DistinctCount(arr)
	expected = 1
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}

	// 测试多个元素
	arr = []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	result = array.DistinctCount(arr)
	expected = 4
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}

	// 测试所有元素相同
	arr = []int{1, 1, 1, 1, 1}
	result = array.DistinctCount(arr)
	expected = 1
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}

}
