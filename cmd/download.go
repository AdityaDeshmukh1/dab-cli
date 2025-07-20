package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
	
	"github.com/adityadeshmukh1/dab-cli/internal/models"
)

func DownloadCommand() *cli.Command {
	return &cli.Command{
		Name:  "download",
		Usage: "Download a track by ID to disk",
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "id", Aliases: []string{"i"}, Required: true},
			&cli.StringFlag{Name: "out", Aliases: []string{"o"}, Usage: "Output filename (optional)"},
		},
		Action: func(c *cli.Context) error {
			trackID := c.Int("id")
			filename := c.String("out")

			token, err := os.ReadFile(".session")
			if err != nil {
				return fmt.Errorf("failed to read session: %v", err)
			}

			// Get streaming URL
			streamReq := fmt.Sprintf("https://dab.yeet.su/api/stream?trackId=%d", trackID)
			req, err := http.NewRequest("GET", streamReq, nil)
			if err != nil {
				return fmt.Errorf("failed to create request: %v", err)
			}
			req.AddCookie(&http.Cookie{Name: "session", Value: string(token)})

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("stream request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("failed to get stream: status code %d", resp.StatusCode)
			}

			var streamData models.StreamResponse
			if err := json.NewDecoder(resp.Body).Decode(&streamData); err != nil {
				return fmt.Errorf("failed to decode stream URL: %v", err)
			}

			if streamData.URL == "" {
				return fmt.Errorf("empty stream URL")
			}

			// Download the audio
			audioResp, err := http.Get(streamData.URL)
			if err != nil {
				return fmt.Errorf("failed to download audio: %v", err)
			}
			defer audioResp.Body.Close()

			if audioResp.StatusCode != http.StatusOK {
				return fmt.Errorf("bad download status: %d", audioResp.StatusCode)
			}

			// Set output file name
			if filename == "" {
				filename = fmt.Sprintf("track_%d.mp3", trackID)
			}
			outPath := filepath.Join(".", filename)

			outFile, err := os.Create(outPath)
			if err != nil {
				return fmt.Errorf("failed to create file: %v", err)
			}
			defer outFile.Close()

			// Write to disk
			_, err = io.Copy(outFile, audioResp.Body)
			if err != nil {
				return fmt.Errorf("failed to save audio: %v", err)
			}

			fmt.Printf("Track downloaded: %s\n", outPath)
			return nil
		},
	}
}

