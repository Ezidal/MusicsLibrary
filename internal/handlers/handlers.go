package handlers

import (
	"LibMusic/internal/config"
	"LibMusic/internal/models"
	"encoding/json"
	"log/slog"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, models.Response{Error: message, Status: code})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

type Storage interface {
	GetAllSongs() ([]models.Song, error)
	AddSong(song models.Song) (int64, error)
	SongExist(song string) error
	DeleteSong(id int) error
	GetText(id int) (string, error)
}

type Handler struct {
	storage Storage
	log     *slog.Logger
	conf    *config.Config
}

func NewHandler(s Storage, l *slog.Logger, conf *config.Config) *Handler {
	return &Handler{storage: s, log: l, conf: conf}
}
