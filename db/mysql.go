package db

import (
	"database/sql"
	"log"
	"regexp"
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

	sql := strings.ToUpper(baseQuery)

	re := regexp.MustCompile("SELECT.*FROM")

	resultQuery := re.ReplaceAllString(sql, "SELECT COUNT(1) c FROM")

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
func (m *MysqlDB) Insert(ch chan []string, stop chan struct{}) func(sql query.SqlInsert) error {

	return func(sql query.SqlInsert) error {

		defer close(stop)

		num := 1000
		res := make([][]any, 0, num)

		commit := func() error {
			if err := m.do(res, sql); err != nil {
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

func (m *MysqlDB) do(data [][]any, sql query.SqlInsert) error {

	sql.AddValues(data...)
	stmt, args := sql.Build()

	defer sql.Clear()

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
func (m *MysqlDB) QueryOne(sql *query.SQLBuilder) func() ([]string, error) {

	query := sql.Build()

	num, err := m.QueryCount(sql)

	if err != nil {
		return func() ([]string, error) {
			return nil, err
		}
	}

	return func() ([]string, error) {

		rows, err := m.Con.Query(query)

		if err != nil {
			return nil, err
		}

		dict := make(map[string]struct{}, num)

		for rows.Next() {
			tmp := ""
			err := rows.Scan(&tmp)
			if err != nil {
				panic(err)
			}
			dict[tmp] = struct{}{}
		}

		return array.MapKeys(dict), nil

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
func (m *MysqlDB) QueryTwo(sql *query.SQLBuilder) func() (map[string]string, error) {

	query := sql.Build()

	num, err := m.QueryCount(sql)

	if err != nil {
		return func() (map[string]string, error) {
			return nil, err
		}
	}
	return func() (map[string]string, error) {

		rows, err := m.Con.Query(query)

		if err != nil {
			return nil, err
		}

		data := make(map[string]string, num)

		for rows.Next() {
			one, two := "", ""
			err := rows.Scan(&one, &two)
			if err != nil {
				return nil, err
			}
			data[one] = two
		}

		return data, nil

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
func (m *MysqlDB) QueryAny(sql *query.SQLBuilder) func() ([][]string, error) {

	query := sql.Build()

	num, err := m.QueryCount(sql)

	if err != nil {
		return func() ([][]string, error) {
			return nil, err
		}
	}

	return func() ([][]string, error) {

		rows, err := m.Con.Query(query)
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

		row := make([]string, lc)
		rowPointers := make([]any, lc)
		for i := range row {
			rowPointers[i] = &row[i]
		}

		for rows.Next() {
			err := rows.Scan(rowPointers...)
			if err != nil {
				log.Fatal(err)
			}

			tmpRow := make([]string, lc)
			copy(tmpRow, row)
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
//	一个函数，接收一个通道，用于发送查询结果的每一行数据，以及可能的错误。
func (m *MysqlDB) QueryAnyIter(sql *query.SQLBuilder) func(ch chan []string) error {

	query := sql.Build()
	return func(ch chan []string) error {

		defer close(ch)

		rows, err := m.Con.Query(query)

		if err != nil {
			return err
		}

		columns, err := rows.Columns()

		if err != nil {
			return err
		}

		lc := len(columns)

		for rows.Next() {
			row := make([]string, lc)
			rowPointers := make([]any, lc)
			for i := range row {
				rowPointers[i] = &row[i]
			}

			err := rows.Scan(rowPointers...)
			if err != nil {
				log.Println(err)
			}

			ch <- row
		}
		return nil

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
func (m *MysqlDB) QueryCount(sql *query.SQLBuilder) (int, error) {

	query := sql.Build()
	rows, err := m.Con.Query(CountQuery(query))
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

func (m *MysqlDB) QueryField(sql *query.SQLBuilder) ([]string, error) {

	q := sql.Copy()

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
func (m *MysqlDB) QueryVector(sql *query.SQLBuilder) func() ([][]string, error) {

	query := sql.Build()

	num, err := m.QueryCount(sql)

	if err != nil {
		return func() ([][]string, error) {
			return nil, err
		}
	}
	return func() ([][]string, error) {

		rows, err := m.Con.Query(query)
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

		row := make([]string, lc)
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
				rawResult[i][rowIndex] = row[i]
			}
			rowIndex++
		}

		return rawResult, nil
	}
}

func (m *MysqlDB) QueryRecord(sql *query.SQLBuilder) func() (*record.Record, error) {

	query := sql.Build()

	num, err := m.QueryCount(sql)

	if err != nil {
		return func() (*record.Record, error) {
			return &record.Record{}, err
		}
	}
	return func() (*record.Record, error) {

		rows, err := m.Con.Query(query)
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
					rawResult[i][rowIndex] = *row[i].(*int)
				case "string":
					rawResult[i][rowIndex] = *row[i].(*string)
				case "float64":
					rawResult[i][rowIndex] = *row[i].(*float64)
				case "bool":
					rawResult[i][rowIndex] = *row[i].(*bool)
				default:
					log.Fatal("Unknown database type: ", ftype[i])
				}

			}
			rowIndex++
		}

		res := record.NewRecord(sql.TableName(), lc)

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
			res[i] = new(string)
		case "integer", "int", "bigint", "smallint", "tinyint":
			res[i] = new(int)
		case "float", "double", "decimal":
			res[i] = new(float64)
		case "date", "time", "datetime", "timestamp":
			res[i] = new(string)
		case "boolean", "bit":
			res[i] = new(bool)
		default:
			log.Printf("Unknown database type: %s", databaseTypeName)

		}
	}

	return res

}
