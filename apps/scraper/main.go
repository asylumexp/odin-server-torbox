package main

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/thoas/go-funk"
)

type Torrent struct {
	Scraper      string   `json:"scraper"`
	Hash         string   `json:"hash"`
	Size         uint64   `json:"size"`
	ReleaseTitle string   `json:"release_title"`
	Magnet       string   `json:"magnet"`
	Url          string   `json:"url"`
	Name         string   `json:"name"`
	Quality      string   `json:"quality"`
	Info         []string `json:"info"`
	Seeds        uint64   `json:"seeds"`
}

type Indexer struct {
	Id    string `xml:"id,attr"`
	Title string `xml:"title"`
	Caps  struct {
		Searching struct {
			Search struct {
				Available       string `xml:"available,attr"`
				SupportedParams string `xml:"supportedParams,attr"`
			} `xml:"search"`
			MovieSearch struct {
				Available       string `xml:"available,attr"`
				SupportedParams string `xml:"supportedParams,attr"`
			} `xml:"movie-search"`
			TvSearch struct {
				Available       string `xml:"available,attr"`
				SupportedParams string `xml:"supportedParams,attr"`
			} `xml:"tv-search"`
		} `xml:"searching"`
		Categories struct {
			Category []struct {
				ID   string `xml:"id,attr"`
				Name string `xml:"name,attr"`
			} `xml:"category"`
		} `xml:"categories"`
	} `xml:"caps"`
}

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		XMLName     xml.Name `xml:"channel"`
		Link        string   `xml:"link"`
		Title       string   `xml:"title"`
		Description string   `xml:"description"`
		Language    string   `xml:"language"`
		Category    string   `xml:"category"`
		Items       []struct {
			XMLName   xml.Name `xml:"item"`
			Title     string   `xml:"title"`
			Size      int64    `xml:"size"`
			Enclosure struct {
				URL    string `xml:"url,attr"`
				Length int64  `xml:"length,attr"`
				Type   string `xml:"type,attr"`
			} `xml:"enclosure"`
			Attrs []struct {
				Name  string `xml:"name,attr"`
				Value string `xml:"value,attr"`
			} `xml:"http://torznab.com/schemas/2015/feed attr"`
		} `xml:"item"`
	} `xml:"channel"`
}

type Payload struct {
	Type          string `json:"type"`
	Title         string `json:"title"`
	Year          string `json:"year"`
	Imdb          string `json:"imdb"`
	Trakt         string `json:"trakt"`
	ShowImdb      string `json:"show_imdb"`
	ShowTvdb      string `json:"show_tvdb"`
	ShowTitle     string `json:"show_title"`
	ShowYear      string `json:"show_year"`
	SeasonNumber  string `json:"season_number"`
	EpisodeImdb   string `json:"episode_imdb"`
	EpisodeTvdb   string `json:"episode_tvdb"`
	EpisodeTitle  string `json:"episode_title"`
	EpisodeNumber string `json:"episode_number"`
}

func main() {
	log.SetLevel(log.DebugLevel)
	jackettUrl := os.Getenv("JACKETT_URL")
	jackettKey := os.Getenv("JACKETT_KEY")

	if jackettUrl == "" || jackettKey == "" {
		log.Error("missing env vars JACKETT_URL and JACKETT_KEY")
		os.Exit(0)
	}
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("odin-craper is up and running!")
	})
	app.Post("/scrape", func(c *fiber.Ctx) error {
		payload := Payload{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		payload.Title = url.QueryEscape(payload.Title)
		payload.EpisodeTitle = url.QueryEscape(payload.EpisodeTitle)
		payload.ShowTitle = url.QueryEscape(payload.ShowTitle)

		log.Info(payload)

		indexers := getIndexerList(payload)
		allTorrents := []Torrent{}

		wg := sync.WaitGroup{}
		for _, indexer := range indexers {
			wg.Add(1)
			go func(indexer Indexer) {
				defer wg.Done()
				t1 := time.Now()
				ts := getTorrents(indexer, payload)
				t2 := time.Now()
				if len(ts) > 0 {
					log.Debug(
						"indexer",
						"id",
						indexer.Id,
						"torrents",
						len(ts),
						"took",
						fmt.Sprintf("%.1fs", t2.Sub(t1).Seconds()),
					)
				}
				allTorrents = append(allTorrents, ts...)
			}(indexer)
		}

		wg.Wait()
		dedupe := dedupe(allTorrents)
		filtered := separateByQuality(dedupe, payload)
		log.Info("torrents", "total", len(allTorrents), "dedupe", len(filtered))
		return c.JSON(filtered)
	})

	app.Listen(":6969")
}

func separateByQuality(torrents []Torrent, payload Payload) []Torrent {
	res := map[string][]Torrent{}

	regexpPatterns := []*regexp.Regexp{
		regexp.MustCompile(
			fmt.Sprintf("s0?%s[.x]?e0?%s", payload.SeasonNumber, payload.EpisodeNumber),
		),
		regexp.MustCompile(
			fmt.Sprintf("Season 0?%s,? ?Episode 0?%s", payload.SeasonNumber, payload.EpisodeNumber),
		),
	}

	for _, q := range []string{"4K", "1080p", "720p", "SD", "CAM"} {
		res[q] = []Torrent{}
	}

	for _, t := range torrents {
		if _, ok := res[t.Quality]; !ok {
			res[t.Quality] = []Torrent{}
		}

		// sort SxEx episodes first
		if payload.Type == "episode" {
			title := strings.ToLower(t.ReleaseTitle)
			for _, pattern := range regexpPatterns {
				if pattern.MatchString(title) {
					res[t.Quality] = append([]Torrent{t}, res[t.Quality]...)
					break
				} else {
					log.Debug(title)
					res[t.Quality] = append(res[t.Quality], t)
				}
			}
		} else {
			res[t.Quality] = append(res[t.Quality], t)
		}
	}

	// if len(res["4K"]) > 20 {
	// 	res["4K"] = res["4K"][:20]
	// }

	// if len(res["1080p"]) > 20 {
	// 	res["1080p"] = res["1080p"][:20]
	// }

	if len(res["720p"]) > 10 {
		res["720p"] = res["720p"][:10]
	}

	if len(res["SD"]) > 10 {
		res["SD"] = res["SD"][:10]
	}

	if len(res["4K"])+len(res["1080p"]) > 30 {

		res["720p"] = []Torrent{}
		res["SD"] = []Torrent{}
	}

	ret := append(res["4K"], res["1080p"]...)
	ret = append(ret, res["720p"]...)
	ret = append(ret, res["SD"]...)
	return ret

}

func getIndexerList(payload Payload) []Indexer {
	cat := "Movies"

	if payload.Type == "episode" {
		cat = "TV"
	}
	var indexers struct {
		Indexers []Indexer `xml:"indexer"`
	}

	jackettUrl := os.Getenv("JACKETT_URL")
	jackettKey := os.Getenv("JACKETT_KEY")

	request := resty.New().
		SetRetryCount(3).
		SetTimeout(time.Second * 30).
		SetRetryWaitTime(time.Second * 2).
		R()
	resp, _ := request.Get(
		fmt.Sprintf(
			"%s/api/v2.0/indexers/type:public/results/torznab/api?apikey=%s&t=indexers&configured=true",
			jackettUrl,
			jackettKey,
		),
	)
	if err := xml.Unmarshal(resp.Body(), &indexers); err != nil {
		log.Error("indexers", err)
	}

	neededIndexers := []Indexer{}
	for _, indexer := range indexers.Indexers {

		for _, category := range indexer.Caps.Categories.Category {
			if category.Name == cat {
				neededIndexers = append(neededIndexers, indexer)
				break
			}
		}
	}
	log.Info("indexers", "cat", cat, "total", len(indexers.Indexers), "needed", len(neededIndexers))
	return neededIndexers
}

func getTorrents(indexer Indexer, payload Payload) []Torrent {

	var rss Rss
	t := "search"
	q := ""
	season := ""
	ep := ""
	traktid := ""
	imdbid := ""
	tvdbid := ""
	tmdbid := ""
	if indexer.Caps.Searching.Search.Available == "yes" {
		q = payload.Title + "+" + payload.Year
		if payload.Type == "episode" {
			q = payload.ShowTitle + "+S" + payload.SeasonNumber + "+E" + payload.EpisodeNumber
		}
	}

	if indexer.Caps.Searching.MovieSearch.Available == "yes" && payload.Type == "movie" {
		t = "movie"
		if strings.Contains(indexer.Caps.Searching.MovieSearch.SupportedParams, "imdbid") {
			imdbid = payload.Imdb
		}
		if strings.Contains(indexer.Caps.Searching.MovieSearch.SupportedParams, "traktid") {
			traktid = payload.Trakt
		}
	}
	if indexer.Caps.Searching.TvSearch.Available == "yes" && payload.Type == "episode" {
		t = "tvsearch"
		if strings.Contains(indexer.Caps.Searching.TvSearch.SupportedParams, "imdbid") {
			imdbid = payload.EpisodeImdb
		}
		if strings.Contains(indexer.Caps.Searching.TvSearch.SupportedParams, "tvdbid") {
			tvdbid = payload.EpisodeTvdb
		}
		if strings.Contains(indexer.Caps.Searching.TvSearch.SupportedParams, "season") {
			q = payload.ShowTitle
			season = payload.SeasonNumber
		}
		if strings.Contains(indexer.Caps.Searching.TvSearch.SupportedParams, "ep") {
			q = payload.ShowTitle
			ep = payload.EpisodeNumber
		}

		if strings.Contains(indexer.Caps.Searching.TvSearch.SupportedParams, "traktid") {
			traktid = payload.Trakt
		}

	}

	query := ""
	if q != "" {
		query = "&q=" + q
	}
	if imdbid != "" {
		query = query + "&imdbid=" + imdbid
	}

	if tvdbid != "" {
		query = query + "&tvdbid=" + tvdbid
	}

	if tmdbid != "" {
		query = query + "&tmdbid=" + tmdbid
	}

	if season != "" {
		query = query + "&season=" + season
	}

	if ep != "" {
		query = query + "&ep=" + ep
	}

	if traktid != "" {
		query = query + "&traktid=" + traktid
	}

	var torrents []Torrent

	request := resty.New().
		// SetRetryCount(1).
		SetTimeout(time.Second * 30).
		// SetRetryWaitTime(time.Second * 2).
		R()

	jackettUrl := os.Getenv("JACKETT_URL")
	jackettKey := os.Getenv("JACKETT_KEY")

	url := fmt.Sprintf(
		"%s/api/v2.0/indexers/%s/results/torznab/api?apikey=%s&t=%s%s",
		jackettUrl,
		indexer.Id,
		jackettKey,
		t,
		query,
	)
	resp, err := request.Get(url)

	if err != nil {
		log.Error("get request", "indexer", indexer.Id, "error", err.Error())
		return torrents
	}

	if err := xml.Unmarshal(resp.Body(), &rss); err != nil {
		if err.Error() != "EOF" {
			log.Error("xml unmarshall", "indexer", indexer.Id, "error", err.Error())
		}
		return torrents
	}

	for _, item := range rss.Channel.Items {
		t := Torrent{}
		t.Scraper = indexer.Id
		t.Size = uint64(item.Size)
		t.ReleaseTitle = item.Title
		infos, q := getInfos(item.Title)
		t.Info = infos
		t.Quality = q
		t.Name = item.Title
		for _, attr := range item.Attrs {
			if attr.Name == "magneturl" {
				t.Magnet = attr.Value
				t.Url = attr.Value
			}
			if attr.Name == "infohash" {
				t.Hash = attr.Value
			}
			if attr.Name == "seeders" {
				c, err := strconv.Atoi(attr.Value)
				if err == nil {
					t.Seeds = uint64(c)
				}
			}
		}
		torrents = append(torrents, t)
	}

	return torrents
}

func dedupe(torrents []Torrent) []Torrent {
	res := []Torrent{}
	hashes := []string{}
	for _, t := range torrents {
		if t.Url != "" && !funk.ContainsString(hashes, t.Hash) {
			res = append(res, t)
			hashes = append(hashes, t.Hash)
		}
	}
	return res
}

func getInfos(title string) ([]string, string) {

	title = strings.ToLower(title)

	res := []string{}
	quality := "SD"
	infoTypes := map[string][]string{
		"AVC":   {"x264", "x 264", "h264", "h 264", "avc"},
		"HEVC":  {"x265", "x 265", "h265", "h 265", "hevc"},
		"XVID":  {"xvid"},
		"DIVX":  {"divx"},
		"MP4":   {"mp4"},
		"WMV":   {"wmv"},
		"MPEG":  {"mpeg"},
		"4K":    {"4k", "2160p", "216o"},
		"1080p": {"1080p", "1o80", "108o", "1o8p"},
		"720p":  {"720", "72o"},
		"REMUX": {"remux", "bdremux"},
		"DV":    {" dv ", "dovi", "dolby vision", "dolbyvision"},
		"HDR": {
			" hdr ",
			"hdr10",
			"hdr 10",
			"uhd bluray 2160p",
			"uhd blu ray 2160p",
			"2160p uhd bluray",
			"2160p uhd blu ray",
			"2160p bluray hevc truehd",
			"2160p bluray hevc dts",
			"2160p bluray hevc lpcm",
			"2160p us bluray hevc truehd",
			"2160p us bluray hevc dts",
		},
		"SDR":      {" sdr"},
		"AAC":      {"aac"},
		"DTS-HDMA": {"hd ma", "hdma"},
		"DTS-HDHR": {"hd hr", "hdhr", "dts hr", "dtshr"},
		"DTS-X":    {"dtsx", " dts x"},
		"ATMOS":    {"atmos"},
		"TRUEHD":   {"truehd", "true hd"},
		"DD+":      {"ddp", "eac3", " e ac3", " e ac 3", "dd+", "digital plus", "digitalplus"},
		"DD": {
			" dd ",
			"dd2",
			"dd5",
			"dd7",
			" ac3",
			" ac 3",
			"dolby digital",
			"dolbydigital",
			"dolby5",
		},
		"MP3":    {"mp3"},
		"WMA":    {" wma"},
		"2.0":    {"2 0 ", "2 0ch", "2ch"},
		"5.1":    {"5 1 ", "5 1ch", "6ch"},
		"7.1":    {"7 1 ", "7 1ch", "8ch"},
		"BLURAY": {"bluray", "blu ray", "bdrip", "bd rip", "brrip", "br rip"},
		"WEB":    {" web ", "webrip", "webdl", "web rip", "web dl", "webmux"},
		"HD-RIP": {" hdrip", " hd rip"},
		"DVDRIP": {"dvdrip", "dvd rip"},
		"HDTV":   {"hdtv"},
		"PDTV":   {"pdtv"},
		"CAM": {
			" cam ", "camrip", "cam rip",
			"hdcam", "hd cam",
			" ts ", " ts1", " ts7",
			"hd ts", "hdts",
			"telesync",
			" tc ", " tc1", " tc7",
			"hd tc", "hdtc",
			"telecine",
			"xbet",
			"hcts", "hc ts",
			"hctc", "hc tc",
			"hqcam", "hq cam",
		},
		"SCR": {"scr ", "screener"},
		"HC": {
			"korsub", " kor ",
			" hc ", "hcsub", "hcts", "hctc", "hchdrip",
			"hardsub", "hard sub",
			"sub hard",
			"hardcode", "hard code",
			"vostfr", "vo stfr",
		},
		"3D": {" 3d"},
	}
	for baseInfo, infoType := range infoTypes {
		for _, info := range infoType {
			if strings.Contains(title, strings.ToLower(baseInfo)) {
				res = append(res, baseInfo)
				break
			}
			if strings.Contains(title, strings.ToLower(info)) {
				res = append(res, baseInfo)
				break
			}
		}
	}

	if funk.Contains(res, "SDR") && funk.Contains(res, "HDR") {
		res = funk.FilterString(res, func(s string) bool {
			return s != "SDR"
		})
	}

	if funk.Contains(res, "DD") && funk.Contains(res, "DD+") {
		res = funk.FilterString(res, func(s string) bool {
			return s != "DD"
		})
	}

	if funk.ContainsString([]string{"2160p", "remux"}, title) &&
		!funk.Contains(res, []string{"HDR", "SDR"}) {
		res = append(res, "HDR")
	}

	if funk.Contains(res, "4K") {
		quality = "4K"
	}
	if funk.Contains(res, "1080p") {
		quality = "1080p"
	}
	if funk.Contains(res, "720p") {
		quality = "720p"
	}

	return res, quality
}
