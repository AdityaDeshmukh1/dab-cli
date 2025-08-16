package store

import (
	"encoding/json"
	"fmt"
	"os"
	"net/http"
	"io"
	"github.com/adityadeshmukh1/dab-cli/internal/models"
)

var SongIDMap = make(map[int]int)

// SetSong stores a track ID for a CLI number
func SetSong(index int, id int) {
	SongIDMap[index] = id
}

// GetSongID returns the real track ID for a CLI number
func GetSongID(index int) (int, bool) {
	id, ok := SongIDMap[index]
	return id, ok
}

// ResetSongs clears the map
func ResetSongs() {
	SongIDMap = make(map[int]int)
}

// SaveToFile saves the map to a JSON file
func SaveToFile(filename string) error {
	data, err := json.Marshal(SongIDMap)
	if err != nil {
		return fmt.Errorf("failed to marshal SongIDMap: %v", err)
	}
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	return nil
}

// LoadFromFile loads the map from a JSON file
func LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	if err := json.Unmarshal(data, &SongIDMap); err != nil {
		return fmt.Errorf("failed to unmarshal SongIDMap: %v", err)
	}
	return nil
}

// Fetch stream URL from API (shared with play.go logic)
func FetchStreamURL(trackID int) (string, error) {
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


