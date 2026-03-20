package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const reelBaseDir = "/tmp/reellearn"

// VideoHandler serves a reel .mp4 file by its absolute path.
// The path must be inside /tmp/reellearn/ to prevent directory traversal.
//
// GET /video?path=/tmp/reellearn/{jobID}/reels/chunk_0000/reel_0000.mp4
func VideoHandler(c *gin.Context) {
	rawPath := c.Query("path")
	if rawPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path query param is required"})
		return
	}

	// Resolve to absolute path and verify it stays inside the allowed base dir
	absPath, err := filepath.Abs(rawPath)
	if err != nil || !strings.HasPrefix(absPath, reelBaseDir+"/") {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid path"})
		return
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "reel not found"})
		return
	}

	c.Header("Content-Type", "video/mp4")
	c.File(absPath)
}
