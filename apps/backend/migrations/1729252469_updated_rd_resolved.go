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

		collection, err := dao.FindCollectionByNameOrId("n434wzwv45po9ib")
		if err != nil {
			return err
		}

		// add
		new_magnet := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "ncpi3amv",
			"name": "magnet",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_magnet)
		collection.Schema.AddField(new_magnet)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("n434wzwv45po9ib")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("ncpi3amv")

		return dao.SaveCollection(collection)
	})
}
