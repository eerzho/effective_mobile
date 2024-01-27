package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Conn struct {
	DB *sql.DB
}

func MustNew(dbUrl string) *Conn {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		panic("failed connection to postgres: " + err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic("failed ping to postgres: " + err.Error())
	}

	return &Conn{DB: db}
}
