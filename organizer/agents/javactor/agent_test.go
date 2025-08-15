package javactor

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getAPIKey(t *testing.T) string {
	t.Helper()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping test because GEMINI_API_KEY is not set")
	}
	return apiKey
}

func TestAgentRun(t *testing.T) {
	apiKey := getAPIKey(t)
	names, err := Run(apiKey, "森泽佳奈")
	require.NoError(t, err)
	assert.Contains(t, names, "森泽佳奈")
	assert.Contains(t, names, "饭冈佳奈子")
	assert.Contains(t, names, "藤原辽子")
}
