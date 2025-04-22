package mteam

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/charleshuang3/autoget/backend/indexers"
	"github.com/charleshuang3/autoget/backend/internal/errors"
	"github.com/rs/zerolog/log"
)

type genDownloadLinkResponse struct {
	Code    interface{} `json:"code"` // maybe string or int
	Message string      `json:"message"`
	Data    string      `json:"data"`
}

func (m *MTeam) Download(id, dir string) (*indexers.DownloadResult, *errors.HTTPStatusError) {
	_, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusBadRequest, "invalid id")
	}

	resp := &genDownloadLinkResponse{}
	er := makeMultipartAPICall(m.config.GetBaseURL(), "/api/torrent/genDlToken", m.config.APIKey, map[string]string{
		"id": id,
	}, resp)
	if er != nil {
		return nil, er
	}

	if resp.Code != "0" {
		log.Error().Any("code", resp.Code).Str("message", resp.Message).Str("API", "/api/torrent/genDlToken").Msg("API error")
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, resp.Message)
	}

	destFilePath := filepath.Join(dir, name+"."+id+".torrent")

	err = downloadFileFromURL(resp.Data, destFilePath)
	if err != nil {
		return nil, errors.NewHTTPStatusError(http.StatusInternalServerError, err.Error())
	}

	return &indexers.DownloadResult{
		TorrentFilePath: destFilePath,
	}, nil
}

// downloadFileFromURL downloads a file from a given URL and saves it to a specified local path.
func downloadFileFromURL(url string, dest string) error {
	// Get the data
	resp, err := http.Get(url)
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
