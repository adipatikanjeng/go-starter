package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func Connect(user, password, dbname, host, port string) (*sql.DB, error) {
	connStr := fmt.Sprintf("%s:%s@/%s",
		user, password, dbname)
	return sql.Open("mysql", connStr)
}
