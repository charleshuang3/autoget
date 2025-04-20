package indexers

import (
	"regexp"
	"time"
)

type IIndexer interface {
	// Name of the indexer.
	Name() string

	// Categories returns indexer's resource categories.
	Categories() ([]*Category, error)

	// List resources in given categories and keyword.
	List(categories []string, keyword string, page, pageSize uint32) ([]Resource, error)

	// Detail of a resource.
	Detail(id string) (Resource, error)

	SetRequestTimeout(timeout time.Duration)
}

type Category struct {
	// Id in format `indexers/{indexer_name}/categories/([0-9a-zA-Z-_]*)`
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	SubCategories []*Category `json:"subCategories,omitempty"`
}

const (
	CategoryIdFormat = "indexers/%s/categories/%s"
)

var (
	categoryIDRegex = regexp.MustCompile(`indexers/([0-9a-zA-Z-_]*)/categories/([0-9a-zA-Z-_]*)`)
)

func VerifyCategoryId(id string) bool {
	matches := categoryIDRegex.FindStringSubmatch(id)
	return len(matches) == 3
}

func ExtractIndexerAndCategoryId(id string) (string, string) {
	matches := categoryIDRegex.FindStringSubmatch(id)
	if len(matches) != 3 {
		return "", ""
	}
	return matches[1], matches[2]
}

type Resource struct {
}
