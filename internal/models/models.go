package models

type Song struct {
	ID          int    `json:"id"`
	Group       string `json:"group" example:"Muse"`
	SongName    string `json:"song" example:"Supermassive Black Hole"`
	ReleaseDate string `json:"releaseDate" example:"2006-07-16"`
	Text        string `json:"text" example:"Текст песни..."`
	Link        string `json:"link" example:"https://youtube.com/watch?v=..."`
}

type SongRequest struct {
	Group    string `json:"group"`
	SongName string `json:"song"`
}

type Response struct {
	Error   string `json:"error,omitempty"`
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
}

type DetailSong struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type PaginatedVerses struct {
	Verses []string `json:"verses"`
	Page   int      `json:"page"`
	Limit  int      `json:"limit"`
	Total  int      `json:"total"`
}
