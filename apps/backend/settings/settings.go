package settings

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

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
	if sets != nil {
		r := RealDebridSettings{}
		if err := sets.UnmarshalJSONField("real_debrid", &r); err == nil {
			return &r
		}
	}
	return nil
}

func (s *Settings) getSettings() *models.Record {
	sets := []*models.Record{}
	s.app.Dao().RecordQuery("settings").All(&sets)
	if len(sets) > 0 {
		return sets[0]
	}
	return nil
}
