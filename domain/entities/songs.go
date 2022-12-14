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
	Duration int     `json:"duration"`
	Album    string  `json:"album"`
	Artwork  string  `json:"artwork"`
	Price    float64 `json:"price"`
	Origin   string  `json:"origin"`
}
