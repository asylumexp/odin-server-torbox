package tmdb

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/odin-movieshow/backend/helpers"
	"github.com/odin-movieshow/backend/settings"
	"github.com/odin-movieshow/backend/types"
	"github.com/thoas/go-funk"

	"github.com/charmbracelet/log"

	resty "github.com/go-resty/resty/v2"
)

type Tmdb struct {
	settings *settings.Settings
	helpers  *helpers.Helpers
}

func New(settings *settings.Settings, helpers *helpers.Helpers) *Tmdb {
	return &Tmdb{settings: settings, helpers: helpers}
}

const (
	TMDB_URL = "https://api.themoviedb.org/3"
)

func (t *Tmdb) PopulateTMDB2(
	k int,
	mux *sync.Mutex,
	objmap []types.TraktItem,
) {
	tmdbKey := os.Getenv("TMDB_KEY")
	resource := "movie"
	tmdbResource := "movie"
	if objmap[k].Type == "show" || objmap[k].Type == "episode" && objmap[k].Show != nil {
		resource = "show"
		tmdbResource = "tv"
	}
	if objmap[k].Season > 0 {
		resource = "season"
		tmdbResource = "tv"
	}
	id := objmap[k].IDs.Tmdb

	if objmap[k].Show != nil {
		id = objmap[k].Show.IDs.Tmdb
	}
	if objmap[k].IDs.Tmdb == 0 {
		return
	}
	cache := t.helpers.ReadTmdbCache(id, resource)
	if cache != nil {
		objmap[k].Tmdb = cache
		return
	}
	request := resty.New().
		SetRetryCount(3).
		SetTimeout(time.Second * 30).
		SetRetryWaitTime(time.Second).
		R()

	url := fmt.Sprintf("%s/%s/%d?api_key=%s&append_to_response=credits,videos,images", TMDB_URL, tmdbResource, id, tmdbKey)
	// log.Debug(url)
	var tmdbObj interface{}
	if _, err := request.SetResult(&tmdbObj).SetHeader("content-type", "application/json").Get(url); err == nil {

		if tmdbObj == nil {
			return
		}
		tmdb := objToTmdb(tmdbObj)

		if (tmdb).Credits != nil {
			if (tmdb).Credits.Crew != nil {
				(tmdb).Credits.Crew = nil
			}
			if (tmdb).Credits.Cast != nil {
				// strip down cast
				cast := (*(tmdb).Credits.Cast)
				castlen := len((cast))
				if castlen > 15 {
					castlen = 15
				}
				castcut := cast[0:castlen]
				tmdb.Credits.Cast = &castcut
			}
		}
		//
		if tmdb.Images != nil &&
			tmdb.Images.Logos != nil {

			for _, l := range *tmdb.Images.Logos {
				if l.Iso_639_1 != nil && *(l.Iso_639_1) == "en" {
					tmdb.LogoPath = l.FilePath
					break
				}
			}
			if tmdb.LogoPath != "" && len(*tmdb.Images.Logos) > 0 {
				tmdb.LogoPath = (*tmdb.Images.Logos)[0].FilePath
			}
			tmdb.Images = nil
		} else {
			tmdb.LogoPath = ""
		}
		tmdbObj := tmdbToObj(tmdb)
		t.helpers.WriteTmdbCache(id, resource, &tmdbObj)
		objmap[k].Tmdb = tmdbObj
	} else {
		log.Error("TMDB", "Response", err)
	}
}

func objToTmdb(obj any) *types.TmdbItem {
	tmdb := types.TmdbItem{}
	ms, err := json.Marshal(obj)
	if err != nil {
		log.Error(err)
		return nil
	}
	err = json.Unmarshal(ms, &tmdb)
	if err != nil {
		log.Error(err)
		return nil
	}
	tmdb.Original = &obj
	return &tmdb
}

func tmdbToObj(tmdb *types.TmdbItem) any {
	var obj any
	ms, err := json.Marshal(tmdb)
	if err == nil {
		err = json.Unmarshal(ms, &obj)
		if err == nil {
			orig := *(tmdb.Original)
			if (orig) == nil {
				return obj
			}
			for k, v := range (orig).(map[string]interface{}) {
				if funk.Contains([]string{"images", "credits"}, k) {
					continue
				}
				obj.(map[string]any)[k] = v
			}
			obj.(map[string]any)["original"] = nil
		}
	}
	return obj
}

func (t *Tmdb) PopulateTMDB(
	k int,
	wg *sync.WaitGroup,
	mux *sync.Mutex,
	objmap []any,
) {
	defer wg.Done()
	// defer mux.Unlock()
	// mux.Lock()
	tmdbKey := os.Getenv("TMDB_KEY")
	resource := "movie"
	tmdbResource := "movie"
	if objmap[k].(map[string]any)["show"] != nil {
		resource = "show"
		tmdbResource = "tv"
	}
	if objmap[k].(map[string]any)["season"] != nil {
		resource = "season"
		tmdbResource = "tv"
	}
	var tmdb any
	if objmap[k].(map[string]any)[resource].(map[string]any)["ids"].(map[string]any)["tmdb"] == nil {
		return
	}
	id := uint(
		objmap[k].(map[string]any)[resource].(map[string]any)["ids"].(map[string]any)["tmdb"].(float64),
	)
	cache := t.helpers.ReadTmdbCache(id, resource)
	if cache != nil {
		objmap[k].(map[string]any)[resource].(map[string]any)["tmdb"] = cache
		return
	}
	request := resty.New().
		SetRetryCount(3).
		SetTimeout(time.Second * 30).
		SetRetryWaitTime(time.Second).
		R()
	if _, err := request.SetResult(&tmdb).SetHeader("content-type", "application/json").Get(fmt.Sprintf("%s/%s/%d?api_key=%s&append_to_response=credits,videos,images", TMDB_URL, tmdbResource, id, tmdbKey)); err == nil {
		log.Debug("tmdb", "resource", resource, "id", id)
		// remove crew
		if tmdb == nil {
			return
		}
		if tmdb.(map[string]any)["credits"] != nil &&
			tmdb.(map[string]any)["credits"].(map[string]any)["crew"] != nil {
			tmdb.(map[string]any)["credits"].(map[string]any)["crew"] = []any{}
		}
		if tmdb.(map[string]any)["credits"] != nil &&
			tmdb.(map[string]any)["credits"].(map[string]any)["cast"].([]any) != nil {
			// strip down cast
			cast := tmdb.(map[string]any)["credits"].(map[string]any)["cast"].([]any)
			castlen := len(cast)
			if castlen > 15 {
				castlen = 15
			}
			tmdb.(map[string]any)["credits"].(map[string]any)["cast"] = cast[0:castlen]
		}

		if tmdb.(map[string]any)["images"] != nil &&
			tmdb.(map[string]any)["images"].(map[string]any)["logos"] != nil {

			for _, l := range tmdb.(map[string]any)["images"].(map[string]any)["logos"].([]any) {
				if l.(map[string]any)["iso_639_1"] != nil &&
					l.(map[string]any)["iso_639_1"].(string) == "en" {
					tmdb.(map[string]any)["logo_path"] = l.(map[string]any)["file_path"]
					break
				}
			}
			tmdb.(map[string]any)["images"] = nil

			tmdb.(map[string]any)["images"] = nil
		} else {
			tmdb.(map[string]any)["logo_path"] = ""
		}

		t.helpers.WriteTmdbCache(id, resource, &tmdb)
		objmap[k].(map[string]any)[resource].(map[string]any)["tmdb"] = tmdb
		// helpers.WriteTMDBImage(tmdb.(map[string]any)["poster_path"].(string))
	} else {
		fmt.Println("TMDB", "Response", err)
	}
}

func (t *Tmdb) GetEpisodes(showId string, seasons []string) *[]any {
	// '/tv/$showId/season/$season?api_key=$key'
	tsets := t.settings.GetTmdb()
	tmdbKey := tsets.Key
	var wg sync.WaitGroup
	res := make([]any, 0)

	for _, s := range seasons {
		wg.Add(1)
		go func(s string) {
			endpoint := fmt.Sprintf("%s/tv/%s/season/%s?api_key=%s", TMDB_URL, showId, s, tmdbKey)

			defer wg.Done()
			var obj any
			request := resty.New().
				SetRetryCount(3).
				SetTimeout(time.Second * 30).
				SetRetryWaitTime(time.Second).
				R()
			if _, err := request.SetResult(&obj).SetHeader("content-type", "application/json").Get(endpoint); err == nil {
				log.Info("tmdb episodes", "show", showId, "season", s)
				res = append(res, obj)
			} else {
				log.Error("tmdb episodes", "error", err)
			}
		}(s)
	}
	wg.Wait()

	return &res
}
