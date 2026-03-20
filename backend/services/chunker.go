package services

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Reel represents a single 15-second clip with its global timestamp within the original video.
type Reel struct {
	Path      string  `json:"path"`
	StartTime float64 `json:"start_time"`
	EndTime   float64 `json:"end_time"`
}

// ProcessVideo splits a downloaded video into intermediate chunks and then into
// 15-second reels. Returns all reels with their global start/end times.
func ProcessVideo(jobDir, videoPath string) ([]Reel, error) {
	duration, err := getDuration(videoPath)
	if err != nil {
		return nil, fmt.Errorf("get duration: %w", err)
	}

	cs := chunkSize(duration)

	chunksDir := filepath.Join(jobDir, "chunks")
	chunks, err := splitSegments(videoPath, chunksDir, "chunk", cs)
	if err != nil {
		return nil, fmt.Errorf("split chunks: %w", err)
	}

	reelsDir := filepath.Join(jobDir, "reels")
	var reels []Reel
	globalStart := 0.0

	for i, chunk := range chunks {
		chunkDur, err := getDuration(chunk)
		if err != nil {
			return nil, fmt.Errorf("get chunk duration: %w", err)
		}

		subDir := filepath.Join(reelsDir, fmt.Sprintf("chunk_%04d", i))
		subReels, err := splitSegments(chunk, subDir, "reel", 15)
		if err != nil {
			return nil, fmt.Errorf("split reels for chunk %d: %w", i, err)
		}

		reelStart := globalStart
		for _, reelPath := range subReels {
			endTime := math.Min(reelStart+15, globalStart+chunkDur)
			reels = append(reels, Reel{
				Path:      reelPath,
				StartTime: math.Round(reelStart*1000) / 1000,
				EndTime:   math.Round(endTime*1000) / 1000,
			})
			reelStart = endTime
		}

		globalStart += chunkDur
	}

	return reels, nil
}

// getDuration returns the duration of a media file in seconds using ffprobe.
func getDuration(filePath string) (float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath,
	)
	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe: %w", err)
	}
	return strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
}

// chunkSize returns the intermediate chunk duration based on total video length.
//
//	>= 3600s (1h+)   → 3600s hourly chunks
//	>= 60s            → 60s minute chunks
//	< 60s             → 15s (video is already short, reels are the chunks)
func chunkSize(duration float64) float64 {
	if duration >= 3600 {
		return 3600
	}
	if duration >= 60 {
		return 60
	}
	return 15
}

// splitSegments uses ffmpeg's segment muxer to cut a file into equal-length pieces.
// Returns sorted file paths of all generated segments.
func splitSegments(inputPath, outputDir, prefix string, segmentDuration float64) ([]string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("mkdir %s: %w", outputDir, err)
	}

	outputPattern := filepath.Join(outputDir, prefix+"_%04d.mp4")

	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-c", "copy",
		"-f", "segment",
		"-segment_time", fmt.Sprintf("%.0f", segmentDuration),
		"-reset_timestamps", "1",
		"-avoid_negative_ts", "make_zero",
		outputPattern,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg: %w\noutput: %s", err, string(out))
	}

	matches, err := filepath.Glob(filepath.Join(outputDir, prefix+"_*.mp4"))
	if err != nil {
		return nil, err
	}
	sort.Strings(matches)
	return matches, nil
}
