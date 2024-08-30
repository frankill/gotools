package db

import (
	"database/sql"
	"log"
	"strings"

	"github.com/frankill/gotools/query"
)

type MysqlDB struct {
	DB
	database string
}

// NewMysqlDB 创建一个新的MysqlDB实例，通过建立与MySQL数据库的连接。
// 参数:
//
//	con - 用于连接到MySQL数据库的连接字符串。
//
// 返回:
//
//	指向新创建的MysqlDB实例的指针。
func NewMysqlDB(dsn string) *MysqlDB {

	db, err := sql.Open("mysql", dsn)

	parts := strings.Split(dsn, "/")
	database := parts[len(parts)-1]

	// 移除可能的查询参数
	database = strings.Split(database, "?")[0]

	if err != nil {
		log.Panic(err)
	}
	return &MysqlDB{
		DB{
			Con: db,
		},
		database,
	}
}

type TableInfo struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
	Comment string
}

func (m *MysqlDB) QueryTableInfo(table string) ([]TableInfo, error) {

	q := query.NewSQLBuilder().From("INFORMATION_SCHEMA.COLUMNS ").
		Eq("TABLE_NAME", table).
		Eq("TABLE_SCHEMA", m.database).
		Select("COLUMN_NAME", "COLUMN_TYPE", "IS_NULLABLE", "COLUMN_KEY", "COLUMN_DEFAULT", "EXTRA", "COLUMN_COMMENT")

	// 准备查询
	rows, err := m.Con.Query(q.Build())
	if err != nil {
		return []TableInfo{}, err
	}
	defer rows.Close()

	var info []TableInfo

	// 遍历结果集
	for rows.Next() {
		var ti TableInfo
		def := sql.NullString{}
		if err := rows.Scan(&ti.Field, &ti.Type, &ti.Null, &ti.Key, &def, &ti.Extra, &ti.Comment); err != nil {
			return []TableInfo{}, err
		}

		ti.Default = def.String

		info = append(info, ti)
	}

	// 检查遍历过程中是否有错误发生
	if err := rows.Err(); err != nil {
		return []TableInfo{}, err
	}

	return info, nil
}
