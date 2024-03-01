package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"backend/helpers"
	"backend/realdebrid"
	"backend/scraper"
	"backend/tmdb"
	"backend/trakt"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/cron"
	"github.com/thoas/go-funk"
)

func getDevice(app *pocketbase.PocketBase, c echo.Context) (*models.Record, error) {
	device := c.Request().Header.Get("Device")
	return app.Dao().FindRecordById("devices", device)
}

func RequireDeviceOrRecordAuth(app *pocketbase.PocketBase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			record, _ := c.Get("authRecord").(*models.Record)
			if record == nil {
				d, _ := getDevice(app, c)
				if d != nil {
					if d.GetBool("verified") {
						u, err := app.Dao().FindRecordById("users", d.Get("user").(string))
						if err == nil {
							c.Set("authRecord", u)
						}
					}
				}
			}

			if c.Get("authRecord") == nil {
				return apis.NewBadRequestError("Verified device code or Auth are required", nil)
			}

			return next(c)
		}
	}
}

func main() {
	l, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err == nil {
		log.SetLevel(l)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	conf := pocketbase.Config{}
	app := pocketbase.NewWithConfig(conf)
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {

		scheduler := cron.New()
		scheduler.MustAdd("hourly", "0 * * * *", func() {
			trakt.RefreshTokens(app)
			realdebrid.RefreshTokens(app)
			trakt.SyncHistory(app)

		})

		// scheduler.MustAdd("daily", "0 0 * * *", func() {
		// 	realdebrid.Cleanup(app)
		// })

		scheduler.Start()

		go func() {
		}()

		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
		e.Router.POST("/scrape", func(c echo.Context) error {
			log.Debug("Scraping")
			data := scraper.GetLinks(apis.RequestInfo(c).Data, app)
			return c.JSON(http.StatusOK, data)
		}, RequireDeviceOrRecordAuth(app))

		e.Router.GET("/device", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]any{"test": "ok"})
		})

		e.Router.GET("/device/verify/:id", func(c echo.Context) error {
			id := c.PathParam("id")
			d, err := app.Dao().FindRecordById("devices", id)
			if err == nil {
				d.Set("verified", true)
				app.Dao().SaveRecord(d)
				log.Info("Device verified", "id", id)
				return c.JSON(http.StatusOK, d)
			}
			return c.JSON(http.StatusNotFound, nil)
		}, apis.RequireGuestOnly())

		e.Router.GET("/sections/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		}, RequireDeviceOrRecordAuth(app))

		e.Router.GET("/user", func(c echo.Context) error {
			u := c.Get("authRecord")
			sections := make(map[string]any)
			err := u.(*models.Record).UnmarshalJSONField("trakt_sections", &sections)

			if err == nil {
				for _, t := range []string{"home", "movies", "shows"} {
					s := sections[t].([]any)
					if s != nil {
						for i := range s {
							sections[t].([]any)[i].(map[string]any)["title"] = helpers.ParseDates(
								sections[t].([]any)[i].(map[string]any)["title"].(string),
							)
							sections[t].([]any)[i].(map[string]any)["url"] = helpers.ParseDates(
								sections[t].([]any)[i].(map[string]any)["url"].(string),
							)
						}
					}

				}

				str, err := json.Marshal(sections)
				if err == nil {
					u.(*models.Record).Set("trakt_sections", string(str))
				}
			}

			return c.JSON(http.StatusOK, u)
		}, RequireDeviceOrRecordAuth(app))

		e.Router.Any("/_trakt/*", func(c echo.Context) error {
			info := apis.RequestInfo(c)

			id := info.AuthRecord.Id

			t := make(map[string]any)
			u, _ := app.Dao().FindRecordById("users", id)
			u.UnmarshalJSONField("trakt_token", &t)
			trakt.Headers = apis.RequestInfo(c).Headers
			trakt.Headers["authorization"] = "Bearer " + t["access_token"].(string)

			trakt.FetchSeasons = true
			trakt.FetchTMDB = true

			jsonData := apis.RequestInfo(c).Data
			url := strings.ReplaceAll(c.Request().URL.String(), "/_trakt", "")
			if strings.Contains(url, "scrobble") {
				go func() {
					trakt.SyncHistory(app)
				}()
			}
			url = helpers.ParseDates(url)
			result, headers, status := trakt.CallEndpoint(
				url,
				c.Request().Method,
				jsonData,
				true,
				app,
			)

			for k, v := range headers {
				if funk.Contains([]string{
					"Content-Encoding",
					"Access-Control-Allow-Origin",
				}, k) {
					continue
				}
				c.Response().Header().Add(k, v[0])
			}
			c.Response().Status = status

			return c.JSON(http.StatusOK, result)
		}, RequireDeviceOrRecordAuth(app))

		e.Router.Any("/_realdebrid/*", func(c echo.Context) error {
			url := strings.ReplaceAll(c.Request().URL.String(), "/_realdebrid", "")
			result, headers, status := realdebrid.CallEndpoint(url, c.Request().Method, nil, app)

			for k, v := range headers {
				if funk.Contains([]string{
					"Content-Encoding",
					"Access-Control-Allow-Origin",
				}, k) {
					continue
				}
				c.Response().Header().Add(k, v[0])
			}
			c.Response().Status = status
			return c.JSON(http.StatusOK, result)
		}, apis.RequireAdminAuth())

		e.Router.GET("/tmdbseasons/:id", func(c echo.Context) error {
			fmt.Println(c.PathParam("id"))
			seasons := c.QueryParam("seasons")
			res := tmdb.GetEpisodes(c.PathParam("id"), strings.Split(seasons, ","), app)
			return c.JSON(http.StatusOK, res)
		}, RequireDeviceOrRecordAuth(app))

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
