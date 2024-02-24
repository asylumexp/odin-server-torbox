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
		new_tmdb := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
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
		}`), new_tmdb)
		collection.Schema.AddField(new_tmdb)

		// add
		new_fanart := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
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
		}`), new_fanart)
		collection.Schema.AddField(new_fanart)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("kog9lz3zq2kj07s")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("r7rf5tvg")

		// remove
		collection.Schema.RemoveField("afgdj4nj")

		return dao.SaveCollection(collection)
	})
}
