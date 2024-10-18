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

		collection, err := dao.FindCollectionByNameOrId("uvaxc6sfjgaxw9c")
		if err != nil {
			return err
		}

		json.Unmarshal([]byte(`[
			"CREATE INDEX ` + "`" + `idx_1ir7EEI` + "`" + ` ON ` + "`" + `tmdb` + "`" + ` (` + "`" + `id` + "`" + `)",
			"CREATE INDEX ` + "`" + `idx_eWT8ulS` + "`" + ` ON ` + "`" + `tmdb` + "`" + ` (` + "`" + `tmdb_id` + "`" + `)"
		]`), &collection.Indexes)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("uvaxc6sfjgaxw9c")
		if err != nil {
			return err
		}

		json.Unmarshal([]byte(`[
			"CREATE INDEX ` + "`" + `idx_1ir7EEI` + "`" + ` ON ` + "`" + `tmdb` + "`" + ` (` + "`" + `id` + "`" + `)"
		]`), &collection.Indexes)

		return dao.SaveCollection(collection)
	})
}
