package indexers

import (
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
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	SubCategories []*Category `json:"subCategories,omitempty"`
}

type Resource struct {
}
