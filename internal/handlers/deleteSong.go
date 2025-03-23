package handlers

import (
	er "LibMusic/internal/logger/err"
	"LibMusic/internal/models"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	// id := r.URL.Query().Get("id")
	path := strings.TrimPrefix(r.URL.Path, "/songs/")
	id := strings.Split(path, "/")[0]
	if id == "" {
		h.log.Error("id is required")
		respondWithError(w, http.StatusBadRequest, "id is required")
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.log.Error("id not a a number", er.Err(err))
		respondWithError(w, http.StatusBadRequest, "id must be a number")
		return
	}
	err = h.storage.DeleteSong(idInt)
	if err != nil {
		h.log.Error("failed to delete song:", er.Err(err))
		respondWithError(w, http.StatusInternalServerError, "failed to delete song")
		return
	}
	respondWithJSON(w, http.StatusOK, models.Response{Status: http.StatusOK, Message: "song deleted" + id})
	h.log.Info("song deleted: " + id)
}
