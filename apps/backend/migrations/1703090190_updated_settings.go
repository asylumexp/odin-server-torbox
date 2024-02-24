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

		// update
		edit_app := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "qcycymzs",
			"name": "app",
			"type": "json",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSize": 2000000
			}
		}`), edit_app)
		collection.Schema.AddField(edit_app)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("kog9lz3zq2kj07s")
		if err != nil {
			return err
		}

		// update
		edit_app := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "qcycymzs",
			"name": "data",
			"type": "json",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSize": 2000000
			}
		}`), edit_app)
		collection.Schema.AddField(edit_app)

		return dao.SaveCollection(collection)
	})
}
