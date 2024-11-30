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
	res, headers, _ := rd.CallEndpoint(fmt.Sprintf("/%s/?limit=1", t), "GET", nil)
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
		res, _, _ := rd.CallEndpoint(fmt.Sprintf("/%s/?limit=200", t), "GET", nil)
		for _, v := range res.([]any) {
			rd.CallEndpoint(
				fmt.Sprintf("/%s/delete/%s", t, v.(map[string]any)["id"].(string)),
				"DELETE",
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
) (any, http.Header, int) {

	var data any
	request := resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(time.Second * 3).
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
		return data, respHeaders, status
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

	return data, respHeaders, status

}

func (rd *RealDebrid) Unrestrict(m string) map[string]any {
	magnet, _, _ := rd.CallEndpoint("/torrents/addMagnet", "POST", map[string]string{
		"host":   "real-debrid.com",
		"magnet": m,
	})

	if magnet == nil || magnet.(map[string]any)["id"] == nil {
		return nil
	}

	magnetId := magnet.(map[string]any)["id"].(string)

	defer rd.CallEndpoint(
		fmt.Sprintf("/torrents/delete/%s", magnetId),
		"DELETE",
		nil,
	)

	rd.CallEndpoint("/torrents/selectFiles/"+magnetId, "POST", map[string]string{
		"files": "all",
	})

	info, _, _ := rd.CallEndpoint("/torrents/info/"+magnetId, "GET", nil)

	if info == nil {
		return nil
	}

	if info.(map[string]any)["links"] == nil {
		return nil
	}

	links := info.(map[string]any)["links"].([]any)

	for _, v := range links {
		log.Debug("realdebrid unrestricted", "link", v.(string))
		u, _, _ := rd.CallEndpoint("/unrestrict/link", "POST", map[string]string{
			"link": v.(string),
		})
		if u == nil {
			continue
		}
		fname := u.(map[string]any)["filename"].(string)

		mimetype := mime.TypeByExtension(fname[strings.LastIndex(fname, "."):])

		isVideo := strings.Contains(
			mimetype,
			"video",
		)

		match, _ := regexp.MatchString("^[Ss]ample[ -_]?[0-9].", fname)

		if !match && isVideo {
			log.Debug("realdebrid unrestricted", "file", fname)
			return u.(map[string]any)
		}
	}

	return nil
}

func (rd *RealDebrid) Cleanup() {
	rd.RemoveByType("downloads")
}
