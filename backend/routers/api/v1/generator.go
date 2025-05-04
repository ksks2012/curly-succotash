package v1

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"curly-succotash/backend/global"
	"curly-succotash/backend/internal/model"
	"curly-succotash/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type Generator struct{}

func NewGenerator() Generator {
	return Generator{}
}

func (g *Generator) Generate(c *gin.Context) {
	var input model.Game
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Invalid input: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// TODO: Save game data
	// result, err := db.Exec("INSERT INTO games (theme, card_count, style) VALUES (?, ?, ?)",
	// 	input.Theme, input.CardCount, input.Style)
	// if err != nil {
	// 	log.Printf("Failed to save game: %s", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save game"})
	// 	return
	// }
	// gameID, _ := result.LastInsertId()

	// Generate cards
	cards, err := service.GenerateCards(c, input)
	if err != nil {
		global.Logger.Errorf(c, "Failed to generate cards: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate cards: " + err.Error()})
		return
	}

	// TODO: Save cards
	// for _, card := range cards {
	// 	_, err = db.Exec("INSERT INTO cards (game_id, name, description, effect) VALUES (?, ?, ?, ?)",
	// 		gameID, card.Name, card.Description, card.Effect)
	// 	if err != nil {
	// 		log.Printf("Failed to save card: %s", err)
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save card"})
	// 		return
	// 	}
	// }

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
		"pdfUrl": fmt.Sprintf("http://localhost:8080/files/%s", filepath.Base(pdfPath)),
	})
}
