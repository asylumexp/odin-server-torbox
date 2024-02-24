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

		collection, err := dao.FindCollectionByNameOrId("8jn6qdaobuh7y0r")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("okqjgflm")

		// add
		new_trakt_id := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "wzppxkyc",
			"name": "trakt_id",
			"type": "number",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"noDecimal": false
			}
		}`), new_trakt_id)
		collection.Schema.AddField(new_trakt_id)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("8jn6qdaobuh7y0r")
		if err != nil {
			return err
		}

		// add
		del_trakt_id := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "okqjgflm",
			"name": "trakt_id",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), del_trakt_id)
		collection.Schema.AddField(del_trakt_id)

		// remove
		collection.Schema.RemoveField("wzppxkyc")

		return dao.SaveCollection(collection)
	})
}
