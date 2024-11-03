package trakt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/odin-movieshow/server/helpers"
	"github.com/odin-movieshow/server/settings"
	"github.com/odin-movieshow/server/tmdb"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/thoas/go-funk"

	"github.com/go-resty/resty/v2"

	"github.com/charmbracelet/log"
)

const (
	TRAKT_URL = "https://api.trakt.tv"
)

// Removes slice element at index(s) and returns new slice
func remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

var Headers = make(map[string]any)
var FetchTMDB = true
var FetchSeasons = true

func normalize(objmap []any, endpoint string) {
	obj := "movie"
	if strings.Contains(endpoint, "/show") {
		obj = "show"
	}

	for i, o := range objmap {
		if objmap[i].(map[string]any)["movie"] == nil && objmap[i].(map[string]any)["show"] == nil {
			m := make(map[string]any)
			m[obj] = o
			objmap[i] = m
		}
	}

}

func removeDuplicates(objmap []any) []any {
	showsSeen := make([]float64, 0)
	toRemove := make([]int, 0)
	for i, o := range objmap {
		id := o.(map[string]any)["show"].(map[string]any)["ids"].(map[string]any)["trakt"].(float64)
		if !funk.ContainsFloat64(showsSeen, id) {
			showsSeen = append(showsSeen, id)
		} else {
			toRemove = append(toRemove, i)
		}
	}

	newmap := make([]any, 0)

	for i, o := range objmap {
		if !funk.ContainsInt(toRemove, i) {
			newmap = append(newmap, o)
		}
	}

	return newmap

}

func removeWatched(objmap []any) []any {
	toRemove := make([]int, 0)
	for i, o := range objmap {
		if o.(map[string]any)["episode"].(map[string]any)["watched"] != nil && o.(map[string]any)["episode"].(map[string]any)["watched"].(bool) == true {
			toRemove = append(toRemove, i)
		}

	}

	newmap := make([]any, 0)

	for i, o := range objmap {
		if !funk.ContainsInt(toRemove, i) {
			newmap = append(newmap, o)
		}
	}

	return newmap

}

func removeSeason0(objmap []any) []any {
	toKeep := []any{}
	for _, o := range objmap {

		if o.(map[string]any)["episode"] != nil && o.(map[string]any)["episode"].(map[string]any)["season"] != nil && o.(map[string]any)["episode"].(map[string]any)["season"].(float64) > 0 {
			toKeep = append(toKeep, o)
		}
	}

	return toKeep

}

func SyncHistory(app *pocketbase.PocketBase) {
	users := []*models.Record{}
	app.Dao().RecordQuery("users").All(&users)
	FetchTMDB = false
	FetchSeasons = false
	var wg sync.WaitGroup
	for _, u := range users {
		records, _ := app.Dao().FindRecordsByFilter("history", "user = {:user}", "-watched_at", 1, 0, dbx.Params{"user": u.Get("id")})
		last_watched := types.DateTime{}
		if len(records) > 0 {
			last_watched = records[0].GetDateTime("watched_at")
		}

		t := make(map[string]any)
		if err := u.UnmarshalJSONField("trakt_token", &t); err != nil {
			continue
		}

		Headers["authorization"] = "Bearer " + t["access_token"].(string)

		wg.Add(1)
		go syncByType(&wg, "movies", last_watched, app, u.Get("id").(string))
		wg.Add(1)
		go syncByType(&wg, "episodes", last_watched, app, u.Get("id").(string))

		wg.Wait()
		log.Info("Done synching trakt history", "user", u.Get("id"))

	}
}

func syncByType(wg *sync.WaitGroup, t string, last_history types.DateTime, app *pocketbase.PocketBase, user string) {
	defer wg.Done()
	limit := 100
	url := "/sync/history/" + t + "?limit=" + fmt.Sprint(limit)
	collection, _ := app.Dao().FindCollectionByNameOrId("history")
	if !last_history.IsZero() {
		url += "&start_at=" + last_history.Time().Add(time.Second*1).Format(time.RFC3339)
	}
	_, headers, _ := CallEndpoint(url, "GET", nil, false, app)
	pages, _ := strconv.Atoi(headers.Get("X-Pagination-Page-Count"))

	for i := 1; i <= pages; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {

			defer wg.Done()
			url += "&page=" + fmt.Sprint(i)

			data, _, _ := CallEndpoint(url, "GET", nil, false, app)
			log.Debug("Synching trakt history", "type", t, "page", fmt.Sprintf("%d/%d", i, pages), "user", user, "count", len(data.([]any)))
			for _, o := range data.([]any) {

				record := models.NewRecord(collection)
				record.Set("watched_at", o.(map[string]any)["watched_at"])
				record.Set("user", user)
				if t == "movies" {
					record.Set("type", "movie")

					record.Set("trakt_id", o.(map[string]any)["movie"].(map[string]any)["ids"].(map[string]any)["trakt"])
					record.Set("data", map[string]any{"genres": o.(map[string]any)["movie"].(map[string]any)["genres"]})
					record.Set("runtime", o.(map[string]any)["movie"].(map[string]any)["runtime"])

				} else if t == "episodes" {
					record.Set("type", "episode")
					record.Set("trakt_id", o.(map[string]any)["episode"].(map[string]any)["ids"].(map[string]any)["trakt"])
					record.Set("show_id", o.(map[string]any)["show"].(map[string]any)["ids"].(map[string]any)["trakt"])
					record.Set("data", map[string]any{"genres": o.(map[string]any)["show"].(map[string]any)["genres"]})
					record.Set("runtime", o.(map[string]any)["episode"].(map[string]any)["runtime"])
				}
				app.Dao().SaveRecord(record)

			}

		}(i, wg)
	}

}

func RefreshTokens(app *pocketbase.PocketBase) {
	records := []*models.Record{}
	app.Dao().RecordQuery("users").All(&records)

	for _, r := range records {
		t := make(map[string]any)
		if err := r.UnmarshalJSONField("trakt_token", &t); err == nil {
			data, _, status := CallEndpoint("/oauth/token", "POST", map[string]any{"grant_type": "refresh_token", "client_id": settings.GetTrakt(app).ClientId, "client_secret": settings.GetTrakt(app).ClientSecret, "code": t["device_code"], "refresh_token": t["refresh_token"]}, false, app)
			if status < 300 && data != nil {
				data.(map[string]any)["device_code"] = t["device_code"]
				r.Set("trakt_token", data)
				app.Dao().Save(r)
				log.Info("trakt refresh token", "user", r.Get("id"))
			}
		}

	}

}

func CallEndpoint(endpoint string, method string, body map[string]any, donorm bool, app *pocketbase.PocketBase) (any, http.Header, int) {

	var objmap any

	request := resty.New().SetRetryCount(3).SetRetryWaitTime(time.Second * 3).R()
	request.SetHeader("trakt-api-version", "2").SetHeader("content-type", "application/json").SetHeader("trakt-api-key", settings.GetTrakt(app).ClientId)
	var respHeaders http.Header
	var status = 200
	for k, v := range Headers {

		if funk.Contains([]string{"Host", "Connection"}, k) {
			continue
		}
		request.SetHeader(k, v.(string))
	}
	if body != nil {
		request.SetBody(body)
	}
	request.Attempt = 3
	var r func(url string) (*resty.Response, error)
	switch method {
	case "POST":
		r = request.Post
	case "PATCH":
		r = request.Patch
	case "PUT":
		r = request.Put
	case "DELETE":
		r = request.Delete
	default:
		r = request.Get

	}

	if !strings.Contains(endpoint, "oauth") {
		if !strings.Contains(endpoint, "extended=") {
			if strings.Contains(endpoint, "?") {
				endpoint += "&"
			} else {
				endpoint += "?"
			}
			endpoint += "extended=full&limit=30"
		}
	}

	// var listError error
	if resp, err := r(fmt.Sprintf("%s%s", TRAKT_URL, endpoint)); err == nil {
		respHeaders = resp.Header()
		status = resp.StatusCode()
		log.Debug("trakt fetch", "url", endpoint, "method", method, "status", status, "body", body, "headers", Headers)
		if status > 299 {
			log.Error("trakt", "fetch", endpoint, "status", status, "res", string(resp.Body()), "body", body, "headers", respHeaders)
		}
		err := json.Unmarshal(resp.Body(), &objmap)
		if err != nil {
			log.Error("trakt", "unmarshal", err)
		}

		switch objmap.(type) {

		case []any:

			if len(objmap.([]any)) == 0 {
				return objmap, respHeaders, status
			}

			if donorm {
				normalize(objmap.([]any), endpoint)
			}

			if (objmap.([]any)[0].(map[string]any)["movie"] != nil || objmap.([]any)[0].(map[string]any)["show"] != nil) && !strings.Contains(endpoint, "sync/history") {
				objmap = GetWatched(objmap.([]any), app)
				if strings.Contains(endpoint, "calendars") {
					objmap = removeSeason0(objmap.([]any))
					objmap = removeWatched(objmap.([]any))
					objmap = removeDuplicates(objmap.([]any))
				}

				// if !strings.Contains(endpoint, "/history") {
				var wg sync.WaitGroup
				var mux sync.Mutex
				if FetchTMDB {
					getTMDB(&wg, &mux, objmap.([]any), app)
				}
				if FetchSeasons {
					// getSeasons(&wg, &mux, objmap.([]any), app)
				}
				wg.Wait()
				objmap = GetWatched(objmap.([]any), app)
				if !strings.Contains(endpoint, "sync/history") {
					objmap = FixEpisodes(objmap)

				}
				// }

			}
		default:

		}

		// }
	} else {
		log.Error("trakt", "endpoint", endpoint, "body", body, "err", err)
	}

	return objmap, respHeaders, status
}

// fixes calendar episodes

func FixEpisodes(result any) []any {
	if funk.IsCollection(result) {

		result = funk.Map(result, func(m any) any {
			var item = ""
			for _, k := range []string{"episode", "movie", "show"} {
				if m.(map[string]any)[k] != nil {
					item = k
					break
				}
			}
			if item == "" {
				return m
			}
			newM := m.(map[string]any)[item].(map[string]any)
			newM["type"] = item
			if item == "episode" {
				if m.(map[string]any)["show"] != nil && newM["show"] == nil {
					newM["show"] = m.(map[string]any)["show"]
				}
				if newM["show"] != nil {
					newM["tmdb"] = newM["show"].(map[string]any)["tmdb"]
				}
			}
			return newM
		})

	}

	return result.([]any)

}

func getTMDB(wg *sync.WaitGroup, mux *sync.Mutex, objmap []any, app *pocketbase.PocketBase) {
	for k := range objmap {
		wg.Add(1)
		go tmdb.PopulateTMDB(k, wg, mux, objmap, app)
	}
}

type Watched struct {
	Plays         int       `json:"plays"`
	LastWatchedAt time.Time `json:"last_watched_at"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
	Movie         struct {
		Title string `json:"title"`
		Year  int    `json:"year"`
		Ids   struct {
			Trakt float64 `json:"trakt"`
			Slug  string  `json:"slug"`
			Imdb  string  `json:"imdb"`
			Tmdb  int     `json:"tmdb"`
		} `json:"ids"`
	} `json:"movie"`

	Show struct {
		Title string `json:"title"`
		Year  int    `json:"year"`
		Ids   struct {
			Trakt  float64 `json:"trakt"`
			Slug   string  `json:"slug"`
			Tvdb   int     `json:"tvdb"`
			Imdb   string  `json:"imdb"`
			Tmdb   int     `json:"tmdb"`
			Tvrage any     `json:"tvrage"`
		} `json:"ids"`
	} `json:"show"`
	Seasons []struct {
		Number   int `json:"number"`
		Episodes []struct {
			Number        int       `json:"number"`
			Plays         int       `json:"plays"`
			LastWatchedAt time.Time `json:"last_watched_at"`
		} `json:"episodes"`
	} `json:"seasons"`
}

func GetWatched(objmap []any, app *pocketbase.PocketBase) []any {

	if len(objmap) == 0 {
		return objmap
	}
	if objmap[0].(map[string]any)["show"] != nil {
		return GetWatchedCalendarEpisodes(objmap, app)
	}

	return GetWatchedMovies(objmap, app)

}

func getHistory(app *pocketbase.PocketBase, htype string) []any {
	records, _ := app.Dao().FindRecordsByFilter("history", "type = {:htype}", "-watched_at", -1, 0, dbx.Params{"htype": htype})
	data := make([]any, 0)
	for _, r := range records {
		item := make(map[string]any)
		item["type"] = r.Get("type")
		item["trakt_id"] = r.Get("trakt_id")
		data = append(data, item)
	}
	return data
}

func GetWatchedCalendarEpisodes(objmap []any, app *pocketbase.PocketBase) []any {
	history := getHistory(app, "episode")
	for i := range objmap {
		tvshow := objmap[i].(map[string]any)["show"].(map[string]any)
		episode := objmap[i].(map[string]any)["episode"]
		if tvshow != nil && episode != nil {
			episode.(map[string]any)["watched"] = false
			tvshow["watched"] = false
			for _, h := range history {
				if h.(map[string]any)["trakt_id"] == episode.(map[string]any)["ids"].(map[string]any)["trakt"] {
					episode.(map[string]any)["watched"] = true
					tvshow["watched"] = true
					break
				}
			}

		}
	}
	newmap := make([]any, 0)

	for _, o := range objmap {
		newmap = append(newmap, o)
	}

	return newmap
}

func GetWatchedMovies(objmap []any, app *pocketbase.PocketBase) []any {
	history := getHistory(app, "movie")
	for i, o := range objmap {
		objmap[i].(map[string]any)["movie"].(map[string]any)["watched"] = false
		for _, h := range history {
			if h.(map[string]any)["trakt_id"] == o.(map[string]any)["movie"].(map[string]any)["ids"].(map[string]any)["trakt"] {
				objmap[i].(map[string]any)["movie"].(map[string]any)["watched"] = true
				break
			}
		}

	}
	newmap := make([]any, 0)

	for _, o := range objmap {
		newmap = append(newmap, o)
	}

	return newmap
}

func GetWatchedSeasonEpisodes(objmap []any, app *pocketbase.PocketBase) []any {
	history := getHistory(app, "episode")
	for _, oseason := range objmap {
		for _, oepisode := range oseason.(map[string]any)["episodes"].([]any) {
			oepisode.(map[string]any)["watched"] = false
			for _, h := range history {
				if h.(map[string]any)["trakt_id"] == oepisode.(map[string]any)["ids"].(map[string]any)["trakt"] {
					oepisode.(map[string]any)["watched"] = true
					break
				}
			}
		}
	}
	newmap := make([]any, 0)
	for _, o := range objmap {
		newmap = append(newmap, o)
	}
	return newmap

}

func GetSeasons(app *pocketbase.PocketBase, id int) any {

	endpoint := fmt.Sprintf("/shows/%d/seasons?extended=full,episodes", id)

	cache := helpers.ReadTraktSeasonCache(app, uint(id))
	if cache != nil {
		return GetWatchedSeasonEpisodes(cache, app)
	}

	result, _, _ := CallEndpoint(endpoint, "GET", nil, false, app)

	helpers.WriteTraktSeasonCache(app, uint(id), &result)

	return GetWatchedSeasonEpisodes(result.([]any), app)
}
