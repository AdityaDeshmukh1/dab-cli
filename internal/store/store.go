package store

import (
	"encoding/json"
	"fmt"
	"os"
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

