package config

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func NewDB() *sql.DB {
	psqlInfo := databaseUrl()
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func databaseUrl() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}
	return "host=localhost user=postgres password=postgres dbname=postgres sslmode=disable"
}
