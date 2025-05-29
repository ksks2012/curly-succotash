package routers

import (
	"net/http"
	"time"

	_ "curly-succotash/backend/docs"
	"curly-succotash/backend/global"
	"curly-succotash/backend/pkg/limiter"
	v1 "curly-succotash/backend/routers/api/v1"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var methodLimiters = limiter.NewMethodLimiter().AddBuckets(
	limiter.LimiterBucketRule{
		Key:          "/auth",
		FillInterval: time.Second,
		Capacity:     10,
		Quantum:      10,
	},
)

func NewRouter() *gin.Engine {
	r := gin.New()
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
	if global.AppSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	}
	if global.AppSetting.RunMode != "release" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	// TODO: middlewares

	r.Static("/files", "./files")

	generator := v1.NewGenerator()

	apiv1 := r.Group("/api/v1")
	// TODO: JWT
	apiv1.Use()
	{
		// Generate game
		apiv1.POST("/generate", generator.Generate)
		apiv1.POST("/game", v1.GenerateGame)

		apiv1.GET("/games", v1.ListGames)
		apiv1.GET("/games/:id", v1.GetGame)
		// TODO:
		apiv1.GET("/generate-pdf/:id", v1.GenerateHTMLPDF)
	}

	return r
}
