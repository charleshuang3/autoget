package category

import (
	"encoding/json"
	"os"
	"testing"

	_ "embed"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed test_res/list_categories.json
	testToCategoriesInput []byte

	//go:embed categories.json
	currentCategories []byte

	apiKey = os.Getenv("MTEAM_API_KEY")
)

func TestToCategoryFile(t *testing.T) {
	categories := &listCategories{}
	err := json.Unmarshal(testToCategoriesInput, categories)
	require.NoError(t, err)

	got := categories.toCategoryFile(false)
	want := &categoryFile{
		CategoryTree: []*categoryWithOrder{
			{
				ID:   "adult",
				Name: "adult",
				SubCategories: []*categoryWithOrder{
					{
						ID:        "115",
						Name:      "AV(有码)",
						Order:     20,
						NumericID: 115,
						SubCategories: []*categoryWithOrder{
							{
								ID:        "410",
								Name:      "AV(有码)/HD Censored",
								Order:     31,
								NumericID: 410,
							},
							{
								ID:        "440",
								Name:      "AV(Gay)/HD",
								Order:     440,
								NumericID: 440,
							},
						},
					},
				},
			},
			{
				ID:   "normal",
				Name: "normal",
				SubCategories: []*categoryWithOrder{
					{
						ID:        "110",
						Name:      "Music",
						Order:     4,
						NumericID: 110,
						SubCategories: []*categoryWithOrder{
							{
								ID:        "434",
								Name:      "Music(无损)",
								Order:     1,
								NumericID: 434,
							},
						},
					},
				},
			},
		},
		CategoryInfos: map[string]*CategoryInfo{
			"110":    {"Music", "normal", []string{"434"}},
			"115":    {"AV(有码)", "adult", []string{"410", "440"}},
			"410":    {"AV(有码)/HD Censored", "adult", []string{"410"}},
			"434":    {"Music(无损)", "normal", []string{"434"}},
			"440":    {"AV(Gay)/HD", "adult", []string{"440"}},
			"adult":  {"adult", "adult", []string{}},
			"normal": {"normal", "normal", []string{}},
		},
	}
	assert.Equal(t, want, got)
}

// If this test failed, that means we need to re-run update_categories.
func TestFetchCategoriesCheckIfUpdated(t *testing.T) {
	if apiKey == "" {
		t.Skip("MTEAM_API_KEY not set")
	}

	want := &categoryFile{}
	err := json.Unmarshal(currentCategories, want)
	require.NoError(t, err)

	got, err := FetchCategories(apiKey, true)
	require.NoError(t, err)

	if d := cmp.Diff(want, got, cmpopts.IgnoreFields(categoryWithOrder{}, "Order", "NumericID")); d != "" {
		t.Errorf("FetchCategories() mismatch (-want +got):\n%s", d)
	}
}
