# gotools

## 简介

`gotools` 提供数据的流处理iter, array,map的隐式循环操作   

## 特性

- **泛型支持**：利用 Go 1.18+ 泛型，确保函数在多种数据类型间通用
- **数据源** ： txt,csv,table,json,msyql,elasticsearch,clickhouse,gob,gzip

## 安装
```bash

	go get github.com/frankill/gotools

```

## 快速示例

```go

package main

import (
	"fmt"

	"github.com/frankill/gotools"
	"github.com/frankill/gotools/iter"
	"github.com/frankill/gotools/operation"
	"github.com/frankill/gotools/query"
)

func main() {

	defer gotools.Clear(1)

	MysqlTest := ""

	q1 := query.NewSQLBuilder().From("test")

	q2 := query.NewSQLBuilder().SQL("select * from test ")

	d1, _ := iter.FromMysql[user](MysqlTest)(q1)

	d2, _ := iter.FromMysql[user](MysqlTest)(q2)

	res := operation.Eq(iter.Collect(d1), iter.Collect(d2))

	fmt.Println(res)
}

type user struct {
	ID    string `mysql:"id"`
	Phone string `mysql:"phone"`
	Name  string `mysql:"name"`
}


```
 

 