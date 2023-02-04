package sql

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
)

func test() {
	cfg := mysql.Config{
		User:   username,
		Passwd: password,
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "jazzrecords",
	}
	connector, err := mysql.NewConnector(&cfg)
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
}
