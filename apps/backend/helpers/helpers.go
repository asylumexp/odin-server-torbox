package helpers

import (
	"fmt"
	"os/user"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

func ArrayContains(alpha []string, str string) bool {

	// iterate using the for loop
	for i := 0; i < len(alpha); i++ {
		// check
		if alpha[i] == str {
			// return true
			return true
		}
	}
	return false
}

func IntArrayContains(alpha []int, str int) bool {

	// iterate using the for loop
	for i := 0; i < len(alpha); i++ {
		// check
		if alpha[i] == str {
			// return true
			return true
		}
	}
	return false
}

func FloatArrayContains(alpha []float64, str float64) bool {

	// iterate using the for loop
	for i := 0; i < len(alpha); i++ {
		// check
		if alpha[i] == str {
			// return true
			return true
		}
	}
	return false
}

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
		record.UnmarshalJSONField("data", &res)
		log.Debug("cache hit", "for", "tmdb", "resource", resource, "id", id)
		return res
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
	log.Info("cache write", "for", "trakt", "resource", "show", "id", id)
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

func ReadTraktSeasonCache(app *pocketbase.PocketBase, id uint) []any {
	record, err := app.Dao().
		FindFirstRecordByFilter("trakt_seasons", "trakt_id = {:id}", dbx.Params{"id": id})
	res := make([]any, 0)
	if err == nil {
		record.UnmarshalJSONField("data", &res)
		log.Debug("cache hit", "for", "trakt", "resource", "show", "id", id)
		return res
	}
	return nil
}

func ParseDates(str string) string {
	now := time.Now()

	re := regexp.MustCompile("::(year|month|day):(\\+|-)?(\\d+)?:")

	matches := re.FindAllStringSubmatch(str, -1)

	for _, match := range matches {
		if len(match) == 4 {
			val := 0
			if v, err := strconv.Atoi(match[3]); err == nil {
				val = v
			}
			if match[2] == "-" {
				val *= -1
			}
			if match[1] == "year" {
				now = now.AddDate(val, 0, 0)
				str = strings.ReplaceAll(str, match[0], fmt.Sprintf("%d", now.Year()))
			} else if match[1] == "month" {
				now = now.AddDate(0, val, 0)
				str = strings.ReplaceAll(str, match[0], fmt.Sprintf("%d", int(now.Month())))
			} else if match[1] == "day" {
				now = now.AddDate(0, 0, val)
				str = strings.ReplaceAll(str, match[0], fmt.Sprintf("%d", now.Day()))
			}
		}

	}

	return str
}
