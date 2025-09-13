package search

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/adityadeshmukh1/dab-cli/internal/store"
)

type SearchResponse struct {
	Tracks []Track `json:"tracks"`
}

type Track struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

// Load that session cookie!
func getSessionToken() (string, error) {
	data, err := os.ReadFile(".session")
	if err != nil {
		return "", fmt.Errorf("could not read session file: %v", err)
	}
	return string(data), nil
}

func Search(query string) ([]Track, error) {
	token, err := getSessionToken()
	if err != nil {
		return nil, err
	}

	encodedQuery := url.QueryEscape(query)
	url := "https://dab.yeet.su/api/search?q=" + encodedQuery
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "session", Value: token})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search request failed: %s", string(body))
	}

	var searchRes SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchRes); err != nil {
		return nil, fmt.Errorf("failed to parse responseL %v", err)
	}

	// Update store
	store.ResetSongs()
	for i, t := range searchRes.Tracks {
		store.SetSong(i+1, t.ID)
	}

	// Save last search
	if err := store.SaveToFile(".dabcli_last_search.json"); err != nil {
		fmt.Printf("Warning: could not save searc results: %v\n", err)
	}

	return searchRes.Tracks, nil
}
