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

		collection, err := dao.FindCollectionByNameOrId("gdnsibibegu5b4y")
		if err != nil {
			return err
		}

		collection.ViewRule = types.Pointer("user.id = @request.auth.id")

		collection.CreateRule = types.Pointer("user.id = @request.auth.id")

		collection.UpdateRule = types.Pointer("user.id = @request.auth.id")

		collection.DeleteRule = types.Pointer("user.id = @request.auth.id")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("gdnsibibegu5b4y")
		if err != nil {
			return err
		}

		collection.ViewRule = nil

		collection.CreateRule = nil

		collection.UpdateRule = nil

		collection.DeleteRule = nil

		return dao.SaveCollection(collection)
	})
}
