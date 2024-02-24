package types

type Torrent struct {
	Scraper      string           `json:"scraper"`
	Hash         string           `json:"hash"`
	Size         uint64           `json:"size"`
	ReleaseTitle string           `json:"release_title"`
	Magnet       string           `json:"magnet"`
	Url          string           `json:"url"`
	Name         string           `json:"name"`
	Quality      string           `json:"quality"`
	Info         []string         `json:"info"`
	Seeds        uint64           `json:"seeds"`
	RealDebrid   []map[string]any `json:"realdebrid"`
}
