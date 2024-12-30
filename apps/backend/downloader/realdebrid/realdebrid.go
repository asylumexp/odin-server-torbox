package realdebrid

import (
	"fmt"
	"mime"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/odin-movieshow/backend/settings"
	"github.com/odin-movieshow/backend/types"

	"github.com/charmbracelet/log"
	"github.com/go-resty/resty/v2"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/thoas/go-funk"
)

type RealDebrid struct {
	app      *pocketbase.PocketBase
	settings *settings.Settings
	Headers  map[string]any
}

func New(app *pocketbase.PocketBase, settings *settings.Settings) *RealDebrid {
	return &RealDebrid{app: app, settings: settings}
}

func (rd *RealDebrid) RemoveByType(t string) {
	var res interface{}
	headers, _ := rd.CallEndpoint(fmt.Sprintf("/%s/?limit=1", t), "GET", nil, &res)
	if res == nil {
		return
	}

	count := 0
	if headers.Get("X-Total-Count") != "" {
		c, err := strconv.Atoi(headers.Get("X-Total-Count"))
		if err == nil {
			count = c
		}
	}

	for count > 0 {

		rd.CallEndpoint(fmt.Sprintf("/%s/?limit=200", t), "GET", nil, &res)
		for _, v := range res.([]any) {
			rd.CallEndpoint(
				fmt.Sprintf("/%s/delete/%s", t, v.(map[string]any)["id"].(string), nil),
				"DELETE",
				nil,
				nil,
			)
		}
		count -= 200
	}

	log.Info("realdebrid cleanup", "type", t, "count", count)
}

func (rd *RealDebrid) RefreshTokens() {
	records := []models.Record{}
	rd.app.Dao().RecordQuery("settings").All(&records)
	if len(records) == 0 {
		return
	}
	data := make(map[string]any)
	r := records[0]

	rdsets := rd.settings.GetRealDebrid()
	request := resty.New().SetRetryCount(3).SetRetryWaitTime(time.Second * 3).R()
	request.SetHeader("Authorization", "Bearer "+rdsets.AccessToken)

	for k, v := range rd.Headers {

		if funk.Contains([]string{"Host", "Connection"}, k) {
			continue
		}
		request.SetHeader(k, v.(string))
	}
	if _, err := request.SetFormData(map[string]string{
		"client_id":     rdsets.ClientId,
		"client_secret": rdsets.ClientSecret,
		"code":          rdsets.RefreshToken,
		"grant_type":    "http://oauth.net/grant_type/device/1.0",
	}).SetResult(&data).Post("https://api.real-debrid.com/oauth/v2/token"); err == nil {
		if data["access_token"] == nil || data["refresh_token"] == nil {
			return
		}
		rdsets.AccessToken = data["access_token"].(string)
		rdsets.RefreshToken = data["refresh_token"].(string)
		r.Set("real_debrid", rdsets)
		log.Info("realdebrid", "token", "refreshed")
		rd.app.Dao().SaveRecord(&r)
	}
}

func (rd *RealDebrid) CallEndpoint(
	endpoint string,
	method string,
	body map[string]string,
	data interface{},
) (http.Header, int) {
	request := resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(time.Second * 1).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return r.StatusCode() == 429
		}).
		R()
	request.SetResult(&data)
	var respHeaders http.Header
	status := 200

	if body != nil {
		request.SetFormData(body)
	}

	rdsets := rd.settings.GetRealDebrid()
	request.SetHeader("Authorization", "Bearer "+rdsets.AccessToken)

	request.Attempt = 3

	var r func(url string) (*resty.Response, error)
	switch method {
	case "POST":
		r = request.Post
	case "PATCH":
		r = request.Patch
	case "PUT":
		r = request.Put
	case "DELETE":
		r = request.Delete
	default:
		r = request.Get

	}

	resp, err := r("https://api.real-debrid.com/rest/1.0" + endpoint)
	if err == nil {
		respHeaders = resp.Header()
		status = resp.StatusCode()
	} else {
		log.Error("realdebrid", "url", endpoint, "error", err)
		return respHeaders, status
	}

	if status > 299 {
		log.Error(
			"realdebrid",
			"status",
			status,
			"url",
			endpoint,
			"data",
			strings.Replace(string(resp.Body()), "\n", "", -1),
		)
	} else {
		log.Debug("realdebrid", "call", fmt.Sprintf("%s %s", method, endpoint))
	}

	return respHeaders, status
}

type Magnet struct {
	Id string `json:"id"`
}

type Info struct {
	Links []string `json:"links"`
}

type Link struct {
	Id         string `json:"id"`
	Filename   string `json:"filename"`
	Filesize   uint   `json:"filesize"`
	Download   string `json:"download"`
	Streamable int    `json:"streamable`
}

func (rd *RealDebrid) Unrestrict(m string) []types.Unrestricted {
	magnet := Magnet{}
	rd.CallEndpoint("/torrents/addMagnet", "POST", map[string]string{
		"host":   "real-debrid.com",
		"magnet": m,
	}, &magnet)

	if magnet.Id == "" {
		return nil
	}

	defer rd.CallEndpoint(
		fmt.Sprintf("/torrents/delete/%s", magnet.Id),
		"DELETE",
		nil,
		nil,
	)

	rd.CallEndpoint("/torrents/selectFiles/"+magnet.Id, "POST", map[string]string{
		"files": "all",
	}, nil)

	info := Info{}
	rd.CallEndpoint("/torrents/info/"+magnet.Id, "GET", nil, &info)

	if len(info.Links) == 0 {
		return nil
	}

	us := []types.Unrestricted{}

	for _, v := range info.Links {
		u := Link{}
		rd.CallEndpoint("/unrestrict/link", "POST", map[string]string{
			"link": v,
		}, &u)
		if u.Filename == "" {
			continue
		}
		fname := u.Filename

		mimetype := mime.TypeByExtension(fname[strings.LastIndex(fname, "."):])

		isVideo := strings.Contains(
			mimetype,
			"video",
		)

		match, _ := regexp.MatchString("^[Ss]ample[ -_]?[0-9].", fname)

		if !match && isVideo {
			log.Debug("realdebrid unrestricted", "file", fname)
			streams := []string{}
			if u.Streamable == 1 {
				streams = append(streams, "https://real-debrid.com/streaming-"+u.Id)
			}
			un := types.Unrestricted{Filename: fname, Filesize: int(u.Filesize), Download: u.Download, Streams: streams}
			us = append(us, un)
		}
	}

	return us
}

func (rd *RealDebrid) Cleanup() {
	rd.RemoveByType("downloads")
}
