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

		// remove
		collection.Schema.RemoveField("cmcphauz")

		// remove
		collection.Schema.RemoveField("r7rf5tvg")

		// remove
		collection.Schema.RemoveField("afgdj4nj")

		// remove
		collection.Schema.RemoveField("wvhuvvrj")

		// remove
		collection.Schema.RemoveField("p0ysoxwr")

		// remove
		collection.Schema.RemoveField("2pryfjne")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("kog9lz3zq2kj07s")
		if err != nil {
			return err
		}

		// add
		del_trakt := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "cmcphauz",
			"name": "trakt",
			"type": "json",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSize": 2000000
			}
		}`), del_trakt); err != nil {
			return err
		}
		collection.Schema.AddField(del_trakt)

		// add
		del_tmdb := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "r7rf5tvg",
			"name": "tmdb",
			"type": "json",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSize": 2000000
			}
		}`), del_tmdb); err != nil {
			return err
		}
		collection.Schema.AddField(del_tmdb)

		// add
		del_fanart := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "afgdj4nj",
			"name": "fanart",
			"type": "json",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSize": 2000000
			}
		}`), del_fanart); err != nil {
			return err
		}
		collection.Schema.AddField(del_fanart)

		// add
		del_jackett := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "wvhuvvrj",
			"name": "jackett",
			"type": "json",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSize": 2000000
			}
		}`), del_jackett); err != nil {
			return err
		}
		collection.Schema.AddField(del_jackett)

		// add
		del_scraper_url := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
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
		}`), del_scraper_url); err != nil {
			return err
		}
		collection.Schema.AddField(del_scraper_url)

		// add
		del_all_debrid := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "2pryfjne",
			"name": "all_debrid",
			"type": "json",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSize": 2000000
			}
		}`), del_all_debrid); err != nil {
			return err
		}
		collection.Schema.AddField(del_all_debrid)

		return dao.SaveCollection(collection)
	})
}
