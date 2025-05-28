package v1

import (
	"curly-succotash/backend/global"
	"curly-succotash/backend/internal/model"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func GeneratePDF(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	var game model.Game
	if err := global.DBEngine.WithContext(ctx).Where("id = ? AND is_del = 0", id).First(&game).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("game not found: %s", err)})
		return
	}
	var cards []model.Card
	if err := global.DBEngine.WithContext(ctx).Where("game_id = ? AND is_del = 0", id).Find(&cards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch cards: %s", err)})
		return
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	for _, card := range cards {
		pdf.Cell(0, 10, fmt.Sprintf("%s (%s)", card.Name, card.Type))
		pdf.Ln(8)
		pdf.MultiCell(0, 8, "Description: "+card.Description, "", "", false)
		pdf.MultiCell(0, 8, "Effect: "+card.Effect, "", "", false)
		pdf.Ln(8)
	}
	outputDir := "./files"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create directory: %s", err)})
		return
	}
	pdfPath := filepath.Join(outputDir, fmt.Sprintf("game_%s.pdf", id))
	if err := pdf.OutputFileAndClose(pdfPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to generate PDF: %s", err)})
		return
	}
	c.File(pdfPath)
}
