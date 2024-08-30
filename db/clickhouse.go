package db

import (
	"crypto/tls"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type CK struct {
	DB
}

func NewCK(user string, pwd string, database string, host ...string) *CK {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: host,
		Auth: clickhouse.Auth{
			Database: database,
			Username: user,
			Password: pwd,
		},
		TLS: &tls.Config{
			InsecureSkipVerify: true,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: time.Second * 30,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug:                true,
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
	})

	return &CK{DB{Con: conn}}
}
