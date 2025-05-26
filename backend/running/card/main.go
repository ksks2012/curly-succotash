package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	_ "github.com/mattn/go-sqlite3"
)

// GameInput defines the structure for user input
type GameInput struct {
	Theme     string `json:"theme" binding:"required"`
	CardCount int    `json:"cardCount" binding:"required,min=1,max=20"`
	Style     string `json:"style" binding:"required"`
}

// Card defines the structure for a generated card
type Card struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Effect      string `json:"effect"`
}

func main() {
	// Initialize database
	db, err := sql.Open("sqlite3", "./games.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS games (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			theme TEXT NOT NULL,
			card_count INTEGER NOT NULL,
			style TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS cards (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			game_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			effect TEXT NOT NULL,
			FOREIGN KEY (game_id) REFERENCES games(id)
		);
	`)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	// Initialize Gin
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// API: Generate game
	r.POST("/api/generate", func(c *gin.Context) {
		var input GameInput
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Printf("Invalid input: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		// Save game data
		result, err := db.Exec("INSERT INTO games (theme, card_count, style) VALUES (?, ?, ?)",
			input.Theme, input.CardCount, input.Style)
		if err != nil {
			log.Printf("Failed to save game: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save game"})
			return
		}
		gameID, _ := result.LastInsertId()

		// Generate cards
		cards, err := generateCards(input)
		if err != nil {
			log.Printf("Failed to generate cards: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate cards: " + err.Error()})
			return
		}

		// Save cards
		for _, card := range cards {
			_, err = db.Exec("INSERT INTO cards (game_id, name, description, effect) VALUES (?, ?, ?, ?)",
				gameID, card.Name, card.Description, card.Effect)
			if err != nil {
				log.Printf("Failed to save card: %s", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save card"})
				return
			}
		}

		// Generate PDF
		pdfPath, err := generatePDF(cards)
		if err != nil {
			log.Printf("Failed to generate PDF: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF: " + err.Error()})
			return
		}

		// Return result
		c.JSON(http.StatusOK, gin.H{
			"cards":  cards,
			"pdfUrl": fmt.Sprintf("http://localhost:8080/files/%s", filepath.Base(pdfPath)),
		})
	})

	// Serve static files for PDF download
	r.Static("/files", "./files")

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// generateCards simulates AI card generation (to be replaced with Hugging Face model)
func generateCards(input GameInput) ([]Card, error) {
	cardTemplates := []Card{
		{Name: "Magic Sword", Description: "A shining sword", Effect: "Deals 2 damage"},
		{Name: "Healing Potion", Description: "Restores vitality", Effect: "Heals 3 health"},
		{Name: "Fireball", Description: "Burns enemies", Effect: "Deals 1 damage to all enemies"},
	}

	cards := make([]Card, input.CardCount)
	for i := 0; i < input.CardCount; i++ {
		template := cardTemplates[rand.Intn(len(cardTemplates))]
		cards[i] = Card{
			ID:          i + 1,
			Name:        fmt.Sprintf("%s %s", input.Theme, template.Name),
			Description: template.Description,
			Effect:      template.Effect,
		}
	}
	return cards, nil
}

// generatePDF creates a PDF with card details
func generatePDF(cards []Card) (string, error) {
	// Ensure files directory exists
	if err := os.MkdirAll("./files", 0755); err != nil {
		return "", fmt.Errorf("failed to create files directory: %s", err)
	}

	// Generate unique PDF path
	pdfPath := fmt.Sprintf("./files/game_%d.pdf", time.Now().UnixNano())

	// Initialize PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set card dimensions (standard poker card: 63.5mm x 88.9mm)
	cardWidth, cardHeight := 63.5, 88.9
	cardsPerRow, cardsPerCol := 3, 2
	margin := 10.0

	for i, card := range cards {
		// Calculate card position
		row := (i / cardsPerRow) % cardsPerCol
		col := i % cardsPerRow
		x := margin + float64(col)*(cardWidth+5)
		y := margin + float64(row)*(cardHeight+5)

		// Draw card border
		pdf.SetFillColor(255, 255, 255) // White background
		pdf.SetDrawColor(0, 0, 0)       // Black border
		pdf.Rect(x, y, cardWidth, cardHeight, "FD")

		// Draw card content
		pdf.SetTextColor(0, 0, 0) // Black text
		pdf.SetFont("Arial", "B", 12)
		pdf.Text(x+5, y+15, card.Name)

		pdf.SetFont("Arial", "", 10)
		pdf.SetXY(x+5, y+20)
		pdf.MultiCell(cardWidth-10, 5, card.Description, "", "L", false)

		pdf.SetXY(x+5, y+40)
		pdf.MultiCell(cardWidth-10, 5, "Effect: "+card.Effect, "", "L", false)

		// Add new page if 6 cards are filled
		if (i+1)%(cardsPerRow*cardsPerCol) == 0 && i < len(cards)-1 {
			pdf.AddPage()
		}
	}

	// Save PDF
	err := pdf.OutputFileAndClose(pdfPath)
	if err != nil {
		return "", fmt.Errorf("failed to save PDF: %s", err)
	}

	// Verify PDF exists
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		return "", fmt.Errorf("PDF file not generated: %s", pdfPath)
	}

	return pdfPath, nil
}
