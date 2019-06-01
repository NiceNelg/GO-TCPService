package mysqlpool

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Model struct {
	db    *sql.DB
	table string
}

func Init(username string, pwd string, host string, port string, database string) (db *sql.DB) {
	db, err := sql.Open("mysql", username+":"+pwd+"@tcp("+host+":"+port+")/"+database)
	if err != nil {
		os.Exit(-1)
	}
	return
}
