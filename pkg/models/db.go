package models

import (
	"github.com/go-pg/pg"
)

var db *pg.DB

func DB(host string, database string, user string, password string) {
	db = pg.Connect(&pg.Options{
		Addr:     host,
		User:     user,
		Database: database,
		Password: password,
	})
}
