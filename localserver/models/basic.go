package models

import (
	"database/sql"
    _ "github.com/mattn/go-sqlite3"
	"fmt"
)

var Db *sql.DB
func init() {
	//load db
	db,err := sql.Open("sqlite3","./db/main")
	Db = db
    if err != nil {
        fmt.Println("Cannot connect to database")
        return
    }
    fmt.Println("Database connected")
}
