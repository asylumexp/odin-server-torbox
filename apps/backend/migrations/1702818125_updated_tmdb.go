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

		collection, err := dao.FindCollectionByNameOrId("uvaxc6sfjgaxw9c")
		if err != nil {
			return err
		}

		// add
		new_tmdb_id := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "0li9qnvg",
			"name": "tmdb_id",
			"type": "number",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"noDecimal": false
			}
		}`), new_tmdb_id)
		collection.Schema.AddField(new_tmdb_id)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("uvaxc6sfjgaxw9c")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("0li9qnvg")

		return dao.SaveCollection(collection)
	})
}
