package db

import (
	"database/sql"
	"log"
	"regexp"
	"strings"

	"github.com/frankill/gotools/array"
	"github.com/frankill/gotools/maps"
	"github.com/frankill/gotools/query"

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
	ModifyFunTemp = array.ToAny[[]string]
)

// MysqlDB 定义了一个与MySQL数据库交互的结构体。
type DB struct {
	Con *sql.DB
}

// Close 关闭MysqlDB实例的底层数据库连接。
// 当MysqlDB实例不再需要使用时，应调用此方法。
func (m *DB) Close() {
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
func (m *DB) Insert(q query.SqlInsert) func(ch chan []any) error {

	return func(ch chan []any) error {

		num := 1000
		res := make([][]any, 0, num)

		commit := func() error {
			if err := m.do(res, q); err != nil {
				return err
			}
			res = res[:0]
			return nil
		}

		for v := range ch {

			if len(v) == 0 {
				continue
			}

			res = append(res, v)

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

func (m *DB) do(data [][]any, query query.SqlInsert) error {

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
func (m *DB) QueryOne(query *query.SQLBuilder) func() ([]string, error) {

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

		return maps.Apply(func(k sql.NullString, v struct{}) string { return k.String }, dict), nil

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
func (m *DB) QueryTwo(query *query.SQLBuilder) func() (map[string]string, error) {

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

		return maps.ApplyBoth(func(k sql.NullString, v sql.NullString) (string, string) { return k.String, v.String }, data), nil

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
func (m *DB) QueryAny(query *query.SQLBuilder) func() ([][]string, error) {

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

			tmpRow := array.Map(func(x ...sql.NullString) string { return x[0].String }, row)
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
func (m *DB) QueryAnyIter(query *query.SQLBuilder) func() (chan []string, chan error) {

	return func() (chan []string, chan error) {
		query_ := query.Build()
		ch := make(chan []string, 100)
		errs := make(chan error, 1)

		go func() {
			defer close(ch)
			defer close(errs)

			rows, err := m.Con.Query(query_)
			if err != nil {
				errs <- err
				return
			}

			columns, err := rows.Columns()

			if err != nil {
				errs <- err
				return
			}
			defer m.Con.Close()

			lc := len(columns)
			for rows.Next() {
				row := make([]sql.NullString, lc)
				rowPointers := make([]any, lc)
				for i := range row {
					rowPointers[i] = &row[i]
				}

				err := rows.Scan(rowPointers...)
				if err != nil {
					errs <- err
					return
				}

				ch <- array.Map(func(i ...sql.NullString) string {
					return i[0].String
				}, row)
			}
		}()

		return ch, errs

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
func (m *DB) QueryCount(query *query.SQLBuilder) (int, error) {

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

func (m *DB) QueryField(query *query.SQLBuilder) ([]string, error) {

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
func (m *DB) QueryVector(query *query.SQLBuilder) func() ([][]string, error) {

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
