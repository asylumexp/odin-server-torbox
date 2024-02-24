package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("jib32sgrokndtt2")
		if err != nil {
			return err
		}

		collection.ListRule = types.Pointer("user.id = @request.auth.id")

		collection.CreateRule = types.Pointer("")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("jib32sgrokndtt2")
		if err != nil {
			return err
		}

		collection.ListRule = nil

		collection.CreateRule = nil

		return dao.SaveCollection(collection)
	})
}
