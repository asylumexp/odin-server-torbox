package main

import (
	"encoding/json"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/odin-movieshow/server/helpers"
	"github.com/odin-movieshow/server/imdb"
	"github.com/odin-movieshow/server/realdebrid"
	"github.com/odin-movieshow/server/scraper"
	"github.com/odin-movieshow/server/tmdb"
	"github.com/odin-movieshow/server/trakt"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/cron"
	"github.com/thoas/go-funk"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func mqttclient() mqtt.Client {
	// mqtt.DEBUG = stdlog.New(os.Stdout, "", 0)
	mqtt.ERROR = stdlog.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().
		AddBroker(os.Getenv("MQTT_URL")).
		SetUsername(os.Getenv("MQTT_USER")).
		SetPassword(os.Getenv("MQTT_PASSWORD"))
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Error("MQTT", "conneced", c.IsConnected())
	} else {
		log.Info("MQTT", "connected", c.IsConnected(), "url", os.Getenv("MQTT_URL"))
	}

	return c
}

func getDevice(app *pocketbase.PocketBase, c echo.Context) (*models.Record, error) {
	device := c.Request().Header.Get("Device")
	return app.Dao().FindRecordById("devices", device)
}

func RequireDeviceOrRecordAuth(app *pocketbase.PocketBase) echo.MiddlewareFunc {
	mut := sync.RWMutex{}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			record, _ := c.Get("authRecord").(*models.Record)
			if record == nil {
				d, _ := getDevice(app, c)
				if d != nil {
					if d.GetBool("verified") {
						u, err := app.Dao().FindRecordById("users", d.Get("user").(string))
						if err == nil {
							mut.Lock()
							c.Set("authRecord", u)
							mut.Unlock()
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
	godotenv.Load()

	log.SetReportCaller(true)
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://308c965810583274884cbc87d1a584de@sentry.dnmc.in/4",
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Error("sentry.Init: %s", err)
	}

	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	sentry.CaptureException(fmt.Errorf("This is a test exception"))
	sentry.CaptureEvent(&sentry.Event{
		Message: "This is a test error event",
		Level:   sentry.LevelError,
	})
	sentry.CaptureMessage("This is a test message")

	l, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err == nil {
		log.SetLevel(l)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if os.Getenv("BACKEND_URL") == "" {
		log.Fatal("BACKEND_URL is required")
		os.Exit(0)
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

		scheduler.MustAdd("daily", "0 4 * * *", func() {
			realdebrid.Cleanup(app)
		})

		scheduler.Start()

		go func() {
		}()

		a, _ := app.Dao().FindAdminByEmail("admin@odin.local")
		a.SetPassword("adminOdin1")
		app.Dao().SaveAdmin(a)

		mq := mqttclient()
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
		e.Router.POST("/scrape", func(c echo.Context) error {
			log.Debug("Scraping")
			data := scraper.GetLinks(apis.RequestInfo(c).Data, app, mq)
			return c.JSON(http.StatusOK, data)
		}, RequireDeviceOrRecordAuth(app))

		e.Router.GET("/imdb/:id", func(c echo.Context) error {
			id := c.PathParam("id")
			return c.JSON(http.StatusOK, imdb.Get(id))
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

		e.Router.GET("/backendurl", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]any{"url": os.Getenv("BACKEND_URL")})
		}, RequireDeviceOrRecordAuth(app))

		e.Router.GET("/mqttconfig", func(c echo.Context) error {
			return c.JSON(
				http.StatusOK,
				map[string]any{
					"url":      os.Getenv("MQTT_URL"),
					"user":     os.Getenv("MQTT_USER"),
					"password": os.Getenv("MQTT_PASSWORD"),
					"port":     443,
				},
			)
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
			url := strings.ReplaceAll(c.Request().URL.String(), "/_trakt", "")
			trakt.Headers = apis.RequestInfo(c).Headers
			// delete passed header of pocketbase

			t := make(map[string]any)
			u, _ := app.Dao().FindRecordById("users", id)
			u.UnmarshalJSONField("trakt_token", &t)
			delete(trakt.Headers, "authorization")

			if t != nil && t["access_token"] != nil {
				trakt.Headers["authorization"] = "Bearer " + t["access_token"].(string)
			}
			if strings.Contains(url, "fresh=true") {
				delete(trakt.Headers, "authorization")
			}

			trakt.FetchSeasons = true
			trakt.FetchTMDB = true

			jsonData := apis.RequestInfo(c).Data

			if strings.Contains(url, "scrobble/stop") {
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
