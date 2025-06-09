package nyaa

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategories(t *testing.T) {
	n := NewClient(&Config{})
	got, err := n.Categories()
	require.Nil(t, err)
	assert.NotEmpty(t, got)
}

func TestDetail(t *testing.T) {
	n := NewClient(&Config{})
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
