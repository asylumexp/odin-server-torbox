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

		collection, err := dao.FindCollectionByNameOrId("gdnsibibegu5b4y")
		if err != nil {
			return err
		}

		// add
		new_last_history := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "ysiazqky",
			"name": "last_history",
			"type": "date",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), new_last_history)
		collection.Schema.AddField(new_last_history)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("gdnsibibegu5b4y")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("ysiazqky")

		return dao.SaveCollection(collection)
	})
}
