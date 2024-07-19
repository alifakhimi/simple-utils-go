package simutils

import (
	"context"
	"database/sql"
	"regexp"

	"github.com/mattn/go-sqlite3"
)

// addRegexpFunction adds the REGEXP function to SQLite
func addRegexpFunction(db *sql.DB) {
	conn, err := db.Conn(context.Background())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	_, err = conn.ExecContext(context.Background(), `PRAGMA case_sensitive_like = true`)
	if err != nil {
		panic(err)
	}

	err = conn.Raw(func(driverConn interface{}) error {
		sqliteConn := driverConn.(*sqlite3.SQLiteConn)
		return sqliteConn.RegisterFunc("regexp", func(re, s string) (bool, error) {
			return regexp.MatchString(re, s)
		}, true)
	})
	if err != nil {
		panic(err)
	}
}
