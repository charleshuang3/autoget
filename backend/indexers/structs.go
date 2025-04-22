package indexers

import (
	"github.com/charleshuang3/autoget/backend/internal/errors"
)

type IIndexer interface {
	// Name of the indexer.
	Name() string

	// Categories returns indexer's resource categories.
	Categories() ([]Category, *errors.HTTPStatusError)

	// List resources in given category and keyword (optional).
	List(category, keyword string, page, pageSize uint32) (*ListResult, *errors.HTTPStatusError)

	// Detail of a resource.
	Detail(id string) (*ResourceDetail, *errors.HTTPStatusError)
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

type Pagination struct {
	Page       uint32 `json:"page"`
	TotalPages uint32 `json:"totalPages"`
	PageSize   uint32 `json:"pageSize"`
	Total      uint32 `json:"total"`
}

type ListResult struct {
	Pagination Pagination         `json:"pagination"`
	Resources  []ListResourceItem `json:"resources"`
}

const (
	Resolution8K    = "8K"
	Resolution4K    = "4K"
	Resolution1080p = "1080p"
	Resolution1080i = "1080i"
	Resolution720p  = "720p"
	ResolutionSD    = "SD" // below 720
)

type VideoDB struct {
	DB     string `json:"db"`
	Link   string `json:"link"`
	Rating string `json:"rating,omitempty"`
}

type ListResourceItem struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Title2     string    `json:"title2,omitempty"`
	Category   string    `json:"category"`
	Size       uint64    `json:"size"`
	Resolution string    `json:"resolution,omitempty"`
	Seeders    uint32    `json:"seeders"`
	Leechers   uint32    `json:"leechers"`
	DBs        []VideoDB `json:"dbs,omitempty"`
	Images     []string  `json:"images,omitempty"`
}

type ResourceDetail struct {
	ListResourceItem

	Mediainfo   string `json:"mediainfo,omitempty"`
	Description string `json:"description,omitempty"`
}
