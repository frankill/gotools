package file

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

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
	Escape byte

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
func NewReader(r io.Reader, sep string, escape byte) *Reader {

	if sep == "" {
		sep = ","
	}

	if escape != '\\' {
		escape = '"'
	}

	return &Reader{
		Comma:  sep,
		r:      bufio.NewReader(r),
		Escape: escape,
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
				i := r.IndexByte(line, '"')
				if i >= 0 {
					// Hit next quote.
					r.recordBuffer = append(r.recordBuffer, line[:i]...)
					line = line[i+quoteLen:]
					pos.col += i + quoteLen
					switch rn := nextRune(line); {
					case rn == '"' && r.Escape == '"':
						// `""` sequence (append quote).
						r.recordBuffer = append(r.recordBuffer, '"')
						line = line[quoteLen:]
						pos.col += quoteLen
					case string(line[:commaLen]) == r.Comma && len(line) >= commaLen:
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

func (r *Reader) IndexByte(s []byte, c byte) int {

	l := len(s)

	if l == 0 {
		return -1
	}

	if s[0] == c {
		return 0
	}

	for i := 1; i < l; i++ {

		if s[i] == c {

			if s[i-1] == '\\' && r.Escape == '\\' {
				continue
			}

			return i
		}
	}
	return -1

}
