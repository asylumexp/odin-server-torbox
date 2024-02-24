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

		// add
		new_runtime := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "ppb4l8ab",
			"name": "runtime",
			"type": "number",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"noDecimal": false
			}
		}`), new_runtime)
		collection.Schema.AddField(new_runtime)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("8jn6qdaobuh7y0r")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("ppb4l8ab")

		return dao.SaveCollection(collection)
	})
}
