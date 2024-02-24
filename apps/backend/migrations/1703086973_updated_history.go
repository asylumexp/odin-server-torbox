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

		collection, err := dao.FindCollectionByNameOrId("jib32sgrokndtt2")
		if err != nil {
			return err
		}

		json.Unmarshal([]byte(`[
			"CREATE INDEX ` + "`" + `idx_SHZ0Jzt` + "`" + ` ON ` + "`" + `history` + "`" + ` (` + "`" + `show_id` + "`" + `)",
			"CREATE INDEX ` + "`" + `idx_WaWH7nI` + "`" + ` ON ` + "`" + `history` + "`" + ` (` + "`" + `trakt_id` + "`" + `)"
		]`), &collection.Indexes)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("jib32sgrokndtt2")
		if err != nil {
			return err
		}

		json.Unmarshal([]byte(`[]`), &collection.Indexes)

		return dao.SaveCollection(collection)
	})
}
