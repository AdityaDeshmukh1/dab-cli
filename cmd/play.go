package cmd 

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"

	"github.com/adityadeshmukh1/dab-cli/internal/models"
)


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


func PlayCommand() *cli.Command {
	return &cli.Command {
		Name: "play",
		Usage: "Play a track by ID",
		Flags: []cli.Flag {
			&cli.IntFlag {Name: "id", Aliases:[]string{"i"}, Required: true},
			&cli.StringFlag{
				Name:    "quality",
				Aliases: []string{"q"},
				Usage:   "Audio quality: low, medium, high",
				Value:   "high", // default
			},
		},
		Action: func  (c *cli.Context) error {
			trackID := c.Int("id")
			codec, format, bitrate := mapQualityToFFmpegFlags(c.String("quality"))

			token, err := os.ReadFile(".session")
			if err != nil {
				return fmt.Errorf("failed to read session: %v", err)
			}

			streamURL := fmt.Sprintf("https://dab.yeet.su/api/stream?trackId=%d", trackID)
			req, err := http.NewRequest("GET", streamURL, nil)
			if err != nil {
				return fmt.Errorf("failed to create request: %v", err)
			}

			req.AddCookie(&http.Cookie {Name: "session", Value: string(token)})

			// Check if the stream is accessible
			resp, err := http.DefaultClient.Do(req)
			if err != nil || resp.StatusCode != http.StatusOK {
				return fmt.Errorf("stream error: %v", err)
			}
			defer resp.Body.Close()

			// Parse the reponse JSON body to get the song URL
			var streamData models.StreamResponse
			if err := json.NewDecoder(resp.Body).Decode(&streamData); err != nil {
				return fmt.Errorf("failed to parse stream JSON: %v", err)
			}

			if streamData.URL == "" {
				return fmt.Errorf("stream URL is empty")
			}

			args := []string{
				"-i", streamData.URL,
			}

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
		},
	}
}
