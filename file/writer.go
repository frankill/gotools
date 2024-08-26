package file

import (
	"bufio"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Writer struct {
	Comma     string // Field delimiter (set to ',' by NewWriter)
	UseCRLF   bool   // True to use \r\n as the line terminator
	w         *bufio.Writer
	useQuote  bool
	escape    byte
	escapeStr string
}

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer, seq string, useQuote bool, escape byte) *Writer {

	if escape != '\\' {
		escape = '"'
	}

	return &Writer{
		Comma:     seq,
		w:         bufio.NewWriter(w),
		useQuote:  useQuote,
		escape:    escape,
		escapeStr: string(escape),
	}
}

func (w *Writer) Write(record []string) error {

	escapeQuote := string(w.escape) + `"`

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
			i := strings.IndexAny(field, "\"\r\n"+w.escapeStr)
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
					_, err = w.w.WriteString(escapeQuote)
				case '\\':
					if w.escape == '\\' {
						_, err = w.w.WriteString(`\\`)
					}
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

func (w *Writer) Flush() error {
	return w.w.Flush()
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
			if w.escape == '\\' && c == '\\' {
				return true
			}
		}
	}

	if strings.ContainsAny(field, w.Comma) || strings.ContainsAny(field, "\"\r\n") {
		return true
	}

	if w.escape == '\\' && strings.ContainsAny(field, `\\`) {
		return true
	}

	r1, _ := utf8.DecodeRuneInString(field)
	return unicode.IsSpace(r1)
}
