package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("n434wzwv45po9ib")
		if err != nil {
			return err
		}

		collection.Name = "rd_resolved"

		json.Unmarshal([]byte(`[
			"CREATE INDEX ` + "`" + `idx_PgBvcoT` + "`" + ` ON ` + "`" + `rd_resolved` + "`" + ` (` + "`" + `resource` + "`" + `)"
		]`), &collection.Indexes)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("n434wzwv45po9ib")
		if err != nil {
			return err
		}

		collection.Name = "resolved"

		json.Unmarshal([]byte(`[
			"CREATE INDEX ` + "`" + `idx_PgBvcoT` + "`" + ` ON ` + "`" + `resolved` + "`" + ` (` + "`" + `resource` + "`" + `)"
		]`), &collection.Indexes)

		return dao.SaveCollection(collection)
	})
}
