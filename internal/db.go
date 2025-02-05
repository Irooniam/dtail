package internal

import (
	"database/sql"
	_ "embed"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type pgDB struct {
	*sql.DB
}

func newPG(db *sql.DB) *pgDB {
	return &pgDB{db}

}

func (p *pgDB) createTriggers() error {
	log.Println("pg trigger")
	return nil
}

func (p *pgDB) getTables() ([]string, error) {
	tables := make([]string, 0)
	sqlstr := `select table_name from information_schema.tables
    		where 
			table_schema not in ('information_schema', 'pg_catalog') and
    			table_type = 'BASE TABLE'`

	rows, err := p.Query(sqlstr)
	if err != nil {
		return tables, err
	}

	defer rows.Close()
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			log.Println(err)
		}

		tables = append(tables, table)
	}

	return tables, nil
}

func NewDB(driver, connUri string) (*sql.DB, error) {
	db, err := sql.Open(driver, connUri)
	if err != nil {
		return &sql.DB{}, err
	}

	return db, nil
}
