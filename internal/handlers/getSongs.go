package handlers

import (
	er "LibMusic/internal/logger/err"
	"net/http"
)

func (h *Handler) GetSongs(w http.ResponseWriter, r *http.Request) {
	// Get all songs from db
	songs, err := h.storage.GetAllSongs()
	if err != nil {
		h.log.Error("failed to get all songs", er.Err(err))
		respondWithError(w, http.StatusInternalServerError, "failed to get all songs")
		return
	}

	respondWithJSON(w, http.StatusOK, songs)
}
