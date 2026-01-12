package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	DB = db

	if err := createTables(); err != nil {
		return nil, err
	}

	return db, nil
}

func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS tracked_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type TEXT NOT NULL,
			external_id TEXT NOT NULL,
			title TEXT NOT NULL,
			year INTEGER,
			poster_url TEXT,
			backdrop_path TEXT,
			overview TEXT,
			genres TEXT,
			path TEXT,
			status TEXT DEFAULT 'wanted',
			added_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(type, external_id)
		)`,
		`CREATE TABLE IF NOT EXISTS seasons (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tracked_item_id INTEGER NOT NULL,
			season_number INTEGER NOT NULL,
			title TEXT,
			overview TEXT,
			poster_path TEXT,
			status TEXT DEFAULT 'wanted',
			FOREIGN KEY (tracked_item_id) REFERENCES tracked_items(id) ON DELETE CASCADE,
			UNIQUE(tracked_item_id, season_number)
		)`,
		`CREATE TABLE IF NOT EXISTS episodes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tracked_item_id INTEGER NOT NULL,
			season_id INTEGER NOT NULL,
			season_number INTEGER NOT NULL,
			episode_number INTEGER NOT NULL,
			title TEXT,
			air_date DATETIME,
			status TEXT DEFAULT 'wanted',
			FOREIGN KEY (tracked_item_id) REFERENCES tracked_items(id) ON DELETE CASCADE,
			FOREIGN KEY (season_id) REFERENCES seasons(id) ON DELETE CASCADE,
			UNIQUE(tracked_item_id, season_number, episode_number)
		)`,
	}

	for _, q := range queries {
		_, err := DB.Exec(q)
		if err != nil {
			return err
		}
	}

	return nil
}
