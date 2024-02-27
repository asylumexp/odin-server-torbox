package realdebrid

import (
	"mime"
	"net/http"
	"regexp"
	"strings"
	"time"

	"backend/helpers"
	"backend/settings"
	"backend/types"

	"github.com/charmbracelet/log"
	"github.com/go-resty/resty/v2"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

var Headers = make(map[string]any)

func RemoveTorrents(app *pocketbase.PocketBase) {
	torrents, _, _ := CallEndpoint("/torrents?limit=300", "GET", nil, app)
	if torrents == nil {
		return
	}
	for _, v := range torrents.([]any) {
		CallEndpoint("/torrents/delete/"+v.(map[string]any)["id"].(string), "DELETE", nil, app)
	}

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

		if helpers.ArrayContains([]string{"Host", "Connection"}, k) {
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
	request := resty.New().SetRetryCount(3).SetRetryWaitTime(time.Second * 3).R()
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

	}

	return data, respHeaders, status

}

func Unrestrict(k int, objmap []types.Torrent, app *pocketbase.PocketBase) {
	item := objmap[k]
	magnet, _, _ := CallEndpoint("/torrents/addMagnet", "POST", map[string]string{
		"host":   "real-debrid.com",
		"magnet": item.Magnet,
	}, app)

	if magnet == nil || magnet.(map[string]any)["id"] == nil {
		return
	}

	magnetId := magnet.(map[string]any)["id"].(string)

	CallEndpoint("/torrents/selectFiles/"+magnetId, "POST", map[string]string{
		"files": "all",
	}, app)

	info, _, _ := CallEndpoint("/torrents/info/"+magnetId, "GET", nil, app)
	if info == nil || info.(map[string]any)["links"] == nil {
		return
	}

	links := info.(map[string]any)["links"].([]any)

	downloads := make([]map[string]any, 0)

	for _, v := range links {
		u, _, _ := CallEndpoint("/unrestrict/link", "POST", map[string]string{
			"link": v.(string),
		}, app)
		if u == nil {
			continue
		}
		fname := u.(map[string]any)["filename"].(string)
		mimetype := mime.TypeByExtension(fname[strings.LastIndex(fname, "."):])
		log.Debug(mimetype, "file", fname)
		isVideo := strings.Contains(
			mimetype,
			"video",
		)
		match, _ := regexp.MatchString("^Sample[ -]?[0-9].", fname)
		if !match && isVideo {
			downloads = append(downloads, u.(map[string]any))
		}
	}

	objmap[k].RealDebrid = downloads
}
