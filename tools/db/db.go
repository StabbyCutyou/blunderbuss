package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const schema = `
CREATE TABLE events (
    application TEXT,
    type TEXT,
    message TEXT,
    context JSON,
    stack_trace TEXT,
    created_at TIMESTAMPTZ
)`

func main() {
	db, err := sqlx.Open("postgres", "user=stabby dbname=bbus sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db.MustExec(schema)
}
