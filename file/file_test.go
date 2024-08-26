package file_test

import (
	"strings"
	"testing"

	"github.com/frankill/gotools/file"
)

func TestReadBasicFields(t *testing.T) {
	input := `a,b,c\n`
	reader := file.NewReader(strings.NewReader(input), ",", '"')
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := [][]string{{"a", "b", "c"}}
	if !compareRecords(records, expected) {
		t.Errorf("expected %v, got %v", expected, records)
	}
}

func TestQuotedFields(t *testing.T) {
	input := `"field1","field2","field3"\n`
	reader := file.NewReader(strings.NewReader(input), ",", '"')
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := [][]string{{"field1", "field2", "field3"}}
	if !compareRecords(records, expected) {
		t.Errorf("expected %v, got %v", expected, records)
	}
}

func TestEscapedQuotesWithDoubleQuote(t *testing.T) {
	input := `"field1","field""2","field3"\n`
	reader := file.NewReader(strings.NewReader(input), ",", '"')
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := [][]string{{"field1", `field"2`, "field3"}}
	if !compareRecords(records, expected) {
		t.Errorf("expected %v, got %v", expected, records)
	}
}

func TestEscapedQuotesWithBackslash(t *testing.T) {
	input := `"field1","field\\\"2","field3"\n`
	reader := file.NewReader(strings.NewReader(input), ",", '\\')
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := [][]string{{"field1", `field\"2`, "field3"}}
	if !compareRecords(records, expected) {
		t.Errorf("expected %v, got %v", expected, records)
	}
}

func TestUnclosedQuoteLazy(t *testing.T) {
	input := `"field1","field2\n`
	reader := file.NewReader(strings.NewReader(input), ",", '"')
	reader.LazyQuotes = true
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := [][]string{{"field1", "field2\n"}}
	if !compareRecords(records, expected) {
		t.Errorf("expected %v, got %v", expected, records)
	}
}

func TestTrailingComma(t *testing.T) {
	input := `"field1","field2",\n`
	reader := file.NewReader(strings.NewReader(input), ",", '"')
	reader.TrailingComma = true
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := [][]string{{"field1", "field2", ""}}
	if !compareRecords(records, expected) {
		t.Errorf("expected %v, got %v", expected, records)
	}
}

func TestCommentLine(t *testing.T) {
	input := `# This is a comment\n"field1","field2"\n`
	reader := file.NewReader(strings.NewReader(input), ",", '"')
	reader.Comment = '#'
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := [][]string{{"field1", "field2"}}
	if !compareRecords(records, expected) {
		t.Errorf("expected %v, got %v", expected, records)
	}
}

func TestFieldWithLineBreak(t *testing.T) {
	input := `"field1","field\n2","field3"\n`
	reader := file.NewReader(strings.NewReader(input), ",", '"')
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := [][]string{{"field1", "field\n2", "field3"}}
	if !compareRecords(records, expected) {
		t.Errorf("expected %v, got %v", expected, records)
	}
}

func TestFieldWithCarriageReturn(t *testing.T) {
	input := `"field1","field\r2","field3"\n`
	reader := file.NewReader(strings.NewReader(input), ",", '"')
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := [][]string{{"field1", "field\r2", "field3"}}
	if !compareRecords(records, expected) {
		t.Errorf("expected %v, got %v", expected, records)
	}
}

func TestFieldWithCommentAfterField(t *testing.T) {
	input := `"field1","field2" # this is a comment\n`
	reader := file.NewReader(strings.NewReader(input), ",", '"')
	reader.Comment = '#'
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := [][]string{{"field1", "field2"}}
	if !compareRecords(records, expected) {
		t.Errorf("expected %v, got %v", expected, records)
	}
}

func TestUnclosedQuoteStrict(t *testing.T) {
	input := `"field1","field2\n`
	reader := file.NewReader(strings.NewReader(input), ",", '"')
	records, err := reader.ReadAll()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	expected := [][]string{{"field1", "field2"}}
	if !compareRecords(records, expected) {
		t.Errorf("expected %v, got %v", expected, records)
	}
}

func TestEmptyFile(t *testing.T) {
	input := ""
	reader := file.NewReader(strings.NewReader(input), ",", '"')
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := [][]string{}
	if !compareRecords(records, expected) {
		t.Errorf("expected %v, got %v", expected, records)
	}
}

func TestFieldsWithEscapedDelimiter(t *testing.T) {
	input := `"field1|field2\\|field3"\n`
	reader := file.NewReader(strings.NewReader(input), "|", '\\')
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := [][]string{{"field1|field2|field3"}}
	if !compareRecords(records, expected) {
		t.Errorf("expected %v, got %v", expected, records)
	}
}

func TestCommentBeforeField(t *testing.T) {
	input := `# this is a comment\n"field1","field2"\n`
	reader := file.NewReader(strings.NewReader(input), ",", '"')
	reader.Comment = '#'
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := [][]string{{"field1", "field2"}}
	if !compareRecords(records, expected) {
		t.Errorf("expected %v, got %v", expected, records)
	}
}

func compareRecords(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !compareFields(a[i], b[i]) {
			return false
		}
	}
	return true
}

func compareFields(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
