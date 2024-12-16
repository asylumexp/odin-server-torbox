package indexer

import (
	"os"

	"github.com/charmbracelet/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/odin-movieshow/backend/common"
	"github.com/odin-movieshow/backend/indexer/jackett"
)

type Indexer struct {
	mqt mqtt.Client
}

func New(mqt mqtt.Client) *Indexer {
	return &Indexer{mqt: mqt}
}

func Index(data common.Payload) {
	jackettUrl := os.Getenv("JACKETT_URL")
	jackettKey := os.Getenv("JACKETT_KEY")
	if jackettUrl == "" || jackettKey == "" {
		log.Error("missing env vars JACKETT_URL and JACKETT_KEY")
		return
	}

	jackett.Search(data)
}
