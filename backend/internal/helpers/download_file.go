package helpers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/anacrolix/torrent/metainfo"
)

// DownloadTorrentFileFromURL downloads a file from a given URL and saves it to a specified local path.
func DownloadTorrentFileFromURL(httpClient *http.Client, url string, dest string) (*metainfo.MetaInfo, *metainfo.Info, error) {
	// Get the data
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("HTTP GET error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("HTTP status error: %d %s", resp.StatusCode, resp.Status)
	}

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	m, err := metainfo.Load(bytes.NewReader(buffer.Bytes()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load metainfo: %w", err)
	}

	info, err := m.UnmarshalInfo()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal info: %w", err)
	}

	// Create the destination file
	out, err := os.Create(dest)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Write the response body to the file
	_, err = io.Copy(out, bytes.NewReader(buffer.Bytes()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to copy response body to file: %w", err)
	}

	return m, &info, nil
}
