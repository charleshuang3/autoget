package mteam

import (
	"encoding/json"
	"os"
	"testing"

	_ "embed"

	"github.com/charleshuang3/autoget/backend/indexers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ indexers.IIndexer = (*MTeam)(nil)

	//go:embed res/list_categories.json
	testToCategoriesInput []byte

	apiKey = os.Getenv("MTEAM_API_KEY")
)

func TestToCategories(t *testing.T) {
	categories := &ListCategories{}
	err := json.Unmarshal(testToCategoriesInput, categories)
	require.NoError(t, err)

	got := categories.toCategories(false)
	want := []*indexers.Category{
		{
			ID:   "adult-root",
			Name: "adult",
			SubCategories: []*indexers.Category{
				{
					ID:   "adult-115",
					Name: "AV(有码)",
					SubCategories: []*indexers.Category{
						{
							ID:   "adult-410",
							Name: "AV(有码)/HD Censored",
						},
						{
							ID:   "adult-440",
							Name: "AV(Gay)/HD",
						},
					},
				},
			},
		},
		{
			ID:   "normal-root",
			Name: "normal",
			SubCategories: []*indexers.Category{
				{
					ID:   "normal-110",
					Name: "Music",
					SubCategories: []*indexers.Category{
						{
							ID:   "normal-434",
							Name: "Music(无损)",
						},
					},
				},
			},
		},
	}
	assert.Equal(t, want, got)
}

func TestCategories(t *testing.T) {
	if apiKey == "" {
		t.Skip("MTEAM_API_KEY not set")
	}

	m := NewMTeam(&Config{
		APIKey: apiKey,
	})
	require.NotNil(t, m)

	got, err := m.Categories()
	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func init() {
	fatalOnUnknownRootCategory = true
}
