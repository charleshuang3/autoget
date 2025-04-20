package scraper

import (
	"time"

	"github.com/gocolly/colly"
)

type Scraper struct {
	C *colly.Collector
}

func (s *Scraper) SetRequestTimeout(timeout time.Duration) {
	s.C.SetRequestTimeout(timeout)
}

func NewScraper() *Scraper {
	return &Scraper{
		C: NewCollyCollector(),
	}
}

func NewCollyCollector() *colly.Collector {
	return colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.IgnoreRobotsTxt(),
	)
}
