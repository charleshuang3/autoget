package nyaa

import (
	"testing"

	"github.com/charleshuang3/autoget/backend/indexers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategories(t *testing.T) {
	n := NewClient(&Config{UseProxy: true}, "", nil)
	got, err := n.Categories()
	require.Nil(t, err)
	assert.NotEmpty(t, got)
	assert.Equal(t, "Anime - English", got[3].Name)
}

func TestDetail(t *testing.T) {
	n := NewClient(&Config{UseProxy: true}, "", nil)
	got, err := n.Detail("1980585", true)
	require.Nil(t, err)

	assert.Equal(t, "1980585", got.ID)
	assert.Equal(t, "[HnY] Bakugan Battle Brawlers 13 SUB - Storm of Passion (854x480 RAW DVD-Rip)(PokePoring Edition).mkv", got.Title)
	assert.Equal(t, "Anime - English", got.Category)
	assert.Equal(t, int64(1749421806), got.CreatedDate)
	assert.Equal(t, uint64(391747993), got.Size)
	assert.NotEmpty(t, got.Description)
	assert.Len(t, got.Files, 1)

	assert.Equal(t, "[HnY] Bakugan Battle Brawlers 13 SUB - Storm of Passion (854x480 RAW DVD-Rip)(PokePoring Edition).mkv", got.Files[0].Name)
	assert.Equal(t, uint64(391747993), got.Files[0].Size)
}

func TestDetailWithComplexFileLists(t *testing.T) {
	n := NewClient(&Config{UseProxy: true}, "", nil)
	got, err := n.Detail("1980395", true)
	require.Nil(t, err)

	assert.Len(t, got.Files, 4)
	assert.Equal(
		t,
		"[tribute]_mcdull_movie_2009_[bd_1920x1080_h265_10bit]/scans/[tribute]_mcdull_movie_2009_cover_[1200dpi_lossless][cf095f5c].jxl",
		got.Files[0].Name,
	)
}

func TestHumanSizeToBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected uint64
	}{
		{
			name:     "Bytes",
			input:    "100 B",
			expected: 100,
		},
		{
			name:     "Kilobytes",
			input:    "1.5 KiB",
			expected: 1536, // 1.5 * 1024
		},
		{
			name:     "Megabytes",
			input:    "2.25 MiB",
			expected: 2359296, // 2.25 * 1024 * 1024
		},
		{
			name:     "Gigabytes",
			input:    "3 GiB",
			expected: 3221225472, // 3 * 1024 * 1024 * 1024
		},
		{
			name:     "Terabytes",
			input:    "0.1 TiB",
			expected: 109951162777, // 0.1 * 1024 * 1024 * 1024 * 1024
		},
		{
			name:     "Zero Bytes",
			input:    "0 B",
			expected: 0,
		},
		{
			name:     "Empty String",
			input:    "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := humanSizeToBytes(tt.input)
			require.Nil(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestList(t *testing.T) {
	n := NewClient(&Config{UseProxy: true}, "", nil)

	tests := []struct {
		name     string
		req      *indexers.ListRequest
		wantPage uint32
	}{
		{
			name:     "Default",
			req:      &indexers.ListRequest{},
			wantPage: 1,
		},
		{
			name:     "Category",
			req:      &indexers.ListRequest{Category: "1_1"},
			wantPage: 1,
		},
		{
			name:     "Keyword",
			req:      &indexers.ListRequest{Keyword: "bakugan"},
			wantPage: 1,
		},
		{
			name:     "Page",
			req:      &indexers.ListRequest{Page: 2},
			wantPage: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := n.List(tt.req)
			require.Nil(t, err)
			assert.Equal(t, tt.wantPage, got.Pagination.Page)
			assert.NotEmpty(t, got.Resources)

			// test on the first item
			firstItem := got.Resources[0]
			assert.NotEmpty(t, firstItem.ID)
			assert.NotEmpty(t, firstItem.Title)
			assert.NotEmpty(t, firstItem.Category)
			assert.NotEmpty(t, firstItem.CreatedDate)
			assert.NotEmpty(t, firstItem.Size)
		})
	}
}

func TestDownload(t *testing.T) {
	dir := t.TempDir()
	n := NewClient(&Config{UseProxy: true}, dir, nil)
	got, err := n.Download("1980585")
	require.Nil(t, err)
	assert.NotEmpty(t, got.TorrentFilePath)
	assert.FileExists(t, got.TorrentFilePath)
}
