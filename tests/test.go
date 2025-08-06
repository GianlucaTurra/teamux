// Package tests
package tests

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func setup() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	schemaSQL, err := os.ReadFile("../sql_scripts/sessions.sql")
	if err != nil {
		log.Fatalf("Failed to read schema SQL: %v", err)
	}
	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		log.Fatalf("Failed to execute schema SQL: %v", err)
	}
	return db
}
