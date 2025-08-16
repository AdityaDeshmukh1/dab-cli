package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"github.com/adityadeshmukh1/dab-cli/internal/store"
)

// Download stream to a file
func downloadToFile(url, filename string) (string, error) {
	// Fetch the audio stream
	audioResp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download audio: %v", err)
	}
	defer audioResp.Body.Close()

	if audioResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad download status: %d", audioResp.StatusCode)
	}

	// Ensure filename
	if filename == "" {
		filename = "track.mp3"
	}
	outPath := filepath.Join(".", filename)

	// Write stream to file
	outFile, err := os.Create(outPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, audioResp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save audio: %v", err)
	}

	return outPath, nil
}

// CLI Command
func DownloadCommand() *cli.Command {
	return &cli.Command{
		Name:  "download",
		Usage: "Download a track from last search",
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "track", Aliases: []string{"t"}, Required: true},
			&cli.StringFlag{Name: "out", Aliases: []string{"o"}, Usage: "Output filename (optional)"},
		},
		Action: func(c *cli.Context) error {
			// Load the last search results
			if err := store.LoadFromFile(".dabcli_last_search.json"); err != nil {
				fmt.Printf("Warning: could not load last search results: %v\n", err)
			}

			trackNumber := c.Int("track")
			trackID, ok := store.GetSongID(trackNumber)
			if !ok {
				return fmt.Errorf("track number %d not found in last search", trackNumber)
			}

			// Fetch stream URL
			url, err := store.FetchStreamURL(trackID)
			if err != nil {
				return err
			}

			// Download audio
			outFile := c.String("out")
			savedFile, err := downloadToFile(url, outFile)
			if err != nil {
				return err
			}

			fmt.Printf("Track downloaded: %s\n", savedFile)
			return nil
		},
	}
}

