package clients

import (
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"store/pkg/confcenter"
	"time"
)

type ClickHouseClient struct {
	driver.Conn
}

//type ClickHouseConfig struct {
//	Addrs    []string
//	Database string
//	Username string
//	Password string
//}

func NewClickHouseClient(conf confcenter.ClickHouseConfig) *ClickHouseClient {
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

	return &ClickHouseClient{
		Conn: conn,
	}
}
