package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type processRequest struct {
	URL string `json:"url" binding:"required"`
}

func ProcessHandler(c *gin.Context) {
	var req processRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "url field is required and must be non-empty",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "received",
		"url":    req.URL,
	})
}
