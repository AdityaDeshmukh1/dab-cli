package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
	"github.com/adityadeshmukh1/dab-cli/internal/models"
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

// Fetch stream URL from API
func fetchStreamURL(trackID int) (string, error) {
	token, err := os.ReadFile(".session")
	if err != nil {
		return "", fmt.Errorf("failed to read session: %v", err)
	}

	streamURL := fmt.Sprintf("https://dab.yeet.su/api/stream?trackId=%d", trackID)
	req, err := http.NewRequest("GET", streamURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.AddCookie(&http.Cookie{Name: "session", Value: string(token)})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("stream request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("stream request failed: %s", string(body))
	}

	var streamData models.StreamResponse
	if err := json.NewDecoder(resp.Body).Decode(&streamData); err != nil {
		return "", fmt.Errorf("failed to parse stream JSON: %v", err)
	}

	if streamData.URL == "" {
		return "", fmt.Errorf("stream URL is empty")
	}

	return streamData.URL, nil
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

// CLI Command
func PlayCommand() *cli.Command {
	return &cli.Command{
		Name:  "play",
		Usage: "Play a track from last search",
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "track", Aliases: []string{"t"}, Required: true},
			&cli.StringFlag{Name: "quality", Aliases: []string{"q"}, Usage: "Audio quality: low, medium, high, flac", Value: "high"},
		},

		Action: func(c *cli.Context) error {
			// Load the last search map
			if err := store.LoadFromFile(".dabcli_last_search.json"); err != nil {
				fmt.Printf("Warning: could not load last search results: %v\n", err)
			}
			trackNumber := c.Int("track")
			trackID, ok := store.GetSongID(trackNumber)
			if !ok {
				return fmt.Errorf("track number %d not found in last search", trackNumber)
			}

			codec, format, bitrate := mapQualityToFFmpegFlags(c.String("quality"))

			url, err := fetchStreamURL(trackID)
			if err != nil {
				return err
			}

			return playStream(url, codec, format, bitrate)
		},
	}
}

