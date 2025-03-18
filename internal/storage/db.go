package storage

import (
	"LibMusic/internal/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(config *config.Config) (*Storage, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.DBUser, config.DBPass, config.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS songs (
		id SERIAL PRIMARY KEY,
		group_name TEXT NOT NULL,
		song_name TEXT NOT NULL,
		release_date DATE,
		text TEXT,
		link TEXT
	);`)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &Storage{db: db}, nil
}

func (s Storage) Get() {}
