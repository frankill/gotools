package maps_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/frankill/gotools/maps"
	"github.com/frankill/gotools/pair"
)

func TestMapKeys(t *testing.T) {
	// 定义测试用例
	testCases := []struct {
		name     string
		inputMap map[int]string
		wantKeys []int
	}{
		{
			name:     "empty map",
			inputMap: map[int]string{},
			wantKeys: []int{},
		},
		{
			name: "map with keys",
			inputMap: map[int]string{
				1: "one",
				2: "two",
				3: "three",
			},
			wantKeys: []int{1, 2, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 调用函数
			gotKeys := maps.Keys(tc.inputMap)

			// 检查结果是否正确
			if len(gotKeys) != len(tc.wantKeys) {
				t.Errorf("MapKeys() = %v, want %v", gotKeys, tc.wantKeys)
			}
			for i, key := range gotKeys {
				if key != tc.wantKeys[i] {
					t.Errorf("MapKeys() = %v, want %v", gotKeys, tc.wantKeys)
					break
				}
			}
		})
	}
}

func TestMpaValues(t *testing.T) {
	// Test case 1
	m1 := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	expected1 := []int{1, 2, 3}
	result1 := maps.Values(m1)
	if len(result1) != len(expected1) {
		t.Errorf("Expected %v, got %v", expected1, result1)
	}

	// Test case 2
	m2 := map[float64]bool{
		1.5: true,
		2.5: false,
		3.5: true,
	}
	expected2 := []bool{true, false, true}
	result2 := maps.Values(m2)
	if len(result2) != len(expected2) {
		t.Errorf("Expected %v, got %v", expected2, result2)
	}

	// Test case 3
	m3 := map[complex128]string{}
	expected3 := []string{}
	result3 := maps.Values(m3)
	if len(result3) != len(expected3) {
		t.Errorf("Expected %v, got %v", expected3, result3)
	}
}
func TestMapFilter(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		filter   func(int, string) bool
		inputMap map[int]string
		wantMap  map[int]string
	}{
		{
			name: "filter out values containing 'a'",
			filter: func(key int, value string) bool {
				return !contains('a', value)
			},
			inputMap: map[int]string{1: "apple", 2: "banana", 3: "cherry"},
			wantMap:  map[int]string{3: "cherry"},
		},
		{
			name: "filter out keys greater than 2",
			filter: func(key int, value string) bool {
				return key <= 2
			},
			inputMap: map[int]string{1: "apple", 2: "banana", 3: "cherry"},
			wantMap:  map[int]string{1: "apple", 2: "banana"},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function under test
			gotMap := maps.Filter(tt.filter, tt.inputMap)

			// Check if the result is as expected
			if len(gotMap) != len(tt.wantMap) {
				t.Errorf("MapFilter() = %v, want %v", gotMap, tt.wantMap)
			}
			for k, v := range tt.wantMap {
				if gotVal, ok := gotMap[k]; !ok || gotVal != v {
					t.Errorf("MapFilter() = %v, want %v", gotMap, tt.wantMap)
					break
				}
			}
		})
	}
}

// Helper function to check if a string contains a character
func contains(ch rune, str string) bool {
	for _, r := range str {
		if r == ch {
			return true
		}
	}
	return false
}

func TestMapApplyValue(t *testing.T) {
	// Define the test cases
	tests := []struct {
		name     string
		f        func(int, string) bool
		m        map[int]string
		expected map[int]bool
	}{
		{
			name: "Test 1",
			f: func(k int, v string) bool {
				return k%2 == 0 && v == "even"
			},
			m: map[int]string{
				1: "odd",
				2: "even",
				3: "odd",
				4: "even",
			},
			expected: map[int]bool{
				1: false,
				2: true,
				3: false,
				4: true,
			},
		},
		{
			name: "Test 2",
			f: func(k int, v string) bool {
				return k > 2 && v == "odd"
			},
			m: map[int]string{
				1: "odd",
				2: "even",
				3: "odd",
				4: "even",
			},
			expected: map[int]bool{
				1: false,
				2: false,
				3: true,
				4: false,
			},
		},
	}

	// Run the tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := maps.ApplyValue(tt.f, tt.m)
			if len(actual) != len(tt.expected) {
				t.Errorf("Expected length %d, but got %d", len(tt.expected), len(actual))
			}
			for k, v := range tt.expected {
				if actual[k] != v {
					t.Errorf("Expected value %v for key %d, but got %v", v, k, actual[k])
				}
			}
		})
	}
}

func TestMapApplyKey(t *testing.T) {
	// Define the test cases
	tests := []struct {
		name     string
		f        func(int, string) string
		m        map[int]string
		expected map[string]string
	}{
		{
			name: "Test 1",
			f: func(k int, v string) string {
				return strconv.Itoa(k) + v
			},
			m: map[int]string{
				1: "a",
				2: "b",
				3: "c",
			},
			expected: map[string]string{
				"1a": "a",
				"2b": "b",
				"3c": "c",
			},
		},
		{
			name: "Test 2",
			f: func(k int, v string) string {
				return strconv.FormatBool(k%2 == 0)
			},
			m: map[int]string{
				1: "a",
				2: "b",
				3: "c",
			},
			expected: map[string]string{
				"false": "c",
				"true":  "b",
			},
		},
	}

	// Run the tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := maps.ApplyKey(tt.f, tt.m)
			if !mapsEqual(actual, tt.expected) {
				t.Errorf("Expected %v, but got %v", tt.expected, actual)
			}
		})
	}
}

// Helper function to check if two maps are equal
func mapsEqual[K, V comparable](m1, m2 map[K]V) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k, v := range m1 {
		if m2[k] != v {
			return false
		}
	}

	return true
}

func TestMapApplyBoth(t *testing.T) {
	tests := []struct {
		name string
		f    func(int, string) (float64, bool)
		m    map[int]string
		want map[float64]bool
	}{
		{
			name: "Test 1",
			f:    func(k int, v string) (float64, bool) { return float64(k), k%2 == 0 },
			m:    map[int]string{1: "1", 2: "2", 3: "3"},
			want: map[float64]bool{1: false, 2: true, 3: false},
		},
		{
			name: "Test 2",
			f:    func(k int, v string) (float64, bool) { return float64(len(v)), k > 2 },
			m:    map[int]string{1: "1", 2: "22", 3: "333"},
			want: map[float64]bool{1: false, 2: false, 3: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maps.ApplyBoth(tt.f, tt.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapApplyBoth() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMapFromArrayWithFun(t *testing.T) {
	// 定义测试用例
	testCases := []struct {
		name     string
		fun      func(int) string
		key      []int
		value    []string
		expected map[string]string
	}{
		{
			name: "test1",
			fun: func(x int) string {
				return strconv.Itoa(x * 2)
			},
			key:      []int{1, 2, 3},
			value:    []string{"a", "b", "c"},
			expected: map[string]string{"2": "a", "4": "b", "6": "c"},
		},
		{
			name: "test2",
			fun: func(x int) string {
				return strconv.Itoa(x * 2)
			},
			key:      []int{4, 5},
			value:    []string{"d", "e"},
			expected: map[string]string{"8": "d", "10": "e"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 调用被测试函数
			result := maps.FromArrayWithFun(tc.fun, tc.key, tc.value)

			// 检查结果是否符合预期
			if len(result) != len(tc.expected) {
				t.Errorf("expected length %d, got %d", len(tc.expected), len(result))
			}
			for k, v := range tc.expected {
				if result[k] != v {
					t.Errorf("expected %v, got %v for key %s", v, result[k], k)
				}
			}
		})
	}
}
func TestMapPopulateSeries(t *testing.T) {
	tests := []struct {
		name   string
		key    []int
		value  []string
		max    int
		expect map[int]string
	}{
		{
			name:   "Test with positive values",
			key:    []int{1, 2, 3},
			value:  []string{"a", "b", "c"},
			max:    5,
			expect: map[int]string{1: "a", 2: "b", 3: "c", 4: "", 5: ""},
		},
		{
			name:   "Test with negative values",
			key:    []int{-1, -2, -3},
			value:  []string{"x", "y", "z"},
			max:    -5,
			expect: map[int]string{-1: "x", -2: "y", -3: "z"},
		},
		// Add more test cases if needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maps.PopulateSeries(tt.key, tt.value, tt.max)
			if !reflect.DeepEqual(result, tt.expect) {
				t.Errorf("Test %s failed. Expected %v, got %v", tt.name, tt.expect, result)
			}
		})
	}
}
func TestMapContains(t *testing.T) {
	// 定义测试用例
	testCases := []struct {
		name     string
		m        map[string]int
		key      []string
		expected []int
	}{
		{
			name:     "Empty map",
			m:        map[string]int{},
			key:      []string{"a", "b", "c"},
			expected: []int{-1, -1, -1},
		},
		{
			name:     "Existing keys",
			m:        map[string]int{"a": 1, "b": 2, "c": 3},
			key:      []string{"a", "b", "c"},
			expected: []int{1, 2, 3},
		},
		{
			name:     "Non-existing keys",
			m:        map[string]int{"a": 1, "b": 2},
			key:      []string{"c", "d"},
			expected: []int{-1, -1},
		},
		{
			name:     "Mixed keys",
			m:        map[string]int{"a": 1, "b": 2, "c": 3},
			key:      []string{"a", "d", "b", "e"},
			expected: []int{1, -1, 2, -1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 调用被测试函数
			result := maps.Contains(tc.m, -1, tc.key...)

			// 检查结果是否符合预期
			if len(result) != len(tc.expected) {
				t.Errorf("Expected length %d, but got %d", len(tc.expected), len(result))
			}
			for i := range result {
				if result[i] != tc.expected[i] {
					t.Errorf("Expected value %d, but got %d at index %d", tc.expected[i], result[i], i)
				}
			}
		})
	}
}
func TestMapConcat(t *testing.T) {
	// Test case 1: Test with empty maps
	m1 := []map[int]string{}
	expected := map[int]string{}
	result := maps.Concat(m1...)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("TestMapConcat failed. Expected %v, got %v", expected, result)
	}

	// Test case 2: Test with single map
	m2 := []map[int]string{{1: "a", 2: "b"}}
	expected = map[int]string{1: "a", 2: "b"}
	result = maps.Concat(m2...)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("TestMapConcat failed. Expected %v, got %v", expected, result)
	}

	// Test case 3: Test with multiple maps
	m3 := []map[int]string{{1: "a"}, {2: "b"}, {1: "c"}}
	expected = map[int]string{1: "c", 2: "b"}
	result = maps.Concat(m3...)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("TestMapConcat failed. Expected %v, got %v", expected, result)
	}
}

func TestMapExists(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		f        func(int, string) bool
		m        map[int]string
		expected bool
	}{
		{
			name: "Test element exists",
			f: func(k int, v string) bool {
				return k == 1 && v == "value1"
			},
			m: map[int]string{
				1: "value1",
				2: "value2",
			},
			expected: true,
		},
		{
			name: "Test element does not exist",
			f: func(k int, v string) bool {
				return k == 3 && v == "value3"
			},
			m: map[int]string{
				1: "value1",
				2: "value2",
			},
			expected: false,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maps.Exists(tt.f, tt.m)
			if result != tt.expected {
				t.Errorf("MapExists() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestMapAll(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		f        func(int, string) bool
		m        map[int]string
		expected bool
	}{
		{
			name: "Test 1",
			f: func(k int, v string) bool {
				return k > 0 && v != ""
			},
			m:        map[int]string{1: "one", 2: "two", 3: "three"},
			expected: true,
		},
		{
			name: "Test 2",
			f: func(k int, v string) bool {
				return k%2 == 0 && v != ""
			},
			m:        map[int]string{1: "one", 2: "two", 3: "three"},
			expected: false,
		},
		{
			name: "Test 3",
			f: func(k int, v string) bool {
				return v == "four"
			},
			m:        map[int]string{1: "one", 2: "two", 3: "three"},
			expected: false,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function under test
			result := maps.All(tt.f, tt.m)

			// Check if the result matches the expected value
			if result != tt.expected {
				t.Errorf("MapAll() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
func TestMapToPairsArray2(t *testing.T) {
	testData := map[string][]int{
		"one":   {1, 2},
		"two":   {3, 4},
		"three": {5},
	}
	expectedPairs := pair.Pair[[]string, []int]{
		First:  []string{"one", "one", "two", "two", "three"},
		Second: []int{1, 2, 3, 4, 5},
	}

	result := maps.ToPairsArray2(testData)
	if !reflect.DeepEqual(result, expectedPairs) {
		t.Errorf("MapToPairsArray2(%v) = %v; want %v", testData, result, expectedPairs)
	}
}
func TestMapToPairsArray(t *testing.T) {
	testData := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	expectedPairs := pair.Pair[[]string, []int]{
		First:  []string{"one", "two", "three"},
		Second: []int{1, 2, 3},
	}

	result := maps.ToPairsArray(testData)
	if !reflect.DeepEqual(result, expectedPairs) {
		t.Errorf("MapToPairsArray(%v) = %v; want %v", testData, result, expectedPairs)
	}
}
func TestMapMerge(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 3, "c": 4}
	m3 := map[string]int{"d": 5}
	expectedMerged := map[string]int{"a": 1, "b": 2, "c": 4, "d": 5}

	result := maps.Merge(m1, m2, m3)
	if !reflect.DeepEqual(result, expectedMerged) {
		t.Errorf("MapMerge(%v, %v) = %v; want %v", m1, m2, result, expectedMerged)
	}
}
func TestMapIntersect(t *testing.T) {
	data1 := map[string][]int{"a": {1, 2, 3}, "b": {5, 6, 7, 8, 9}, "c": {4}}
	data2 := map[string][]int{"a": {3, 4, 5}, "b": {7}, "c": {8, 9, 10}}

	result := maps.Intersect(data1, data2)

	expectedFirst := [][]string{{"a", "a"}, {"b", "a"}, {"c", "a"}, {"b", "b"}, {"b", "c"}}
	expectedSecond := [][]int{{3}, {5}, {4}, {7}, {8, 9}}

	if !reflect.DeepEqual(result.First, expectedFirst) {
		t.Errorf("First arrays do not match. Expected %v, got %v", expectedFirst, result.First)
	}

	if !reflect.DeepEqual(result.Second, expectedSecond) {
		t.Errorf("Second arrays do not match. Expected %v, got %v", expectedSecond, result.Second)
	}
}
