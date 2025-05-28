package v1

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"curly-succotash/backend/global"
	"curly-succotash/backend/internal/model"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gin-gonic/gin"
)

// GenerateHTMLPDF generates a PDF using HTML and wkhtmltopdf
// GenerateHTMLPDF handles the HTTP request to generate a PDF file from the HTML representation
// of a board game and its associated cards. It fetches the game and its cards from the database,
// renders an HTML template using the game and card data, converts the HTML to a PDF using wkhtmltopdf,
// saves the PDF to the server, and serves the generated PDF file as a response.
//
// API
// @Summary      Generate PDF for a board game
// @Description  Generates a PDF file containing the board game's details and its cards, and returns the PDF file.
// @Tags         games
// @Produce      application/pdf
// @Param        id   path      string  true  "Game ID"
// @Success      200  {file}    file    "PDF file"
// @Failure      404  {object}  map[string]string  "game not found"
// @Failure      500  {object}  map[string]string  "internal server error"
// @Router       /api/v1/games/{id}/pdf [get]
func GenerateHTMLPDF(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	// Fetch game
	var game model.Game
	if err := global.DBEngine.WithContext(ctx).Where("id = ? AND is_del = 0", id).First(&game).Error; err != nil {
		global.Logger.Errorf(c, "game not found: %s", err)
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("game not found: %s", err)})
		return
	}

	// Fetch cards
	var cards []model.Card
	if err := global.DBEngine.WithContext(ctx).Where("game_id = ? AND is_del = 0", id).Find(&cards).Error; err != nil {
		global.Logger.Errorf(c, "failed to fetch cards: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch cards: %s", err)})
		return
	}

	// HTML template
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
		<style>
			@page { margin: 10mm; }
			body { font-family: Arial, sans-serif; }
			.card { 
				width: 250px; 
				height: 350px; 
				border: 2px solid black; 
				border-radius: 8px; 
				padding: 10px; 
				margin: 10px; 
				float: left; 
				background-color: #f9fafb;
				box-shadow: 2px 2px 5px rgba(0,0,0,0.2);
			}
			.card-title { font-size: 18px; font-weight: bold; color: #1f2937; }
			.card-type { font-size: 14px; color: #6b7280; }
			.card-desc, .card-effect { 
				font-size: 12px; 
				margin-top: 8px; 
				color: #374151; 
				line-height: 1.4; 
				max-height: 120px; 
				overflow: hidden; 
			}
			.page-break { clear: both; page-break-after: always; }
			.header { margin-bottom: 20px; }
		</style>
	</head>
	<body>
		<div class="header">
			<h1 class="text-2xl font-bold mb-4">Board Game: {{.Game.Theme}}</h1>
			<p class="text-base mb-4">Story: {{.Game.Description}}</p>
			<h2 class="text-xl font-bold mb-4">Cards (Game ID: {{.Game.ID}})</h2>
		</div>
		{{range $i, $card := .Cards}}
			<div class="card">
				<div class="card-title">{{.Name}}</div>
				<div class="card-type">({{.Type | title}})</div>
				<div class="card-desc">Description: {{.Description}}</div>
				<div class="card-effect">Effect: {{.Effect}}</div>
			</div>
			{{if eq (mod $i 4) 3}}<div class="page-break"></div>{{end}}
		{{end}}
	</body>
	</html>`

	// Parse template
	t, err := template.New("pdf").Funcs(template.FuncMap{
		"mod": func(i, n int) int { return i % n },
		"title": func(s string) string {
			if len(s) == 0 {
				return s
			}
			return string(s[0]-32) + s[1:]
		},
	}).Parse(tmpl)
	if err != nil {
		global.Logger.Errorf(c, "failed to parse template: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to parse template: %s", err)})
		return
	}

	// Render HTML
	data := struct {
		Game  model.Game
		Cards []model.Card
	}{game, cards}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		global.Logger.Errorf(c, "failed to render template: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to render template: %s", err)})
		return
	}

	// Create PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		global.Logger.Errorf(c, "failed to create PDF generator: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create PDF generator: %s", err)})
		return
	}
	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(buf.Bytes())))
	pdfg.MarginTop.Set(10)
	pdfg.MarginBottom.Set(10)
	pdfg.MarginLeft.Set(10)
	pdfg.MarginRight.Set(10)

	// Generate PDF
	if err := pdfg.Create(); err != nil {
		global.Logger.Errorf(c, "failed to generate PDF: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to generate PDF: %s", err)})
		return
	}

	// Save PDF
	outputDir := "./files"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		global.Logger.Errorf(c, "failed to create directory: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create directory: %s", err)})
		return
	}
	pdfPath := filepath.Join(outputDir, fmt.Sprintf("game_%s.pdf", id))
	if err := pdfg.WriteFile(pdfPath); err != nil {
		global.Logger.Errorf(c, "failed to save PDF: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to save PDF: %s", err)})
		return
	}

	// Serve PDF
	c.File(pdfPath)
}
