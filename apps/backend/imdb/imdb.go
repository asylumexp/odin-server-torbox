package imdb

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-resty/resty/v2"
)

func Get(id string) map[string]interface{} {
	request := resty.New().
		SetRetryCount(3).
		SetTimeout(time.Second * 30).
		SetRetryWaitTime(time.Second).
		R()

	res, err := request.Get("https://www.imdb.com/title/" + id)
	r := make(map[string]interface{})
	if err == nil {
		b := string(res.Body())
		start := "<script type=\"application/ld+json\">"
		end := "</script>"
		startIndex := strings.Index(b, start)
		endIndex := strings.Index(b[startIndex+len(start):], end)
		str := b[startIndex+len(start) : endIndex+startIndex+len(start)]
		if err := json.Unmarshal([]byte(str), &r); err != nil {
			log.Error("imdb", "id", id, "err", err)
		} else {
			log.Info("imdb", "id", id, "parsed", "success")
		}
	}
	return r
}
