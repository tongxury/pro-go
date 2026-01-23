package mysqlz

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	DSN string
}

type Client struct {
	*sql.DB
}

func NewClient(conf Config) (*Client, error) {

	db, err := sql.Open("mysql", conf.DSN)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	_, err = db.Query("select version()")
	if err != nil {
		return nil, err
	}

	return &Client{DB: db}, nil
}
