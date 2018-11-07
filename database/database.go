package database

import (
	"github.com/go-pg/pg"
)

var (
	connection *pg.DB
)

// Init connects to the DB and initiate
func Init() error {
	connection = pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
		Addr:     "127.0.0.1:5432",
	})

	return nil
}

func Get() *pg.DB {
	return connection
}
