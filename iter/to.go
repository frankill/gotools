package iter

import (
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/frankill/gotools/array"
	"github.com/frankill/gotools/db"
	"github.com/frankill/gotools/file"
	"github.com/frankill/gotools/query"
	"github.com/olivere/elastic/v7"
	"github.com/xuri/excelize/v2"
)

// ToGzip 方法将通道中的数据写入 gzip 文件
// 参数:
//
// path - gzip 文件路径
// append - 是否追加到文件末尾（true）还是覆盖文件（false）
//
// 返回:
//
// 一个函数，接受一个通道作为参数，写入通道中的数据到 gzip 文件中，并返回错误信息。
func ToGzip(path string, append bool) func(ch chan string) error {
	return func(ch chan string) error {
		var f *os.File
		var err error

		if append {
			f, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		} else {
			f, err = os.Create(path)
		}
		if err != nil {
			return err
		}
		defer f.Close()

		gz, err := gzip.NewWriterLevel(f, gzip.BestSpeed)
		if err != nil {
			return err
		}
		defer gz.Close()

		for t := range ch {
			_, err = gz.Write([]byte(t + "\n"))
			if err != nil {
				return err
			}
		}

		return nil
	}
}

// ToJson 方法将通道中的数据写入到 JSON 文件中
// 参数:
//
// path - 文件路径
// append - 是否追加写入文件
// ch - 一个通道，用于接收待写入的数据。
//
// 返回:
//
// 一个函数，用于执行 JSON 文件写入操作。
func ToJson[T any](path string, append bool) func(ch chan T) error {

	return func(ch chan T) error {

		if path == "" {
			return errors.New("path is empty")
		}

		var f *os.File
		var err error

		if append {
			f, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		} else {
			f, err = os.Create(path)
		}
		if err != nil {
			return err
		}
		defer f.Close()

		// 创建 JSON 编码器
		encoder := json.NewEncoder(f)

		for t := range ch {

			if err := encoder.Encode(t); err != nil {
				return err
			}
		}

		return nil
	}

}

// ToTxt 方法将通道中的数据写入到文本文件中
// 参数:
//
// path - 文件路径
// append - 是否追加写入文件
// ch - 一个通道，用于接收待写入的数据。
//
// 返回:
//
// 一个函数，用于执行文件写入操作。
func ToTxt(path string, append bool) func(ch chan string) error {

	return func(ch chan string) error {

		if path == "" {
			return errors.New("path is empty")
		}

		var f *os.File
		var err error

		if append {
			f, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		} else {
			f, err = os.Create(path)
		}

		if err != nil {
			return err
		}
		defer f.Close()

		for t := range ch {

			if _, err := f.WriteString(t + "\n"); err != nil {
				return err
			}
		}

		return nil
	}

}

// ToMysqlInset 方法将通道中的数据插入到 MySQL 数据库中
// 参数:
//
// con mysql 连接字符串
//
//	q - *SQLBuilder 类型的结构体，用于构建 SQL 查询语句。
//	ch - 一个通道，用于接收待插入的数据。
//
// 返回:
//
// 一个函数， 用于执行数据库插入操作。
func ToMysqlInset(con string, q query.SqlInsert) func(ch chan []string) error {

	return func(ch chan []string) error {

		con := db.NewMysqlDB(con)

		defer con.Close()

		err := con.Insert(q)(ch)
		if err != nil {
			return err
		}

		return nil
	}

}

// ToCKInsert 方法将通道中的数据插入到 ClickHouse 数据库中
// 参数:
//
// ck *db.CKinfo - *db.CKinfo 类型的结构体，用于构建 ck 客户端信息
//
//	q - *SQLBuilder 类型的结构体，用于构建 SQL 查询语句。
//	ch - 一个通道，用于接收待插入的数据。
//
// 返回:
//
// 一个函数， 用于执行数据库插入操作。
func ToCK(ck *db.CKinfo, q query.SqlInsert) func(ch chan []string) error {

	return func(ch chan []string) error {

		con := db.NewCK(ck)

		defer con.Close()

		err := con.Insert(q)(ch)
		if err != nil {
			return err
		}

		return nil
	}

}

// ToElasticSearch 将数据写入 ElasticSearch。
// 参数:
//
// client: ElasticSearch 客户端。
// index: ElasticSearch 索引名称。
// ctype: ElasticSearch 写入类型 index create
//
//	data: 需要写入的数据。
//
// 返回:
//
//	error: 错误信息，如果写入失败。
func ToElasticSearch[T any](client *elastic.Client) func(ch chan db.ElasticBluk[T]) error {
	return func(ch chan db.ElasticBluk[T]) error {

		con := db.NewElasticSearchClient[T](client)

		err := con.BulkInsert()(ch)
		if err != nil {
			return err
		}

		return nil
	}
}

// ToCsv 将通道中的数据写入指定的 CSV 文件。
// 参数:
//   - path: CSV 文件的路径，字符串类型。
//   - ch: 一个通道，通道中的每个值是一个字符串切片（[]string），表示 CSV 文件中的一行数据。
//
// 函数功能:
//   - 从通道中读取数据，并将数据写入指定的 CSV 文件。
//   - 使用一个 goroutine 执行写入操作，并通过 stop 通道同步写入完成。
func ToCsv(path string, append bool, header ...string) func(ch chan []string) error {

	return func(ch chan []string) error {

		var f *os.File
		var err error

		if append {
			f, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
		} else {
			f, err = os.Create(path)
			if err != nil {
				return err
			}
		}

		defer f.Close()

		writer := csv.NewWriter(f)
		defer writer.Flush()

		if len(header) > 0 {
			if err := writer.Write(header); err != nil {
				return err
			}
		}

		for row := range ch {
			if err := writer.Write(row); err != nil {
				return err
			}
		}

		return nil
	}

}

type TableField struct {
	Path     string
	Seq      string
	UseQuote bool
	Header   []string
	Append   bool
	Escape   byte
}

func T(path string) *TableField {

	return &TableField{
		Path:     path,
		Seq:      ",",
		UseQuote: true,
		Header:   []string{},
		Append:   false,
		Escape:   '"',
	}
}

func (t *TableField) SetEscape(escape byte) *TableField {
	t.Escape = escape

	return t
}

func (t *TableField) SetHeader(header ...string) *TableField {
	t.Header = header

	return t
}

func (t *TableField) SetSeq(seq string) *TableField {
	t.Seq = seq

	return t
}

func (t *TableField) SetUseQuote(useQuote bool) *TableField {
	t.UseQuote = useQuote

	return t
}

func (t *TableField) SetAppend(append bool) *TableField {
	t.Append = append

	return t
}

func (t *TableField) SetPath(path string) *TableField {
	t.Path = path

	return t
}

// ToTable 将通道中的数据写入指定分隔符 文件。
// 参数:
//
// t *TableField: 表格配置
// ch: 一个通道，通道中的每个值是一个字符串切片（[]string），表示表格文件中的一行数据。
//
// 返回:
//
//	error
func ToTable(t *TableField) func(ch chan []string) error {

	return func(ch chan []string) error {

		var f *os.File
		var err error

		if t.Append {
			f, err = os.OpenFile(t.Path, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
		} else {
			f, err = os.Create(t.Path)
			if err != nil {
				return err
			}
		}
		defer f.Close()

		if t.Seq == "" {
			return errors.New("seq cannot be empty")
		}

		writer := file.NewWriter(f, t.Seq, t.UseQuote, t.Escape)

		defer writer.Flush()

		if len(t.Header) > 0 {
			if err := writer.Write(t.Header); err != nil {
				return err
			}
		}

		for record := range ch {
			err := writer.Write(record)

			if err != nil {
				return err
			}
		}
		return nil
	}

}

type ExcelField struct {
	Path   string
	Sheet  string
	Header []string
	Exist  bool
	Append bool
}

func (e *ExcelField) SetHeader(header ...string) *ExcelField {
	e.Header = header

	return e
}

func (e *ExcelField) SetSheet(sheet string) *ExcelField {

	e.Sheet = sheet

	return e
}

func (e *ExcelField) SetPath(path string) *ExcelField {

	e.Path = path

	return e
}

func (e *ExcelField) SetExist(append bool) *ExcelField {

	e.Exist = append

	return e
}

func (e *ExcelField) SetAppend(append bool) *ExcelField {

	e.Append = append

	return e
}

func E(path string) *ExcelField {

	return &ExcelField{
		Path:   path,
		Exist:  false,
		Sheet:  "",
		Header: []string{},
		Append: false,
	}
}

func (e *ExcelField) getRow(f *excelize.File, sheet string) (int, error) {
	rows, err := f.GetRows(sheet)

	return len(rows) + 1, err
}
func ToExcel(e *ExcelField) func(ch chan []string) error {

	return func(ch chan []string) error {

		var f *excelize.File
		var err error
		var startRow int

		if e.Exist {
			f, err = excelize.OpenFile(e.Path)
			if err != nil {
				return err
			}
			if e.Sheet == "" {
				e.Sheet = "Sheet" + strconv.FormatInt(int64(f.SheetCount+1), 10)
			}

		} else {
			f = excelize.NewFile()
			if e.Sheet == "" {
				e.Sheet = "Sheet1"
			}
		}

		startRow = 1

		index, _ := f.GetSheetIndex(e.Sheet)

		if e.Append && index >= 0 {
			startRow, err = e.getRow(f, e.Sheet)

			if err != nil {
				return err
			}
		}

		if !e.Append && index >= 0 {
			f.DeleteSheet(e.Sheet)
			if _, err := f.NewSheet(e.Sheet); err != nil {
				return err
			}
		}

		if index == -1 {
			if _, err := f.NewSheet(e.Sheet); err != nil {
				return err
			}
		}

		if len(e.Header) > 0 {
			data := array.ArrayToAny(e.Header)
			f.SetSheetRow(e.Sheet, fmt.Sprintf("A%d", startRow), &data)
			startRow++
		}
		for strSlice := range ch {
			data := array.ArrayToAny(strSlice)
			f.SetSheetRow(e.Sheet, fmt.Sprintf("A%d", startRow), &data)
			startRow++
		}

		if err := f.SaveAs(e.Path); err != nil {
			return err
		}

		return nil
	}

}
