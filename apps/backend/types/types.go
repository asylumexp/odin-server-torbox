package types

type Torrent struct {
	Scraper      string         `json:"scraper"`
	Hash         string         `json:"hash"`
	Size         uint64         `json:"size"`
	ReleaseTitle string         `json:"release_title"`
	Magnet       string         `json:"magnet"`
	Url          string         `json:"url"`
	Name         string         `json:"name"`
	Quality      string         `json:"quality"`
	Info         []string       `json:"info"`
	Seeds        uint64         `json:"seeds"`
	Links        []Unrestricted `json:"links"`
}

type Unrestricted struct {
	Filename string `json:"filename"`
	Filesize int    `json:"filesize"`
	Download string `json:"download"`
}
