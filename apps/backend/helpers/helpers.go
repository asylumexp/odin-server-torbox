package helpers

import (
	"fmt"
	"os/user"
	"time"

	"github.com/charmbracelet/log"
	"github.com/odin-movieshow/backend/types"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

type Helpers struct {
	app *pocketbase.PocketBase
}

func New(app *pocketbase.PocketBase) *Helpers {
	return &Helpers{app: app}
}

func (h *Helpers) GetHomeDir() string {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return currentUser.HomeDir
}

func (h *Helpers) ReadTmdbCache(id uint, resource string) interface{} {
	record, err := h.app.Dao().
		FindFirstRecordByFilter("tmdb", "tmdb_id = {:id} && type = {:type}", dbx.Params{"id": id, "type": resource})
	res := make(map[string]any)
	if err == nil {
		err := record.UnmarshalJSONField("data", &res)
		date := record.GetDateTime("updated")
		now := time.Now()
		if err == nil {
			if date.Time().Before(now.AddDate(0, 0, -1)) {
				return nil
			}
			log.Debug("cache hit", "for", "tmdb", "resource", resource, "id", id)

			return res
		}
	}
	return nil
}

func (h *Helpers) WriteTmdbCache(id uint, resource string, data *interface{}) {
	if data == nil {
		return
	}
	log.Info("cache write", "for", "tmdb", "resource", resource, "id", id)
	record, err := h.app.Dao().
		FindFirstRecordByFilter("tmdb", "tmdb_id = {:id} && type = {:type}", dbx.Params{"id": id, "type": resource})

	if err == nil {
		record.Set("data", &data)
		h.app.Dao().SaveRecord(record)
	} else {

		collection, _ := h.app.Dao().FindCollectionByNameOrId("tmdb")
		record := models.NewRecord(collection)
		record.Set("data", &data)
		record.Set("tmdb_id", id)
		record.Set("type", resource)
		h.app.Dao().SaveRecord(record)
	}
}

func (h *Helpers) WriteTraktSeasonCache(id uint, data *interface{}) {
	if data == nil {
		return
	}
	log.Info("cache write", "for", "trakt", "resource", "show_seasons", "id", id)
	record, err := h.app.Dao().
		FindFirstRecordByFilter("trakt_seasons", "trakt_id = {:id}", dbx.Params{"id": id})

	if err == nil {
		record.Set("data", &data)
		h.app.Dao().SaveRecord(record)
	} else {
		collection, _ := h.app.Dao().FindCollectionByNameOrId("trakt_seasons")
		record := models.NewRecord(collection)
		record.Set("data", &data)
		record.Set("trakt_id", id)
		h.app.Dao().SaveRecord(record)
	}
}

func (h *Helpers) ReadRDCache(resource string, magnet string) *types.Torrent {
	record, err := h.app.Dao().
		FindFirstRecordByFilter("rd_resolved", "magnet = {:magnet}", dbx.Params{"magnet": magnet})
	var res types.Torrent
	if err == nil {
		err := record.UnmarshalJSONField("data", &res)
		date := record.GetDateTime("updated")
		now := time.Now().Add(time.Duration((-8) * time.Hour))
		if err == nil {
			if date.Time().Before(now) {
				return nil
			}
			log.Debug("cache hit", "for", "RD", "resource", resource)
			return &res
		}
	}
	return nil
}

func (h *Helpers) ReadRDCacheByResource(resource string) []types.Torrent {
	records, err := h.app.Dao().
		FindRecordsByFilter("rd_resolved", "resource = {:resource}", "id", -1, 0, dbx.Params{"resource": resource})
	res := make([]types.Torrent, 0)
	if err == nil {
		for _, record := range records {
			var r types.Torrent
			date := record.GetDateTime("updated")
			now := time.Now().Add(time.Duration((-8) * time.Hour))
			// add 1 hour to date
			if date.Time().Before(now) {
				continue
			}
			err := record.UnmarshalJSONField("data", &r)
			if err == nil {
				res = append(res, r)
			}
		}
	}
	return res
}

func (h *Helpers) WriteRDCache(resource string, magnet string, data interface{}) {
	log.Info("cache write", "for", "RD", "resource", resource)
	record, err := h.app.Dao().
		FindFirstRecordByFilter("rd_resolved", "magnet = {:magnet}", dbx.Params{"magnet": magnet})

	if err == nil {
		record.Set("data", &data)
		h.app.Dao().SaveRecord(record)
	} else {
		collection, _ := h.app.Dao().FindCollectionByNameOrId("rd_resolved")
		record := models.NewRecord(collection)
		record.Set("data", &data)
		record.Set("magnet", magnet)
		record.Set("resource", resource)
		h.app.Dao().SaveRecord(record)
	}
}

func (h *Helpers) ReadTraktSeasonCache(id uint) []any {
	record, err := h.app.Dao().
		FindFirstRecordByFilter("trakt_seasons", "trakt_id = {:id}", dbx.Params{"id": id})
	res := make([]any, 0)
	if err == nil {
		err := record.UnmarshalJSONField("data", &res)
		date := record.GetDateTime("updated")
		now := time.Now()

		if err == nil {
			if date.Time().Before(now.AddDate(0, 0, -1)) {
				return nil
			}
			log.Debug("cache hit", "for", "trakt", "resource", "show_seasons", "id", id)
			return res
		}
	}
	return nil
}
