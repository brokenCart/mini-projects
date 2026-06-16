package urlshortener

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// GetDB opens a connection to the SQLite database with the given name and returns the database handle.
func GetDB(dbName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// CreateTable creates a table named "mappings" in the SQLite database if it doesn't already exist.
// The table has two columns: "path" (TEXT, PRIMARY KEY) and "url" (TEXT, NOT NULL).
func CreateTable(db *sql.DB) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS mappings (
			path TEXT PRIMARY KEY,
			url TEXT NOT NULL
	);`
	_, err := db.Exec(createTableSQL)
	return err
}
