package db

import (
	"time"

	"gorm.io/gorm"
)

const (
	StoreMaxDays = 30
)

type DownloadState uint

const (
	DownloadStarted DownloadState = iota
	DownloadSeeding
	DownloadStopped
	DownloadDeleted
)

type MoveState uint

const (
	UnMoved MoveState = iota
	Moved
	Organized
)

type DownloadStatus struct {
	ID        string `gorm:"primarykey"` // hash
	CreatedAt time.Time
	UpdatedAt time.Time

	Downloader       string
	DownloadProgress int32 // in x/1000
	State            DownloadState

	UploadHistories map[string]int64 `gorm:"serializer:json"`

	ResIndexer string
	ResTitle   string
	ResTitle2  string
	Category   string
	FileList   []string `gorm:"serializer:json"`

	MoveState MoveState
}

func (s *DownloadStatus) AddToday(b int64) {
	t := time.Now().Format("2006-01-02")
	s.UploadHistories[t] = b
}

func (s *DownloadStatus) GetXDayBefore(x int) (int64, bool) {
	t := time.Now().AddDate(0, 0, -x).Format("2006-01-02")
	n, ok := s.UploadHistories[t]
	return n, ok
}

func (s *DownloadStatus) CleanupHistory() {
	now := time.Now()

	for k := range s.UploadHistories {
		t, _ := time.Parse("2006-01-02", k)
		if now.Sub(t).Hours() > StoreMaxDays*24 {
			delete(s.UploadHistories, k)
		}
	}
}

func GetUnfinishedDownloadStatusByDownloader(db *gorm.DB, downloader string) ([]DownloadStatus, error) {
	var ss []DownloadStatus
	err := db.Where("downloader = ?", downloader).Where("state = ?", DownloadStarted).Find(&ss).Error
	return ss, err
}

func GetFinishedUnmoveedDownloadStatusByDownloader(db *gorm.DB, downloader string) ([]DownloadStatus, error) {
	var ss []DownloadStatus
	err := db.Where("downloader = ?", downloader).Where("state >= ?", DownloadSeeding).Where("move_state = ?", UnMoved).Find(&ss).Error
	return ss, err
}

func GetStoppedMovedDownloadStatusByDownloader(db *gorm.DB, downloader string) ([]DownloadStatus, error) {
	var ss []DownloadStatus
	err := db.Where("downloader = ?", downloader).Where("state = ?", DownloadStopped).Where("move_state >= ?", Moved).Find(&ss).Error
	return ss, err
}

func GetDownloadStatus(db *gorm.DB, hash string) (*DownloadStatus, error) {
	s := &DownloadStatus{}
	err := db.First(s, "id = ?", hash).Error
	return s, err
}

func SaveDownloadStatus(db *gorm.DB, s *DownloadStatus) error {
	return db.Save(s).Error
}

func RemoveDownloadStatus(db *gorm.DB, id string) error {
	return db.Where("id = ?", id).Delete(&DownloadStatus{}).Error
}

func UpdateDownloadStateForStatuses(db *gorm.DB, ids []string, state DownloadState) error {
	return db.Model(&DownloadStatus{}).Where("id IN ?", ids).Update("state", state).Error
}
