package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dataSourceName string) (*sql.DB, error) {

	_, err := os.Stat(dataSourceName)
	dbExists := err == nil || os.IsExist(err)

	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	if !dbExists {
		content, err := os.ReadFile("./pkg/database/data.sql")
		if err != nil {
			return nil, err
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return nil, err
		}
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil

}
