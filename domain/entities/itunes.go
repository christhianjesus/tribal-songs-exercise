package entities

type ITunesResponse struct {
	Results []ITunesSong `json:"results"`
}

type ITunesSong struct {
	ID       int     `json:"trackId"`
	Name     string  `json:"trackName"`
	Artist   string  `json:"artistName"`
	Duration int     `json:"trackTimeMillis"`
	Album    string  `json:"collectionName"`
	Artwork  string  `json:"artworkUrl100"`
	Price    float64 `json:"trackPrice"`
	Origin   string  `json:"-"`
}
