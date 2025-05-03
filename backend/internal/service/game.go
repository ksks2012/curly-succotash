package service

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"curly-succotash/backend/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

// generateCards simulates AI card generation (to be replaced with Hugging Face model)
func GenerateCards(c *gin.Context, input model.Game) ([]model.Card, error) {
	cardTemplates := []model.Card{
		{Name: "Magic Sword", Description: "A shining sword", Effect: "Deals 2 damage"},
		{Name: "Healing Potion", Description: "Restores vitality", Effect: "Heals 3 health"},
		{Name: "Fireball", Description: "Burns enemies", Effect: "Deals 1 damage to all enemies"},
	}

	cards := make([]model.Card, input.CardCount)
	for i := 0; i < input.CardCount; i++ {
		template := cardTemplates[rand.Intn(len(cardTemplates))]
		cards[i] = model.Card{
			GameID:      i + 1,
			Name:        fmt.Sprintf("%s %s", input.Theme, template.Name),
			Description: template.Description,
			Effect:      template.Effect,
		}
	}
	return cards, nil
}

// generatePDF creates a PDF with card details
func GeneratePDF(c *gin.Context, cards []model.Card) (string, error) {
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
