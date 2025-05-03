package routers

import (
	"net/http"
	"time"

	"curly-succotash/backend/global"

	v1 "curly-succotash/backend/routers/api/v1"

	"curly-succotash/backend/pkg/limiter"

	"github.com/gin-gonic/gin"
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
	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
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
	}

	return r
}
