package prefetcheddata

import (
	"maps"
	"slices"
	"strings"

	"github.com/charleshuang3/autoget/backend/indexers"
)

var (
	Categories = map[string]indexers.Category{
		"0_0": {ID: "0_0", Name: "All categories", SubCategories: []indexers.Category{
			{ID: "1_0", Name: "Anime", SubCategories: []indexers.Category{
				{ID: "1_1", Name: "Anime - AMV"},
				{ID: "1_2", Name: "Anime - English"},
				{ID: "1_3", Name: "Anime - Non-English"},
				{ID: "1_4", Name: "Anime - Raw"},
			}},
			{ID: "2_0", Name: "Audio", SubCategories: []indexers.Category{
				{ID: "2_1", Name: "Audio - Lossless"},
				{ID: "2_2", Name: "Audio - Lossy"},
			}},
			{ID: "3_0", Name: "Literature", SubCategories: []indexers.Category{
				{ID: "3_1", Name: "Literature - English"},
				{ID: "3_2", Name: "Literature - Non-English"},
				{ID: "3_3", Name: "Literature - Raw"},
			}},
			{ID: "4_0", Name: "Live Action", SubCategories: []indexers.Category{
				{ID: "4_1", Name: "Live Action - English"},
				{ID: "4_2", Name: "Live Action - Idol/PV"},
				{ID: "4_3", Name: "Live Action - Non-English"},
				{ID: "4_4", Name: "Live Action - Raw"},
			}},
			{ID: "5_0", Name: "Pictures", SubCategories: []indexers.Category{
				{ID: "5_1", Name: "Pictures - Graphics"},
				{ID: "5_2", Name: "Pictures - Photos"},
			}},
			{ID: "6_0", Name: "Software", SubCategories: []indexers.Category{
				{ID: "6_1", Name: "Software - Apps"},
				{ID: "6_2", Name: "Software - Games"},
			}},
		}},
	}

	CategoriesList = slices.SortedFunc(maps.Values(Categories), func(a, b indexers.Category) int {
		return strings.Compare(a.ID, b.ID)
	})
)
