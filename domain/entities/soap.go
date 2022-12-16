package entities

type SearchLyricResponse struct {
	Results []SearchLyricResult `xml:"SearchLyricResult>SearchLyricResult"`
}

type SearchLyricResult struct {
	ID       int     `xml:"TrackId"`
	Name     string  `xml:"Song"`
	Artist   string  `xml:"Artist"`
	Duration int     `xml:"-"`
	Album    string  `xml:"-"`
	Artwork  string  `xml:"-"`
	Price    float64 `xml:"-"`
	Origin   string  `xml:"-"`
}
