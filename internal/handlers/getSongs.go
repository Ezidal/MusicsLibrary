package handlers

import (
	er "LibMusic/internal/logger/err"
	"LibMusic/internal/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// GetSongs godoc
// @Summary Получить список песен
// @Description Возвращает список песен с возможностью фильтрации и пагинации
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Фильтр по названию группы"
// @Param song query string false "Фильтр по названию песни"
// @Param releaseDate query string false "Фильтр по дате выпуска (формат: YYYY-MM-DD)"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество песен на странице" default(10)
// @Success 200 {array} models.Song
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /songs [get]
func (h *Handler) GetSongs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit < 1 {
		limit = 10
	}

	group := query.Get("group")
	song := query.Get("song")
	releaseDate := query.Get("releaseDate")

	offset := (page - 1) * limit

	baseQuery := "SELECT id, group_name, song_name, release_date, text, link FROM songs"
	filters := []string{}
	args := []any{}

	if group != "" {
		filters = append(filters, "group_name = $"+strconv.Itoa(len(args)+1))
		args = append(args, group)
	}
	if song != "" {
		filters = append(filters, "song_name = $"+strconv.Itoa(len(args)+1))
		args = append(args, song)
	}
	if releaseDate != "" {
		filters = append(filters, "release_date = $"+strconv.Itoa(len(args)+1))
		args = append(args, releaseDate)
	}

	finalQuery := baseQuery
	if len(filters) > 0 {
		finalQuery += " WHERE " + strings.Join(filters, " AND ")
	}
	finalQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := h.storage.Custom(finalQuery, args...)
	if err != nil {
		h.log.Error("Error getting songs", er.Err(err))
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	defer rows.Close()

	songs := []models.Song{}
	for rows.Next() {
		var song models.Song
		err := rows.Scan(
			&song.ID,
			&song.Group,
			&song.SongName,
			&song.ReleaseDate,
			&song.Text,
			&song.Link,
		)
		if err != nil {
			h.log.Error("Error scanning song", er.Err(err))
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		songs = append(songs, song)
	}
	if len(songs) == 0 {
		h.log.Info("No songs found")
		respondWithError(w, http.StatusNotFound, "No songs found")
		return
	}
	respondWithJSON(w, http.StatusOK, songs)
	h.log.Info("Songs sent")
}
