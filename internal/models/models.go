package models

type Song struct {
	Id          int    `json:"id"`
	Group       string `json:"group"`
	SongName    string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type SongRequest struct {
	Group    string `json:"group"`
	SongName string `json:"song"`
}

type Response struct {
	Error  string `json:"error",omitempty`
	Status int    `json:"status"`
}

type DetailSong struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
