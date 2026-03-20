package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"reel_learn/backend/services"
)

type processRequest struct {
	YoutubeURL string `json:"youtube_url" binding:"required"`
}

type processResponse struct {
	Reels []services.Reel `json:"reels"`
}

func ProcessHandler(c *gin.Context) {
	var req processRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "youtube_url field is required",
		})
		return
	}

	jobDir, videoPath, err := services.Download(req.YoutubeURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "download failed: " + err.Error(),
		})
		return
	}

	reels, err := services.ProcessVideo(jobDir, videoPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "processing failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, processResponse{Reels: reels})
}
