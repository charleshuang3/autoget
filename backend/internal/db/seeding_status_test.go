package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreSeedingStatus(t *testing.T) {
	db, err := SqliteForTest()
	require.NoError(t, err)

	want := &SeedingStatus{
		ID: "1",
		UploadHistories: map[string]int64{
			"2025-06-04": 100000,
			"2025-06-05": 100001,
		},
	}

	db.Create(want)

	got := &SeedingStatus{}
	db.First(got, "id = ?", want.ID)

	want.CreatedAt = got.CreatedAt
	want.UpdatedAt = got.UpdatedAt

	assert.Equal(t, want, got)
}

func TestSeedingStatus_AddToday(t *testing.T) {
	s := &SeedingStatus{
		UploadHistories: make(map[string]int64),
	}
	today := time.Now().Format("2006-01-02")
	amount := int64(12345)

	s.AddToday(amount)

	assert.Contains(t, s.UploadHistories, today)
	assert.Equal(t, amount, s.UploadHistories[today])

	// Test adding again for the same day
	s.AddToday(amount + 100)
	assert.Equal(t, amount+100, s.UploadHistories[today])
}

func TestSeedingStatus_GetXDayBefore(t *testing.T) {
	s := &SeedingStatus{
		UploadHistories: make(map[string]int64),
	}

	// Add some historical data
	dayMinus1 := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	dayMinus5 := time.Now().AddDate(0, 0, -5).Format("2006-01-02")
	s.UploadHistories[dayMinus1] = 100
	s.UploadHistories[dayMinus5] = 500

	// Test existing history
	n, ok := s.GetXDayBefore(1)
	assert.True(t, ok)
	assert.Equal(t, int64(100), n)

	n, ok = s.GetXDayBefore(5)
	assert.True(t, ok)
	assert.Equal(t, int64(500), n)

	// Test non-existing history
	n, ok = s.GetXDayBefore(2)
	assert.False(t, ok)
	assert.Equal(t, int64(0), n)
}

func TestSeedingStatus_CleanupHistory(t *testing.T) {
	s := &SeedingStatus{
		UploadHistories: make(map[string]int64),
	}

	now := time.Now()

	// Add entries older than storeMaxDays
	oldDate1 := now.AddDate(0, 0, -(StoreMaxDays + 1)).Format("2006-01-02")
	oldDate2 := now.AddDate(0, 0, -(StoreMaxDays + 5)).Format("2006-01-02")
	s.UploadHistories[oldDate1] = 100
	s.UploadHistories[oldDate2] = 200

	// Add entries within storeMaxDays
	recentDate1 := now.AddDate(0, 0, -5).Format("2006-01-02")
	recentDate2 := now.AddDate(0, 0, -(StoreMaxDays - 1)).Format("2006-01-02")
	s.UploadHistories[recentDate1] = 300
	s.UploadHistories[recentDate2] = 400

	s.CleanupHistory()

	// Verify old entries are removed
	assert.NotContains(t, s.UploadHistories, oldDate1)
	assert.NotContains(t, s.UploadHistories, oldDate2)

	// Verify recent entries are kept
	assert.Contains(t, s.UploadHistories, recentDate1)
	assert.Equal(t, int64(300), s.UploadHistories[recentDate1])
	assert.Contains(t, s.UploadHistories, recentDate2)
	assert.Equal(t, int64(400), s.UploadHistories[recentDate2])
}
