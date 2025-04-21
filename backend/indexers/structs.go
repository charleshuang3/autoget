package indexers

import (
	"time"
)

type IIndexer interface {
	// Name of the indexer.
	Name() string

	// Categories returns indexer's resource categories.
	Categories() ([]Category, error)

	// List resources in given categories and keyword.
	List(categories []string, keyword string, page, pageSize uint32) ([]Resource, error)

	// Detail of a resource.
	Detail(id string) (Resource, error)

	SetRequestTimeout(timeout time.Duration)
}

type IndexerBasicInfo struct {
	Name_   string
	Private bool
}

func NewIndexerBasicInfo(name string, private bool) *IndexerBasicInfo {
	return &IndexerBasicInfo{
		Name_:   name,
		Private: private,
	}
}

func (info *IndexerBasicInfo) Name() string {
	return info.Name_
}

type Category struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	SubCategories []Category `json:"subCategories,omitempty"`
}

type Resource struct {
}
