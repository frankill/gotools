package db

import (
	"crypto/tls"

	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
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
