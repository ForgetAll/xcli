package model

import (
	"database/sql"
	"os/user"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

var DbConn *sql.DB

func InitDb() {
	var err error
	usr, err := user.Current()
	checkError(err)
	dbPath := usr.HomeDir + "/.xcli_db"
	DbConn, err = sql.Open("sqlite3", dbPath)
	checkError(err)
}

func Release() {
	if DbConn != nil {
		_ = DbConn.Close()
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
