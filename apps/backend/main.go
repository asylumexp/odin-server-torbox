package main

import (
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
)

func getDevice(app *pocketbase.PocketBase, c echo.Context) (*models.Record, error) {
	device := c.Request().Header.Get("Device")
	return app.Dao().FindRecordById("devices", device)
}

func RequireDeviceOrRecordAuth(app *pocketbase.PocketBase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			record, _ := c.Get("authRecord").(*models.Record)
			d, _ := getDevice(app, c)
			verified := false
			if d != nil {
				verified = d.GetBool("verified")
			}

			if !verified && record == nil {
				return apis.NewBadRequestError("Verified device code or Auth are required", nil)
			}

			return next(c)
		}
	}
}

func main() {
	log.SetLevel(log.DebugLevel)
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

		scheduler.Start()

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
			// id := c.PathParam("id")
			// d, _ := app.Dao().FindRecordById("devices", id)
			// if d != nil {
			// 	d.Set("verified", true)
			// 	d.Save()
			// 	return c.JSON(http.StatusOK, d)
			// }
			return c.JSON(http.StatusNotFound, nil)
		})

		e.Router.GET("/sections/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		}, RequireDeviceOrRecordAuth(app))

		e.Router.Any("/_trakt/*", func(c echo.Context) error {
			info := apis.RequestInfo(c)
			id := ""
			if info.AuthRecord == nil {
				d, _ := getDevice(app, c)
				if d != nil {
					id = d.Get("user").(string)
				}
			} else {
				id = info.AuthRecord.Id
			}

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
				if helpers.ArrayContains([]string{
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
				if helpers.ArrayContains([]string{
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

	// go func() {
	// 	for {
	// 		// fmt.Println("hello world")
	// 		time.Sleep(1 * time.Second)
	// 		res := make([]User, 0)
	// 		err := app.Dao().DB().NewQuery("SELECT * FROM users").All(&res)
	// 		if err == nil {
	// 			// fmt.Println(res[0].Id)
	// 		}
	// 	}
	// }()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
