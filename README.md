# gotools

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
	
	iter.Walk(
		func(x int) {
			println(x)
		},
	)(
		iter.Filter(func(x int) bool {
			for i := 5; i*i <= x; i += 6 {
				if x%i == 0 || x%(i+2) == 0 {
					return false
				}
			}
			return true
		})(
			iter.Filter(func(x int) bool {
				if x <= 1 {
					return false
				}
				if x <= 3 {
					return true
				}
				if x%3 == 0 {
					return false
				}
				if x%2 == 0 {
					return false
				}
				return true
			})(
				iter.FromArray(Identity[int])(array.ArraySeq(1, 100, 1)),
			),
		),
	)
 

}


```

```go

package main

import (
	"fmt"
	"strconv"

	"github.com/frankill/gotools/array"
	"github.com/frankill/gotools/iter"
)

func main() {

	input := iter.FromArray(func(x int) []string { return []string{fmt.Sprintf("%d", x)} })(array.ArraySeq(1, 100, 1))

	pipe := iter.NewPipeline[[]string]()

	pipe.SetStart(func() chan []string { return input })

	pipe.SetEnd(iter.ToCsv("pipe.csv"))

	pipe.AddStep(iter.Map(func(x []string) []string { return append(x, "test") }))

	pipe.AddStep(iter.Filter(func(x []string) bool {
		tmp, _ := strconv.Atoi(x[0])
		return tmp%2 == 0
	}))

	pipe.Run()

}


```