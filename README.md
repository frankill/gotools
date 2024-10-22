# gotools

## 简介

`gotools` 提供数据的流处理iter, array的隐式循环操作，etl的读取  

## 特性

- **泛型支持**：利用 Go 1.18+ 泛型，确保函数在多种数据类型间通用。
- **效率** ： 凑合用吧，操作简化就行
- **数据源** ： txt,csv,table,json, msyql , elasticsearch,clickhouse,gob,gzip

## 安装
```bash
go get github.com/frankill/gotools
```

## 快速示例

```go

 package main

import (
	"fmt"
	"strconv"

	"github.com/frankill/gotools/iter"
)

func compare(x, y []string) bool {
	xn, _ := strconv.Atoi(x[0])
	yn, _ := strconv.Atoi(y[0])
	return xn < yn
}

func main() {

	data := []struct {
		Number int
		Label  string
	}{
		{2, "test"},
		{4, "test"},
		{6, "test"},
		{6, "test1"},
		{10, "test"},
		{12, "test"},
		{13, "test"},
	}

	iter.ToCsv("pipe.csv", false)(
		iter.FromArray2(func(x struct {
			Number int
			Label  string
		}) []string {
			return []string{strconv.Itoa(x.Number), x.Label}
		})(data),
	)

	data = []struct {
		Number int
		Label  string
	}{
		{1, "test"},
		{6, "test"},
		{10, "test1"},
		{14, "test13"},
		{15, "test1"},
		{17, "1"},
	}

	iter.ToCsv("pipe_1.csv", false)(
		iter.FromArray2(func(x struct {
			Number int
			Label  string
		}) []string {
			return []string{strconv.Itoa(x.Number), x.Label}
		})(data),
	)

	ch1, _ := iter.FromCsv("pipe.csv")(false)
	ch2, _ := iter.FromCsv("pipe_1.csv")(false)

	iter.Walk(func(x []string) { fmt.Println(x) })(
		iter.Unique(func(x, y []string) bool {
			return x[0] == y[0]
		})(
			iter.Merge(compare)(
				iter.SortS(compare)(ch2),
				iter.SortS(compare)(ch1),
			),
		),
	)
}





```
 

 