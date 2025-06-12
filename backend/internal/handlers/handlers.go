package handlers

import (
	"slices"

	"github.com/charleshuang3/autoget/backend/downloaders"
	"github.com/charleshuang3/autoget/backend/indexers"
	"github.com/charleshuang3/autoget/backend/indexers/mteam"
	"github.com/charleshuang3/autoget/backend/indexers/nyaa"
	"github.com/charleshuang3/autoget/backend/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Service struct {
	config *config.Config

	indexers    map[string]indexers.IIndexer
	downloaders map[string]downloaders.IDownloader
}

func NewService(config *config.Config, db *gorm.DB, downloaders map[string]downloaders.IDownloader) *Service {
	s := &Service{
		config:      config,
		indexers:    map[string]indexers.IIndexer{},
		downloaders: downloaders,
	}

	if config.MTeam != nil {
		s.indexers["m-team"] = mteam.NewMTeam(config.MTeam)
	}
	if config.Nyaa != nil {
		s.indexers["nyaa"] = nyaa.NewClient(config.Nyaa)
	}
	if config.Sukebei != nil {
		s.indexers["sukebei"] = nyaa.NewClient(config.Sukebei)
	}

	return s
}

func (s *Service) SetupRouter(router *gin.RouterGroup) {
	router.GET("/indexers", s.listIndexers)
	router.GET("/indexers/:indexer/categories", s.indexerCategories)
	router.GET("/indexers/:indexer/resources", s.indexerListResources)
	router.GET("/indexers/:indexer/resources/:resource", s.indexerResourceDetail)
	router.GET("/indexers/:indexer/resources/:resource/download", s.indexerDownload)

	router.GET("/downloaders", s.listDownloaders)
}

func (s *Service) listIndexers(c *gin.Context) {
	resp := []string{}
	for k := range s.indexers {
		resp = append(resp, k)
	}
	slices.Sort(resp)
	c.JSON(200, resp)
}

func (s *Service) indexerCategories(c *gin.Context) {
	indexerName := c.Param("indexer")
	indexer, ok := s.indexers[indexerName]
	if !ok {
		c.JSON(404, gin.H{"error": "Indexer not found"})
		return
	}

	categories, err := indexer.Categories()
	if err != nil {
		c.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	c.JSON(200, categories)
}

type ListRequest struct {
	Category  string   `form:"category"`
	Keyword   string   `form:"keyword"`
	Page      uint32   `form:"page"`
	PageSize  uint32   `form:"pageSize"`
	Free      bool     `form:"free"`
	Standards []string `form:"standards"`
}

func (s *Service) indexerListResources(c *gin.Context) {
	indexerName := c.Param("indexer")
	indexer, ok := s.indexers[indexerName]
	if !ok {
		c.JSON(404, gin.H{"error": "Indexer not found"})
		return
	}

	req := &ListRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	lreq := &indexers.ListRequest{
		Category:  req.Category,
		Keyword:   req.Keyword,
		Page:      req.Page,
		PageSize:  req.PageSize,
		Free:      req.Free,
		Standards: req.Standards,
	}

	listResult, err := indexer.List(lreq)
	if err != nil {
		c.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	c.JSON(200, listResult)
}

func (s *Service) indexerResourceDetail(c *gin.Context) {
	indexerName := c.Param("indexer")
	indexer, ok := s.indexers[indexerName]
	if !ok {
		c.JSON(404, gin.H{"error": "Indexer not found"})
		return
	}

	resourceID := c.Param("resource")
	detail, err := indexer.Detail(resourceID, true)
	if err != nil {
		c.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	c.JSON(200, detail)
}

func (s *Service) indexerDownload(c *gin.Context) {
	// TODO
}

type listDownloadersRespItem struct {
	TorrentsDir string `json:"torrents_dir"`
	DownloadDir string `json:"download_dir"`
}

type listDownloadersResp struct {
	Map map[string]listDownloadersRespItem `json:"downloaders"`
}

func (s *Service) listDownloaders(c *gin.Context) {
	m := map[string]listDownloadersRespItem{}
	for name, dl := range s.downloaders {
		m[name] = listDownloadersRespItem{
			TorrentsDir: dl.TorrentsDir(),
			DownloadDir: dl.DownloadDir(),
		}
	}
	c.JSON(200, listDownloadersResp{Map: m})
}
