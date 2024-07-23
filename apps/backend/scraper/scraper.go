package scraper

import (
	"fmt"
	"sync"

	"github.com/odin-movieshow/server/helpers"
	"github.com/odin-movieshow/server/realdebrid"
	"github.com/odin-movieshow/server/settings"
	"github.com/odin-movieshow/server/types"
	"github.com/thoas/go-funk"

	"github.com/charmbracelet/log"
	"github.com/go-resty/resty/v2"
	"github.com/pocketbase/pocketbase"
)

func GetLinks(data map[string]any, app *pocketbase.PocketBase) []types.Torrent {
	mux := sync.Mutex{}
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

				q1s := funk.Filter(allTorrentsUnrestricted, func(t types.Torrent) bool {
					return t.Quality == "4K" && len(t.RealDebrid) > 0
				}).([]types.Torrent)

				q2s := funk.Filter(allTorrentsUnrestricted, func(t types.Torrent) bool {
					return t.Quality == "1080p" && len(t.RealDebrid) > 0
				}).([]types.Torrent)

				q3s := funk.Filter(allTorrentsUnrestricted, func(t types.Torrent) bool {
					return t.Quality == "720p" && len(t.RealDebrid) > 0
				}).([]types.Torrent)

				q4s := funk.Filter(allTorrentsUnrestricted, func(t types.Torrent) bool {
					return t.Quality == "SD" && len(t.RealDebrid) > 0
				}).([]types.Torrent)

				if k.Quality == "1080p" {
					if len(q2s) > 20 {
						continue
					}
				}

				if k.Quality == "720p" {
					if len(q1s)+len(q2s) > 30 {
						continue
					}
					if len(q3s) > 10 {
						continue
					}
				}
				if k.Quality == "SD" {
					if len(q1s)+len(q2s) > 30 {
						continue
					}
					if len(q4s) > 10 {
						continue
					}
				}
				mux.Lock()
				allTorrentsUnrestricted = append(
					allTorrentsUnrestricted,
					realdebrid.Unrestrict(k, app),
				)
				mux.Unlock()
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
