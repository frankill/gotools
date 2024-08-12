package file

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/frankill/gotools/array"
	"github.com/xuri/excelize/v2"
)

// WriteToExcelStringSliceChannel 从一个字符串切片的通道读取数据，并将其写入指定的 Excel 文件。
// 参数:
//
//	ch - 包含字符串切片的通道，每个切片代表一行数据。
//	stop - 用于停止写入过程的信号通道。
//	filename - 输出的 Excel 文件的路径。
//	sheet - Excel 文件中的工作表名称，默认为 "Sheet1"。
//
// 返回:
//
//	如果在写入过程中遇到错误，则返回错误。
func WriteToExcelStringSliceChannel(ch chan []string, stop chan struct{}, filename string, sheet string) error {

	defer close(stop)
	file := excelize.NewFile()

	if sheet == "" {
		sheet = "Sheet1"
	}
	stream, err := file.NewStreamWriter(sheet)
	if err != nil {
		return err
	}
	defer stream.Flush()
	num := 1

	for strSlice := range ch {
		stream.SetRow(fmt.Sprintf("A%d", num), array.ArrayToAny(strSlice))
		num++
	}

	if err := file.SaveAs(filename); err != nil {
		return err
	}

	return nil
}

// WriteToExcel 将二维字符串数组 data 写入指定的 Excel 文件。
// 参数:
//
//	data - 二维字符串数组，其中每一行代表 Excel 中的一行数据。
//	filename - 输出的 Excel 文件的路径。
//	sheet - Excel 文件中的工作表名称，默认为 "Sheet1"。
//
// 返回:
//
//	如果在写入过程中遇到错误，则返回错误。
func WriteToExcel(data [][]string, filename string, sheet string) error {

	file := excelize.NewFile()

	if sheet == "" {
		sheet = "Sheet1"
	}
	stream, err := file.NewStreamWriter(sheet)
	if err != nil {
		return err
	}
	defer stream.Flush()
	num := 1

	for _, row := range data {
		stream.SetRow(fmt.Sprintf("A%d", num), array.ArrayToAny(row))
		num++
	}

	// 保存文件
	if err := file.SaveAs(filename); err != nil {
		return err
	}

	return nil
}

// WriteToCSVStringSliceChannel 从一个字符串切片的通道读取数据，并将其写入指定的 CSV 文件。
// 参数:
//
//	ch - 包含字符串切片的通道，每个切片代表一行数据。
//	stop - 用于停止写入过程的信号通道。
//	filename - 输出的 CSV 文件的路径。
//
// 返回:
//
//	如果在写入过程中遇到错误，则返回错误。
func WriteToCSVStringSliceChannel(ch chan []string, stop chan struct{}, filename string) error {

	defer close(stop)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for row := range ch {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// WriteToCSV 将二维字符串数组 data 写入指定的 CSV 文件。
// 参数:
//
//	data - 二维字符串数组，其中每一行代表 CSV 中的一行数据。
//	filename - 输出的 CSV 文件的路径。
//
// 返回:
//
//	如果在写入过程中遇到错误，则返回错误。
func WriteToCSV(data [][]string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range data {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func WriteToTable(data [][]string, filename string, seq string, useQuote bool) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := NewWriter(file, seq, useQuote)

	return writer.WriteAll(data)
}

func WriteToTableSliceChannel(ch chan []string, stop chan struct{}, filename string, seq string, useQuote bool) error {

	defer close(stop)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := NewWriter(file, seq, useQuote)

	for record := range ch {
		err := writer.Write(record)
		if err != nil {
			return err
		}
	}
	return nil
}

type Writer struct {
	Comma    string // Field delimiter (set to ',' by NewWriter)
	UseCRLF  bool   // True to use \r\n as the line terminator
	w        *bufio.Writer
	useQuote bool
}

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer, seq string, useQuote bool) *Writer {
	return &Writer{
		Comma:    seq,
		w:        bufio.NewWriter(w),
		useQuote: useQuote,
	}
}

func (w *Writer) Write(record []string) error {

	for n, field := range record {
		if n > 0 {
			for _, c := range w.Comma {
				if _, err := w.w.WriteRune(c); err != nil {
					return err
				}
			}
		}

		// If we don't have to have a quoted field then just
		// write out the field and continue to the next field.
		if !w.fieldNeedsQuotes(field) {
			if _, err := w.w.WriteString(field); err != nil {
				return err
			}
			continue
		}

		if err := w.w.WriteByte('"'); err != nil {
			return err
		}
		for len(field) > 0 {
			// Search for special characters.
			i := strings.IndexAny(field, "\"\r\n")
			if i < 0 {
				i = len(field)
			}

			// Copy verbatim everything before the special character.
			if _, err := w.w.WriteString(field[:i]); err != nil {
				return err
			}
			field = field[i:]

			// Encode the special character.
			if len(field) > 0 {
				var err error
				switch field[0] {
				case '"':
					_, err = w.w.WriteString(`""`)
				case '\r':
					if !w.UseCRLF {
						err = w.w.WriteByte('\r')
					}
				case '\n':
					if w.UseCRLF {
						_, err = w.w.WriteString("\r\n")
					} else {
						err = w.w.WriteByte('\n')
					}
				}
				field = field[1:]
				if err != nil {
					return err
				}
			}
		}
		if err := w.w.WriteByte('"'); err != nil {
			return err
		}
	}
	var err error
	if w.UseCRLF {
		_, err = w.w.WriteString("\r\n")
	} else {
		err = w.w.WriteByte('\n')
	}
	return err
}

func (w *Writer) Flush() {
	w.w.Flush()
}

func (w *Writer) Error() error {
	_, err := w.w.Write(nil)
	return err
}

func (w *Writer) WriteAll(records [][]string) error {
	for _, record := range records {
		err := w.Write(record)
		if err != nil {
			return err
		}
	}
	return w.w.Flush()
}

func (w *Writer) fieldNeedsQuotes(field string) bool {

	if !w.useQuote {
		return false
	}

	if field == "" {
		return false
	}

	if field == `\.` {
		return true
	}

	if len(w.Comma) == 1 && w.Comma[0] < utf8.RuneSelf {
		for i := 0; i < len(field); i++ {
			c := field[i]
			if c == '\n' || c == '\r' || c == '"' || c == byte(w.Comma[0]) {
				return true
			}
		}
	}

	if strings.ContainsAny(field, w.Comma) || strings.ContainsAny(field, "\"\r\n") {
		return true
	}

	r1, _ := utf8.DecodeRuneInString(field)
	return unicode.IsSpace(r1)
}
