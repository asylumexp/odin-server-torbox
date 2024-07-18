package realdebrid

import (
	"fmt"
	"mime"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/odin-movieshow/server/settings"

	"github.com/charmbracelet/log"
	"github.com/go-resty/resty/v2"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/thoas/go-funk"
)

var Headers = make(map[string]any)

func RemoveByType(app *pocketbase.PocketBase, t string) {
	res, headers, _ := CallEndpoint(fmt.Sprintf("/%s/?limit=1", t), "GET", nil, app)
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
		res, _, _ := CallEndpoint(fmt.Sprintf("/%s/?limit=200", t), "GET", nil, app)
		for _, v := range res.([]any) {
			CallEndpoint(
				fmt.Sprintf("/%s/delete/%s", t, v.(map[string]any)["id"].(string)),
				"DELETE",
				nil,
				app,
			)
		}
		count -= 200
	}

	log.Info("realdebrid cleanup", "type", t, "count", count)

}

func RefreshTokens(app *pocketbase.PocketBase) {
	records := []models.Record{}
	app.Dao().RecordQuery("settings").All(&records)
	if len(records) == 0 {
		return
	}
	data := make(map[string]any)
	r := records[0]

	rd := settings.GetRealDebrid(app)
	request := resty.New().SetRetryCount(3).SetRetryWaitTime(time.Second * 3).R()
	request.SetHeader("Authorization", "Bearer "+rd.AccessToken)

	for k, v := range Headers {

		if funk.Contains([]string{"Host", "Connection"}, k) {
			continue
		}
		request.SetHeader(k, v.(string))
	}
	if _, err := request.SetFormData(map[string]string{
		"client_id":     rd.ClientId,
		"client_secret": rd.ClientSecret,
		"code":          rd.RefreshToken,
		"grant_type":    "http://oauth.net/grant_type/device/1.0",
	}).SetResult(&data).Post("https://api.real-debrid.com/oauth/v2/token"); err == nil {
		rd.AccessToken = data["access_token"].(string)
		rd.RefreshToken = data["refresh_token"].(string)
		r.Set("real_debrid", rd)
		log.Info("realdebrid", "token", "refreshed")
		app.Dao().SaveRecord(&r)
	}

}

func CallEndpoint(
	endpoint string,
	method string,
	body map[string]string,
	app *pocketbase.PocketBase,
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

	rd := settings.GetRealDebrid(app)
	request.SetHeader("Authorization", "Bearer "+rd.AccessToken)

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

func Unrestrict(m string, app *pocketbase.PocketBase) any {
	magnet, _, _ := CallEndpoint("/torrents/addMagnet", "POST", map[string]string{
		"host":   "real-debrid.com",
		"magnet": m,
	}, app)

	if magnet == nil || magnet.(map[string]any)["id"] == nil {
		return nil
	}

	magnetId := magnet.(map[string]any)["id"].(string)

	defer CallEndpoint(
		fmt.Sprintf("/torrents/delete/%s", magnetId),
		"DELETE",
		nil,
		app,
	)

	CallEndpoint("/torrents/selectFiles/"+magnetId, "POST", map[string]string{
		"files": "all",
	}, app)

	info, _, _ := CallEndpoint("/torrents/info/"+magnetId, "GET", nil, app)

	if info == nil {
		return nil
	}

	if info.(map[string]any)["links"] == nil {
		return nil
	}

	links := info.(map[string]any)["links"].([]any)

	for _, v := range links {
		log.Debug("realdebrid unrestricted", "link", v.(string))
		u, _, _ := CallEndpoint("/unrestrict/link", "POST", map[string]string{
			"link": v.(string),
		}, app)
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
			return u
		}
	}

	return nil
}

func Cleanup(app *pocketbase.PocketBase) {
	RemoveByType(app, "downloads")
}
