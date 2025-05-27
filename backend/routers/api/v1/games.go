package v1

import (
	"context"
	"curly-succotash/backend/global"
	"curly-succotash/backend/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListGames(c *gin.Context) {
	var games []model.Game
	ctx := context.Background()
	if err := global.DBEngine.WithContext(ctx).Where("is_del = 0").Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, games)
}
