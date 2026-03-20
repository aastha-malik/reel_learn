package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"reel_learn/backend/handlers"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"POST", "GET", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	r.POST("/process", handlers.ProcessHandler)
	r.GET("/video", handlers.VideoHandler)

	r.Run(":8080")
}
