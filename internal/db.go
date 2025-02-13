package internal

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

//go:embed trigger.sql
var TRIGGER_SQL string

//go:embed notify.sql
var NOTIFY_SQL string

type pgDB struct {
	*sql.DB
}

func newPG(db *sql.DB) *pgDB {
	return &pgDB{db}

}

func (p *pgDB) createTriggers(table string) error {
	sql := fmt.Sprintf(TRIGGER_SQL, table, table)
	_, err := p.Exec(sql)
	if err != nil {
		return err
	}

	//now re/create notification trigger
	sql = fmt.Sprintf(NOTIFY_SQL, table, table, table)
	_, err = p.Exec(sql)
	if err != nil {
		return err
	}

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
