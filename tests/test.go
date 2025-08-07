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
	createTable("../sql_scripts/sessions.sql", db)
	createTable("../sql_scripts/windows.sql", db)
	createTable("../sql_scripts/sessions_windows.sql", db)
	return db
}

func createTable(scriptFile string, db *sql.DB) {
	schemaSQL, err := os.ReadFile(scriptFile)
	if err != nil {
		log.Fatalf("Failed to read schema SQL: %v", err)
	}
	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		log.Fatalf("Failed to execute schema SQL: %v", err)
	}
}
