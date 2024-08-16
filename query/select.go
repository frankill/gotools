package query

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/olivere/elastic/v7"
)

func ExtractTableName(sql string) string {
	re := regexp.MustCompile(`\bFROM\s+(\w+)`)
	matches := re.FindStringSubmatch(sql)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// SQLBuilder 结构体用于构建 SQL 查询
type SQLBuilder struct {
	selectFields []string
	tableName    string
	whereClauses []string
	groupBy      []string
	orderBy      []string
	limit        int
	offset       int
}

func (sb *SQLBuilder) TableName() string {
	return sb.tableName
}

func (sb *SQLBuilder) Clear() {

	sb.selectFields = nil
	sb.tableName = ""
	sb.whereClauses = nil
	sb.groupBy = nil
	sb.orderBy = nil
	sb.limit = 0
	sb.offset = 0

}

// copy struct

func (sb *SQLBuilder) Copy() *SQLBuilder {
	newSb := &SQLBuilder{}
	newSb.selectFields = make([]string, len(sb.selectFields))
	copy(newSb.selectFields, sb.selectFields)
	newSb.tableName = sb.tableName
	newSb.whereClauses = make([]string, len(sb.whereClauses))
	copy(newSb.whereClauses, sb.whereClauses)
	newSb.groupBy = make([]string, len(sb.groupBy))
	copy(newSb.groupBy, sb.groupBy)
	newSb.orderBy = make([]string, len(sb.orderBy))
	copy(newSb.orderBy, sb.orderBy)
	newSb.limit = sb.limit
	newSb.offset = sb.offset
	return newSb
}

// NewSQLBuilder 返回一个新的 SQLBuilder 实例
func NewSQLBuilder() *SQLBuilder {
	return &SQLBuilder{}
}

// TransformSelect 方法用于对 SELECT 字段进行转换处理
func (sb *SQLBuilder) TransformSelect(transformFunc func(string) string) *SQLBuilder {
	transformedFields := make([]string, len(sb.selectFields))
	for i, field := range sb.selectFields {
		transformedFields[i] = transformFunc(field)
	}
	sb.selectFields = transformedFields
	return sb
}

// Select 方法用于指定选择的字段
func (sb *SQLBuilder) Select(fields ...string) *SQLBuilder {
	sb.selectFields = append(sb.selectFields, fields...)
	return sb
}

// From 方法用于指定查询的表名
func (sb *SQLBuilder) From(tableName any) *SQLBuilder {

	switch v := tableName.(type) {
	case string:
		sb.tableName = v
	case *SQLBuilder:
		sb.tableName = v.Build()
	default:
		panic(fmt.Sprintf("Unsupported type %T for table name", v))
	}

	return sb
}

// Where 方法用于添加 WHERE 子句条件
func (sb *SQLBuilder) Where(clauses ...string) *SQLBuilder {
	sb.whereClauses = append(sb.whereClauses, clauses...)
	return sb
}

// equal 方法用语添加 等于的条件到where中
func (sb *SQLBuilder) Eq(field string, value any) *SQLBuilder {
	switch v := value.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		sb.whereClauses = append(sb.whereClauses, fmt.Sprintf("%s = %v", field, v))
	case string:
		sb.whereClauses = append(sb.whereClauses, fmt.Sprintf("%s = '%s'", field, v))
	}
	return sb
}

// unequal 方法用语添加 不等于的条件到where中
func (sb *SQLBuilder) Uneq(field string, value any) *SQLBuilder {
	switch v := value.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		sb.whereClauses = append(sb.whereClauses, fmt.Sprintf("%s != %v", field, v))
	case string:
		sb.whereClauses = append(sb.whereClauses, fmt.Sprintf("%s != '%s'", field, v))
	}
	return sb
}

// gt 方法用语添加 大于的条件到where中
func (sb *SQLBuilder) Gt(field string, value any) *SQLBuilder {
	switch v := value.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		sb.whereClauses = append(sb.whereClauses, fmt.Sprintf("%s > %v", field, v))
	case string:
		sb.whereClauses = append(sb.whereClauses, fmt.Sprintf("%s > '%s'", field, v))
	}
	return sb
}

// gte 方法用语添加 大于等于的条件到where中
func (sb *SQLBuilder) Gte(field string, value any) *SQLBuilder {
	switch v := value.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		sb.whereClauses = append(sb.whereClauses, fmt.Sprintf("%s >= %v", field, v))
	case string:
		sb.whereClauses = append(sb.whereClauses, fmt.Sprintf("%s >= '%s'", field, v))
	}
	return sb
}

// lt 方法用语添加 小于的条件到where中
func (sb *SQLBuilder) Lt(field string, value any) *SQLBuilder {
	switch v := value.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		sb.whereClauses = append(sb.whereClauses, fmt.Sprintf("%s < %v", field, v))
	case string:
		sb.whereClauses = append(sb.whereClauses, fmt.Sprintf("%s < '%s'", field, v))
	}
	return sb
}

// lte 方法用语添加 小于等于的条件到where中
func (sb *SQLBuilder) Lte(field string, value any) *SQLBuilder {
	switch v := value.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		sb.whereClauses = append(sb.whereClauses, fmt.Sprintf("%s <= %v", field, v))
	case string:
		sb.whereClauses = append(sb.whereClauses, fmt.Sprintf("%s <= '%s'", field, v))
	}
	return sb
}

// And 方法用于添加 AND 条件
func (sb *SQLBuilder) And(clauses ...string) *SQLBuilder {
	if len(clauses) == 1 {
		sb.whereClauses = append(sb.whereClauses, clauses[0])
	} else {
		andClauses := fmt.Sprintf("(%s)", strings.Join(clauses, " AND "))
		sb.whereClauses = append(sb.whereClauses, andClauses)
	}
	return sb
}

// Or 方法用于添加 OR 条件
func (sb *SQLBuilder) Or(clauses ...string) *SQLBuilder {
	if len(clauses) == 1 {
		sb.whereClauses = append(sb.whereClauses, clauses[0])
	} else {
		orClauses := fmt.Sprintf("(%s)", strings.Join(clauses, " OR "))
		sb.whereClauses = append(sb.whereClauses, orClauses)
	}
	return sb
}

// In 方法用于生成 field IN (subquery) 或 field IN (value1, value2, ...) 的条件语句
func (sb *SQLBuilder) In(field string, values any) *SQLBuilder {
	switch v := values.(type) {
	case []string:
		sb.whereClauses = append(sb.whereClauses, fmt.Sprintf("%s IN ('%s')", field, strings.Join(v, "', '")))
	case []int:
		var strValues []string
		for _, val := range v {
			strValues = append(strValues, fmt.Sprintf("%d", val))
		}
		sb.whereClauses = append(sb.whereClauses, fmt.Sprintf("%s IN (%s)", field, strings.Join(strValues, ", ")))
	case string:
		// Check if the string is a valid subquery
		if isValidSubquery(v) {
			sb.whereClauses = append(sb.whereClauses, fmt.Sprintf("%s IN (%s)", field, v))
		} else {
			sb.whereClauses = append(sb.whereClauses, fmt.Sprintf("%s IN ('%s')", field, v))
		}
	case *SQLBuilder:
		sb.whereClauses = append(sb.whereClauses, fmt.Sprintf("%s IN (%s)", field, v.Build()))
	default:
		panic(fmt.Sprintf("Unsupported type %T for IN clause", v))
	}
	return sb
}

// isValidSubquery 检查一个字符串是否是有效的子查询
func isValidSubquery(s string) bool {
	// 简单示例：假设有效的子查询以 "SELECT" 开头
	return strings.HasPrefix(strings.ToUpper(strings.TrimSpace(s)), "SELECT")
}

// Between 方法用于生成 field BETWEEN lower AND upper 的条件语句
func (sb *SQLBuilder) Between(field string, lower, upper any) *SQLBuilder {
	switch v := lower.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		sb.whereClauses = append(sb.whereClauses, fmt.Sprintf("%s BETWEEN %v AND %v", field, lower, upper))
	case string:
		sb.whereClauses = append(sb.whereClauses, fmt.Sprintf("%s BETWEEN '%s' AND '%s'", field, lower, upper))
	default:
		// Add cases for other types as needed
		panic(fmt.Sprintf("Unsupported type %T for BETWEEN clause", v))
	}
	return sb
}

// GroupBy 方法用于指定 GROUP BY 子句
func (sb *SQLBuilder) GroupBy(fields ...string) *SQLBuilder {
	sb.groupBy = append(sb.groupBy, fields...)
	return sb
}

// OrderBy 方法用于指定 ORDER BY 子句
func (sb *SQLBuilder) OrderBy(fields ...string) *SQLBuilder {
	sb.orderBy = append(sb.orderBy, fields...)
	return sb
}

// Limit 方法用于指定 LIMIT 子句
func (sb *SQLBuilder) Limit(limit int) *SQLBuilder {
	sb.limit = limit
	return sb
}

// Offset 方法用于指定 OFFSET 子句
func (sb *SQLBuilder) Offset(offset int) *SQLBuilder {
	sb.offset = offset
	return sb
}

// Build 方法用于构建最终的 SQL 查询语句
func (sb *SQLBuilder) Build() string {
	var sql strings.Builder

	// 构建 SELECT 子句
	sql.WriteString("SELECT ")
	if len(sb.selectFields) > 0 {
		sql.WriteString(strings.Join(sb.selectFields, ", "))
	} else {
		sql.WriteString("*")
	}

	// 构建 FROM 子句
	if sb.tableName != "" {
		sql.WriteString(" FROM ")
		sql.WriteString(sb.tableName)
	}

	// 构建 WHERE 子句
	if len(sb.whereClauses) > 0 {
		sql.WriteString(" WHERE ")
		sql.WriteString(strings.Join(sb.whereClauses, " AND "))
	}

	// 构建 GROUP BY 子句
	if len(sb.groupBy) > 0 {
		sql.WriteString(" GROUP BY ")
		sql.WriteString(strings.Join(sb.groupBy, ", "))
	}

	// 构建 ORDER BY 子句
	if len(sb.orderBy) > 0 {
		sql.WriteString(" ORDER BY ")
		sql.WriteString(strings.Join(sb.orderBy, ", "))
	}

	// 构建 LIMIT 子句
	if sb.limit > 0 {
		sql.WriteString(fmt.Sprintf(" LIMIT %d", sb.limit))
	}

	// 构建 OFFSET 子句
	if sb.offset > 0 {
		sql.WriteString(fmt.Sprintf(" OFFSET %d", sb.offset))
	}

	return sql.String()
}

type EsQuery struct {
	querys []elastic.Query
	typ    string
}

func NewQuery() *EsQuery {
	return &EsQuery{
		typ: "",
	}
}

func NewFilterQuery() *EsQuery {
	return &EsQuery{
		typ:    "filter",
		querys: make([]elastic.Query, 0),
	}
}

func NewMustQuery() *EsQuery {
	return &EsQuery{
		typ:    "must",
		querys: make([]elastic.Query, 0),
	}
}

func NewMustNotQuery() *EsQuery {
	return &EsQuery{
		typ:    "must_not",
		querys: make([]elastic.Query, 0),
	}
}

func NewShouldQuery() *EsQuery {
	return &EsQuery{
		typ:    "should",
		querys: make([]elastic.Query, 0),
	}
}

func (q *EsQuery) Eq(field string, value any) *EsQuery {
	q.querys = append(q.querys, elastic.NewTermQuery(field, value))
	return q
}

func (q *EsQuery) In(field string, values ...any) *EsQuery {
	q.querys = append(q.querys, elastic.NewTermsQuery(field, values...))
	return q
}

func (q *EsQuery) NotIn(field string, values ...any) *EsQuery {
	q.querys = append(q.querys, elastic.NewBoolQuery().MustNot(elastic.NewTermsQuery(field, values...)))
	return q
}

func (q *EsQuery) Gt(field string, value any) *EsQuery {
	q.querys = append(q.querys, elastic.NewRangeQuery(field).Gt(value))
	return q
}

func (q *EsQuery) Gte(field string, value any) *EsQuery {
	q.querys = append(q.querys, elastic.NewRangeQuery(field).Gte(value))
	return q
}

func (q *EsQuery) Lt(field string, value any) *EsQuery {
	q.querys = append(q.querys, elastic.NewRangeQuery(field).Lt(value))
	return q
}

func (q *EsQuery) Lte(field string, value any) *EsQuery {
	q.querys = append(q.querys, elastic.NewRangeQuery(field).Lte(value))
	return q
}

func (q *EsQuery) Neq(field string, value any) *EsQuery {
	q.querys = append(q.querys, elastic.NewBoolQuery().MustNot(elastic.NewTermQuery(field, value)))
	return q
}

func (q *EsQuery) Like(field string, value string) *EsQuery {
	q.querys = append(q.querys, elastic.NewRegexpQuery(field, value))
	return q
}

func (q *EsQuery) Wildcard(field string, value string) *EsQuery {
	q.querys = append(q.querys, elastic.NewWildcardQuery(field, value))
	return q
}

func (q *EsQuery) Prefix(field string, value string) *EsQuery {
	q.querys = append(q.querys, elastic.NewPrefixQuery(field, value))
	return q
}

func (q *EsQuery) Fuzzy(field string, value string) *EsQuery {
	q.querys = append(q.querys, elastic.NewFuzzyQuery(field, value))
	return q
}

func (q *EsQuery) NotLike(field string, value string) *EsQuery {
	q.querys = append(q.querys, elastic.NewBoolQuery().MustNot(elastic.NewRegexpQuery(field, value)))
	return q
}

func (q *EsQuery) NotPrefix(field string, value string) *EsQuery {
	q.querys = append(q.querys, elastic.NewBoolQuery().MustNot(elastic.NewPrefixQuery(field, value)))
	return q
}

func (q *EsQuery) NotWildcard(field string, value string) *EsQuery {
	q.querys = append(q.querys, elastic.NewBoolQuery().MustNot(elastic.NewWildcardQuery(field, value)))
	return q
}

func (q *EsQuery) Script(script string, params map[string]interface{}, lang string) *EsQuery {
	scriptQuery := elastic.NewScript(script).Params(params).Lang(lang)
	q.querys = append(q.querys, elastic.NewScriptQuery(scriptQuery))
	return q
}

func (q *EsQuery) ScriptID(scriptID string) *EsQuery {
	q.querys = append(q.querys, elastic.NewScriptQuery(elastic.NewScriptStored(scriptID)))
	return q
}

func (q *EsQuery) Exists(field string) *EsQuery {
	q.querys = append(q.querys, elastic.NewExistsQuery(field))
	return q
}

func (q *EsQuery) NotExists(field string) *EsQuery {
	q.querys = append(q.querys, elastic.NewBoolQuery().MustNot(elastic.NewExistsQuery(field)))
	return q
}

func (q *EsQuery) Between(field string, lower, upper any) *EsQuery {
	q.querys = append(q.querys, elastic.NewRangeQuery(field).Gte(lower).Lte(upper))
	return q
}

func (q *EsQuery) Should(clauses *EsQuery, should int) *EsQuery {
	q.querys = append(q.querys, elastic.NewBoolQuery().Should(clauses.querys...).MinimumNumberShouldMatch(should))
	return q
}

func (q *EsQuery) Must(clauses *EsQuery) *EsQuery {
	q.querys = append(q.querys, elastic.NewBoolQuery().Must(clauses.querys...))
	return q
}

func (q *EsQuery) MustNot(clauses *EsQuery) *EsQuery {
	q.querys = append(q.querys, elastic.NewBoolQuery().MustNot(clauses.querys...))
	return q
}

func (q *EsQuery) Filter(clauses *EsQuery) *EsQuery {
	q.querys = append(q.querys, elastic.NewBoolQuery().Filter(clauses.querys...))
	return q
}

func (q *EsQuery) Where(clauses *EsQuery) *EsQuery {
	q.querys = append(q.querys, clauses.querys...)
	return q
}

func (q *EsQuery) Build() elastic.Query {

	if q.typ == "filter" {
		return elastic.NewBoolQuery().Filter(q.querys...)
	} else if q.typ == "must" {
		return elastic.NewBoolQuery().Must(q.querys...)
	} else if q.typ == "must_not" {
		return elastic.NewBoolQuery().MustNot(q.querys...)
	} else if q.typ == "should" {
		return elastic.NewBoolQuery().Should(q.querys...)
	}
	return nil
}

func (q *EsQuery) Source() string {
	source, _ := q.Build().Source()
	b, _ := json.Marshal(source)
	return string(b)
}