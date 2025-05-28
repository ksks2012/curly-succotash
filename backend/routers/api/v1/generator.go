package v1

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"curly-succotash/backend/global"
	"curly-succotash/backend/internal/model"
	"curly-succotash/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// Generator handles game generation requests
type Generator struct{}

func NewGenerator() Generator {
	return Generator{}
}

// Generate processes a game generation request. It performs the following steps:
// 1. Parses and validates the input JSON to create a Game object.
// 2. Sets default values for the Game object and saves it to the database.
// 3. Generates cards based on the game data and saves them to the database.
// 4. Generates a PDF file containing the card details.
// 5. Returns the generated cards and the URL to the PDF file as a JSON response.
//
// If any step fails, an appropriate error message is returned in the response.
//
// Parameters:
// - c: The Gin context, which provides request and response handling.
//
// Example Response:
// HTTP 200 OK
//
//	{
//	  "cards": [...],
//	  "pdfUrl": "http://localhost:8080/files/generated.pdf"
//	}
//
// Possible Errors:
// - HTTP 400 Bad Request: Invalid input JSON.
// - HTTP 500 Internal Server Error: Database save failure, card generation failure, or PDF generation failure.
func (g *Generator) Generate(c *gin.Context) {
	var input model.Game
	if err := c.ShouldBindJSON(&input); err != nil {
		global.Logger.Errorf(c, "Invalid input: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Set default values for Model fields
	input.CreatedBy = "system"
	input.ModifiedBy = "system"
	input.CreatedAt = time.Now()

	// Save game data using GORM
	if err := global.DBEngine.Create(&input).Error; err != nil {
		global.Logger.Errorf(c, "Failed to save game: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save game: " + err.Error()})
		return
	}
	gameID := input.ID

	// Generate cards
	cards, err := service.GenerateCards(c, input)
	if err != nil {
		global.Logger.Errorf(c, "Failed to generate cards: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate cards: " + err.Error()})
		return
	}

	// Save cards using GORM
	for _, card := range cards {
		dbCard := model.Card{
			GameID:      gameID,
			Name:        card.Name,
			Description: card.Description,
			Effect:      card.Effect,
			Model: model.Model{
				CreatedBy:  "system",
				ModifiedBy: "system",
				CreatedOn:  uint32(time.Now().Unix()),
				ModifiedOn: uint32(time.Now().Unix()),
				DeletedOn:  0,
				IsDel:      0},
		}
		if err := global.DBEngine.Create(&dbCard).Error; err != nil {
			global.Logger.Errorf(c, "Failed to save card: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save card: " + err.Error()})
			return
		}
	}

	// Generate PDF
	pdfPath, err := service.GeneratePDF(c, cards)
	if err != nil {
		global.Logger.Errorf(c, "Failed to generate PDF: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF: " + err.Error()})
		return
	}

	// Return result
	c.JSON(http.StatusOK, gin.H{
		"cards":  cards,
		"pdfUrl": fmt.Sprintf("http://localhost:%s/%s/%s", global.ServerSetting.HttpPort, global.StoragePathSetting.PDFFoldar, filepath.Base(pdfPath)),
	})
}
