package data

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB = nil

func OpenDatabase(username, password, database string) (err error) {
	db, err = sql.Open("mysql", username+":"+password+"@/"+database+"?parseTime=true")
	return err
}

func CloseDatabase() {
	db.Close()
}

type databaseError struct {
	s string
}

func (e databaseError) Error() string {
	return e.s
}
