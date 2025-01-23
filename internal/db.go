package internal

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB() {
	log.Println("new db")
}
