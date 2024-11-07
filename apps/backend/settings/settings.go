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

type Settings struct {
	app *pocketbase.PocketBase
}

func New(app *pocketbase.PocketBase) *Settings {
	return &Settings{app: app}
}

func (s *Settings) GetRealDebrid() *RealDebridSettings {
	sets := s.getSettings()
	if s != nil {
		r := RealDebridSettings{}
		if err := sets.UnmarshalJSONField("real_debrid", &r); err == nil {
			return &r
		}
	}
	return nil
}

func (s *Settings) GetTrakt() *TraktSettings {
	sets := s.getSettings()
	if s != nil {
		t := TraktSettings{}

		if err := sets.UnmarshalJSONField("trakt", &t); err == nil {
			return &t
		}
	}

	return nil

}

func (s *Settings) GetScraperUrl() string {
	sets := s.getSettings()
	if s != nil {
		return sets.Get("scraper_url").(string)
	}
	return ""

}

func (s *Settings) GetTmdb() *TmdbSettings {
	sets := s.getSettings()
	if s != nil {
		t := TmdbSettings{}
		if err := sets.UnmarshalJSONField("tmdb", &t); err == nil {
			return &t
		}
	}
	return nil
}

func (s *Settings) GetJackett() *JackettSettings {
	sets := s.getSettings()
	if sets != nil {
		j := JackettSettings{}
		if err := sets.UnmarshalJSONField("jackett", &j); err == nil {
			return &j
		}

	}
	return nil
}

func (s *Settings) getSettings() *models.Record {
	sets := []*models.Record{}
	s.app.Dao().RecordQuery("settings").All(&s)
	if len(sets) > 0 {
		return sets[0]
	}
	return nil
}
