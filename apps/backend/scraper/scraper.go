package scraper

import (
	"fmt"
	"sync"

	"github.com/odin-movieshow/server/helpers"
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
	chunks := helpers.Chunk(allTorrents)
	allTorrentsUnrestricted := []types.Torrent{}
	for _, c := range chunks {
		wg.Add(1)
		go func(torrents []types.Torrent) {
			defer wg.Done()
			for _, k := range torrents {
				allTorrentsUnrestricted = append(
					allTorrentsUnrestricted,
					realdebrid.Unrestrict(k, app),
				)
			}
		}(c)

	}
	wg.Wait()

	filtered := []types.Torrent{}
	for _, t := range allTorrentsUnrestricted {
		if len(t.RealDebrid) > 0 {
			filtered = append(filtered, t)
		}
	}
	log.Info("scrape done", "unrestricted", len(filtered))
	return filtered
}
