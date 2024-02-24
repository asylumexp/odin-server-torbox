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
			"id": "gy7oibx0bg0num2",
			"created": "2023-12-21 16:39:00.621Z",
			"updated": "2023-12-21 16:39:00.621Z",
			"name": "devices",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "17ici2x4",
					"name": "user",
					"type": "relation",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"collectionId": "_pb_users_auth_",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": null
					}
				},
				{
					"system": false,
					"id": "iy8rm0eq",
					"name": "token",
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
					"id": "xqkkkixs",
					"name": "verified",
					"type": "bool",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {}
				}
			],
			"indexes": [
				"CREATE INDEX ` + "`" + `idx_Q19a20R` + "`" + ` ON ` + "`" + `devices` + "`" + ` (` + "`" + `user` + "`" + `)",
				"CREATE INDEX ` + "`" + `idx_JuaFpI7` + "`" + ` ON ` + "`" + `devices` + "`" + ` (` + "`" + `token` + "`" + `)",
				"CREATE INDEX ` + "`" + `idx_7TjpMBM` + "`" + ` ON ` + "`" + `devices` + "`" + ` (\n  ` + "`" + `user` + "`" + `,\n  ` + "`" + `token` + "`" + `\n)"
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

		collection, err := dao.FindCollectionByNameOrId("gy7oibx0bg0num2")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
