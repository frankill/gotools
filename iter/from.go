package iter

import (
	"bufio"
	"compress/gzip"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/frankill/gotools/array"
	"github.com/frankill/gotools/db"
	"github.com/frankill/gotools/file"
	"github.com/frankill/gotools/query"
	"github.com/olivere/elastic/v7"
	"github.com/xuri/excelize/v2"
)

// FromGzip 方法从 gzip 文件中读取数据并发送到通道中
// 参数:
//
// path - gzip 文件路径
// skip - 跳过的行数
//
// 返回:
//
// 一个通道，用于接收从 gzip 文件中读取的数据。
func FromGzip(path string) func(skip int) chan string {
	return func(skip int) chan string {
		ch := make(chan string, BufferSize)

		go func() {
			defer close(ch)

			f, err := os.Open(path)
			if err != nil {
				log.Panicln(err)
			}
			defer f.Close()

			gz, err := gzip.NewReader(f)
			if err != nil {
				log.Panicln(err)
			}
			defer gz.Close()

			scanner := bufio.NewScanner(gz)

			for i := 0; i < skip; i++ {
				scanner.Scan()
			}

			for scanner.Scan() {

				ch <- scanner.Text()
			}
		}()

		return ch
	}
}

// FromJson 方法从 JSON 文件中读取数据并发送到通道中
// 参数:
//
// path - 文件路径
// skip - 跳过的行数
//
// 返回:
//
// 一个通道，用于接收从文件中读取的数据。
func FromJson[T any](path string) func(skip int) chan T {

	return func(skip int) chan T {

		ch := make(chan T, BufferSize)

		go func() {
			defer close(ch)

			f, err := os.Open(path)
			if err != nil {
				log.Panicln(err)
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)

			for i := 0; i < skip; i++ {
				scanner.Scan()
			}

			for scanner.Scan() {

				var t T
				err = json.Unmarshal(scanner.Bytes(), &t)
				if err != nil {
					log.Panicln(err)
				}

				ch <- t
			}
		}()

		return ch
	}
}

// FromTxt 方法从文本文件中读取数据并发送到通道中
// 参数:
//
// path - 文件路径
// skip - 跳过的行数
//
// 返回:
//
// 一个通道，用于接收从文件中读取的数据。
func FromTxt(path string) func(skip int) chan string {

	return func(skip int) chan string {

		ch := make(chan string, BufferSize)

		go func() {
			defer close(ch)

			f, err := os.Open(path)
			if err != nil {
				log.Panicln(err)
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)

			for i := 0; i < skip; i++ {
				scanner.Scan()
			}

			for scanner.Scan() {

				ch <- scanner.Text()
			}
		}()

		return ch

	}
}

// FromMysqlQuery 从 MySQL 数据库中执行查询并返回数据通道。
// 参数:
//
//	query: 执行的 SQL 查询语句。
//
// 返回:
//
//	chan interface{}: 查询结果数据通道。
//	error: 错误信息，如果查询失败。
//
// 示例:
//
//	resultChan, err := FromMysqlQuery("SELECT * FROM table WHERE condition")
//	if err != nil {
//	    log.Fatal(err)
//	}
func FromMysqlQuery[T any](con string) func(query *query.SQLBuilder) chan T {

	return func(query *query.SQLBuilder) chan T {

		query_ := query.Build()

		ch := make(chan T, 100)

		db, err := sql.Open("mysql", con)
		if err != nil {
			log.Panicln(err)
		}

		rows, err := db.Query(query_)
		if err != nil {
			log.Panicln(err)
		}

		columns, err := rows.Columns()
		if err != nil {
			log.Panicln(err)
		}

		columnValues := make([]interface{}, len(columns))
		for i := range columnValues {
			columnValues[i] = new(sql.RawBytes)
		}

		go func() {

			defer close(ch)
			defer db.Close()

			for rows.Next() {
				if err := rows.Scan(columnValues...); err != nil {
					fmt.Println("Scan error:", err)
					continue
				}

				instance := new(T)
				v := reflect.ValueOf(instance).Elem()
				fieldMap := make(map[string]reflect.Value)
				t := reflect.TypeOf(instance).Elem()
				for i := 0; i < t.NumField(); i++ {
					field := t.Field(i)
					if tag, ok := field.Tag.Lookup("sql"); ok {
						fieldMap[tag] = v.Field(i)
					}
				}

				for i, column := range columns {
					if field, ok := fieldMap[column]; ok {
						rawBytes := columnValues[i].(*sql.RawBytes)
						if err := convertToGoType(field, *rawBytes); err != nil {
							fmt.Println("Convert error:", err)
							continue
						}
					}
				}

				ch <- *instance
			}
		}()

		return ch
	}
}

// Converts raw database bytes to Go native types, handling NULL values and missing fields.
func convertToGoType(field reflect.Value, value []byte) error {

	// Convert non-NULL values
	switch field.Kind() {
	case reflect.String:
		field.SetString(string(value))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(string(value), 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(string(value), 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uintValue)
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(string(value), 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatValue)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(string(value))
		if err != nil {
			return err
		}
		field.SetBool(boolValue)
	default:
		return fmt.Errorf("unsupported type: %s", field.Kind())
	}
	return nil
}

// FromMysqlStr 从 MySQL 数据库中执行查询并返回数据通道，
// 使用 SQL 查询字符串。
// 参数:
//
//	queryStr: SQL 查询字符串。
//
// 返回:
//
//	chan []string: 查询结果数据通道。

func FromMysqlStr(con string) func(query *query.SQLBuilder) chan []string {

	return func(query *query.SQLBuilder) chan []string {

		query_ := query.Build()
		ch := make(chan []string, 100)

		db, err := sql.Open("mysql", con)
		if err != nil {
			log.Panicln(err)
		}

		rows, err := db.Query(query_)
		if err != nil {
			log.Panicln(err)
		}

		columns, err := rows.Columns()

		if err != nil {
			log.Panicln(err)
		}

		go func() {
			defer close(ch)
			defer db.Close()

			lc := len(columns)
			for rows.Next() {
				row := make([]sql.NullString, lc)
				rowPointers := make([]any, lc)
				for i := range row {
					rowPointers[i] = &row[i]
				}

				err := rows.Scan(rowPointers...)
				if err != nil {
					log.Fatalln(err)
				}

				ch <- array.ArrayMap(func(i ...sql.NullString) string {
					return i[0].String
				}, row)
			}
		}()

		return ch
	}
}

// FromElasticSearch 从 ElasticSearch 中读取数据。
// 参数:
//
//	client: ElasticSearch 客户端。
//	index: ElasticSearch 索引名称。
//	 query: ElasticSearch 查询请求。
//
// 返回:
//
//	chan interface{}: 查询结果数据通道。
func FromElasticSearch[T any](client *elastic.Client) func(index string, query elastic.Query) chan db.ElasticBluk[T] {
	return func(index string, query elastic.Query) chan db.ElasticBluk[T] {

		con := db.NewElasticSearchClient[T](client)

		ch, err := con.QueryAnyIter(index, query)

		if err != nil {
			log.Panicln(err)
		}

		return ch

	}
}

// FromArray 将输入切片 `a` 中的每个元素应用函数 `f`，并返回一个包含结果的通道。
// 参数:
//   - f: 一个函数，接受切片中的元素 `T` 类型，并返回 `U` 类型的结果。
//   - a: 一个 `T` 类型的切片，将对其每个元素应用函数 `f`。
//
// 返回:
//   - 一个 `U` 类型的通道，通道中的值是对切片 `a` 中的每个元素应用函数 `f` 的结果。
func FromArray[T any, U any](f func(x T) U) func(a []T) chan U {

	return func(a []T) chan U {
		ch := make(chan U, BufferSize)

		go func() {
			defer close(ch)
			for _, v := range a {
				ch <- f(v)
			}

		}()
		return ch
	}

}

// FromMap 将一个映射转换为一个通道，每个通道中的元素是映射中的键值对。
// 参数:
//   - m: 一个映射，键类型为 K，值类型为 V。
//
// 返回:
//   - 一个通道，每个通道中的元素是包含一个键值对的切片，类型为 []array.Pair[K, V]。
//
// 函数功能:
//   - 遍历映射 m，将每个键值对包装在一个切片中，然后将这些切片逐个发送到通道 ch 中。
//   - 每个通道中的元素都是一个包含单个键值对的切片。
//   - 当所有键值对都被发送到通道后，关闭通道。
func FromMap[K comparable, V any](m map[K]V) func() chan array.Pair[K, V] {

	return func() chan array.Pair[K, V] {
		ch := make(chan array.Pair[K, V], BufferSize)

		go func() {
			defer close(ch)
			for k, v := range m {
				ch <- array.Pair[K, V]{
					First:  k,
					Second: v,
				}
			}
		}()

		return ch
	}
}

// FromCsv 从指定的 CSV 文件路径读取数据，并将其以切片的形式发送到通道。
// 参数:
//   - path: CSV 文件的路径，字符串类型。
//
// 返回:
//   - 一个通道，通道中的值是读取的 CSV 文件中的每一行数据，每一行数据是一个字符串切片（[]string）。
func FromCsv(path string) func(header bool) chan []string {

	return func(header bool) chan []string {
		ch := make(chan []string, BufferSize)

		go func() {
			defer close(ch)

			f, err := os.Open(path)
			if err != nil {
				log.Panicln(err)
			}
			defer f.Close()

			reader := csv.NewReader(f)

			if header {
				_, err := reader.Read()

				if err == io.EOF {
					return
				}

				if err != nil {
					log.Panicln(err)
				}
			}

			for {

				row, err := reader.Read()

				if err == io.EOF {
					break
				}

				if err != nil {
					break
				}
				ch <- row

			}

		}()
		return ch
	}
}

// FromTable 从指定的表格文件路径读取数据，并将其以切片的形式发送到通道。
// 参数:
//
//   - path: 表格文件的路径，字符串类型。
//   - header: 是否包含表头，布尔类型。
//   - seq: 表格文件的分隔符，字符串类型。
//
// 返回:
//   - 一个通道，通道中的值是读取的表格文件中的每一行数据，每一行数据是一个字符串切片（[]string）。
func FromTable(path string) func(header bool, seq string) chan []string {

	return func(header bool, seq string) chan []string {
		ch := make(chan []string, BufferSize)

		go func() {
			defer close(ch)

			f, err := os.Open(path)
			if err != nil {
				log.Panicln(err)
			}
			defer f.Close()

			reader := file.NewReader(f, seq)

			if header {
				_, err := reader.Read()

				if err == io.EOF {
					return
				}
				if err != nil {
					log.Panicln(err)
				}
			}

			for {
				record, err := reader.Read()

				if err == io.EOF {
					return
				}

				if err != nil {
					return
				}
				ch <- record
			}
		}()
		return ch
	}
}

// FromExcel 从指定的 Excel 文件路径读取数据，并将其以切片的形式发送到通道。
// 参数:
//   - path: Excel 文件的路径，字符串类型。
//   - sheet: 工作表的名称，字符串类型。
//   - header: 是否包含表头，布尔类型。
//
// 返回:
//   - 一个通道，通道中的值是读取的 Excel 文件中的每一行数据，每一行数据是一个字符串切片（[]string）。
func FromExcel(path string) func(sheet string, header bool) chan []string {

	return func(sheet string, header bool) chan []string {
		ch := make(chan []string, BufferSize)

		go func() {
			defer close(ch)

			f, err := excelize.OpenFile(path)
			if err != nil {
				log.Panicln(err)
			}

			if sheet == "" {
				sheet = "Sheet1"
			}

			rows, err := f.Rows(sheet)
			if err != nil {
				log.Panicln(err)
			}

			if header {
				rows.Next()
			}
			for rows.Next() {
				row, err := rows.Columns()
				if err != nil {
					log.Panicln(err)
				}
				ch <- row
			}

		}()
		return ch
	}
}
