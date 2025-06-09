package prefetcheddata

import "github.com/charleshuang3/autoget/backend/indexers"

var (
	Categories = map[string]indexers.Category{
		"0_0": {ID: "0_0", Name: "All categories"},
		"1_0": {ID: "1_0", Name: "Anime"},
		"1_1": {ID: "1_1", Name: "Anime - AMV"},
		"1_2": {ID: "1_2", Name: "Anime - English"},
		"1_3": {ID: "1_3", Name: "Anime - Non-English"},
		"1_4": {ID: "1_4", Name: "Anime - Raw"},
		"2_0": {ID: "2_0", Name: "Audio"},
		"2_1": {ID: "2_1", Name: "Audio - Lossless"},
		"2_2": {ID: "2_2", Name: "Audio - Lossy"},
		"3_0": {ID: "3_0", Name: "Literature"},
		"3_1": {ID: "3_1", Name: "Literature - English"},
		"3_2": {ID: "3_2", Name: "Literature - Non-English"},
		"3_3": {ID: "3_3", Name: "Literature - Raw"},
		"4_0": {ID: "4_0", Name: "Live Action"},
		"4_1": {ID: "4_1", Name: "Live Action - English"},
		"4_2": {ID: "4_2", Name: "Live Action - Idol/PV"},
		"4_3": {ID: "4_3", Name: "Live Action - Non-English"},
		"4_4": {ID: "4_4", Name: "Live Action - Raw"},
		"5_0": {ID: "5_0", Name: "Pictures"},
		"5_1": {ID: "5_1", Name: "Pictures - Graphics"},
		"5_2": {ID: "5_2", Name: "Pictures - Photos"},
		"6_0": {ID: "6_0", Name: "Software"},
		"6_1": {ID: "6_1", Name: "Software - Apps"},
		"6_2": {ID: "6_2", Name: "Software - Games"},
	}
)
