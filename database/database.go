package database

import (
	"github.com/go-pg/pg"
	"os"
)

var (
	connection *pg.DB
)

// Init connects to the DB and initiate
func Init() error {
	addr := os.Getenv("DB_ADDR")
	if addr == "" {
		addr = "127.0.0.1:5432"
	}

	connection = pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
		Addr:     addr,
	})

	return nil
}

func Get() *pg.DB {
	return connection
}
