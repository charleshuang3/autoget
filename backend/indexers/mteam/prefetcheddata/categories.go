package prefetcheddata

import (
	"sort"
	"strconv"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

const (
	categoryAdult   = "adult"
	categoryNormal  = "normal"
	categoryGayPorn = "440"

	baseURL = "https://api.m-team.cc"
)

var (
	rootCategories = map[string]string{
		"100": categoryNormal, // Movie
		"105": categoryNormal, // TV Series
		"444": categoryNormal, // Documentary
		"110": categoryNormal, // Music
		"443": categoryNormal, // edu
		"447": categoryNormal, // Game
		"449": categoryNormal, // Anime
		"450": categoryNormal, // Others
		"115": categoryAdult,  // AV Censored
		"120": categoryAdult,  // AV Uncensored
		"445": categoryAdult,  // IV
		"446": categoryAdult,  // HCG
	}
)

type listCategories struct {
	Data struct {
		List []struct {
			CreatedDate      string `json:"createdDate"`
			LastModifiedDate string `json:"lastModifiedDate"`
			ID               string `json:"id"`
			Order            string `json:"order"`
			NameChs          string `json:"nameChs"`
			NameCht          string `json:"nameCht"`
			NameEng          string `json:"nameEng"`
			Image            string `json:"image"`
			Parent           string `json:"parent"`
		} `json:"list"`

		// We don't use following fields because they don't contains
		// all subcategories. For example the parent of tvshow(105).
		Adult  []string `json:"adult"`
		Movie  []string `json:"movie"`
		Music  []string `json:"music"`
		Tvshow []string `json:"tvshow"`

		// We don't use following fields
		Waterfall []string `json:"waterfall"`
	} `json:"data"`

	// We don't use following fields
	Code    string `json:"code"`
	Message string `json:"message"`
}

// categoryWithOrder has same json definition with indexers.Category.
type categoryWithOrder struct {
	ID            string               `json:"id"`
	Name          string               `json:"name"`
	SubCategories []*categoryWithOrder `json:"subCategories,omitempty"`
	Order         int
	NumericID     int
}

type CategoryInfo struct {
	Name       string   `json:"name"`
	Mode       string   `json:"mode"`
	Categories []string `json:"categories"` // You can not search resources on "115" but need to includes all sub.
}

type categoryFile struct {
	CategoryTree  []*categoryWithOrder     `json:"category_tree"`
	CategoryInfos map[string]*CategoryInfo `json:"categories"`
}

func (l *listCategories) toCategoryFile(excludeGayContent bool) *categoryFile {
	adultRoot := &categoryWithOrder{
		ID:   categoryAdult,
		Name: categoryAdult,
	}
	normalRoot := &categoryWithOrder{
		ID:   categoryNormal,
		Name: categoryNormal,
	}
	roots := []*categoryWithOrder{
		adultRoot,
		normalRoot,
	}

	categories := map[string]*categoryWithOrder{
		categoryAdult:  adultRoot,
		categoryNormal: normalRoot,
	}

	for _, cat := range l.Data.List {
		if excludeGayContent && cat.ID == categoryGayPorn {
			continue
		}
		id, err := strconv.Atoi(cat.ID)
		if err != nil {
			log.Fatal().Msgf("Category ID is not a number: %s", cat.ID)
		}
		order, err := strconv.Atoi(cat.Order)
		if err != nil {
			log.Fatal().Msgf("Category Order is not a number: id = %s, order = %s", cat.ID, cat.Order)
		}

		categories[cat.ID] = &categoryWithOrder{
			ID:        cat.ID,
			Name:      cat.NameChs,
			Order:     order,
			NumericID: id,
		}
	}

	for _, cat := range l.Data.List {
		if excludeGayContent && cat.ID == categoryGayPorn {
			continue
		}
		parent := cat.Parent
		if parent == "" {
			var ok bool
			parent, ok = rootCategories[cat.ID]
			if !ok {
				log.Fatal().Msgf("Got unknown root category: %s %s", cat.ID, cat.NameChs)
			}
		}

		p, ok := categories[parent]
		if !ok {
			log.Fatal().Msgf("Category %s has unknown parent %s", cat.ID, parent)
		}

		p.SubCategories = append(p.SubCategories, categories[cat.ID])
	}

	sortSubCategories(adultRoot)
	sortSubCategories(normalRoot)

	categoryInfos := map[string]*CategoryInfo{}
	categoryInfo(adultRoot, categoryInfos, categoryAdult)
	categoryInfo(normalRoot, categoryInfos, categoryNormal)

	return &categoryFile{
		CategoryTree:  roots,
		CategoryInfos: categoryInfos,
	}
}

func sortSubCategories(category *categoryWithOrder) {
	sort.SliceStable(category.SubCategories, func(i, j int) bool {
		if category.SubCategories[i].Order != category.SubCategories[j].Order {
			return category.SubCategories[i].Order < category.SubCategories[j].Order
		}
		return category.SubCategories[i].NumericID < category.SubCategories[j].NumericID
	})

	for _, sub := range category.SubCategories {
		sortSubCategories(sub)
	}
}

func categoryInfo(categories *categoryWithOrder, m map[string]*CategoryInfo, mode string) {
	subs := []string{}
	if categories.Name != categoryAdult && categories.Name != categoryNormal {
		for _, sub := range categories.SubCategories {
			subs = append(subs, sub.ID)
		}
		if len(subs) == 0 {
			subs = append(subs, categories.ID)
		}
	}

	m[categories.ID] = &CategoryInfo{
		Name:       categories.Name,
		Mode:       mode,
		Categories: subs,
	}

	for _, sub := range categories.SubCategories {
		categoryInfo(sub, m, mode)
	}
}

func FetchCategories(apiKey string, excludeGayContent bool) (*categoryFile, error) {
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodPost, baseURL+"/api/torrent/categoryList", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var categories listCategories
	if err := json.NewDecoder(resp.Body).Decode(&categories); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w, body: %s", err, resp.Body)
	}

	return categories.toCategoryFile(excludeGayContent), nil
}
