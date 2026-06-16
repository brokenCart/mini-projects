package urlshortener

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func GetDB(dbName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTable(db *sql.DB) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS mappings (
			path TEXT PRIMARY KEY,
			url TEXT NOT NULL
	);`
	_, err := db.Exec(createTableSQL)
	return err
}

func GetData(db *sql.DB) (*sql.Rows, error) {
	getDataSQL := `SELECT path, url FROM mappings`
	rows, err := db.Query(getDataSQL)
	return rows, err
}
