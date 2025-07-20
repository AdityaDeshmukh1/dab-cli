package cmd 

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
	
	"github.com/adityadeshmukh1/dab-cli/internal/models"
)

func PlayCommand() *cli.Command {
	return &cli.Command {
		Name: "play",
		Usage: "Play a track by ID",
		Flags: []cli.Flag {
			&cli.IntFlag {Name: "id", Aliases:[]string{"i"}, Required: true},
		},
		Action: func  (c *cli.Context) error {
			trackID := c.Int("id")

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
			
			// Pipe the audio directly to mpv
			cmd := exec.Command("mpv", "--no-terminal", "--quiet", streamData.URL)
			
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			return cmd.Run()

			
		},
	}
}
