package models

import (
	"github.com/go-pg/pg"
)

var db *pg.DB

func DB(host string, database string) {
	db = pg.Connect(&pg.Options{
		Addr:     host,
		User:     "postgres",
		Database: database,
	})
}
