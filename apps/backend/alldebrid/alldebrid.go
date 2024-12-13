package alldebrid

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

type AllDebrid struct {
	app      *pocketbase.PocketBase
	settings *settings.Settings
	Headers  map[string]any
}

type Magnet struct {
	Magnet string `json:"magnet"`
	Hash   string `json:"hash"`
	Name   string `json:"name"`
	Size   uint32 `json:"size"`
	Ready  bool   `json:"ready"`
	ID     uint32 `json:"id"`
}

type FileNode struct {
	N string     `json:"n"`
	S string     `json:"s"`
	L string     `json:"l"`
	E []FileNode `json:"e"`
}

func New(app *pocketbase.PocketBase, settings *settings.Settings) *AllDebrid {
	return &AllDebrid{app: app, settings: settings}
}

func (ad *AllDebrid) RemoveByType(t string) {
	var res interface{}
	headers, _ := ad.CallEndpoint(fmt.Sprintf("/%s/?limit=1", t), "GET", nil, &res)
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
		ad.CallEndpoint(fmt.Sprintf("/%s/?limit=200", t), "GET", nil, &res)
		for _, v := range res.([]any) {
			ad.CallEndpoint(
				fmt.Sprintf("/%s/delete/%s", t, v.(map[string]any)["id"].(string)),
				"DELETE",
				nil,
				nil,
			)
		}
		count -= 200
	}

	log.Info("realdebrid cleanup", "type", t, "count", count)
}

func (ad *AllDebrid) RefreshTokens() {
	records := []models.Record{}
	ad.app.Dao().RecordQuery("settings").All(&records)
	if len(records) == 0 {
		return
	}
	data := make(map[string]any)
	r := records[0]

	rdsets := ad.settings.GetAllDebrid()
	request := resty.New().SetRetryCount(3).SetRetryWaitTime(time.Second * 3).R()
	request.SetHeader("Authorization", "Bearer "+rdsets.AccessToken)

	for k, v := range ad.Headers {

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
		ad.app.Dao().SaveRecord(&r)
	}
}

func (ad *AllDebrid) CallEndpoint(
	endpoint string,
	method string,
	body map[string]string,
	res interface{},
) (http.Header, int) {
	request := resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(time.Second * 1).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return r.StatusCode() == 429
		}).
		R()
	request.SetResult(res)
	var respHeaders http.Header
	status := 200

	if body != nil {
		request.SetFormData(body)
	}

	rdsets := ad.settings.GetAllDebrid()
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

	resp, err := r("https://api.alldebrid.com/v4" + endpoint)
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
			strings.ReplaceAll(string(resp.Body()), "\n", ""),
		)
	} else {
		log.Debug("realdebrid", "call", fmt.Sprintf("%s %s", method, endpoint))
	}

	return respHeaders, status
}

func getLinks(files []FileNode) []string {
	links := []string{}
	for _, v := range files {
		if v.L != "" {
			links = append(links, v.L)
			continue
		}
		if len(v.E) > 0 {
			links = append(links, getLinks(v.E)...)
		}
	}
	return links
}

func (ad *AllDebrid) Unrestrict(m string) interface{} {
	var res struct {
		Data struct {
			Magnets []Magnet `json:"magnets"`
		} `json:"data"`
	}

	ad.CallEndpoint("/magnet/upload?magnets[]="+m, "GET", nil, &res)

	if len(res.Data.Magnets) == 0 {
		return nil
	}

	magnetId := strconv.Itoa(int(res.Data.Magnets[0].ID))
	defer ad.CallEndpoint(
		fmt.Sprintf("/magnet/delete?id[]=%s", magnetId),
		"GET",
		nil,
		nil,
	)

	var files struct {
		Data struct {
			Magnets []struct {
				Files []FileNode `json:"files"`
			} `json:"magnets"`
		} `json:"data"`
	}

	ad.CallEndpoint("/magnet/files?id="+magnetId, "GET", nil, &files)

	// info, _, _ := ad.CallEndpoint("/torrents/info/"+magnetId., "GET", nil,)
	//
	// if info == nil {
	// 	return nil
	// }
	//
	// if info.(map[string]any)["links"] == nil {
	// 	return nil
	// }
	//
	//

	links := getLinks(files.Data.Magnets[0].Files)

	for _, v := range links {
		var u struct {
			Data struct {
				Link     string `json:"link"`
				Filename string `json:"filename"`
				Filesize int    `json:"filesize"`
			} `json:"data"`
		}
		ad.CallEndpoint("/link/unlock&link="+v, "GET", nil, &u)
		fname := u.Data.Filename

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

func (ad *AllDebrid) Cleanup() {
	ad.RemoveByType("downloads")
}
