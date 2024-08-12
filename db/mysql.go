package db

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/frankill/gotools/array"
	"github.com/frankill/gotools/query"
	"github.com/frankill/gotools/record"
	_ "github.com/go-sql-driver/mysql"
)

// CountQuery 将一个 SQL 查询语句转换为计数查询语句。
// 参数:
//
//	baseQuery - 原始的 SQL 查询语句。
//
// 返回:
//
//	转换后的计数查询语句。
func CountQuery(baseQuery string) string {

	query := strings.ToUpper(baseQuery)

	re := regexp.MustCompile("SELECT.*FROM")

	resultQuery := re.ReplaceAllString(query, "SELECT COUNT(1) c FROM")

	return strings.ToLower(resultQuery)
}

var (
	ModifyFunTemp = array.ArrayToAny[[]string]
)

// MysqlDB 定义了一个与MySQL数据库交互的结构体。
type MysqlDB struct {
	Con *sql.DB
}

// NewMysqlDB 创建一个新的MysqlDB实例，通过建立与MySQL数据库的连接。
// 参数:
//
//	con - 用于连接到MySQL数据库的连接字符串。
//
// 返回:
//
//	指向新创建的MysqlDB实例的指针。
func NewMysqlDB(con string) *MysqlDB {

	db, err := sql.Open("mysql", con)

	if err != nil {
		log.Panic(err)
	}
	return &MysqlDB{
		Con: db,
	}
}

// Close 关闭MysqlDB实例的底层数据库连接。
// 当MysqlDB实例不再需要使用时，应调用此方法。
func (m *MysqlDB) Close() {
	m.Con.Close()
}

// Insert 方法从一个通道接收数据，并批量插入到 MySQL 数据库中。
// 参数:
//
//	ch - 一个通道，用于接收待插入的数据。
//	stop - 一个信号通道，用于控制插入操作的停止。
//
// 返回:
//
//	一个函数，接受一个 SqlInsert 类型的参数，执行数据库插入操作。
func (m *MysqlDB) Insert(ch chan []string) func(query query.SqlInsert) error {

	return func(query query.SqlInsert) error {

		num := 1000
		res := make([][]any, 0, num)

		commit := func() error {
			if err := m.do(res, query); err != nil {
				return err
			}
			res = res[:0]
			return nil
		}

		for v := range ch {

			if len(v) == 0 {
				continue
			}

			res = append(res, ModifyFunTemp(v))

			if len(res) == num {
				err := commit()
				if err != nil {
					return err
				}
			}
		}

		if len(res) > 0 {
			err := commit()
			if err != nil {
				return err
			}
		}
		return nil
	}

}

func (m *MysqlDB) do(data [][]any, query query.SqlInsert) error {

	query.AddValues(data...)
	stmt, args := query.Build()

	defer query.Clear()

	tj, err := m.Con.Begin()

	defer func() {
		if err != nil {
			tj.Rollback()
		}
	}()

	if err != nil {
		return err
	}

	smt, err := tj.Prepare(stmt)
	if err != nil {
		return err
	}

	_, err = smt.Exec(args...)
	if err != nil {
		return err
	}
	err = tj.Commit()

	return err
}

// QueryOne 方法执行一个 SQL 查询语句，并返回查询结果的第一列数据。
// 参数:
//
//	sql - *SQLBuilder 类型的结构体，用于构建 SQL 查询语句。
//
// 返回:
//
//	一个函数，无参数，返回查询结果的第一列数据和可能的错误。
func (m *MysqlDB) QueryOne(query *query.SQLBuilder) func() ([]string, error) {

	query_ := query.Build()

	num, err := m.QueryCount(query)

	if err != nil {
		return func() ([]string, error) {
			return nil, err
		}
	}

	return func() ([]string, error) {

		rows, err := m.Con.Query(query_)

		if err != nil {
			return nil, err
		}

		dict := make(map[sql.NullString]struct{}, num)

		for rows.Next() {
			var tmp sql.NullString
			err := rows.Scan(&tmp)
			if err != nil {
				panic(err)
			}
			dict[tmp] = struct{}{}
		}

		return array.MapApply(func(k sql.NullString, v struct{}) string { return k.String }, dict), nil

	}

}

// QueryTwo 方法执行一个 SQL 查询语句，并返回查询结果的两列数据作为键值对的映射。
// 参数:
//
//	sql - *SQLBuilder 类型的结构体，用于构建 SQL 查询语句。
//
// 返回:
//
//	一个函数，无参数，返回查询结果的两列数据作为键值对的映射和可能的错误。
func (m *MysqlDB) QueryTwo(query *query.SQLBuilder) func() (map[string]string, error) {

	query_ := query.Build()

	num, err := m.QueryCount(query)

	if err != nil {
		return func() (map[string]string, error) {
			return nil, err
		}
	}
	return func() (map[string]string, error) {

		rows, err := m.Con.Query(query_)

		if err != nil {
			return nil, err
		}

		data := make(map[sql.NullString]sql.NullString, num)

		for rows.Next() {
			var one, two sql.NullString
			err := rows.Scan(&one, &two)
			if err != nil {
				return nil, err
			}
			data[one] = two
		}

		return array.MapApplyBoth(func(k sql.NullString, v sql.NullString) (string, string) { return k.String, v.String }, data), nil

	}

}

// QueryAny 方法执行一个 SQL 查询语句，并返回查询结果的所有列数据。
// 参数:
//
//	sql - *SQLBuilder 类型的结构体，用于构建 SQL 查询语句。
//
// 返回:
//
//	一个函数，无参数，返回查询结果的所有列数据和可能的错误。
func (m *MysqlDB) QueryAny(query *query.SQLBuilder) func() ([][]string, error) {

	query_ := query.Build()

	num, err := m.QueryCount(query)

	if err != nil {
		return func() ([][]string, error) {
			return nil, err
		}
	}

	return func() ([][]string, error) {

		rows, err := m.Con.Query(query_)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		lc := len(columns)

		rawResult := make([][]string, 0, num)

		row := make([]sql.NullString, lc)
		rowPointers := make([]any, lc)
		for i := range row {
			rowPointers[i] = &row[i]
		}

		for rows.Next() {
			err := rows.Scan(rowPointers...)
			if err != nil {
				log.Fatalln(err)
			}

			tmpRow := array.ArrayMap(func(x ...sql.NullString) string { return x[0].String }, row)
			rawResult = append(rawResult, tmpRow)

		}

		return rawResult, nil

	}

}

// QueryAnyIter 方法执行一个 SQL 查询语句，并通过一个通道逐行返回查询结果的所有列数据。
// 参数:
//
//	sql - *SQLBuilder 类型的结构体，用于构建 SQL 查询语句。
//
// 返回:
//
//	一个函数，返回一个通道，用于发送查询结果的每一行数据，以及可能的错误。
func (m *MysqlDB) QueryAnyIter(query *query.SQLBuilder) func() (chan []string, error) {

	query_ := query.Build()
	return func() (chan []string, error) {
		ch := make(chan []string, 100)

		rows, err := m.Con.Query(query_)

		if err != nil {
			return nil, err
		}

		columns, err := rows.Columns()

		if err != nil {
			return nil, err
		}

		go func() {
			defer close(ch)

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

		return ch, nil

	}

}

// QueryCount 方法执行一个 SQL 查询语句，并返回查询结果的行数。
// 参数:
//
//	sql - *SQLBuilder 类型的结构体，用于构建 SQL 查询语句。
//
// 返回:
//
//	查询结果的行数和可能的错误。
func (m *MysqlDB) QueryCount(query *query.SQLBuilder) (int, error) {

	query_ := query.Build()
	rows, err := m.Con.Query(CountQuery(query_))
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		err := rows.Scan(&rowCount)
		if err != nil {
			return 0, err
		}
	}

	return rowCount, nil

}

func (m *MysqlDB) QueryField(query *query.SQLBuilder) ([]string, error) {

	q := query.Copy()

	rows, err := m.Con.Query(q.Eq("1", 2).Build())
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	return rows.Columns()

}

// QueryVector 方法执行一个 SQL 查询语句，并返回查询结果的每一列数据。
// 参数:
//
//	sql - *SQLBuilder 类型的结构体，用于构建 SQL 查询语句。
//
// 返回:
//
//	一个函数，无参数，返回查询结果的每一列数据和可能的错误。
//	查询结果的每一列数据是一个二维切片，其中每一行是一个一维切片，代表查询结果的一列数据。
func (m *MysqlDB) QueryVector(query *query.SQLBuilder) func() ([][]string, error) {

	query_ := query.Build()

	num, err := m.QueryCount(query)

	if err != nil {
		return func() ([][]string, error) {
			return nil, err
		}
	}
	return func() ([][]string, error) {

		rows, err := m.Con.Query(query_)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		lc := len(columns)

		rawResult := make([][]string, lc)
		for i := range rawResult {
			rawResult[i] = make([]string, num)
		}

		row := make([]sql.NullString, lc)
		rowPointers := make([]any, lc)
		for i := range row {
			rowPointers[i] = &row[i]
		}

		rowIndex := 0
		for rows.Next() {
			err := rows.Scan(rowPointers...)
			if err != nil {
				log.Fatal(err)
			}

			for i := 0; i < lc; i++ {
				rawResult[i][rowIndex] = row[i].String
			}
			rowIndex++
		}

		return rawResult, nil
	}
}

func (m *MysqlDB) QueryRecord(query *query.SQLBuilder) func() (*record.Record, error) {

	query_ := query.Build()

	num, err := m.QueryCount(query)

	if err != nil {
		return func() (*record.Record, error) {
			return &record.Record{}, err
		}
	}
	return func() (*record.Record, error) {

		rows, err := m.Con.Query(query_)
		if err != nil {
			return &record.Record{}, err
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			return &record.Record{}, err
		}

		lc := len(columns)

		rawResult := make([][]any, lc)
		for i := range rawResult {
			rawResult[i] = make([]any, num)
		}

		field, _ := rows.ColumnTypes()
		ftype := extrcaType(field)
		row := makeRow(field)

		rowIndex := 0
		for rows.Next() {
			err := rows.Scan(row...)
			if err != nil {
				log.Fatal(err)
			}

			for i := 0; i < lc; i++ {
				switch ftype[i] {
				case "int":
					rawResult[i][rowIndex] = int((*row[i].(*sql.NullInt64)).Int64)
				case "string":
					rawResult[i][rowIndex] = (*row[i].(*sql.NullString)).String
				case "float64":
					rawResult[i][rowIndex] = (*row[i].(*sql.NullFloat64)).Float64
				case "bool":
					rawResult[i][rowIndex] = (*row[i].(*sql.NullBool)).Bool
				default:
					log.Fatal("Unknown database type: ", ftype[i])
				}

			}
			rowIndex++
		}

		res := record.NewRecord(query.TableName(), lc)

		for i := 0; i < lc; i++ {
			res.AddField(columns[i], rawResult[i]...)
		}

		return res, nil
	}
}

func extrcaType(columns []*sql.ColumnType) []string {

	res := make([]string, len(columns))

	for i := 0; i < len(columns); i++ {

		databaseTypeName := strings.ToLower(columns[i].DatabaseTypeName())

		switch databaseTypeName {
		case "varchar", "text", "char", "enum", "set":
			res[i] = "string"
		case "integer", "int", "bigint", "smallint", "tinyint":
			res[i] = "int"
		case "float", "double", "decimal":
			res[i] = "float64"
		case "date", "time", "datetime", "timestamp":
			res[i] = "string"
		case "boolean", "bit":
			res[i] = "bool"
		default:
			log.Printf("Unknown database type: %s", databaseTypeName)

		}
	}

	return res
}
func makeRow(columns []*sql.ColumnType) []any {

	res := make([]any, len(columns))

	for i := 0; i < len(columns); i++ {

		databaseTypeName := strings.ToLower(columns[i].DatabaseTypeName())

		switch databaseTypeName {
		case "varchar", "text", "char", "enum", "set":
			res[i] = new(sql.NullString)
		case "integer", "int", "bigint", "smallint", "tinyint":
			res[i] = new(sql.NullInt32)
		case "float", "double", "decimal":
			res[i] = new(sql.NullFloat64)
		case "date", "time", "datetime", "timestamp":
			res[i] = new(sql.NullString)
		case "boolean", "bit":
			res[i] = new(sql.NullString)
		default:
			log.Printf("Unknown database type: %s", databaseTypeName)

		}
	}

	return res

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

// NewMysqlQuery creates a query function for MySQL that scans results into the specified type T.
func NewMysqlQuery[T any](con string) func(query *query.SQLBuilder) (chan T, error) {

	return func(query *query.SQLBuilder) (chan T, error) {

		query_ := query.Build()

		ch := make(chan T, 100)

		db, err := sql.Open("mysql", con)
		if err != nil {
			return nil, err
		}

		rows, err := db.Query(query_)
		if err != nil {
			return nil, err
		}

		columns, err := rows.Columns()
		if err != nil {
			return nil, err
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

		return ch, nil
	}
}

func NewMysqlQueryStr(con string) func(query *query.SQLBuilder) (chan []string, error) {

	return func(query *query.SQLBuilder) (chan []string, error) {

		query_ := query.Build()
		ch := make(chan []string, 100)

		db, err := sql.Open("mysql", con)
		if err != nil {
			return nil, err
		}

		rows, err := db.Query(query_)
		if err != nil {
			return nil, err
		}

		columns, err := rows.Columns()

		if err != nil {
			return nil, err
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

		return ch, nil
	}
}
