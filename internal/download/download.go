package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

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

func Download(trackNumber int) bool {
	// Load the last search results
	if err := store.LoadFromFile(".dabcli_last_search.json"); err != nil {
		fmt.Printf("Warning: could not load last search results: %v\n", err)
	}

	trackID, ok := store.GetSongID(trackNumber)
	if !ok {
		fmt.Printf("track number %d not found in last search", trackNumber)
		return false
	}

	// Fetch stream URL
	url, err := store.FetchStreamURL(trackID)
	if err != nil {
		fmt.Printf("Error in fetching the stream URL: %v", err)
		return false
	}

	// Download audio
	outFile := "song.mp3"
	savedFile, err := downloadToFile(url, outFile)
	if err != nil {
		fmt.Print(err)
		return false
	}

	fmt.Printf("Track downloaded: %s\n", savedFile)
	return true
}
