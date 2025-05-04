package sqlitestorage

import (
	"context"
	"errors"
	"fmt"

	"curly-succotash/backend/internal/model"
	"curly-succotash/backend/pkg/setting"

	"gorm.io/gorm"
)

// ErrDatabasePathRequired indicates database path is not given
var ErrDatabasePathRequired = errors.New("database path is required")

// SQLiteStorageEngine implements the StorageEngine interface using GORM
type SQLiteStorageEngine struct {
	DB *gorm.DB
}

// NewSQLiteStorageEngine creates a new SQLite-based storage engine
func NewSQLiteStorageEngine(databaseSetting *setting.DatabaseSettingS) (*SQLiteStorageEngine, error) {
	if databaseSetting.Path == "" {
		return nil, ErrDatabasePathRequired
	}

	// Initialize GORM database
	db, err := model.NewDBEngine(databaseSetting)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GORM database: %s", err)
	}

	return &SQLiteStorageEngine{
		DB: db,
	}, nil
}

// Open ensures the database connection is ready
func (eng *SQLiteStorageEngine) Open() error {
	if eng.DB == nil {
		return errors.New("database not initialized")
	}
	// GORM connection is already opened in NewDBEngine
	return nil
}

// Close closes the database connection
func (eng *SQLiteStorageEngine) Close() error {
	if eng.DB == nil {
		return nil
	}
	sqlDB, err := eng.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying SQL DB: %s", err)
	}
	err = sqlDB.Close()
	eng.DB = nil
	return err
}

// FetchMetaInt64 retrieves an int64 value from the meta table
func (eng *SQLiteStorageEngine) FetchMetaInt64(ctx context.Context, metaKey string, defaultValue int64) (int64, error) {
	var meta model.Meta
	err := eng.DB.Where("key = ?", metaKey).First(&meta).Error
	if err == gorm.ErrRecordNotFound {
		return defaultValue, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to fetch meta value for key %s: %s", metaKey, err)
	}
	return meta.Value, nil
}
