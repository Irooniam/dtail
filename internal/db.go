package internal

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DataStore struct {
	*sql.DB
}

func NewDB(driver, connUri string) (*DataStore, error) {
	db, err := sql.Open(driver, connUri)
	if err != nil {
		return &DataStore{}, err
	}

	return &DataStore{db}, nil
}
