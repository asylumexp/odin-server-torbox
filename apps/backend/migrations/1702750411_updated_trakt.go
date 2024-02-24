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

		// update
		edit_token := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "morbfnbp",
			"name": "token",
			"type": "json",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSize": 2000000
			}
		}`), edit_token)
		collection.Schema.AddField(edit_token)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("gdnsibibegu5b4y")
		if err != nil {
			return err
		}

		// update
		edit_token := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "morbfnbp",
			"name": "data",
			"type": "json",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSize": 2000000
			}
		}`), edit_token)
		collection.Schema.AddField(edit_token)

		return dao.SaveCollection(collection)
	})
}
