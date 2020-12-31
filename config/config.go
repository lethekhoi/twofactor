package config

import (
	"database/sql"

	_ "github.com/godror/godror"
)

func GetDB() (*sql.DB, error) {
	db, err := sql.Open("godror", "dip/dip@localhost:1521/orcl")

	if err != nil {
		return nil, err
	}
	return db, nil
}
