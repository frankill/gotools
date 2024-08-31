package db

import (
	"crypto/tls"

	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/frankill/gotools/query"
)

type CK struct {
	DB
}

type CKinfo struct {
	HostPost []string
	Database string
	User     string
	Pwd      string
	Tls      *tls.Config
	Settings map[string]string
	Debug    bool
}

func NewCKLoc(ck *CKinfo) (driver.Conn, error) {

	c := clickhouse.Options{
		Addr: ck.HostPost,
		Auth: clickhouse.Auth{
			Database: ck.Database,
			Username: ck.User,
			Password: ck.Pwd,
		},

		DialTimeout: time.Second * 30,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug:                ck.Debug,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "my-app", Version: "0.1"},
			},
		},
	}

	if ck.Tls != nil {
		c.TLS = ck.Tls
	}

	if ck.Settings != nil {
		for k, v := range ck.Settings {
			c.Settings[k] = v
		}
	}

	conn, err := clickhouse.Open(&c)

	return conn, err

}

func NewCK(ck *CKinfo) *CK {

	c := clickhouse.Options{
		Addr: ck.HostPost,
		Auth: clickhouse.Auth{
			Database: ck.Database,
			Username: ck.User,
			Password: ck.Pwd,
		},

		DialTimeout: time.Second * 30,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug:                ck.Debug,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "my-app", Version: "0.1"},
			},
		},
	}
	conn := clickhouse.OpenDB(&c)

	return &CK{DB{Con: conn}}
}
func (m *CK) Insert(q query.SqlInsert) func(ch chan []any) error {

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

func (m *CK) do(data [][]any, query query.SqlInsert) error {

	query.AddValues()
	stmt, _ := query.Build()

	defer query.Clear()

	tj, err := m.Con.Begin()

	if err != nil {
		return err
	}

	smt, err := tj.Prepare(stmt)
	if err != nil {
		return err
	}
	for _, v := range data {

		_, err = smt.Exec(v...)
		if err != nil {
			return err
		}
	}

	err = tj.Commit()

	return err
}
