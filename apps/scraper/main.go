package main

import (
	"os"

	"github.com/odin-movieshow/scraper/jackett"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
)

func main() {
	log.SetLevel(log.DebugLevel)
	jackettUrl := os.Getenv("JACKETT_URL")
	jackettKey := os.Getenv("JACKETT_KEY")

	if jackettUrl == "" || jackettKey == "" {
		log.Error("missing env vars JACKETT_URL and JACKETT_KEY")
		os.Exit(0)
	}
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("odin-craper is up and running!")
	})
	app.Post("/scrape", jackett.Search)

	app.Listen(":6969")
}
