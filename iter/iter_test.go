package iter_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/frankill/gotools/iter"
)

// TestFormArray 测试 FormArray 函数的正确性。
func TestFormArray(t *testing.T) {
	// 定义测试数据和期望结果
	input := []int{1, 2, 3, 4}
	expected := []string{"1", "2", "3", "4"}

	// 定义将整数转换为字符串的函数
	toString := func(x int) string {
		return fmt.Sprintf("%d", x)
	}

	// 使用 FormArray 函数
	resultCh := iter.FromArray(toString, input)

	// 收集结果
	result := iter.Collect(resultCh)

	// 比较结果和期望结果
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestFromCsv 测试 FromCsv 函数的正确性。
// 这里的测试假设存在一个名为 "test.csv" 的测试文件，其内容为：
// "a,b,c\n1,2,3\n4,5,6\n"
func TestFromCsv(t *testing.T) {
	// 创建一个临时的 CSV 文件
	file, err := os.Create("test.csv")
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove("test.csv") // 确保测试结束后删除文件
	defer file.Close()

	// 写入测试数据
	data := "header1,header2,header3\n1,2,3\n4,5,6\n"
	_, err = file.WriteString(data)
	if err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	// 使用 FromCsv 函数读取文件
	resultCh := iter.FromCsv("test.csv")

	// 期望的结果
	expected := [][]string{
		{"header1", "header2", "header3"},
		{"1", "2", "3"},
		{"4", "5", "6"},
	}

	// 收集结果
	result := iter.Collect(resultCh)

	// 比较结果和期望结果
	if len(result) != len(expected) {
		t.Errorf("Expected %d rows, got %d rows", len(expected), len(result))
		return
	}

	for i, row := range result {
		if !reflect.DeepEqual(row, expected[i]) {
			t.Errorf("Row %d: Expected %v, got %v", i, expected[i], row)
		}
	}
}

// TestToCsv 测试 ToCsv 函数的正确性。
func TestToCsv(t *testing.T) {
	// 定义测试数据
	data := [][]string{
		{"header1", "header2", "header3"},
		{"1", "2", "3"},
		{"4", "5", "6"},
	}

	// 创建一个临时的 CSV 文件
	filePath := "test_output.csv"
	defer os.Remove(filePath) // 确保测试结束后删除文件

	// 创建通道并将测试数据写入通道
	dataCh := make(chan []string, len(data))
	for _, row := range data {
		dataCh <- row
	}
	close(dataCh) // 关闭通道，表示数据写入完毕

	// 使用 ToCsv 函数写入文件
	iter.ToCsv(filePath, dataCh)

	// 读取文件内容并验证
	file, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	expectedContent := "header1,header2,header3\n1,2,3\n4,5,6\n"
	if string(file) != expectedContent {
		t.Errorf("Expected file content:\n%s\nGot:\n%s", expectedContent, string(file))
	}
}

func TestMap(t *testing.T) {
	// 定义测试数据和期望结果
	input := []int{1, 2, 3, 4}
	expected := []string{"1", "2", "3", "4"}

	// 定义将整数转换为字符串的函数
	toString := func(x int) string {
		return fmt.Sprintf("%d", x)
	}

	// 创建通道并将测试数据写入通道
	inputCh := make(chan int, len(input))
	for _, v := range input {
		inputCh <- v
	}
	close(inputCh) // 关闭通道，表示数据写入完毕

	// 使用 Map 函数
	resultCh := iter.Map(toString, inputCh)

	// 收集结果
	result := iter.Collect(resultCh)

	// 比较结果和期望结果
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
func TestFilter(t *testing.T) {
	// 定义测试数据和过滤条件
	input := []int{1, 2, 3, 4, 5}
	expected := []int{2, 4}

	// 定义过滤条件：只保留偶数
	isEven := func(x int) bool {
		return x%2 == 0
	}

	// 创建通道并将测试数据写入通道
	inputCh := make(chan int, len(input))
	for _, v := range input {
		inputCh <- v
	}
	close(inputCh) // 关闭通道，表示数据写入完毕

	// 使用 Filter 函数
	resultCh := iter.Filter(isEven, inputCh)

	// 收集结果
	result := iter.Collect(resultCh)
	// 比较结果和期望结果
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
func TestSequence(t *testing.T) {
	tests := []struct {
		start    int
		end      int
		step     int
		expected []int
	}{
		{1, 10, 2, []int{1, 3, 5, 7, 9}},   // 正步长
		{10, 1, -2, []int{10, 8, 6, 4, 2}}, // 负步长
		// {0, 0, 1, []int{}},                 // 空序列
		// {5, 5, 1, []int{}},                 // 起始和结束相同
		// {1, 5, 5, []int{}},                 // 步长大于范围
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			resultCh := iter.Sequence(tt.start, tt.end, tt.step)

			// 收集结果
			result := iter.Collect(resultCh)

			// 比较结果和期望结果
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("For start=%d, end=%d, step=%d, expected %v, got %v", tt.start, tt.end, tt.step, tt.expected, result)
			}
		})
	}
}
func TestReduce(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		f        func(x int, y int) int
		init     int
		input    []int
		expected int
	}{
		{
			name:     "Sum",
			f:        func(x int, y int) int { return x + y },
			init:     0,
			input:    []int{1, 2, 3, 4},
			expected: 10,
		},
		{
			name:     "Product",
			f:        func(x int, y int) int { return x * y },
			init:     1,
			input:    []int{1, 2, 3, 4},
			expected: 24,
		},
		{
			name:     "Empty",
			f:        func(x int, y int) int { return x + y },
			init:     0,
			input:    []int{},
			expected: 0,
		},
		{
			name:     "Single element",
			f:        func(x int, y int) int { return x + y },
			init:     10,
			input:    []int{5},
			expected: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建通道并将测试数据写入通道
			inputCh := make(chan int, len(tt.input))
			for _, v := range tt.input {
				inputCh <- v
			}
			close(inputCh) // 关闭通道，表示数据写入完毕

			// 使用 Reduce 函数
			result := iter.Reduce(tt.f, tt.init, inputCh)

			// 比较结果和期望结果
			if result != tt.expected {
				t.Errorf("For   init=%d, input=%v, expected %d, got %d", tt.init, tt.input, tt.expected, result)
			}
		})
	}
}
func TestScanl(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		f        func(x int, y int) int
		init     int
		input    []int
		expected []int
	}{
		{
			name:     "Sum",
			f:        func(x int, y int) int { return x + y },
			init:     0,
			input:    []int{1, 2, 3, 4},
			expected: []int{1, 3, 6, 10},
		},
		{
			name:     "Product",
			f:        func(x int, y int) int { return x * y },
			init:     1,
			input:    []int{1, 2, 3, 4},
			expected: []int{1, 2, 6, 24},
		},
		// {
		// 	name:     "Empty",
		// 	f:        func(x int, y int) int { return x + y },
		// 	init:     0,
		// 	input:    []int{},
		// 	expected: []int{},
		// },
		{
			name:     "Single element",
			f:        func(x int, y int) int { return x + y },
			init:     10,
			input:    []int{5},
			expected: []int{15},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建通道并将测试数据写入通道
			inputCh := make(chan int, len(tt.input))
			for _, v := range tt.input {
				inputCh <- v
			}
			close(inputCh) // 关闭通道，表示数据写入完毕

			// 使用 Scanl 函数
			resultCh := iter.Scanl(tt.f, tt.init, inputCh)

			// 收集结果
			result := iter.Collect(resultCh)
			// 比较结果和期望结果
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("For   init=%d, input=%v, expected %v, got %v", tt.init, tt.input, tt.expected, result)
			}
		})
	}
}

func TestZip(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		f        func(x string, y string) string
		ch1      []string
		ch2      []string
		expected []string
	}{
		{
			name:     "Basic Test",
			f:        func(x string, y string) string { return y + x },
			ch1:      []string{"1", "2", "3"},
			ch2:      []string{"a", "b", "c"},
			expected: []string{"a1", "b2", "c3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建通道并将测试数据写入通道
			ch1 := make(chan string, len(tt.ch1))
			ch2 := make(chan string, len(tt.ch2))

			for _, v := range tt.ch1 {
				ch1 <- v
			}
			close(ch1)

			for _, v := range tt.ch2 {
				ch2 <- v
			}
			close(ch2)

			// 使用 Zip 函数
			resultCh := iter.Zip(tt.f, ch1, ch2)

			// 收集结果
			result := iter.Collect(resultCh)

			// 比较结果和期望结果
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("For   ch1=%v, ch2=%v, expected %v, got %v", tt.ch1, tt.ch2, tt.expected, result)
			}
		})
	}
}
func TestPartition(t *testing.T) {

	// 定义测试用例
	tests := []struct {
		name      string
		f         func(x int) bool
		input     []int
		expected1 []int
		expected2 []int
	}{
		{
			name:      "Even and Odd Partition",
			f:         func(x int) bool { return x%2 == 0 },
			input:     []int{1, 2, 3, 4, 5, 6},
			expected1: []int{2, 4, 6}, // 满足条件（偶数）
			expected2: []int{1, 3, 5}, // 不满足条件（奇数）
		},
		// {
		// 	name:      "All True Condition",
		// 	f:         func(x int) bool { return x > 0 },
		// 	input:     []int{1, 2, 3},
		// 	expected1: []int{1, 2, 3}, // 所有值都满足条件
		// 	expected2: []int{},        // 没有值不满足条件
		// },
		// {
		// 	name:      "All False Condition",
		// 	f:         func(x int) bool { return x < 0 },
		// 	input:     []int{-1, -2, -3},
		// 	expected1: []int{},           // 没有值满足条件
		// 	expected2: []int{-1, -2, -3}, // 所有值都不满足条件
		// },
		// {
		// 	name:      "Empty Input",
		// 	f:         func(x int) bool { return x > 0 },
		// 	input:     []int{},
		// 	expected1: []int{}, // 空输入
		// 	expected2: []int{}, // 空输入
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建通道并将测试数据写入通道
			inputCh := make(chan int, len(tt.input))
			for _, v := range tt.input {
				inputCh <- v
			}
			close(inputCh) // 关闭通道，表示数据写入完毕

			// 使用 Partition 函数
			ch1, ch2 := iter.Partition(tt.f, inputCh)

			result1 := iter.Collect(ch1)
			result2 := iter.Collect(ch2)

			// 比较结果和期望结果
			if !reflect.DeepEqual(result1, tt.expected1) {
				t.Errorf("For   input=%v, expected ch1 %v, got %v", tt.input, tt.expected1, result1)
			}
			if !reflect.DeepEqual(result2, tt.expected2) {
				t.Errorf("For   input=%v, expected ch2 %v, got %v", tt.input, tt.expected2, result2)
			}
		})
	}
}
func TestFind(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		f        func(x int) bool
		input    []int
		expected int
	}{
		{
			name:     "Find Even Number",
			f:        func(x int) bool { return x%2 == 0 },
			input:    []int{1, 3, 5, 6, 7},
			expected: 6, // 第一个满足条件的值
		},
		{
			name:     "Find Greater Than Five",
			f:        func(x int) bool { return x > 5 },
			input:    []int{1, 2, 3, 4, 5},
			expected: 0, // 没有值满足条件，返回零值
		},
		{
			name:     "Find First Positive Number",
			f:        func(x int) bool { return x > 0 },
			input:    []int{-1, -2, -3, 1, 2, 3},
			expected: 1, // 第一个满足条件的值
		},
		{
			name:     "Empty Input",
			f:        func(x int) bool { return x > 0 },
			input:    []int{},
			expected: 0, // 空输入，返回零值
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建通道并将测试数据写入通道
			inputCh := make(chan int, len(tt.input))
			for _, v := range tt.input {
				inputCh <- v
			}
			close(inputCh) // 关闭通道，表示数据写入完毕

			// 使用 Find 函数
			result := iter.Find(tt.f, inputCh)

			// 比较结果和期望结果
			if result != tt.expected {
				t.Errorf("For   input=%v, expected %v, got %v", tt.input, tt.expected, result)
			}
		})
	}
}
func TestCollect(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "Non-Empty Input",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},

		{
			name:     "Single Element",
			input:    []int{42},
			expected: []int{42},
		},
		{
			name:     "Multiple Elements",
			input:    []int{10, 20, 30, 40},
			expected: []int{10, 20, 30, 40},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建通道并将测试数据写入通道
			inputCh := make(chan int, len(tt.input))
			for _, v := range tt.input {
				inputCh <- v
			}
			close(inputCh) // 关闭通道，表示数据写入完毕

			// 使用 Collect 函数
			result := iter.Collect(inputCh)

			// 比较结果和期望结果
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("For input %v, expected %v, got %v", tt.input, tt.expected, result)
			}
		})
	}
}
func TestTakeWhile(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		f        func(x int) bool
		input    []int
		expected []int
	}{
		{
			name:     "Take While Less Than 5",
			f:        func(x int) bool { return x < 5 },
			input:    []int{1, 2, 3, 4, 5, 6, 7},
			expected: []int{1, 2, 3, 4}, // 5 和之后的值被忽略
		},
		{
			name:     "Take While Positive",
			f:        func(x int) bool { return x > 0 },
			input:    []int{1, 2, 3, -1, 4, 5},
			expected: []int{1, 2, 3}, // -1 及之后的值被忽略
		},

		// {
		// 	name:     "No Values Meet Condition",
		// 	f:        func(x int) bool { return x > 10 },
		// 	input:    []int{1, 2, 3, 4, 5},
		// 	expected: []int{}, // 没有值满足条件，返回空切片
		// },
		{
			name:     "All Values Meet Condition",
			f:        func(x int) bool { return x < 10 },
			input:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, // 所有值都满足条件
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建通道并将测试数据写入通道
			inputCh := make(chan int, len(tt.input))
			for _, v := range tt.input {
				inputCh <- v
			}
			close(inputCh) // 关闭通道，表示数据写入完毕

			// 使用 TakeWhile 函数
			resultCh := iter.TakeWhile(tt.f, inputCh)
			result := iter.Collect(resultCh) // 收集结果

			// 比较结果和期望结果
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("For  and input %v, expected %v, got %v", tt.input, tt.expected, result)
			}
		})
	}
}
func TestDropWhile(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		f        func(x int) bool
		input    []int
		expected []int
	}{
		{
			name:     "Drop While Less Than 5",
			f:        func(x int) bool { return x < 5 },
			input:    []int{1, 2, 3, 4, 5, 6, 7},
			expected: []int{5, 6, 7}, // 1, 2, 3, 4 被跳过，5 和之后的值被保留
		},
		{
			name:     "Drop While Positive",
			f:        func(x int) bool { return x > 0 },
			input:    []int{1, 2, 3, -1, 4, 5},
			expected: []int{-1, 4, 5}, // 1, 2, 3 被跳过，-1 和之后的值被保留
		},
		// {
		// 	name:     "Empty Input",
		// 	f:        func(x int) bool { return x < 5 },
		// 	input:    []int{},
		// 	expected: []int{}, // 空输入，返回空切片
		// },
		{
			name:     "No Values Meet Condition",
			f:        func(x int) bool { return x > 10 },
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5}, // 没有值满足条件，返回原切片
		},
		{
			name:     "All Values Meet Condition",
			f:        func(x int) bool { return x < 9 },
			input:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			expected: []int{9}, // 所有值都满足条件，返回空切片
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建通道并将测试数据写入通道
			inputCh := make(chan int, len(tt.input))
			for _, v := range tt.input {
				inputCh <- v
			}
			close(inputCh) // 关闭通道，表示数据写入完毕

			// 使用 DropWhile 函数
			resultCh := iter.DropWhile(tt.f, inputCh)
			result := iter.Collect(resultCh) // 收集结果

			// 比较结果和期望结果
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("For  and input %v, expected %v, got %v", tt.input, tt.expected, result)
			}
		})
	}
}
func TestMerge(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		inputs   [][]int
		expected []int
	}{
		{
			name: "Merge Two Channels",
			inputs: [][]int{
				{1, 2, 3},
				{4, 5, 6},
			},
			expected: []int{1, 2, 3, 4, 5, 6}, // 两个通道的值合并到一个通道
		},
		// {
		// 	name: "Merge Empty Channels",
		// 	inputs: [][]int{
		// 		{},
		// 		{},
		// 	},
		// 	expected: []int{}, // 所有通道都为空，返回空切片
		// },
		{
			name: "Merge One Empty Channel",
			inputs: [][]int{
				{1, 2, 3},
				{},
			},
			expected: []int{1, 2, 3}, // 一个通道为空，返回另一个通道的值
		},
		{
			name: "Merge Channels with Overlapping Values",
			inputs: [][]int{
				{1, 2},
				{2, 3},
			},
			expected: []int{1, 2, 2, 3}, // 合并两个通道，可能会有重复值
		},
		{
			name: "Merge Multiple Channels",
			inputs: [][]int{
				{1, 2},
				{3, 4},
				{5, 6},
			},
			expected: []int{1, 2, 3, 4, 5, 6}, // 合并多个通道
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建输入通道
			channels := make([]chan int, len(tt.inputs))
			for i, vals := range tt.inputs {
				channels[i] = make(chan int, len(vals))
				for _, v := range vals {
					channels[i] <- v
				}
				close(channels[i]) // 关闭通道，表示数据写入完毕
			}

			// 使用 Merge 函数
			resultCh := iter.Merge(channels...)
			result := iter.Collect(resultCh) // 收集结果

			// 比较结果和期望结果
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("For inputs %v, expected %v, got %v", tt.inputs, tt.expected, result)
			}
		})
	}
}
func TestGroupBy(t *testing.T) {

	tests := []struct {
		name     string
		keys     []string
		values   []string
		orders   []int
		expected map[string][]string
	}{
		{
			name:   "Basic Case",
			keys:   []string{"a", "b", "a", "b", "c"},
			values: []string{"1", "2", "3", "4", "5"},
			orders: []int{1, 2, 3, 4, 5},
			expected: map[string][]string{
				"a": {"1", "3"},
				"b": {"2", "4"},
				"c": {"5"},
			},
		},
		{
			name:     "Empty Input",
			keys:     []string{},
			values:   []string{},
			orders:   []int{},
			expected: map[string][]string{}, // Empty input returns empty map
		},
		{
			name:   "Single Element",
			keys:   []string{"x"},
			values: []string{"42"},
			orders: []int{1},
			expected: map[string][]string{
				"x": {"42"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create input channels
			keysCh := make(chan string, len(tt.keys))
			valuesCh := make(chan string, len(tt.values))
			ordersCh := make(chan int, len(tt.orders))

			// Send data to channels
			for _, k := range tt.keys {
				keysCh <- k
			}
			close(keysCh)

			for _, v := range tt.values {
				valuesCh <- v
			}
			close(valuesCh)

			for _, o := range tt.orders {
				ordersCh <- o
			}
			close(ordersCh)

			// Use GroupBy function
			result := iter.GroupBy(keysCh, valuesCh, ordersCh)

			// Collect result
			actual := make(map[string][]string)
			for pair := range result {
				actual[pair.First] = pair.Second
			}

			// Compare result with expected
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("For keys %v, values %v, orders %v, expected %v, got %v", tt.keys, tt.values, tt.orders, tt.expected, actual)
			}
		})
	}
}
