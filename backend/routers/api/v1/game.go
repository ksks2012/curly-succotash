package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"curly-succotash/backend/global"
	"curly-succotash/backend/internal/ai"
	"curly-succotash/backend/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GenerateGameRequest defines the request payload for generating a game
type GenerateGameRequest struct {
	Theme       string `json:"theme" binding:"required"`
	CardCount   int    `json:"cardCount" binding:"required,min=10,max=100"`
	Style       string `json:"style" binding:"required"`
	Description string `json:"description"`
}

// GameResponse defines the response for game queries
type GameResponse struct {
	ID          uint32       `json:"id"`
	Theme       string       `json:"theme"`
	CardCount   int          `json:"card_count"`
	Style       string       `json:"style"`
	Description string       `json:"description"`
	Cards       []model.Card `json:"cards"`
}

type cardResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Effect      string `json:"effect"`
}

// GenerateGame generates a new board game using Gemini AI
func GenerateGame(c *gin.Context) {
	var req GenerateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	tx := global.DBEngine.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Initialize Gemini client
	aiClient, err := ai.NewGeminiClient()
	if err != nil {
		global.Logger.Errorf(ctx, "failed to initialize AI client: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to initialize AI client: %s", err)})
		return
	}
	defer aiClient.Close()

	// Generate game description (story background)
	prompt := fmt.Sprintf(global.StoryPromptTemplate, req.Theme)
	story, err := aiClient.GenerateContent(prompt)
	if err != nil {
		global.Logger.Errorf(ctx, "failed to generate story: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to generate story: %s", err)})
		return
	}
	if req.Description != "" {
		story = req.Description // Override with user input if provided
	}
	fmt.Println("Generated story: %s", story)
	global.Logger.Infof(ctx, "Generated story: %s", story)

	// Create game entry
	game := model.Game{
		Theme:       req.Theme,
		CardCount:   req.CardCount,
		Style:       req.Style,
		Description: story,
		CreatedAt:   time.Now(),
		Model: model.Model{
			CreatedBy:  "system",
			ModifiedBy: "system",
			CreatedOn:  uint32(time.Now().Unix()),
			ModifiedOn: uint32(time.Now().Unix()),
		},
	}
	if err := tx.Create(&game).Error; err != nil {
		global.Logger.Errorf(ctx, "failed to create game: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create game: %s", err)})
		return
	}

	// Generate cards (roles, events, items)
	cards, err := generateCards(c, tx, aiClient, game.ID, req.CardCount, story)
	if err != nil {
		global.Logger.Errorf(ctx, "failed to generate cards: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to generate cards: %s", err)})
		return
	}
	for _, card := range cards {
		if err := tx.Create(&card).Error; err != nil {
			global.Logger.Errorf(ctx, "failed to create card: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create card: %s", err)})
			return
		}
	}

	// Save plot points and objective in meta table
	metas := []model.Meta{
		{Key: fmt.Sprintf("game_%d_plot_points", game.ID), Value: 0},
		{Key: fmt.Sprintf("game_%d_main_objective_completed", game.ID), Value: 0},
	}
	for _, meta := range metas {
		if err := tx.Create(&meta).Error; err != nil {
			global.Logger.Errorf(ctx, "failed to create meta: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create meta: %s", err)})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"game_id": game.ID,
		"message": "Game generated successfully",
	})
}

// generateCards creates AI-generated role, event, and item cards
func generateCards(c *gin.Context, tx *gorm.DB, aiClient *ai.GeminiClient, gameID uint32, cardCount int, story string) ([]model.Card, error) {
	ctx := c.Request.Context()

	global.Logger.Info(ctx, "Generating cards")
	cards := []model.Card{}

	// Role cards
	rolePrompt := fmt.Sprintf(global.RolePrompt, 4, story)

	roleText, err := aiClient.GenerateContent(rolePrompt)
	if err != nil {
		global.Logger.Errorf(ctx, "Role generation error: %v", err)
		return nil, fmt.Errorf("failed to generate role: %s", err)
	}

	var role []cardResponse
	if err := json.Unmarshal([]byte(roleText), &role); err != nil {
		global.Logger.Errorf(ctx, "JSON parse error: %v", err)
		return nil, fmt.Errorf("failed to parse role JSON: %s", err)
	}

	for _, r := range role {
		cards = append(cards, model.Card{
			GameID:      gameID,
			Type:        "role",
			Name:        r.Name,
			Description: r.Description,
			Effect:      r.Effect,
		})
	}

	// Event and Item cards
	remaining := cardCount - len(cards)
	eventPrompt := fmt.Sprintf(global.EventPrompt, remaining, story)
	eventText, err := aiClient.GenerateContent(eventPrompt)
	if err != nil {
		global.Logger.Errorf(ctx, "Event generation error: %v", err)
		return nil, fmt.Errorf("failed to generate event: %s", err)
	}

	var event []cardResponse
	if err := json.Unmarshal([]byte(eventText), &event); err != nil {
		global.Logger.Errorf(ctx, "JSON parse error: %v", err)
		return nil, fmt.Errorf("failed to parse event JSON: %s", err)
	}

	for _, e := range event {
		cards = append(cards, model.Card{
			GameID:      gameID,
			Type:        "event",
			Name:        e.Name,
			Description: e.Description,
			Effect:      e.Effect,
		})
	}

	return cards, nil
}

// GetGame retrieves a stored game by ID
func GetGame(c *gin.Context) {
	id := c.Param("id")
	var game model.Game
	ctx := context.Background()

	if err := global.DBEngine.WithContext(ctx).Where("id = ? AND is_del = 0", id).First(&game).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("game not found: %s", err)})
		return
	}

	var cards []model.Card
	if err := global.DBEngine.WithContext(ctx).Where("game_id = ? AND is_del = 0", id).Find(&cards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch cards: %s", err)})
		return
	}

	c.JSON(http.StatusOK, GameResponse{
		ID:          game.ID,
		Theme:       game.Theme,
		CardCount:   game.CardCount,
		Style:       game.Style,
		Description: game.Description,
		Cards:       cards,
	})
}
