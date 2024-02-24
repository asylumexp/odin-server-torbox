package settings

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

type TmdbSettings struct {
	Key string `json:"key"`
}

type TraktSettings struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type JackettSettings struct {
	Url string `json:"url"`
	Key string `json:"key"`
}

type RealDebridSettings struct {
	AccessToken  string `json:"access_token"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

func GetRealDebrid(app *pocketbase.PocketBase) *RealDebridSettings {
	s := getSettings(app)
	if s != nil {
		r := RealDebridSettings{}
		if err := s.UnmarshalJSONField("real_debrid", &r); err == nil {
			return &r
		}
	}
	return nil
}

func GetTrakt(app *pocketbase.PocketBase) *TraktSettings {
	s := getSettings(app)
	if s != nil {
		t := TraktSettings{}

		if err := s.UnmarshalJSONField("trakt", &t); err == nil {
			return &t
		}
	}

	return nil

}

func GetScraperUrl(app *pocketbase.PocketBase) string {
	s := getSettings(app)
	if s != nil {
		return s.Get("scraper_url").(string)
	}
	return ""

}

func GetTmdb(app *pocketbase.PocketBase) *TmdbSettings {
	s := getSettings(app)
	if s != nil {
		t := TmdbSettings{}
		if err := s.UnmarshalJSONField("tmdb", &t); err == nil {
			return &t
		}
	}
	return nil
}

func GetJackett(app *pocketbase.PocketBase) *JackettSettings {
	s := getSettings(app)
	if s != nil {
		j := JackettSettings{}
		if err := s.UnmarshalJSONField("jackett", &j); err == nil {
			return &j
		}

	}
	return nil
}

func getSettings(app *pocketbase.PocketBase) *models.Record {
	s := []*models.Record{}
	app.Dao().RecordQuery("settings").All(&s)
	if len(s) > 0 {
		return s[0]
	}
	return nil
}
