package sukebei

import (
	"testing"

	"github.com/charleshuang3/autoget/backend/indexers"
	"github.com/charleshuang3/autoget/backend/indexers/nyaa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategories(t *testing.T) {
	n := NewClient(&nyaa.Config{UseProxy: true}, "", nil)
	got, err := n.Categories()
	require.Nil(t, err)
	assert.NotEmpty(t, got)
	assert.Equal(t, "Real Life - Videos", got[9].Name)
}

func TestList(t *testing.T) {
	n := NewClient(&nyaa.Config{UseProxy: true}, "", nil)

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
			req:      &indexers.ListRequest{Keyword: "中文字幕"},
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
	n := NewClient(&nyaa.Config{UseProxy: true}, dir, nil)
	got, err := n.Download("4322631")
	require.Nil(t, err)
	assert.NotEmpty(t, got.TorrentFilePath)
	assert.FileExists(t, got.TorrentFilePath)
}

func TestDetail(t *testing.T) {
	n := NewClient(&nyaa.Config{UseProxy: true}, "", nil)
	got, err := n.Detail("4322631", true)
	require.Nil(t, err)

	assert.Equal(t, "4322631", got.ID)
	assert.NotEmpty(t, got.Title)
	assert.NotEmpty(t, got.Category)
	assert.NotEmpty(t, got.CreatedDate)
	assert.NotEmpty(t, got.Size)
	assert.NotEmpty(t, got.Description)
	assert.NotEmpty(t, got.Files)
}
