package entities

type SearchParams struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
}

type Song struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Artist   string  `json:"artist"`
	Duration int     `json:"duration,omitempty"`
	Album    string  `json:"album,omitempty"`
	Artwork  string  `json:"artwork,omitempty"`
	Price    float64 `json:"price,omitempty"`
	Origin   string  `json:"origin"`
}
