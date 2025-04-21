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
	m := NewMTeam(&Config{
		APIKey: apiKey,
	})
	require.NotNil(t, m)

	got, err := m.Categories()
	require.NoError(t, err)
	assert.Len(t, got, 2)
}
