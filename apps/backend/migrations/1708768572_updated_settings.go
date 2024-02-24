package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("kog9lz3zq2kj07s")
		if err != nil {
			return err
		}

		// add
		new_scraper_url := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "p0ysoxwr",
			"name": "scraper_url",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_scraper_url)
		collection.Schema.AddField(new_scraper_url)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("kog9lz3zq2kj07s")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("p0ysoxwr")

		return dao.SaveCollection(collection)
	})
}
