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
		new_trakt_clientId := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "hrdh9lv5",
			"name": "trakt_clientId",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_trakt_clientId)
		collection.Schema.AddField(new_trakt_clientId)

		// add
		new_trakt_clientSecret := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "gnl3hohb",
			"name": "trakt_clientSecret",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_trakt_clientSecret)
		collection.Schema.AddField(new_trakt_clientSecret)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("kog9lz3zq2kj07s")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("hrdh9lv5")

		// remove
		collection.Schema.RemoveField("gnl3hohb")

		return dao.SaveCollection(collection)
	})
}
