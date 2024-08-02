# purrr

## 简介

`gotools` 总结了日常使用的一些数据处理逻辑，旨在为开发者提供一系列泛型高级函数，以简化数据处理流程并增强代码的表达力。本库借鉴了函数式编程理念，通过提供类型灵活的高阶函数，使数据操作更加高效和优雅。

## 特性

- **泛型支持**：利用 Go 1.18+ 泛型，确保函数在多种数据类型间通用。
- **效率** ： 凑合用吧，操作简化就行
 

## 安装
```bash
go get github.com/frankill/gotools
```

## 快速示例

```go

package main

import (
	"fmt"

	"github.com/frankill/gotools/array"
	"github.com/frankill/gotools/iter"
)

func main() {
	data := map[string]int{"a": 1, "b": 2, "c": 3}
	doubled := array.MapApply(func(k string, v int) int { return v * 2 }, data)
	fmt.Println(doubled) // 输出: [2, 4, 6]

	arr1 := []any{1, 2, 3}
	arr2 := []any{"a", "b", "c"}
	combined := array.ArrayMap(func(x ...any) string {
		return fmt.Sprintf("%d-%s", x[0].(int), x[1].(string))
	}, arr1, arr2)
	fmt.Println(combined) // 输出: ["1-a", "2-b", "3-c"]

	arr3 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	index := array.MatchZero([]int{4, 2, 6}, arr3)

	fmt.Println(index)

	input := iter.FromArray(func(x int) []string { return []string{fmt.Sprintf("%d", x)} }, array.ArraySeq(1, 100, 1))

	pipe := iter.NewPipeline[[]string]()

	pipe.SetStart(input)

	pipe.SetEnd(iter.ToCsv("pipe.csv"))

	pipe.AddStep(iter.Map(func(x []string) []string { return append(x, "test") }))

	pipe.AddStep(iter.Filter(func(x []string) bool {
		tmp, _ := strconv.Atoi(x[0])
		return tmp%2 == 0
	}))

	pipe.Compute()

}


```