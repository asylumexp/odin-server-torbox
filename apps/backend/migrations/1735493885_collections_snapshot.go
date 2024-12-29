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
		jsonData := `[
			{
				"id": "_pb_users_auth_",
				"created": "2023-12-15 18:30:15.247Z",
				"updated": "2023-12-24 21:24:28.471Z",
				"name": "users",
				"type": "auth",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "users_name",
						"name": "name",
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
						"id": "users_avatar",
						"name": "avatar",
						"type": "file",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"mimeTypes": [
								"image/jpeg",
								"image/png",
								"image/svg+xml",
								"image/gif",
								"image/webp"
							],
							"thumbs": null,
							"maxSelect": 1,
							"maxSize": 5242880,
							"protected": false
						}
					},
					{
						"system": false,
						"id": "ebrrsidx",
						"name": "trakt_sections",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					},
					{
						"system": false,
						"id": "hpy6zodi",
						"name": "trakt_token",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					}
				],
				"indexes": [],
				"listRule": "id = @request.auth.id",
				"viewRule": "id = @request.auth.id",
				"createRule": null,
				"updateRule": "id = @request.auth.id",
				"deleteRule": "id = @request.auth.id",
				"options": {
					"allowEmailAuth": true,
					"allowOAuth2Auth": true,
					"allowUsernameAuth": true,
					"exceptEmailDomains": null,
					"manageRule": null,
					"minPasswordLength": 8,
					"onlyEmailDomains": null,
					"onlyVerified": false,
					"requireEmail": false
				}
			},
			{
				"id": "kog9lz3zq2kj07s",
				"created": "2023-12-15 18:32:16.845Z",
				"updated": "2024-12-25 07:47:18.193Z",
				"name": "settings",
				"type": "base",
				"system": false,
				"schema": [
					{
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
					},
					{
						"system": false,
						"id": "wfnmmcoj",
						"name": "real_debrid",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id != \"\"",
				"viewRule": "@request.auth.id != \"\"",
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "uvaxc6sfjgaxw9c",
				"created": "2023-12-15 18:34:04.439Z",
				"updated": "2024-11-03 08:06:04.441Z",
				"name": "tmdb",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "txaj4w1q",
						"name": "data",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					},
					{
						"system": false,
						"id": "5e3mgwxm",
						"name": "type",
						"type": "select",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"movie",
								"show"
							]
						}
					},
					{
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
					}
				],
				"indexes": [
					"CREATE INDEX ` + "`" + `idx_1ir7EEI` + "`" + ` ON ` + "`" + `tmdb` + "`" + ` (` + "`" + `id` + "`" + `)",
					"CREATE INDEX ` + "`" + `idx_eWT8ulS` + "`" + ` ON ` + "`" + `tmdb` + "`" + ` (` + "`" + `tmdb_id` + "`" + `)",
					"CREATE INDEX ` + "`" + `idx_7QPyCnz` + "`" + ` ON ` + "`" + `tmdb` + "`" + ` (` + "`" + `type` + "`" + `)"
				],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "jib32sgrokndtt2",
				"created": "2023-12-20 14:59:04.811Z",
				"updated": "2023-12-20 16:23:08.540Z",
				"name": "history",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "xfytwctq",
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
						"id": "axymvzt6",
						"name": "data",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					},
					{
						"system": false,
						"id": "vv5en8cb",
						"name": "type",
						"type": "select",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"movie",
								"episode"
							]
						}
					},
					{
						"system": false,
						"id": "293urao6",
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
					},
					{
						"system": false,
						"id": "spgro5mw",
						"name": "watched_at",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "09ahphxh",
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
					},
					{
						"system": false,
						"id": "xeh8twko",
						"name": "show_id",
						"type": "number",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"noDecimal": false
						}
					}
				],
				"indexes": [
					"CREATE INDEX ` + "`" + `idx_SHZ0Jzt` + "`" + ` ON ` + "`" + `history` + "`" + ` (` + "`" + `show_id` + "`" + `)",
					"CREATE INDEX ` + "`" + `idx_WaWH7nI` + "`" + ` ON ` + "`" + `history` + "`" + ` (` + "`" + `trakt_id` + "`" + `)"
				],
				"listRule": "user.id = @request.auth.id",
				"viewRule": "user.id = @request.auth.id",
				"createRule": "",
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "gy7oibx0bg0num2",
				"created": "2023-12-21 16:39:00.621Z",
				"updated": "2023-12-21 17:11:57.884Z",
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
					},
					{
						"system": false,
						"id": "u20rxpwn",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"indexes": [
					"CREATE INDEX ` + "`" + `idx_Q19a20R` + "`" + ` ON ` + "`" + `devices` + "`" + ` (` + "`" + `user` + "`" + `)",
					"CREATE INDEX ` + "`" + `idx_JuaFpI7` + "`" + ` ON ` + "`" + `devices` + "`" + ` (` + "`" + `token` + "`" + `)",
					"CREATE INDEX ` + "`" + `idx_7TjpMBM` + "`" + ` ON ` + "`" + `devices` + "`" + ` (\n  ` + "`" + `user` + "`" + `,\n  ` + "`" + `token` + "`" + `\n)"
				],
				"listRule": "user.id = @request.auth.id",
				"viewRule": "user.id = @request.auth.id",
				"createRule": "user.id = @request.auth.id",
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "7n09aifcpshv459",
				"created": "2024-02-22 10:21:47.021Z",
				"updated": "2024-10-18 09:43:52.826Z",
				"name": "trakt_seasons",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "hjqjynor",
						"name": "data",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					},
					{
						"system": false,
						"id": "uaatrqia",
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
					}
				],
				"indexes": [
					"CREATE INDEX ` + "`" + `idx_O92766T` + "`" + ` ON ` + "`" + `trakt_seasons` + "`" + ` (` + "`" + `trakt_id` + "`" + `)"
				],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "n434wzwv45po9ib",
				"created": "2024-10-18 09:42:49.313Z",
				"updated": "2024-10-18 11:54:29.844Z",
				"name": "rd_resolved",
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
						"name": "data",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					},
					{
						"system": false,
						"id": "ncpi3amv",
						"name": "magnet",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"indexes": [
					"CREATE INDEX ` + "`" + `idx_PgBvcoT` + "`" + ` ON ` + "`" + `rd_resolved` + "`" + ` (` + "`" + `resource` + "`" + `)"
				],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			}
		]`

		collections := []*models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collections); err != nil {
			return err
		}

		return daos.New(db).ImportCollections(collections, true, nil)
	}, func(db dbx.Builder) error {
		return nil
	})
}
