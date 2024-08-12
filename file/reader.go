package file

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"

	"github.com/xuri/excelize/v2"
)

// ReadFromCsvSliceChannel 从 CSV 文件中读取数据，并将每一行数据作为字符串切片发送到通道。
// 参数:
//
//	filename - CSV 文件的路径。
//	ch - 用于接收字符串切片的通道。
//
// 返回:
//
//	如果在读取过程中遇到错误，则返回错误。
func ReadFromCsvSliceChannel(filename string, ch chan []string) error {

	defer close(ch)

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)

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
	return nil
}

// ReadFromCsv 一次性读取整个 CSV 文件的内容，并返回一个二维字符串数组。
// 参数:
//
//	filename - CSV 文件的路径。
//
// 返回:
//
//	一个二维字符串数组，其中每一行代表 CSV 文件中的一行数据。
//	如果在读取过程中遇到错误，则返回错误。
func ReadFromCsv(filename string) ([][]string, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)

	return reader.ReadAll()
}

// ReadFromExcel 从 Excel 文件中读取指定工作表的所有行，并返回一个二维字符串数组。
// 参数:
//
//	filename - Excel 文件的路径。
//	sheet - 工作表的名称。
//
// 返回:
//
//	一个二维字符串数组，其中每一行代表 Excel 工作表中的一行数据。
//	如果在读取过程中遇到错误，则返回错误。
func ReadFromExcel(filename string, sheet string) ([][]string, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	return f.GetRows(sheet)
}

// ReadFromExcelSliceChannel 从 Excel 文件中读取指定工作表的每一行数据，并通过通道发送出去。
// 参数:
//
//	filename - Excel 文件的路径。
//	sheet - 工作表的名称。
//	ch - 用于发送字符串切片的通道。
//
// 返回:
//
//	如果在读取过程中遇到错误，则返回错误。
func ReadFromExcelSliceChannel(filename string, sheet string, ch chan []string) error {

	defer close(ch)

	f, err := excelize.OpenFile(filename)
	if err != nil {
		return err
	}
	rows, err := f.Rows(sheet)
	if err != nil {
		return err
	}
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return err
		}
		ch <- row
	}
	return nil
}

func ReadFromTable(filename string, seq string) ([][]string, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := NewReader(f, seq)

	return reader.ReadAll()
}

func ReadFromTableSliceChannel(filename string, seq string, ch chan []string) error {

	defer close(ch)

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := NewReader(f, seq)

	for {
		record, err := reader.Read()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}
		ch <- record
	}
}

// A ParseError is returned for parsing errors.
// Line and column numbers are 1-indexed.
type ParseError struct {
	StartLine int   // Line where the record starts
	Line      int   // Line where the error occurred
	Column    int   // Column (1-based byte index) where the error occurred
	Err       error // The actual error
}

func (e *ParseError) Error() string {
	if e.Err == ErrFieldCount {
		return fmt.Sprintf("record on line %d: %v", e.Line, e.Err)
	}
	if e.StartLine != e.Line {
		return fmt.Sprintf("record on line %d; parse error on line %d, column %d: %v", e.StartLine, e.Line, e.Column, e.Err)
	}
	return fmt.Sprintf("parse error on line %d, column %d: %v", e.Line, e.Column, e.Err)
}

func (e *ParseError) Unwrap() error { return e.Err }

// These are the errors that can be returned in [ParseError.Err].
var (
	ErrBareQuote  = errors.New("bare \" in non-quoted-field")
	ErrQuote      = errors.New("extraneous or missing \" in quoted-field")
	ErrFieldCount = errors.New("wrong number of fields")

	// Deprecated: ErrTrailingComma is no longer used.
	ErrTrailingComma = errors.New("extra delimiter at end of line")
)

type Reader struct {
	Comma string

	Comment rune

	FieldsPerRecord int

	LazyQuotes bool

	TrimLeadingSpace bool

	ReuseRecord bool

	TrailingComma bool

	r *bufio.Reader

	numLine int

	offset int64

	rawBuffer []byte

	recordBuffer []byte

	fieldIndexes []int

	fieldPositions []position

	lastRecord []string
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader, sep string) *Reader {
	return &Reader{
		Comma: sep,
		r:     bufio.NewReader(r),
	}
}

func (r *Reader) Read() (record []string, err error) {
	if r.ReuseRecord {
		record, err = r.readRecord(r.lastRecord)
		r.lastRecord = record
	} else {
		record, err = r.readRecord(nil)
	}
	return record, err
}

func (r *Reader) FieldPos(field int) (line, column int) {
	if field < 0 || field >= len(r.fieldPositions) {
		panic("out of range index passed to FieldPos")
	}
	p := &r.fieldPositions[field]
	return p.line, p.col
}

func (r *Reader) InputOffset() int64 {
	return r.offset
}

// pos holds the position of a field in the current line.
type position struct {
	line, col int
}

func (r *Reader) ReadAll() (records [][]string, err error) {
	for {
		record, err := r.readRecord(nil)
		if err == io.EOF {
			return records, nil
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
}

func (r *Reader) readLine() ([]byte, error) {
	line, err := r.r.ReadSlice('\n')

	if err == bufio.ErrBufferFull {
		r.rawBuffer = append(r.rawBuffer[:0], line...)
		for err == bufio.ErrBufferFull {
			line, err = r.r.ReadSlice('\n')
			r.rawBuffer = append(r.rawBuffer, line...)
		}
		line = r.rawBuffer
	}
	readSize := len(line)
	if readSize > 0 && err == io.EOF {
		err = nil
		// For backwards compatibility, drop trailing \r before EOF.
		if line[readSize-1] == '\r' {
			line = line[:readSize-1]
		}
	}
	r.numLine++
	r.offset += int64(readSize)
	// Normalize \r\n to \n on all input lines.
	if n := len(line); n >= 2 && line[n-2] == '\r' && line[n-1] == '\n' {
		line[n-2] = '\n'
		line = line[:n-1]
	}
	return line, err
}

// lengthNL reports the number of bytes for the trailing \n.
func lengthNL(b []byte) int {
	if len(b) > 0 && b[len(b)-1] == '\n' {
		return 1
	}
	return 0
}

// nextRune returns the next rune in b or utf8.RuneError.
func nextRune(b []byte) rune {
	r, _ := utf8.DecodeRune(b)
	return r
}

func (r *Reader) readRecord(dst []string) ([]string, error) {

	// Read line (automatically skipping past empty lines and any comments).
	var line []byte
	var errRead error
	for errRead == nil {
		line, errRead = r.readLine()
		if r.Comment != 0 && nextRune(line) == r.Comment {
			line = nil
			continue // Skip comment lines
		}
		if errRead == nil && len(line) == lengthNL(line) {
			line = nil
			continue // Skip empty lines
		}
		break
	}
	if errRead == io.EOF {
		return nil, errRead
	}

	if r.Comma == "" {
		return []string{string(line[:len(line)-1])}, nil
	}

	// Parse each field in the record.
	var err error
	const quoteLen = len(`"`)
	commaLen := len(r.Comma)
	commabyte := []byte(r.Comma)
	recLine := r.numLine
	r.recordBuffer = r.recordBuffer[:0]
	r.fieldIndexes = r.fieldIndexes[:0]
	r.fieldPositions = r.fieldPositions[:0]
	pos := position{line: r.numLine, col: 1}
parseField:
	for {
		if r.TrimLeadingSpace {
			i := bytes.IndexFunc(line, func(r rune) bool {
				return !unicode.IsSpace(r)
			})
			if i < 0 {
				i = len(line)
				pos.col -= lengthNL(line)
			}
			line = line[i:]
			pos.col += i
		}
		if len(line) == 0 || line[0] != '"' {
			// Non-quoted string field
			i := bytes.Index(line, commabyte)
			field := line
			if i >= 0 {
				field = field[:i]
			} else {
				field = field[:len(field)-lengthNL(field)]
			}
			// Check to make sure a quote does not appear in field.
			if !r.LazyQuotes {
				if j := bytes.IndexByte(field, '"'); j >= 0 {
					col := pos.col + j
					err = &ParseError{StartLine: recLine, Line: r.numLine, Column: col, Err: ErrBareQuote}
					break parseField
				}
			}
			r.recordBuffer = append(r.recordBuffer, field...)
			r.fieldIndexes = append(r.fieldIndexes, len(r.recordBuffer))
			r.fieldPositions = append(r.fieldPositions, pos)
			if i >= 0 {
				line = line[i+commaLen:]
				pos.col += i + commaLen
				continue parseField
			}
			break parseField
		} else {
			// Quoted string field
			fieldPos := pos
			line = line[quoteLen:]
			pos.col += quoteLen
			for {
				i := bytes.IndexByte(line, '"')
				if i >= 0 {
					// Hit next quote.
					r.recordBuffer = append(r.recordBuffer, line[:i]...)
					line = line[i+quoteLen:]
					pos.col += i + quoteLen
					switch rn := nextRune(line); {
					case rn == '"':
						// `""` sequence (append quote).
						r.recordBuffer = append(r.recordBuffer, '"')
						line = line[quoteLen:]
						pos.col += quoteLen
					case string(line[:commaLen]) == r.Comma:
						// `",` sequence (end of field).
						line = line[commaLen:]
						pos.col += commaLen
						r.fieldIndexes = append(r.fieldIndexes, len(r.recordBuffer))
						r.fieldPositions = append(r.fieldPositions, fieldPos)
						continue parseField
					case lengthNL(line) == len(line):
						// `"\n` sequence (end of line).
						r.fieldIndexes = append(r.fieldIndexes, len(r.recordBuffer))
						r.fieldPositions = append(r.fieldPositions, fieldPos)
						break parseField
					case r.LazyQuotes:
						// `"` sequence (bare quote).
						r.recordBuffer = append(r.recordBuffer, '"')
					default:
						// `"*` sequence (invalid non-escaped quote).
						err = &ParseError{StartLine: recLine, Line: r.numLine, Column: pos.col - quoteLen, Err: ErrQuote}
						break parseField
					}
				} else if len(line) > 0 {
					// Hit end of line (copy all data so far).
					r.recordBuffer = append(r.recordBuffer, line...)
					if errRead != nil {
						break parseField
					}
					pos.col += len(line)
					line, errRead = r.readLine()
					if len(line) > 0 {
						pos.line++
						pos.col = 1
					}
					if errRead == io.EOF {
						errRead = nil
					}
				} else {
					// Abrupt end of file (EOF or error).
					if !r.LazyQuotes && errRead == nil {
						err = &ParseError{StartLine: recLine, Line: pos.line, Column: pos.col, Err: ErrQuote}
						break parseField
					}
					r.fieldIndexes = append(r.fieldIndexes, len(r.recordBuffer))
					r.fieldPositions = append(r.fieldPositions, fieldPos)
					break parseField
				}
			}
		}
	}
	if err == nil {
		err = errRead
	}

	str := string(r.recordBuffer)
	dst = dst[:0]
	if cap(dst) < len(r.fieldIndexes) {
		dst = make([]string, len(r.fieldIndexes))
	}
	dst = dst[:len(r.fieldIndexes)]
	var preIdx int
	for i, idx := range r.fieldIndexes {
		dst[i] = str[preIdx:idx]
		preIdx = idx
	}

	// Check or update the expected fields per record.
	if r.FieldsPerRecord > 0 {
		if len(dst) != r.FieldsPerRecord && err == nil {
			err = &ParseError{
				StartLine: recLine,
				Line:      recLine,
				Column:    1,
				Err:       ErrFieldCount,
			}
		}
	} else if r.FieldsPerRecord == 0 {
		r.FieldsPerRecord = len(dst)
	}
	return dst, err
}
