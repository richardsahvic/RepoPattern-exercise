package datasource

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func InitConnection() *sqlx.DB {
	db, err := sqlx.Connect("mysql", "dev:dev@(localhost:3306)/myapp?parseTime=true")
	if err != nil {
		log.Fatalln("Failed to connect to database,    ", err)
	}
	return db
}
