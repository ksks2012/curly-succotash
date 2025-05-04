package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Model defines the common fields for models in this migration
type Model struct {
	ID         uint32 `gorm:"primaryKey" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `gorm:"index" json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
}

// Game20250503 represents a board game entry for this migration
type Game20250503 struct {
	Model
	Theme     string    `gorm:"type:text;not null" json:"theme"`
	CardCount int       `gorm:"column:card_count;not null" json:"card_count"`
	Style     string    `gorm:"type:text;not null" json:"style"`
	CreatedAt time.Time `gorm:"type:datetime;not null" json:"created_at"`
}

func (Game20250503) TableName() string {
	return "games"
}

// Card20250503 represents a card entry for this migration
type Card20250503 struct {
	Model
	GameID      int    `gorm:"not null;index" json:"game_id"`
	Name        string `gorm:"type:text;not null" json:"name"`
	Description string `gorm:"type:text;not null" json:"description"`
	Effect      string `gorm:"type:text;not null" json:"effect"`
}

func (Card20250503) TableName() string {
	return "cards"
}

// Meta20250503 represents a key-value pair for this migration
type Meta20250503 struct {
	Key   string `gorm:"primaryKey" json:"key"`
	Value int64  `gorm:"not null" json:"value"`
}

func (Meta20250503) TableName() string {
	return "meta"
}

var createTables = &gormigrate.Migration{
	ID: "20250503120000_create_tables",
	Migrate: func(tx *gorm.DB) error {
		// Create games, cards, and meta tables
		if err := tx.Migrator().AutoMigrate(&Game20250503{}, &Card20250503{}, &Meta20250503{}); err != nil {
			return err
		}
		// Add foreign key for cards.game_id (MySQL-specific)
		if tx.Dialector.Name() == "mysql" {
			return tx.Exec("ALTER TABLE cards ADD CONSTRAINT fk_cards_game_id FOREIGN KEY (game_id) REFERENCES games(id)").Error
		}
		return nil
	},
	Rollback: func(tx *gorm.DB) error {
		// Drop tables in reverse order
		return tx.Migrator().DropTable("cards", "games", "meta")
	},
}
