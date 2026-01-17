package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date TEXT NOT NULL,
	title VARCHAR(255) NOT NULL,
	comment TEXT,
	repeat VARCHAR(128)
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);`

func Init(dbFile string) error {
	conn, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	db = conn

	_, err = db.Exec(schema)
	if err != nil {
		return err
	}

	return nil
}
