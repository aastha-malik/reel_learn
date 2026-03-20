package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type VideoMetadata struct {
	Title     string `json:"title"`
	Channel   string `json:"channel"`
	Duration  int    `json:"duration"`
	Thumbnail string `json:"thumbnail"`
}

func MetadataHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url query param is required"})
		return
	}

	out, err := exec.Command("yt-dlp", "--dump-json", "--no-playlist", url).Output()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("yt-dlp failed: %v", err)})
		return
	}

	var raw map[string]any
	if err := json.Unmarshal(out, &raw); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse metadata"})
		return
	}

	meta := VideoMetadata{
		Title:     stringField(raw, "title"),
		Channel:   stringField(raw, "channel"),
		Thumbnail: stringField(raw, "thumbnail"),
	}

	if d, ok := raw["duration"].(float64); ok {
		meta.Duration = int(d)
	}

	c.JSON(http.StatusOK, meta)
}

func stringField(m map[string]any, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
