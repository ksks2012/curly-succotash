package migrations

import "github.com/go-gormigrate/gormigrate/v2"

func GetMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		createTables,
		AddGameInfo,
		// NOTE: Add future migrations here
	}
}
