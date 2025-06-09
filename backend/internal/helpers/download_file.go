package helpers

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownloadFileFromURL downloads a file from a given URL and saves it to a specified local path.
func DownloadFileFromURL(httpClient *http.Client, url string, dest string) error {
	// Get the data
	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP GET error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status error: %d %s", resp.StatusCode, resp.Status)
	}

	// Create the destination file
	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Write the response body to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy response body to file: %w", err)
	}

	return nil
}
