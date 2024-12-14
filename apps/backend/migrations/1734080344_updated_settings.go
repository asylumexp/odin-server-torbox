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
		new_all_debrid := &schema.SchemaField{}
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
		}`), new_all_debrid); err != nil {
			return err
		}
		collection.Schema.AddField(new_all_debrid)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("kog9lz3zq2kj07s")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("2pryfjne")

		return dao.SaveCollection(collection)
	})
}
