package scraper

import (
	"fmt"
	"sync"

	"github.com/odin-movieshow/server/realdebrid"
	"github.com/odin-movieshow/server/settings"
	"github.com/odin-movieshow/server/types"

	"github.com/charmbracelet/log"
	"github.com/go-resty/resty/v2"
	"github.com/pocketbase/pocketbase"
)

func GetLinks(data map[string]any, app *pocketbase.PocketBase) []types.Torrent {

	j := settings.GetJackett(app)

	if j == nil {
		log.Error("jackett", "error", "no settings")
		return []types.Torrent{}
	}

	allTorrents := []types.Torrent{}

	res := resty.New().
		R().
		SetBody(data).
		SetHeader("Content-Type", "application/json").
		SetResult(&allTorrents)

	_, err := res.Post(fmt.Sprintf("%s/scrape", settings.GetScraperUrl(app)))
	if err != nil {
		log.Error("scrape", err)
		return []types.Torrent{}
	}

	wg := sync.WaitGroup{}

	for k := range allTorrents {

		defer wg.Done()
		realdebrid.Unrestrict(k, allTorrents, app)
	}

	filtered := []types.Torrent{}
	unrestricted := 0
	for _, t := range allTorrents {
		if len(t.RealDebrid) > 0 {
			filtered = append(filtered, t)
			unrestricted += len(t.RealDebrid)
		}
	}
	log.Info("scrape done", "unrestricted", unrestricted)
	return allTorrents
}
