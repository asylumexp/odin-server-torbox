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

		collection, err := dao.FindCollectionByNameOrId("jib32sgrokndtt2")
		if err != nil {
			return err
		}

		// add
		new_season := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "hpieoqru",
			"name": "season",
			"type": "number",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"noDecimal": false
			}
		}`), new_season)
		collection.Schema.AddField(new_season)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("jib32sgrokndtt2")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("hpieoqru")

		return dao.SaveCollection(collection)
	})
}
