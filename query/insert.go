package query

import (
	"fmt"
	"log"
	"strings"

	"github.com/frankill/gotools/array"
)

type SqlInsert interface {
	AddValues(vals ...[]any)
	Build() (string, []any)
	Clear()
}

type CKInsert struct {
	TableName string
	Columns   []string
}

func NewCKInsert(tableName string) *CKInsert {
	return &CKInsert{
		TableName: tableName,
	}
}

func (c *CKInsert) AddColumn(col ...string) *CKInsert {
	c.Columns = append(c.Columns, col...)
	return c
}

func (c *CKInsert) Clear() {

}

func (c *CKInsert) AddValues(vals ...[]any) {

}

func (c *CKInsert) Build() (string, []any) {
	return fmt.Sprintf("INSERT INTO %s (%s) ", c.TableName, strings.Join(c.Columns, ", ")), nil
}

type MysqlInsert struct {
	TableName     string
	Columns       []string
	UpdateColumns []string
	InsertValues  [][]any
	IsUpdate      bool
	IsReplace     bool
	IsIgnore      bool
}

func NewMysqlInsert(tableName string, Replace bool, ignore bool) *MysqlInsert {
	return &MysqlInsert{
		TableName: tableName,
		IsReplace: Replace,
		IsIgnore:  ignore,
	}
}

func (m *MysqlInsert) AddColumn(col ...string) *MysqlInsert {
	m.Columns = append(m.Columns, col...)
	return m
}

func (m *MysqlInsert) AddUpdateColumn(col ...string) *MysqlInsert {
	if m.IsReplace {
		log.Println("Warning: Replace insert type does not support update columns.")
		return m
	}

	if m.IsIgnore {
		log.Println("Warning: Ignore insert type does not support update columns.")
		return m
	}
	m.UpdateColumns = append(m.UpdateColumns, col...)
	m.IsUpdate = true

	return m
}

func (m *MysqlInsert) AddValues(vals ...[]any) {
	m.InsertValues = append(m.InsertValues, vals...)
}

func (m *MysqlInsert) AddValue(vals ...any) *MysqlInsert {
	m.InsertValues = append(m.InsertValues, vals)
	return m
}

func (m *MysqlInsert) Build() (string, []any) {

	if len(m.Columns) < len(m.UpdateColumns) {
		log.Println("Warning: Insert columns count is less than update columns count.")
		return "", nil
	}

	var builder strings.Builder
	if m.IsReplace {
		builder.WriteString("REPLACE")
	} else {
		builder.WriteString("INSERT")
	}

	if m.IsIgnore {
		builder.WriteString(" IGNORE")
	}

	if (len(m.Columns) == 1 && m.Columns[0] == "*") || len(m.Columns) == 0 {
		builder.WriteString(fmt.Sprintf(" INTO %s VALUES ", m.TableName))
	} else {
		builder.WriteString(fmt.Sprintf(" INTO `%s` (`%s`) VALUES ", m.TableName, strings.Join(m.Columns, "`, `")))
	}

	var perch string

	if len(m.Columns) > 1 {
		perch = "(" + strings.Join(array.Repeat("?", len(m.Columns)), ", ") + ")"
	} else if len(m.InsertValues) > 0 {
		perch = "(" + strings.Join(array.Repeat("?", len(m.InsertValues[0])), ", ") + ")"
	}

	if perch != "" {
		builder.WriteString(strings.Join(array.Repeat(perch, max(len(m.InsertValues), 1)), ", "))
	}

	if m.IsUpdate && len(m.UpdateColumns) > 0 {
		builder.WriteString(" ON DUPLICATE KEY UPDATE ")

		var updateStrings []string
		for _, col := range m.UpdateColumns {
			updateStrings = append(updateStrings, fmt.Sprintf("`%s` = VALUES(`%s`)", col, col))
		}

		builder.WriteString(strings.Join(updateStrings, ", "))
	}

	return builder.String(), array.Concat(m.InsertValues...)
}

func (m *MysqlInsert) Clear() {
	m.InsertValues = [][]any{}
}
