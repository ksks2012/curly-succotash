package model

import (
	"context"
	"fmt"
	"time"

	"curly-succotash/backend/global"
	"curly-succotash/backend/migrations"
	"curly-succotash/backend/pkg/setting"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Model defines the common fields for all models
type Model struct {
	ID         uint32 `gorm:"primaryKey" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `gorm:"index" json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
}

func (Model) TableName() string {
	return "model"
}

// Game represents a board game entry
type Game struct {
	Model
	Theme       string    `gorm:"type:text;not null" json:"theme"`
	CardCount   int       `gorm:"column:card_count;not null" json:"card_count"`
	Style       string    `gorm:"type:text;not null" json:"style"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"type:datetime;not null" json:"created_at"`
}

func (Game) TableName() string {
	return "games"
}

// Card represents a card entry
type Card struct {
	Model
	GameID      int    `gorm:"not null;index" json:"game_id"`
	Type        string `gorm:"type:text;not null" json:"type"` // Added: role, event, item
	Name        string `gorm:"type:text;not null" json:"name"`
	Description string `gorm:"type:text;not null" json:"description"`
	Effect      string `gorm:"type:text;not null" json:"effect"`
}

// TableName specifies the table name for Card
func (Card) TableName() string {
	return "cards"
}

// Meta represents a key-value pair in the meta table
type Meta struct {
	Key   string `gorm:"primaryKey" json:"key"`
	Value int64  `gorm:"not null" json:"value"`
}

// TableName specifies the table name for Meta
func (Meta) TableName() string {
	return "meta"
}

// NewDBEngine initializes the database engine
func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// Open database connection based on DBType
	switch databaseSetting.DBType {
	case "sqlite3":
		db, err = gorm.Open(sqlite.Open(databaseSetting.Path+"?_foreign_keys=on"), &gorm.Config{})
	case "mysql", "mariadb":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			databaseSetting.UserName,
			databaseSetting.Password,
			databaseSetting.Host,
			databaseSetting.DBName,
		)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database type: %s", databaseSetting.DBType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %s", err)
	}

	// Set GORM configurations
	if global.AppSetting.RunMode == "debug" {
		db = db.Debug()
	}

	// Set connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %s", err)
	}
	sqlDB.SetMaxIdleConns(databaseSetting.MaxIdleConns)
	sqlDB.SetMaxOpenConns(databaseSetting.MaxOpenConns)

	// Register callbacks
	db.Callback().Create().Before("gorm:create").Register("app:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Before("gorm:update").Register("app:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Before("gorm:delete").Register("app:soft_delete", softDeleteCallback)

	// Apply migrations
	if err := applyMigrations(db); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to apply migrations: %s", err)
	}

	return db, nil
}

// applyMigrations runs gormigrate migrations
func applyMigrations(db *gorm.DB) error {
	ctx := context.Background()
	global.Logger.Infof(ctx, "Applying database migrations")

	// Run migrations from migrations package
	migrator := gormigrate.New(db, gormigrate.DefaultOptions, migrations.GetMigrations())
	if err := migrator.Migrate(); err != nil {
		return fmt.Errorf("failed to run migrations: %s", err)
	}
	global.Logger.Infof(ctx, "Successfully applied migrations")

	// Verify table existence
	var tableCount int64
	query := "SELECT count(*) FROM sqlite_master WHERE type='table' AND name IN ('games', 'cards', 'meta')"
	if db.Dialector.Name() == "mysql" {
		query = "SELECT count(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name IN ('games', 'cards', 'meta')"
	}
	if err := db.Raw(query).Scan(&tableCount).Error; err != nil {
		return fmt.Errorf("failed to verify tables: %s", err)
	}
	if tableCount != 3 {
		return fmt.Errorf("expected 3 tables, found %d", tableCount)
	}
	global.Logger.Infof(ctx, "Verified tables exist: games, cards, meta")

	return nil
}

// updateTimeStampForCreateCallback sets CreatedOn and ModifiedOn on create
func updateTimeStampForCreateCallback(db *gorm.DB) {
	if db.Error != nil {
		return
	}
	nowTime := uint32(time.Now().Unix())
	if _, ok := db.Statement.Schema.FieldsByName["CreatedOn"]; ok {
		if db.Statement.ReflectValue.FieldByName("CreatedOn").IsZero() {
			db.Statement.SetColumn("CreatedOn", nowTime)
		}
	}
	if _, ok := db.Statement.Schema.FieldsByName["ModifiedOn"]; ok {
		if db.Statement.ReflectValue.FieldByName("ModifiedOn").IsZero() {
			db.Statement.SetColumn("ModifiedOn", nowTime)
		}
	}
}

// updateTimeStampForUpdateCallback sets ModifiedOn on update
func updateTimeStampForUpdateCallback(db *gorm.DB) {
	if db.Error != nil {
		return
	}
	if _, ok := db.Statement.Context.Value("gorm:update_column").(bool); !ok {
		db.Statement.SetColumn("ModifiedOn", uint32(time.Now().Unix()))
	}
}

// softDeleteCallback implements soft delete
func softDeleteCallback(db *gorm.DB) {
	if db.Error != nil {
		return
	}
	if db.Statement.Schema == nil {
		return
	}
	if !db.Statement.Unscoped {
		if _, ok := db.Statement.Schema.FieldsByName["DeletedOn"]; ok {
			if _, ok := db.Statement.Schema.FieldsByName["IsDel"]; ok {
				now := uint32(time.Now().Unix())
				db.Statement.AddClause(clause.Set{{Column: clause.Column{Name: "deleted_on"}, Value: now}})
				db.Statement.AddClause(clause.Set{{Column: clause.Column{Name: "is_del"}, Value: 1}})
			}
		}
	}
}
