package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Download fetches a YouTube video using yt-dlp and returns the job directory
// and path to the downloaded mp4 file.
func Download(youtubeURL string) (jobDir string, videoPath string, err error) {
	jobDir = filepath.Join("/tmp/reellearn", newID())

	if err = os.MkdirAll(jobDir, 0755); err != nil {
		return "", "", fmt.Errorf("create job dir: %w", err)
	}

	outputTemplate := filepath.Join(jobDir, "video.%(ext)s")

	cmd := exec.Command("yt-dlp",
		"--no-playlist",
		"--merge-output-format", "mp4",
		"-o", outputTemplate,
		youtubeURL,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("yt-dlp failed: %w\n%s", err, string(out))
	}

	matches, err := filepath.Glob(filepath.Join(jobDir, "video.*"))
	if err != nil || len(matches) == 0 {
		return "", "", fmt.Errorf("downloaded file not found in %s", jobDir)
	}

	return jobDir, matches[0], nil
}

func newID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
