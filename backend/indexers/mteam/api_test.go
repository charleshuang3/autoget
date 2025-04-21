package mteam

import (
	"os"
	"testing"

	"github.com/charleshuang3/autoget/backend/indexers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ indexers.IIndexer = (*MTeam)(nil)

	apiKey = os.Getenv("MTEAM_API_KEY")
)

func TestCategories(t *testing.T) {
	if apiKey == "" {
		t.Skip("MTEAM_API_KEY not set")
	}

	m := NewMTeam(&Config{
		APIKey: apiKey,
	})
	require.NotNil(t, m)

	got, err := m.Categories()
	require.Nil(t, err)
	assert.Len(t, got, 2)
}

func TestList(t *testing.T) {
	if apiKey == "" {
		t.Skip("MTEAM_API_KEY not set")
	}

	m := NewMTeam(&Config{
		APIKey: apiKey,
	})
	require.NotNil(t, m)

	tests := []struct {
		name     string
		category string
		keyword  string
		page     uint32
		pageSize uint32
	}{
		{
			name:     "normal root",
			category: "normal",
			keyword:  "",
			page:     1,
			pageSize: 2,
		},
		{
			name:     "adult root",
			category: "adult",
			keyword:  "",
			page:     1,
			pageSize: 2,
		},
		{
			name:     "normal sub",
			category: "434",
			keyword:  "",
			page:     1,
			pageSize: 2,
		},
		{
			name:     "adult sub",
			category: "433",
			keyword:  "",
			page:     1,
			pageSize: 2,
		},
		{
			name:     "adult mid level",
			category: "115",
			keyword:  "",
			page:     1,
			pageSize: 2,
		},
		{
			name:     "normal with keyword",
			category: "normal",
			keyword:  "地狱",
			page:     1,
			pageSize: 2,
		},
		{
			name:     "adult with keyword",
			category: "adult",
			keyword:  "sone",
			page:     1,
			pageSize: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := m.List(tc.category, tc.keyword, tc.page, tc.pageSize)
			require.Nil(t, err)
			assert.NotEmpty(t, got.Resources)
		})
	}
}
