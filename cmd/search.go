package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

type SearchResponse struct {
	Tracks []Track `json:"tracks"`
}

type Track struct {
	ID 			int `json:"id"`
	Title 	string `json:"title"`
	Artist 	string `json:"artist"`
}

// Load that session cookie!

func getSessionToken() (string, error) {
	data, err := os.ReadFile(".session")
	if err != nil {
		return "", fmt.Errorf("could not read session file: %v", err)
	}
	return string(data), nil
}

func SearchCommand() *cli.Command {
	return &cli.Command {
		Name: 	"search", 
		Usage: "Search for tracks, albums, or artists",
		Flags: 	[]cli.Flag {
			&cli.StringFlag{Name: "query", Aliases: []string{"q"}, Required:true},
		},
		Action: func(c *cli.Context) error {
			query := c.String("query")
			token, err := getSessionToken()
			if err != nil {
				return err
			}

			url := "https://dab.yeet.su/api/search?q=" + query
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return fmt.Errorf("failed to create search request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(&http.Cookie{Name: "session", Value: token})

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("search request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("search failed: %s", string(body))
			}

			var searchRes SearchResponse
			if err := json.NewDecoder(resp.Body).Decode(&searchRes); err != nil {
				return fmt.Errorf("failed to parse response: %v", err)
			}

			if len(searchRes.Tracks) == 0 {
				fmt.Println("No tracks found.")
				return nil
			}

			fmt.Println("Search Results:")
			for i, t := range searchRes.Tracks {
				fmt.Printf("%2d. %s - %s (ID: %d)\n", i+1, t.Title, t.Artist, t.ID)
			}

			return nil
		},
	}
}
