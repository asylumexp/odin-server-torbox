package scraper

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/odin-movieshow/backend/alldebrid"
	"github.com/odin-movieshow/backend/helpers"
	"github.com/odin-movieshow/backend/realdebrid"
	"github.com/odin-movieshow/backend/settings"
	"github.com/odin-movieshow/backend/types"
	"github.com/thoas/go-funk"

	"github.com/charmbracelet/log"
	"github.com/go-resty/resty/v2"
	"github.com/pocketbase/pocketbase"
)

type Scraper struct {
	app        *pocketbase.PocketBase
	settings   *settings.Settings
	helpers    *helpers.Helpers
	realdebrid *realdebrid.RealDebrid
	alldebrid  *alldebrid.AllDebrid
}

func New(
	app *pocketbase.PocketBase,
	settings *settings.Settings,
	helpers *helpers.Helpers,
	realdebrid *realdebrid.RealDebrid,
	alldebrid *alldebrid.AllDebrid,
) *Scraper {
	return &Scraper{app: app, settings: settings, helpers: helpers, realdebrid: realdebrid, alldebrid: alldebrid}
}

func (s *Scraper) GetLinks(data map[string]any, mqt mqtt.Client) {
	// mux := sync.Mutex{}
	j := s.settings.GetJackett()

	if j == nil {
		log.Error("jackett", "error", "no settings")
		return
	}

	topic := "odin-movieshow/" + data["type"].(string)
	indexertopic := "odin-movieshow/indexer/" + data["type"].(string)
	if data["type"] == "episode" {
		topic += "/" + data["episode_trakt"].(string)
		indexertopic += "/" + data["episode_trakt"].(string)
	} else {
		topic += "/" + data["trakt"].(string)
		indexertopic += "/" + data["trakt"].(string)
	}

	log.Debug("MQTT", "indexer topic", indexertopic)
	log.Debug("MQTT", "result topic", topic)
	torrentQueue := make(chan types.Torrent)

	allTorrentsUnrestricted := s.helpers.ReadRDCacheByResource(topic)
	for _, u := range allTorrentsUnrestricted {
		cstr, _ := json.Marshal(u)
		mqt.Publish(topic, 0, false, cstr)
	}

	if token := mqt.Subscribe(indexertopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		newTorrents := []types.Torrent{}
		json.Unmarshal(msg.Payload(), &newTorrents)
		go func() {
			for _, t := range newTorrents {
				if t.Magnet != "" {
					torrentQueue <- t
				}
			}
		}()
		// fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	}); token.Wait() &&
		token.Error() != nil {
		log.Error("mqtt-subscribe-indexer", "error", token.Error())
	}

	res := resty.New().SetTimeout(15*time.Minute).
		R().
		SetBody(data).
		SetHeader("Content-Type", "application/json")

	go func() {
		done := []string{}
		for {
			select {
			case k := <-torrentQueue:
				if !funk.Contains(done, k.Magnet) && k.Quality != "720p" && k.Quality != "SD" &&
					k.Quality != "CAM" {
					s.unrestrict(k, mqt, topic)
					done = append(done, k.Magnet)
				}
			}
		}
	}()

	_, err := res.Post(fmt.Sprintf("%s/scrape", s.settings.GetScraperUrl()))
	if err != nil {
		log.Error("scrape", err)
		return
	}

	<-torrentQueue
	log.Warn("DONE")
}

func (s *Scraper) unrestrict(
	k types.Torrent,
	mqt mqtt.Client,
	topic string,
) {
	cache := s.helpers.ReadRDCache(topic, k.Magnet)
	if cache != nil {
		cstr, _ := json.Marshal(cache)
		mqt.Publish(topic, 0, false, cstr)
		return
	}
	return
	us := s.alldebrid.Unrestrict(k.Magnet)
	// us := s.realdebrid.Unrestrict(k.Magnet)
	if len(us) == 0 {
		return
	}
	k.Links = us
	log.Info(k.ReleaseTitle)
	s.helpers.WriteRDCache(topic, k.Magnet, k)
	kstr, _ := json.Marshal(k)
	mqt.Publish(topic, 0, false, kstr)
}
