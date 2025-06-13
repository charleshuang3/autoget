package indexers

import (
	"strings"
	"testing"
)

func TestRenderRSSResult(t *testing.T) {
	tests := []struct {
		name                   string
		indexer                string
		downloadStarted        []string
		downloadPendingToStart []string
		expectedSubstrings     []string
		notExpectedSubstrings  []string
	}{
		{
			name:    "Both lists populated",
			indexer: "TestIndexer",
			downloadStarted: []string{
				"Item A",
				"Item B",
			},
			downloadPendingToStart: []string{
				"Item C",
				"Item D",
			},
			expectedSubstrings: []string{
				"# TestIndexer RSS",
				"## Download Started",
				"- Item A",
				"- Item B",
				"## Download Pending to Start",
				"- Item C",
				"- Item D",
			},
			notExpectedSubstrings: []string{},
		},
		{
			name:            "DownloadStarted empty",
			indexer:         "TestIndexer",
			downloadStarted: []string{},
			downloadPendingToStart: []string{
				"Item C",
				"Item D",
			},
			expectedSubstrings: []string{
				"# TestIndexer RSS",
				"## Download Pending to Start",
				"- Item C",
				"- Item D",
			},
			notExpectedSubstrings: []string{
				"## Download Started",
			},
		},
		{
			name:    "DownloadPendingToStart empty",
			indexer: "TestIndexer",
			downloadStarted: []string{
				"Item A",
				"Item B",
			},
			downloadPendingToStart: []string{},
			expectedSubstrings: []string{
				"# TestIndexer RSS",
				"## Download Started",
				"- Item A",
				"- Item B",
			},
			notExpectedSubstrings: []string{
				"## Download Pending to Start",
			},
		},
		{
			name:                   "Both lists empty",
			indexer:                "TestIndexer",
			downloadStarted:        []string{},
			downloadPendingToStart: []string{},
			expectedSubstrings: []string{
				"# TestIndexer RSS",
			},
			notExpectedSubstrings: []string{
				"## Download Started",
				"## Download Pending to Start",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RenderRSSResult(tt.indexer, tt.downloadStarted, tt.downloadPendingToStart)
			if err != nil {
				t.Fatalf("RenderRSSResult returned an error: %v", err)
			}

			for _, expected := range tt.expectedSubstrings {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected to find substring:\n%q\nBut got:\n%q", expected, result)
				}
			}

			for _, notExpected := range tt.notExpectedSubstrings {
				if strings.Contains(result, notExpected) {
					t.Errorf("Expected NOT to find substring:\n%q\nBut found:\n%q", notExpected, result)
				}
			}
		})
	}
}
