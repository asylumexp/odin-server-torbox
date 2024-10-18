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

		// update
		edit_data := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "kkyx6hfp",
			"name": "data",
			"type": "json",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSize": 2000000
			}
		}`), edit_data)
		collection.Schema.AddField(edit_data)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("n434wzwv45po9ib")
		if err != nil {
			return err
		}

		// update
		edit_data := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "kkyx6hfp",
			"name": "links",
			"type": "json",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSize": 2000000
			}
		}`), edit_data)
		collection.Schema.AddField(edit_data)

		return dao.SaveCollection(collection)
	})
}
