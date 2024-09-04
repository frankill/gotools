package iter

import (
	"bufio"
	"compress/gzip"
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/gob"
	"encoding/json"
	"errors"
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

func ErrorCH(ch chan error) {

	for err := range ch {
		log.Println(err)
	}
}

// FromGzip 方法从 gzip 文件中读取数据并发送到通道中
// 参数:
//
// path - gzip 文件路径
// skip - 跳过的行数
//
// 返回:
//
// 一个通道，用于接收从 gzip 文件中读取的数据。
// 一个通道，用于接收错误信息
func FromGzip(path string) func(skip int) (chan string, chan error) {
	return func(skip int) (chan string, chan error) {
		ch := make(chan string, BufferSize)
		errs := make(chan error, 1)

		go func() {
			defer close(ch)
			defer close(errs)

			f, err := os.Open(path)
			if err != nil {
				errs <- err
				return
			}
			defer f.Close()

			gz, err := gzip.NewReader(f)
			if err != nil {
				errs <- err
				return
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

		return ch, errs
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
// 一个通道，用于接收错误信息
func FromJson[T any](path string) func(skip int) (chan T, chan error) {

	return func(skip int) (chan T, chan error) {

		ch := make(chan T, BufferSize)
		errs := make(chan error, 1)

		go func() {
			defer close(ch)
			defer close(errs)

			f, err := os.Open(path)
			if err != nil {
				errs <- err
				return
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
					errs <- err
					return
				}

				ch <- t
			}
		}()

		return ch, errs
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
// error: 错误信息
func FromTxt(path string) func(skip int) (chan string, chan error) {

	return func(skip int) (chan string, chan error) {

		ch := make(chan string, BufferSize)
		errs := make(chan error, 1)

		go func() {
			defer close(ch)
			defer close(errs)

			f, err := os.Open(path)
			if err != nil {
				errs <- err
				return
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

		return ch, errs

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
//	chan db.ElasticBluk[T]: 数据通道。
//	chan error: 错误通道。
func FromElasticSearch[T any](client *elastic.Client) func(index string, query any) (chan db.ElasticBluk[T], chan error) {
	return func(index string, query any) (chan db.ElasticBluk[T], chan error) {

		con := db.NewElasticSearchClient[T](client)

		return con.QueryAnyIter(index, query)

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
//   - 一个通道，通道中的值是读取 CSV 文件时发生的错误。
func FromCsv(path string) func(header bool) (chan []string, chan error) {

	return func(header bool) (chan []string, chan error) {
		ch := make(chan []string, BufferSize)
		errs := make(chan error, 1)

		go func() {
			defer close(ch)
			defer close(errs)

			f, err := os.Open(path)
			if err != nil {
				errs <- err
				return
			}
			defer f.Close()

			reader := csv.NewReader(f)

			if header {
				_, err := reader.Read()

				if err == io.EOF {
					return
				}

				if err != nil {
					errs <- err
					return
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
		return ch, errs
	}
}

// FromTable 从指定的文本文件路径读取数据，并将其以切片的形式发送到通道。
// 参数:
//
//   - path: 文本文件的路径，字符串类型。
//   - header: 是否包含表头，布尔类型。
//   - seq: 文本文件的分隔符，字符串类型。
//
// 返回:
//   - 一个通道，通道中的值是读取的表格文件中的每一行数据，每一行数据是一个字符串切片（[]string）。
//   - 一个通道，通道中的值是读取表格文件时发生的错误。
func FromTable(path string) func(header bool, seq string, escape byte) (chan []string, chan error) {

	return func(header bool, seq string, escape byte) (chan []string, chan error) {
		ch := make(chan []string, BufferSize)
		errs := make(chan error, 1)

		go func() {
			defer close(ch)
			defer close(errs)

			f, err := os.Open(path)
			if err != nil {
				errs <- err
				return
			}
			defer f.Close()

			reader := file.NewReader(f, seq, escape)

			if header {
				_, err := reader.Read()

				if err == io.EOF {
					return
				}
				if err != nil {
					errs <- err
					return
				}
			}

			for {
				record, err := reader.Read()

				if err == io.EOF {
					return
				}

				if err != nil {
					errs <- err
					return
				}
				ch <- record
			}
		}()
		return ch, errs
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
//   - 一个通道，通道中的值是读取的 Excel 文件中的每一行数据，每一行数据是一个字符串切片（[]string）。
func FromExcel(path string) func(sheet string, header bool) (chan []string, chan error) {

	return func(sheet string, header bool) (chan []string, chan error) {
		ch := make(chan []string, BufferSize)
		errs := make(chan error, 1)

		go func() {
			defer close(ch)
			defer close(errs)

			f, err := excelize.OpenFile(path)
			if err != nil {
				errs <- err
				return
			}

			if sheet == "" {
				sheet = "Sheet1"
			}

			rows, err := f.Rows(sheet)
			if err != nil {
				errs <- err
				return
			}

			if header {
				rows.Next()
			}
			for rows.Next() {
				row, err := rows.Columns()
				if err != nil {
					errs <- err
					return
				}
				ch <- row
			}

		}()
		return ch, errs
	}
}

type ftype struct {
	t reflect.Value
	n int
}

// FromMysqlQuery 从 MySQL 数据库中执行查询并返回数据通道。
// 参数:
//
//	query: *query.SQLBuilder - 查询语句
//
// 返回:
//
//	chan T: 查询结果数据通道
//	error: 错误信息，如果查询失败。
func FromMysql[T any](con string) func(query *query.SQLBuilder) (chan T, chan error) {

	return func(query *query.SQLBuilder) (chan T, chan error) {

		query_ := query.Build()

		ch := make(chan T, BufferSize)
		errs := make(chan error, 1)

		go func() {

			defer close(ch)
			defer close(errs)

			con, err := sql.Open("mysql", con)
			if err != nil {
				errs <- err
				return
			}
			defer con.Close()

			rows, err := con.Query(query_)
			if err != nil {
				errs <- err
				return
			}

			columns, err := rows.Columns()
			if err != nil {
				errs <- err
				return
			}

			columnValues := make([]interface{}, len(columns))
			for i := range columnValues {
				columnValues[i] = new(sql.RawBytes)
			}

			instance := new(T)
			v := reflect.ValueOf(instance).Elem()
			fieldMap := make(map[string]ftype, len(columns))
			t := reflect.TypeOf(instance).Elem()
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				if tag, ok := field.Tag.Lookup("mysql"); ok {
					var tmp ftype
					tmp.t = v.Field(i)
					tmp.n = fieldType(tmp.t)
					fieldMap[tag] = tmp
				}
			}

			for rows.Next() {
				if err := rows.Scan(columnValues...); err != nil {
					errs <- err
					return
				}

				for i, column := range columns {
					if field, ok := fieldMap[column]; ok {
						rawBytes := columnValues[i].(*sql.RawBytes)
						if err := convertToGoType(field, *rawBytes); err != nil {
							errs <- err
							return
						}
					}
				}

				ch <- *instance
			}
		}()

		return ch, errs
	}
}

func FromCK[T any](ck *db.CKinfo) func(query *query.SQLBuilder) (chan T, chan error) {

	return func(query *query.SQLBuilder) (chan T, chan error) {

		query_ := query.Build()

		ch := make(chan T, BufferSize)
		errs := make(chan error, 1)

		go func() {

			defer close(ch)
			defer close(errs)

			con, err := db.NewCKLoc(ck)
			if err != nil {
				errs <- err
				return
			}
			defer con.Close()

			rows, err := con.Query(context.Background(), query_)

			if err != nil {
				errs <- err
				return
			}
			for rows.Next() {
				var instance T
				if err := rows.ScanStruct(&instance); err != nil {
					errs <- err
					return
				}
				ch <- instance
			}

		}()

		return ch, errs
	}
}

func fieldType(field reflect.Value) int {

	switch field.Kind() {
	case reflect.String:
		return 1
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return 2
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return 3
	case reflect.Float32, reflect.Float64:
		return 4
	case reflect.Bool:
		return 5
	default:
		return 0
	}
}

// Converts raw database bytes to Go native types, handling NULL values and missing fields.
func convertToGoType(field ftype, value []byte) error {

	switch field.n {
	case 1:
		field.t.SetString(string(value))
	case 2:
		intValue, err := strconv.ParseInt(string(value), 10, 64)
		if err != nil {
			return err
		}
		field.t.SetInt(intValue)
	case 3:
		intValue, err := strconv.ParseUint(string(value), 10, 64)
		if err != nil {
			return err
		}
		field.t.SetUint(intValue)
	case 4:
		intValue, err := strconv.ParseFloat(string(value), 64)
		if err != nil {
			return err
		}
		field.t.SetFloat(intValue)
	case 5:
		intValue, err := strconv.ParseBool(string(value))
		if err != nil {
			return err
		}
		field.t.SetBool(intValue)
	default:
		return errors.New("unknown field type")
	}

	return nil
}

// FromMysqlStr 从 MySQL 数据库中执行查询并返回数据通道，
// 使用 SQL 查询字符串。
// 参数:
//
//	query - *SQLBuilder 类型的结构体，用于构建 SQL 查询语句。
//
// 返回:
//
//	chan []string: 查询结果数据通道。
//	error: 错误信息，如果查询失败。
func FromMysqlStr(con string) func(query *query.SQLBuilder) (chan []string, chan error) {

	return func(query *query.SQLBuilder) (chan []string, chan error) {

		return db.NewMysqlDB(con).QueryAnyIter(query)()

	}
}

func FromCKStr(ck *db.CKinfo) func(query *query.SQLBuilder) (chan []string, chan error) {

	return func(query *query.SQLBuilder) (chan []string, chan error) {

		return db.NewCK(ck).QueryAnyIter(query)()
	}
}

// FromGob 从指定的 gob 文件路径读取数据，并通过通道返回。
func FromGob[T any](path string, del bool) (chan T, chan error) {
	ch := make(chan T, BufferSize)
	errs := make(chan error, 1) // 只有一个错误通道缓冲区

	file, err := os.Open(path)
	if err != nil {
		errs <- err
		close(errs)
		return ch, errs
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)

	go func() {
		defer func() {
			if del {
				os.Remove(path)
			}
		}()
		defer close(ch)
		defer close(errs)
		for {
			var instance T
			if err := decoder.Decode(&instance); err != nil {
				if err.Error() == "EOF" {
					// 文件结束，正常退出
					return
				}
				errs <- err
				return
			}
			ch <- instance
		}
	}()

	return ch, errs
}
