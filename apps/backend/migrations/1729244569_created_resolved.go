package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "n434wzwv45po9ib",
			"created": "2024-10-18 09:42:49.313Z",
			"updated": "2024-10-18 09:42:49.313Z",
			"name": "resolved",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "pyyh9opk",
					"name": "resource",
					"type": "text",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
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
				}
			],
			"indexes": [
				"CREATE INDEX ` + "`" + `idx_PgBvcoT` + "`" + ` ON ` + "`" + `resolved` + "`" + ` (` + "`" + `resource` + "`" + `)"
			],
			"listRule": null,
			"viewRule": null,
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("n434wzwv45po9ib")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
