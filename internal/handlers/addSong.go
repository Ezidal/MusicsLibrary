package handlers

import (
	er "LibMusic/internal/logger/err"
	"LibMusic/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// AddSong godoc
// @Summary Добавить новую песню
// @Description Добавляет новую песню в библиотеку, используя данные из внешнего API
// @Tags songs
// @Accept json
// @Produce json
// @Param input body models.SongRequest true "Данные для добавления песни"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /songs [post]
func (h *Handler) AddSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	var songReq models.SongRequest
	var songResp models.DetailSong
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&songReq)
	if err != nil {
		h.log.Error("failed to decode song", er.Err(err))
		respondWithError(w, http.StatusBadRequest, "failed to decode song")
		return
	}

	// Check if song already exist
	if h.storage.SongExist(songReq.SongName) != nil {
		respondWithError(w, http.StatusBadRequest, "Song already exist")
		return
	}
	h.log.Debug("songReq: " + songReq.Group + " " + songReq.SongName)

	params := url.Values{}
	params.Add("group", songReq.Group)
	params.Add("song", songReq.SongName)

	fullURL := fmt.Sprintf("%s?%s", h.conf.ExternalApiUrl, params.Encode())
	resp, err := http.Get(fullURL)
	if err != nil {
		h.log.Error("failed to get song info", er.Err(err))
		respondWithError(w, http.StatusInternalServerError, "failed to get song info")
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.log.Error("failed to read response body", er.Err(err))
		respondWithError(w, http.StatusInternalServerError, "failed to read response body")
		return
	}
	// h.log.Debug("response body: " + string(body))
	err = json.Unmarshal(body, &songResp)
	if err != nil {
		h.log.Error("failed to unmarshal response body", er.Err(err))
		respondWithError(w, http.StatusInternalServerError, "failed to unmarshal response body")
		return
	}

	date, err := time.Parse("02.01.2006", songResp.ReleaseDate)
	if err != nil {
		h.log.Error("Failed to parse date:", er.Err(err))
	}
	song.ReleaseDate = date.Format("2006-01-02")
	song.Link = songResp.Link
	song.Text = songResp.Text
	song.Group = songReq.Group
	song.SongName = songReq.SongName

	id, err := h.storage.AddSong(song)
	if err != nil {
		h.log.Error("failed to add song", er.Err(err))
		respondWithError(w, http.StatusInternalServerError, "failed to add song")
		return
	}
	idStr := strconv.Itoa(int(id))
	respondWithJSON(w, http.StatusOK, models.Response{Status: http.StatusOK, Message: "Song added, id: " + idStr})
	h.log.Info("Song added, id: " + idStr)

}
