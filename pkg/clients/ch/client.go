package ch

import (
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"time"
)

type Client struct {
	driver.Conn
}

type Config struct {
	Addrs    []string
	Database string
	Username string
	Password string
}

func NewClient(conf Config) *Client {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: conf.Addrs,
		Auth: clickhouse.Auth{
			Database: conf.Database,
			Username: conf.Username,
			Password: conf.Password,
		},
		ReadTimeout:  60 * time.Second,
		MaxOpenConns: 100,
		MaxIdleConns: 50,
	})

	if err != nil {
		panic(err)
	}

	return &Client{
		Conn: conn,
	}
}
