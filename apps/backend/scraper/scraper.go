package scraper

import (
	"fmt"
	"sync"

	"backend/realdebrid"
	"backend/settings"
	"backend/types"

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
		wg.Add(1)

		go func(k int) {
			defer wg.Done()
			realdebrid.Unrestrict(k, allTorrents, app)
		}(k)
	}
	wg.Wait()

	return allTorrents
}
