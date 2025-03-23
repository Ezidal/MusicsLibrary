package storage

import (
	"LibMusic/internal/config"
	"LibMusic/internal/models"
	"database/sql"
	"fmt"
	"slices"

	_ "github.com/lib/pq"
)

type Storage struct {
	db        *sql.DB
	arrString []string
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

	stmt, err = db.Prepare(`SELECT song_name FROM songs;`)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer rows.Close()
	var arrString []string
	for rows.Next() {
		var songName string
		err := rows.Scan(&songName)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		arrString = append(arrString, songName)
	}

	return &Storage{db: db, arrString: arrString}, nil
}

func (s *Storage) GetAllSongs() ([]models.Song, error) {
	rows, err := s.db.Query("SELECT * FROM songs")
	if err != nil {
		return nil, fmt.Errorf("failed to get songs: %w", err)
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		err := rows.Scan(&song.Id, &song.Group, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			return nil, fmt.Errorf("failed to scan songs: %w", err)
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func (s *Storage) AddSong(song models.Song) error {
	s.arrString = append(s.arrString, song.SongName)
	_, err := s.db.Exec("INSERT INTO songs (group_name, song_name, release_date, text, link) VALUES ($1, $2, $3, $4, $5)",
		song.Group, song.SongName, song.ReleaseDate, song.Text, song.Link)
	if err != nil {
		return fmt.Errorf("failed to add song: %w", err)
	}

	return nil
}

func (s *Storage) SongExist(song string) error {
	if slices.Contains(s.arrString, song) {
		return fmt.Errorf("song already exists")
	}
	return nil
}
