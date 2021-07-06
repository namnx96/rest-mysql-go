package db

import (
	"database/sql"
	"fmt"
)

const (
	Host     = "localhost"
	User     = "root"
	Password = "123456"
	Schema   = "test_db"
)

var (
	DB *sql.DB
)

func Open() error {
	var err error
	dbConn := fmt.Sprintf("%s:%s@tcp(%s)/%s", User, Password, Host, Schema)

	DB, err = sql.Open("mysql", dbConn)
	return err
}

func Close() error {
	return DB.Close()
}
