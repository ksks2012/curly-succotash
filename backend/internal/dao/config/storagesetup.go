package config

import (
	"errors"

	interfaces "curly-succotash/backend/interfaces"
	sqlitestorage "curly-succotash/backend/internal/dao/dbversion/sqlite"
	"curly-succotash/backend/pkg/setting"
)

// StorageSetup contains storage type and storage instance
type StorageSetup struct {
	Type     string
	Instance interfaces.StorageEngine
}

func (s *StorageSetup) NewDBEngine(databaseSetting *setting.DatabaseSettingS) (err error) {
	switch databaseSetting.DBType {
	case "pxc":
		// TODO:
		// s.Instance, err = setupMySQLRoundRobinStorageEngine(databaseSetting)
	case "mysql", "mariadb":
		// TODO:
		// s.Instance, err = setupMySQLStorageEngine(databaseSetting)
	case "sqlite3":
		s.Instance, err = sqlitestorage.NewSQLiteStorageEngine(databaseSetting)
	default:
		err = errors.New("unknown storage engine type: " + databaseSetting.DBType)
	}
	if err != nil {
		s.Instance = nil
		return err
	}
	s.Type = databaseSetting.DBType
	return nil
}
