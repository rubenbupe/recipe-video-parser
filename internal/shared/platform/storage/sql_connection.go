package storage

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Connection struct {
	Db *sql.DB
}

var connections map[string]*Connection = map[string]*Connection{}

func CreateConnection(name string, cfg *Dbconfig) (*Connection, error) {
	if len(cfg.Database) < 3 || cfg.Database[len(cfg.Database)-3:] != ".db" {
		cfg.Database += ".db"
	}
	dsn := cfg.Database
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	connection := &Connection{
		Db: db,
	}

	connections[name] = connection

	return connection, nil
}

type Dbconfig struct {
	Database string        `default:"app"`
	Timeout  time.Duration `default:"5s"`
}
