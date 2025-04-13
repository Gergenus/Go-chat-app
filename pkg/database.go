package pkg

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DBwraper struct {
	DB *sql.DB
}

func InitDB(user, password, host, port, dbname, sslmode string) DBwraper {
	connect := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbname, sslmode)
	conn, err := sql.Open("postgres", connect)
	if err != nil {
		log.Fatal("initdb", err)
	}
	return DBwraper{DB: conn}
}
