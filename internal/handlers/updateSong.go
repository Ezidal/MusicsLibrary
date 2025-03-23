package handlers

import (
	er "LibMusic/internal/logger/err"
	"LibMusic/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// UpdateSong godoc
// @Summary Обновить данные песни
// @Description Обновляет данные песни по её ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param input body models.Song true "Данные для обновления песни"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /songs/{id} [put]
func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	path := strings.TrimPrefix(r.URL.Path, "/songs/")
	id := strings.Split(path, "/")[0]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing song ID")
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.log.Error("Invalid song ID", er.Err(err))
		respondWithError(w, http.StatusBadRequest, "Invalid song ID")
		return
	}
	err = json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		h.log.Error("Failed to decode request body", er.Err(err))
		respondWithError(w, http.StatusBadRequest, "Failed to decode request body")
		return
	}

	err = h.storage.UpdateSong(idInt, song)
	if err != nil {
		h.log.Error("Failed to update song", er.Err(err))
		respondWithError(w, http.StatusInternalServerError, "Failed to update song")
		return
	}

	respondWithJSON(w, http.StatusOK, models.Response{Status: http.StatusOK, Message: "Song updated by ID: " + id})
	h.log.Info("Song updated by ID: " + id)

}
