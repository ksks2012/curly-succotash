package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Model20250520AddGameInfo defines the common fields (unchanged)
type Model20250520AddGameInfo struct {
	ID         uint32 `gorm:"primaryKey" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `gorm:"index" json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
}

// Game20250520AddGameInfo adds the Description field
type Game20250520AddGameInfo struct {
	Model20250520AddGameInfo
	Theme       string    `gorm:"type:text;not null" json:"theme"`
	CardCount   int       `gorm:"column:card_count;not null" json:"card_count"`
	Style       string    `gorm:"type:text;not null" json:"style"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"type:datetime;not null" json:"created_at"`
}

type Card20250520AddGameInfo struct {
	Model20250520AddGameInfo
	GameID      int    `gorm:"not null;index" json:"game_id"`
	Type        string `gorm:"type:text;not null" json:"type"` // Added: role, event, item
	Name        string `gorm:"type:text;not null" json:"name"`
	Description string `gorm:"type:text;not null" json:"description"`
	Effect      string `gorm:"type:text;not null" json:"effect"`
}

// TableName specifies the table name for Game20250520AddGameInfo
func (Game20250520AddGameInfo) TableName() string {
	return "games"
}

// TableName specifies the table name for Game20250520AddGameInfo
func (Card20250520AddGameInfo) TableName() string {
	return "cards"
}

var AddGameInfo = &gormigrate.Migration{
	ID: "20250520200000_add_game_info",
	Migrate: func(tx *gorm.DB) error {
		// Add Description column
		if err := tx.Migrator().AutoMigrate(&Game20250520AddGameInfo{}, &Card20250520AddGameInfo{}); err != nil {
			return err
		}
		return nil
	},
	Rollback: func(tx *gorm.DB) error {
		// Drop Description column
		return tx.Migrator().DropColumn(&Game20250520AddGameInfo{}, "description")
	},
}
