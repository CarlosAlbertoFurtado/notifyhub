package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter(nh *NotificationHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestLogger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	api := r.Group("/api")
	{
		notifications := api.Group("/notifications")
		{
			notifications.POST("/send", nh.Send)
			notifications.GET("", nh.List)
			notifications.GET("/stats", nh.Stats)
			notifications.GET("/:id", nh.GetByID)
		}
	}

	return r
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		gin.DefaultWriter.Write([]byte(
			c.Request.Method + " " + c.Request.URL.Path + " " +
				c.Writer.Header().Get("Status") + " " +
				latency.String() + "\n",
		))
	}
}
