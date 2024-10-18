package helpers

import (
	"fmt"
	"os/user"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/odin-movieshow/server/types"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

func GetHomeDir() string {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return currentUser.HomeDir
}

func ReadTmdbCache(app *pocketbase.PocketBase, id uint, resource string) interface{} {
	record, err := app.Dao().
		FindFirstRecordByFilter("tmdb", "tmdb_id = {:id} && type = {:type}", dbx.Params{"id": id, "type": resource})
	res := make(map[string]any)
	if err == nil {
		err := record.UnmarshalJSONField("data", &res)
		date := record.GetDateTime("updated")
		now := time.Now()
		if err == nil {
			if date.Time().Before(now.AddDate(0, 0, -1)) {
				return nil
			}
			log.Debug("cache hit", "for", "tmdb", "resource", resource, "id", id)

			return res
		}
	}
	return nil
}

func WriteTmdbCache(app *pocketbase.PocketBase, id uint, resource string, data *interface{}) {
	log.Info("cache write", "for", "tmdb", "resource", resource, "id", id)
	record, err := app.Dao().
		FindFirstRecordByFilter("tmdb", "tmdb_id = {:id} && type = {:type}", dbx.Params{"id": id, "type": resource})

	if err == nil {
		record.Set("data", &data)
		app.Dao().SaveRecord(record)
	} else {
		collection, _ := app.Dao().FindCollectionByNameOrId("tmdb")
		record := models.NewRecord(collection)
		record.Set("data", &data)
		record.Set("tmdb_id", id)
		record.Set("type", resource)
		app.Dao().SaveRecord(record)
	}

}

func WriteTraktSeasonCache(app *pocketbase.PocketBase, id uint, data *interface{}) {
	log.Info("cache write", "for", "trakt", "resource", "show_seasons", "id", id)
	record, err := app.Dao().
		FindFirstRecordByFilter("trakt_seasons", "trakt_id = {:id}", dbx.Params{"id": id})

	if err == nil {
		record.Set("data", &data)
		app.Dao().SaveRecord(record)
	} else {
		collection, _ := app.Dao().FindCollectionByNameOrId("trakt_seasons")
		record := models.NewRecord(collection)
		record.Set("data", &data)
		record.Set("trakt_id", id)
		app.Dao().SaveRecord(record)
	}

}

func ReadRDCache(app *pocketbase.PocketBase, resource string, magnet string) *types.Torrent {
	record, err := app.Dao().
		FindFirstRecordByFilter("rd_resolved", "magnet = {:magnet}", dbx.Params{"magnet": magnet})
	var res types.Torrent
	if err == nil {
		err := record.UnmarshalJSONField("data", &res)
		date := record.GetDateTime("updated")
		now := time.Now().Add(time.Duration((-2) * time.Hour))
		if err == nil {
			if date.Time().Before(now) {
				return nil
			}
			log.Debug("cache hit", "for", "RD", "resource", resource)
			return &res
		}
	}
	return nil
}

func ReadRDCacheByResource(app *pocketbase.PocketBase, resource string) []types.Torrent {
	records, err := app.Dao().FindRecordsByFilter("rd_resolved", "resource = {:resource}", "id", -1, 0, dbx.Params{"resource": resource})
	res := make([]types.Torrent, 0)
	if err == nil {
		for _, record := range records {
			var r types.Torrent
			date := record.GetDateTime("updated")
			now := time.Now().Add(time.Duration((-2) * time.Hour))
			// add 1 hour to date
			if date.Time().Before(now) {
				continue
			}
			err := record.UnmarshalJSONField("data", &r)
			if err == nil {
				res = append(res, r)
			}
		}

	}
	return res
}

func WriteRDCache(app *pocketbase.PocketBase, resource string, magnet string, data interface{}) {
	log.Info("cache write", "for", "RD", "resource", resource)
	record, err := app.Dao().
		FindFirstRecordByFilter("rd_resolved", "magnet = {:magnet}", dbx.Params{"magnet": magnet})

	if err == nil {
		record.Set("data", &data)
		app.Dao().SaveRecord(record)
	} else {
		collection, _ := app.Dao().FindCollectionByNameOrId("rd_resolved")
		record := models.NewRecord(collection)
		record.Set("data", &data)
		record.Set("magnet", magnet)
		record.Set("resource", resource)
		app.Dao().SaveRecord(record)
	}
}

func ReadTraktSeasonCache(app *pocketbase.PocketBase, id uint) []any {
	record, err := app.Dao().
		FindFirstRecordByFilter("trakt_seasons", "trakt_id = {:id}", dbx.Params{"id": id})
	res := make([]any, 0)
	if err == nil {
		err := record.UnmarshalJSONField("data", &res)
		date := record.GetDateTime("updated")
		now := time.Now()

		if err == nil {
			if date.Time().Before(now.AddDate(0, 0, -1)) {
				return nil
			}
			log.Debug("cache hit", "for", "trakt", "resource", "show_seasons", "id", id)
			return res
		}
	}
	return nil
}

func ParseDates(str string) string {

	re := regexp.MustCompile("::(year|month|day):(\\+|-)?(\\d+)?:")

	matches := re.FindAllStringSubmatch(str, -1)
	now := time.Now()
	for _, match := range matches {

		yearVal := 0
		monthVal := 0
		dayVal := 0
		if len(match) == 4 {
			val := 0
			if v, err := strconv.Atoi(match[3]); err == nil {
				val = v
			}
			if match[2] == "-" {
				val *= -1
			}
			if match[1] == "year" {

				yearVal = val
				str = strings.ReplaceAll(str, match[0], "#year#")
			} else if match[1] == "month" {
				monthVal = val
				str = strings.ReplaceAll(str, match[0], "#month#")
			} else if match[1] == "day" {
				dayVal = val
				str = strings.ReplaceAll(str, match[0], "#day#")
			}
		}
		now = now.AddDate(yearVal, monthVal, dayVal)
		str = strings.ReplaceAll(str, "#year#", fmt.Sprintf("%d", now.Year()))
		str = strings.ReplaceAll(str, "#month#", fmt.Sprintf("%d", now.Month()))
		str = strings.ReplaceAll(str, "#day#", fmt.Sprintf("%d", now.Day()))

	}

	re2 := regexp.MustCompile("::monthdays::")

	matches2 := re2.FindAllStringSubmatch(str, -1)
	dinm := daysInMonth(now)
	for _, match := range matches2 {
		str = strings.ReplaceAll(str, match[0], fmt.Sprintf("%d", dinm))
	}

	return str
}

func daysInMonth(t time.Time) int {
	t = time.Date(t.Year(), t.Month(), 32, 0, 0, 0, 0, time.UTC)
	daysInMonth := 32 - t.Day()
	days := make([]int, daysInMonth)
	for i := range days {
		days[i] = i + 1
	}

	d := days[len(days)-1]
	d += 1
	return d
}

func Chunk[T any](slice []T) [][]T {
	size := ItemsPerChunk(slice)
	var chunks [][]T
	for i := 0; i < len(slice); {
		// Clamp the last chunk to the slice bound as necessary.
		end := size
		if l := len(slice[i:]); l < size {
			end = l
		}

		// Set the capacity of each chunk so that appending to a chunk does not
		// modify the original slice.
		chunks = append(chunks, slice[i:i+end:i+end])
		i += end
	}

	return chunks
}

func ItemsPerChunk[T any](slice []T) int {
	batchSize := runtime.NumCPU() / 2
	if len(slice) < batchSize {
		batchSize = 1
	}
	return len(slice) / batchSize
}
