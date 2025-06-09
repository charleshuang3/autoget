package prefetcheddata

import (
	"maps"
	"slices"
	"strings"

	"github.com/charleshuang3/autoget/backend/indexers"
)

var (
	Categories = map[string]indexers.Category{
		"0_0": {ID: "0_0", Name: "All categories"},
		"1_0": {ID: "1_0", Name: "Art"},
		"1_1": {ID: "1_1", Name: "Art - Anime"},
		"1_2": {ID: "1_2", Name: "Art - Doujinshi"},
		"1_3": {ID: "1_3", Name: "Art - Games"},
		"1_4": {ID: "1_4", Name: "Art - Manga"},
		"1_5": {ID: "1_5", Name: "Art - Pictures"},
		"2_0": {ID: "2_0", Name: "Real Life"},
		"2_1": {ID: "2_1", Name: "Real Life - Pictures"},
		"2_2": {ID: "2_2", Name: "Real Life - Videos"},
	}

	CategoriesList = slices.SortedFunc(maps.Values(Categories), func(a, b indexers.Category) int {
		return strings.Compare(a.ID, b.ID)
	})
)
