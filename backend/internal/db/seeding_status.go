package db

import (
	"time"
)

const (
	StoreMaxDays = 30
)

type SeedingStatus struct {
	ID        string `gorm:"primarykey"` // download/hash
	CreatedAt time.Time
	UpdatedAt time.Time

	UploadHistories map[string]int64 `gorm:"serializer:json"`
}

func (s *SeedingStatus) AddToday(b int64) {
	t := time.Now().Format("2006-01-02")
	s.UploadHistories[t] = b
}

func (s *SeedingStatus) GetXDayBefore(x int) (int64, bool) {
	t := time.Now().AddDate(0, 0, -x).Format("2006-01-02")
	n, ok := s.UploadHistories[t]
	return n, ok
}

func (s *SeedingStatus) CleanupHistory() {
	now := time.Now()

	for k := range s.UploadHistories {
		t, _ := time.Parse("2006-01-02", k)
		if now.Sub(t).Hours() > StoreMaxDays*24 {
			delete(s.UploadHistories, k)
		}
	}
}
