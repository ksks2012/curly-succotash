package v1

import (
	"context"
	"curly-succotash/backend/global"
	"curly-succotash/backend/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListGames handles the GET request to retrieve a list of games.
//
// @Summary      List games
// @Description  Retrieves all games that are not marked as deleted.
// @Tags         games
// @Produce      json
// @Success      200  {array}   model.Game
// @Failure      500  {object}  gin.H{"error": string}
// @Router       /api/v1/games [get]
func ListGames(c *gin.Context) {
	var games []model.Game
	ctx := context.Background()
	if err := global.DBEngine.WithContext(ctx).Where("is_del = 0").Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, games)
}
