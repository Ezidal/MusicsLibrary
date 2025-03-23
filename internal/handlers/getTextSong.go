package handlers

import (
	er "LibMusic/internal/logger/err"
	"LibMusic/internal/models"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) GetSongText(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/songs/")
	parts := strings.Split(path, "/")
	id := parts[0]
	if id == "" {
		h.log.Error("Missing song ID")
		respondWithError(w, http.StatusBadRequest, "Missing song ID")
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.log.Error("Invalid song ID", er.Err(err))
		respondWithError(w, http.StatusBadRequest, "Invalid song ID")
		return
	}

	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit < 1 {
		limit = 2
	}

	text, err := h.storage.GetText(idInt)
	if err != nil {
		h.log.Error("Failed to get song text:", er.Err(err))
		respondWithError(w, http.StatusInternalServerError, "Failed to get song text")
		return
	}
	verses := strings.Split(text, "\n\n")
	total := len(verses)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginated := models.PaginatedVerses{
		Verses: verses[start:end],
		Page:   page,
		Limit:  limit,
		Total:  total,
	}

	respondWithJSON(w, http.StatusOK, paginated)
	h.log.Info("Get song text by ID: " + id)
}
