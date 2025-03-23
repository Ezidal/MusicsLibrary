package handlers

import (
	er "LibMusic/internal/logger/err"
	"LibMusic/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

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

	respondWithJSON(w, http.StatusOK, models.Response{Status: http.StatusOK, Message: "Song updated" + id})
	h.log.Info("Song updated by ID: " + id)

}
