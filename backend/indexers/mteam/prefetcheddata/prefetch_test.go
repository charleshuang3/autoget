package prefetcheddata

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

// If this test failed, that means we need to re-run cmd/update.
func TestFetchCheckIfUpdated(t *testing.T) {
	if apiKey == "" {
		t.Skip("MTEAM_API_KEY not set")
	}

	want := &prefetched{}
	err := json.Unmarshal(dataJSON, want)
	require.NoError(t, err)

	got, err := FetchAll(apiKey, true)
	require.NoError(t, err)

	if d := cmp.Diff(want, got, cmpopts.IgnoreFields(categoryWithOrder{}, "Order", "NumericID")); d != "" {
		t.Errorf("FetchCategories() mismatch (-want +got):\n%s", d)
	}
}
