package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var DB *sql.DB

type Database struct {
	*sql.DB
}

func (d *Database) New(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	DB = db
	return db, nil

}
