package play

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/adityadeshmukh1/dab-cli/internal/store"
)

// Map quality to FFmpeg flags
func mapQualityToFFmpegFlags(q string) (codec, format, bitrate string) {
	switch q {
	case "low":
		return "libmp3lame", "mp3", "96k"
	case "medium":
		return "libmp3lame", "mp3", "160k"
	case "high":
		return "libmp3lame", "mp3", "256k"
	case "flac":
		return "flac", "flac", ""
	default:
		return "libmp3lame", "mp3", "192k"
	}
}

// Play stream via FFmpeg → MPV
func playStream(url, codec, format, bitrate string) error {
	args := []string{"-i", url}

	if codec != "" {
		args = append(args, "-c:a", codec)
	}
	if bitrate != "" {
		args = append(args, "-b:a", bitrate)
	}

	args = append(args, "-f", format, "pipe:1")

	ffmpeg := exec.Command("ffmpeg", args...)
	mpv := exec.Command("mpv", "-")

	r, w := io.Pipe()
	ffmpeg.Stdout = w
	mpv.Stdin = r

	ffmpeg.Stderr = os.Stderr
	mpv.Stdout = os.Stdout
	mpv.Stderr = os.Stderr

	if err := ffmpeg.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %v", err)
	}
	if err := mpv.Start(); err != nil {
		return fmt.Errorf("failed to start mpv: %v", err)
	}

	ffmpeg.Wait()
	w.Close()
	mpv.Wait()
	return nil
}

func Play(trackNumber int, quality string) error {
	// Load the last search map
	if err := store.LoadFromFile(".dabcli_last_search.json"); err != nil {
		return fmt.Errorf("could not load last search results: %v", err)
	}

	trackID, ok := store.GetSongID(trackNumber)
	if !ok {
		return fmt.Errorf("track number %d not found in last search", trackNumber)
	}

	codec, format, bitrate := mapQualityToFFmpegFlags(quality)

	url, err := store.FetchStreamURL(trackID)
	if err != nil {
		return fmt.Errorf("failed to fetch stream URL: %v", err)
	}

	return playStream(url, codec, format, bitrate)
}
