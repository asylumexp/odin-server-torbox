package scraper

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/odin-movieshow/server/helpers"
	"github.com/odin-movieshow/server/realdebrid"
	"github.com/odin-movieshow/server/settings"
	"github.com/odin-movieshow/server/types"
	"github.com/thoas/go-funk"

	"github.com/charmbracelet/log"
	"github.com/go-resty/resty/v2"
	"github.com/pocketbase/pocketbase"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func mqttclient() mqtt.Client {
	// mqtt.DEBUG = log.New(os.Stdout, "", 0)
	// mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().
		AddBroker("wss://mqtt.dnmc.in/mqtt").
		SetUsername("mqtt").
		SetPassword("mqtt9040!")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)

	return c
}

func GetLinks(data map[string]any, app *pocketbase.PocketBase, mqtt mqtt.Client) []types.Torrent {
	// mux := sync.Mutex{}
	j := settings.GetJackett(app)

	if j == nil {
		log.Error("jackett", "error", "no settings")
		return []types.Torrent{}
	}

	allTorrents := []types.Torrent{}

	topic := "odin-movieshow/" + data["type"].(string)
	if data["episode_trakt"] != nil {
		topic = topic + "/" + data["episode_trakt"].(string)
	}
	if data["trakt"] != nil {
		topic = topic + "/" + data["trakt"].(string)
	}

	log.Info(topic)
	allTorrentsUnrestricted := helpers.ReadRDCacheByResource(app, topic)
	for _, u := range allTorrentsUnrestricted {
		cstr, _ := json.Marshal(u)
		mqtt.Publish(topic, 0, false, cstr)
	}

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

	for _, k := range allTorrents {

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

		cache := helpers.ReadRDCache(app, topic, k.Magnet)
		if cache != nil {
			// allTorrentsUnrestricted = append(allTorrentsUnrestricted, *cache)
			// cstr, _ := json.Marshal(cache)
			// mqtt.Publish(topic, 0, false, cstr)
			continue
		}

		u := realdebrid.Unrestrict(k.Magnet, app)
		k.RealDebrid = append(k.RealDebrid, u)

		if len(k.RealDebrid) > 0 {
			allTorrentsUnrestricted = append(allTorrentsUnrestricted, k)
			helpers.WriteRDCache(app, topic, k.Magnet, k)
			kstr, _ := json.Marshal(k)
			mqtt.Publish(topic, 0, false, kstr)
		}
		// mux.Unlock()
	}

	// }
	// wg.Wait()

	// filtered := []types.Torrent{}
	// for _, t := range allTorrentsUnrestricted {
	// 	if len(t.RealDebrid) > 0 {
	// 		filtered = append(filtered, t)
	// 	}
	// }
	// log.Info("scrape done", "unrestricted", len(filtered))

	return allTorrentsUnrestricted
}
