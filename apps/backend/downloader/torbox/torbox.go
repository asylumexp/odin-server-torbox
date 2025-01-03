package torbox

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-resty/resty/v2"
	"github.com/odin-movieshow/backend/settings"
	"github.com/odin-movieshow/backend/types"
	"github.com/pocketbase/pocketbase"
)

type baseResponse[T any] struct {
    Success bool   `json:"success"`
    Error   string `json:"error"`
    Detail  string `json:"detail"`
    Data    T      `json:"data"`
}

type torrentInfoResult struct {
    ID            int    `json:"id"`
    AuthID        string `json:"auth_id"`
    Hash          string `json:"hash"`
    Name          string `json:"name"`
    Magnet        string `json:"magnet"`
    DownloadState string `json:"download_state"`
	Files []torrentFile `json:"files"`
}

type torrentAddResult struct {
    Hash      string `json:"hash"`
    TorrentID int    `json:"torrent_id"`
    AuthID    string `json:"auth_id"`
}

type torrentFile struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Size int64  `json:"size"`
    // ... etc.
}


type TorBox struct {
    app      *pocketbase.PocketBase
    settings *settings.Settings
    bearerToken string
    Headers     map[string]any
}

func New(app *pocketbase.PocketBase, settings *settings.Settings) *TorBox {
    token := os.Getenv("TORBOX_KEY")
    return &TorBox{
        app:         app,
        settings:    settings,
        bearerToken: token,
    }
}


func (tb *TorBox) callEndpoint(endpoint string, method string, body map[string]string, res interface{}) (http.Header, int) {
    baseURL := "https://api.torbox.app/v1/api/"
    fullURL := baseURL + endpoint

    client := resty.New().
        SetRetryCount(3).
        SetRetryWaitTime(time.Second * 1).
        AddRetryCondition(func(r *resty.Response, err error) bool {
            return r.StatusCode() == 429
        }).
        R()

    if res != nil {
        client.SetResult(res)
    }

    if body != nil {
        client.SetFormData(body)
    }

    if tb.bearerToken != "" {
        client.SetHeader("Authorization", "Bearer "+tb.bearerToken)
    }

    var fn func(url string) (*resty.Response, error)
    switch method {
    case "POST":
        fn = client.Post
    case "PATCH":
        fn = client.Patch
    case "PUT":
        fn = client.Put
    case "DELETE":
        fn = client.Delete
    default:
        fn = client.Get
    }

    resp, err := fn(fullURL)
    status := 0
    headers := http.Header{}
    if err != nil {
        log.Error("torbox", "url", endpoint, "error", err)
        return headers, status
    }

    headers = resp.Header()
    status = resp.StatusCode()

    if status > 299 {
        log.Error("torbox", 
            "status", status,
            "url", endpoint,
            "data", strings.ReplaceAll(string(resp.Body()), "\n", ""),
        )
    } else {
        log.Debug("torbox", "call", fmt.Sprintf("%s %s", method, endpoint))
    }

    return headers, status
}


func (tb *TorBox) AddMagnet(magnet string) *torrentAddResult {
    form := map[string]string{
        "magnet":    magnet,
        "seed":      "3",
        "allow_zip": "false",
        "name":      "",
    }

    var response baseResponse[torrentAddResult]
    _, status := tb.callEndpoint("torrents/createtorrent", "POST", form, &response)
    if status > 299 {
        return nil
    }
    if !response.Success {
        log.Error("torbox AddMagnet error", "error", response.Error, "detail", response.Detail)
        return nil
    }
    return &response.Data
}

func (tb *TorBox) GetCurrentTorrents() []torrentInfoResult {
    var response baseResponse[[]torrentInfoResult]
    _, status := tb.callEndpoint("torrents/mylist?bypass_cache=false", "GET", nil, &response)
    if status > 299 {
        return nil
    }
    if !response.Success {
        log.Error("torbox GetCurrentTorrents error", "error", response.Error, "detail", response.Detail)
        return nil
    }
    return response.Data
}

func (tb *TorBox) RequestDownloadLink(torrentID, fileID int, zip bool) string {
    query := fmt.Sprintf("torrents/requestdl?torrent_id=%d&file_id=%d&zip_link=%v", torrentID, fileID, zip)

    var response baseResponse[string]
    _, status := tb.callEndpoint(query, "GET", nil, &response)
    if status > 299 {
        return ""
    }
    if !response.Success {
        log.Error("torbox RequestDownloadLink error", "error", response.Error, "detail", response.Detail)
        return ""
    }
    return response.Data
}

func (tb *TorBox) ControlTorrent(hash, action string) bool {
    body := map[string]string{
        "hash":      hash,
        "operation": action,
    }

    var response baseResponse[any]
    _, status := tb.callEndpoint("torrents/controltorrent", "POST", body, &response)
    if status > 299 {
        return false
    }
    if !response.Success {
        log.Error("torbox ControlTorrent error", "error", response.Error, "detail", response.Detail)
        return false
    }
    return true
}

func (tb *TorBox) Unrestrict(magnet string) []types.Unrestricted {
    // 1) Add the magnet
    addResult := tb.AddMagnet(magnet)
    if addResult == nil || addResult.TorrentID == 0 {
        return nil
    }
    torrentID := addResult.TorrentID

    maxChecks := 5
    interval := time.Second * 2
    ready := false
    for i := 0; i < maxChecks; i++ {
        torrents := tb.GetCurrentTorrents()
        for _, t := range torrents {
            if t.ID == torrentID && (t.DownloadState == "finished" || t.DownloadState == "seeding") {
                ready = true
                break
            }
        }
        if ready {
            break
        }
        time.Sleep(interval)
    }
    if !ready {
        return nil
    }

    var singleTorrent torrentInfoResult
    {
        torrList := tb.GetCurrentTorrents()
        for _, t := range torrList {
            if t.ID == torrentID {
                singleTorrent = t
                break
            }
        }
    }
    var result []types.Unrestricted
    re := regexp.MustCompile("^[Ss]ample[ -_]?[0-9].")
    for _, file := range singleTorrent.Files {
        downloadLink := tb.RequestDownloadLink(torrentID, file.ID, false)
        if downloadLink == "" {
            continue
        }

        guessedName := singleTorrent.Name
        mimeType := mime.TypeByExtension(guessedName[strings.LastIndex(guessedName, "."):])
        isVideo := strings.Contains(mimeType, "video")
        match := re.MatchString(guessedName)
        if !match && isVideo {
            un := types.Unrestricted{
                Filename:  guessedName,
                Filesize:  0,
                Download:  downloadLink,
            }
            result = append(result, un)
        }
    }

    return result
}
