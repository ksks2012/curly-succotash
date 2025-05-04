package model

import (
	"context"
	"fmt"
	"time"

	"curly-succotash/backend/global"
	"curly-succotash/backend/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

// Model defines the common fields for all models
type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
}

// Game represents a board game entry
type Game struct {
	Model
	Theme     string    `gorm:"type:text;not null" json:"theme"`
	CardCount int       `gorm:"column:card_count;not null" json:"card_count"`
	Style     string    `gorm:"type:text;not null" json:"style"`
	CreatedAt time.Time `gorm:"type:datetime;not null" json:"created_at"`
}

// Card represents a card entry
type Card struct {
	Model
	GameID      int    `gorm:"not null;index" json:"game_id"`
	Name        string `gorm:"type:text;not null" json:"name"`
	Description string `gorm:"type:text;not null" json:"description"`
	Effect      string `gorm:"type:text;not null" json:"effect"`
}

// Meta represents a key-value pair in the meta table
type Meta struct {
	Key   string `gorm:"primary_key" json:"key"`
	Value int64  `gorm:"not null" json:"value"`
}

// NewDBEngine initializes the database engine
func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// Open database connection based on DBType
	switch databaseSetting.DBType {
	case "sqlite3":
		db, err = gorm.Open("sqlite3", databaseSetting.Path+"?_foreign_keys=on")
	case "mysql", "mariadb":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			databaseSetting.UserName,
			databaseSetting.Password,
			databaseSetting.Host,
			databaseSetting.DBName,
		)
		db, err = gorm.Open("mysql", dsn)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", databaseSetting.DBType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %s", err)
	}

	// Set GORM configurations
	if global.AppSetting.RunMode == "debug" {
		db.LogMode(true)
	}
	db.SingularTable(true)
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)

	// Apply migrations
	if err := applyMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to apply migrations: %s", err)
	}

	return db, nil
}

// applyMigrations runs GORM AutoMigrate for models
func applyMigrations(db *gorm.DB) error {
	ctx := context.Background()
	global.Logger.Infof(ctx, "Applying database migrations")

	// AutoMigrate tables
	if err := db.AutoMigrate(&Game{}, &Card{}, &Meta{}).Error; err != nil {
		return fmt.Errorf("failed to auto-migrate tables: %s", err)
	}
	global.Logger.Infof(ctx, "Successfully migrated tables: game, card, meta")

	// Verify table existence
	type TableCount struct {
		Count int
	}
	var result TableCount
	err := db.Raw("SELECT count(*) AS count FROM sqlite_master WHERE type='table' AND name IN ('game', 'card', 'meta')").Scan(&result).Error
	if err != nil {
		return fmt.Errorf("failed to verify tables: %s", err)
	}
	if result.Count != 3 {
		return fmt.Errorf("expected 3 tables, found %d", result.Count)
	}
	global.Logger.Infof(ctx, "Verified tables exist: game, card, meta")

	return nil
}

// updateTimeStampForCreateCallback sets CreatedOn and ModifiedOn on create
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				_ = createTimeField.Set(nowTime)
			}
		}
		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				_ = modifyTimeField.Set(nowTime)
			}
		}
	}
}

// updateTimeStampForUpdateCallback sets ModifiedOn on update
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		_ = scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

// deleteCallback implements soft delete
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")
		isDelField, hasIsDelField := scope.FieldByName("IsDel")
		if !scope.Search.Unscoped && hasDeletedOnField && hasIsDelField {
			now := time.Now().Unix()
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v,%v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(now),
				scope.Quote(isDelField.DBName),
				scope.AddToVars(1),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

// addExtraSpaceIfExist adds a space if the string is non-empty
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
