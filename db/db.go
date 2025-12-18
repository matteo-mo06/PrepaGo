package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	driver   = "mysql"
	host     = "localhost"
	port     = 3306
	user     = "root"
	password = "nyat@.pqsdfgLCL0"
	dbname   = "exam_api"
)

var Conn *sql.DB

func NewDB() *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		user, password, host, port, dbname)

	conn, err := sql.Open(driver, dsn)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connected to database!")
	return conn
}