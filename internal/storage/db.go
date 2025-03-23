package storage

import (
	"LibMusic/internal/config"
	er "LibMusic/internal/logger/err"
	"LibMusic/internal/models"
	"database/sql"
	"fmt"
	"log/slog"
	"slices"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Storage struct {
	db        *sql.DB
	arrString []string
}

func Migrate(db *sql.DB, config *config.Config, log *slog.Logger) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Debug("failed to create migrate driver %s", er.Err(err))
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		config.DBName,
		driver,
	)
	if err != nil {
		log.Debug("failed to initialize migrate %s", er.Err(err))
		return fmt.Errorf("failed to initialize migrate: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Debug("failed to apply migrations %s", er.Err(err))
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Debug("failed to get migration version %s", er.Err(err))
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	log.Info("Database migrations applied. Current version: %d, dirty: %v", version, dirty)
	return nil
}

func New(config *config.Config, log *slog.Logger) (*Storage, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.DBUser, config.DBPass, config.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	// stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS songs (
	// 	id SERIAL PRIMARY KEY,
	// 	group_name TEXT NOT NULL,
	// 	song_name TEXT NOT NULL,
	// 	release_date DATE,
	// 	text TEXT,
	// 	link TEXT
	// );`)

	// if err != nil {
	// 	return nil, fmt.Errorf("%w", err)
	// }
	// _, err = stmt.Exec()
	// if err != nil {
	// 	return nil, fmt.Errorf("%w", err)
	// }

	err = Migrate(db, config, log)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	stmt, err := db.Prepare(`SELECT song_name FROM songs;`)
	if err != nil {
		log.Debug("failed to prepare select songs %s", er.Err(err))
		return nil, fmt.Errorf("%w", err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Debug("failed to query select songs %s", er.Err(err))
		return nil, fmt.Errorf("%w", err)
	}
	defer rows.Close()
	var arrString []string
	for rows.Next() {
		var songName string
		err := rows.Scan(&songName)
		if err != nil {
			log.Debug("failed to scan songs %s", er.Err(err))
			return nil, fmt.Errorf("%w", err)
		}
		arrString = append(arrString, songName)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

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

func (s *Storage) AddSong(song models.Song) (int64, error) {
	s.arrString = append(s.arrString, song.SongName)

	query := `
		INSERT INTO songs (group_name, song_name, release_date, text, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id int64
	err := s.db.QueryRow(query, song.Group, song.SongName, song.ReleaseDate, song.Text, song.Link).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to add song: %w", err)
	}

	return id, nil
}

func (s *Storage) SongExist(song string) error {
	if slices.Contains(s.arrString, song) {
		return fmt.Errorf("song already exists")
	}
	return nil
}

func removeElementByValue(slice []string, value string) []string {
	for i, v := range slice {
		if v == value {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice // Если элемент не найден, возвращаем исходный срез
}

func (s *Storage) DeleteSong(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid song ID: %d", id)
	}
	var song string
	err := s.db.QueryRow("SELECT song_name FROM songs WHERE id = $1", id).Scan(&song)
	if err != nil {
		return fmt.Errorf("failed to delete song: %w", err)
	}
	s.arrString = removeElementByValue(s.arrString, song)

	result, err := s.db.Exec("DELETE FROM songs WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete song: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no song found with ID %d", id)
	}

	return nil
}

func (s *Storage) GetText(id int) (string, error) {
	var text string
	err := s.db.QueryRow("SELECT text FROM songs WHERE id = $1", id).Scan(&text)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("song with ID %d not found", id)
		}
		return "", fmt.Errorf("failed to get song text: %w", err)
	}
	return text, nil
}

func (s *Storage) UpdateSong(id int, song models.Song) error {
	_, err := s.db.Exec(`UPDATE songs 
		SET group_name = $1, song_name = $2, release_date = $3, text = $4, link = $5 
		WHERE id = $6`, song.Group, song.SongName, song.ReleaseDate, song.Text, song.Link, id)
	if err != nil {
		return fmt.Errorf("failed to update song: %w", err)
	}
	return nil
}

func (s *Storage) Custom(query string, args ...any) (*sql.Rows, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	return rows, nil

}
