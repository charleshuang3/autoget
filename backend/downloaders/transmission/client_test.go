package transmission

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/charleshuang3/autoget/backend/downloaders/config"
	"github.com/charleshuang3/autoget/backend/internal/db"
	"github.com/hekmon/transmissionrpc/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type requestPayload struct {
	Method    string      `json:"method"`
	Arguments interface{} `json:"arguments,omitempty"`
	Tag       int         `json:"tag,omitempty"`
}

type torrentGetResults struct {
	Torrents []transmissionrpc.Torrent `json:"torrents"`
}

type answerPayload struct {
	Arguments interface{} `json:"arguments"`
	Result    string      `json:"result"`
	Tag       int         `json:"tag"`
}

// fake transmission rpc server
type fakeTransmission struct {
	reqs []*requestPayload
	resp []any
}

func (f *fakeTransmission) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	req := &requestPayload{}
	json.NewDecoder(r.Body).Decode(req)
	f.reqs = append(f.reqs, req)

	w.Header().Set("Content-Type", "application/json")

	body := f.resp[0]
	f.resp = f.resp[1:]

	resp := &answerPayload{
		Arguments: body,
		Result:    "success",
		Tag:       req.Tag,
	}

	json.NewEncoder(w).Encode(resp)
}

func newTorrent(id int64, hash string, status transmissionrpc.TorrentStatus, uploaded int64) transmissionrpc.Torrent {
	return transmissionrpc.Torrent{
		ID:           &id,
		HashString:   &hash,
		Status:       &status,
		UploadedEver: &uploaded,
	}
}

func TestCheckDailySeeding(t *testing.T) {
	fake := &fakeTransmission{}

	serv := httptest.NewServer(http.HandlerFunc(fake.ServeHTTP))

	httpClient = &http.Client{}
	t.Cleanup(func() {
		httpClient = http.DefaultClient
		serv.Close()
	})

	d, err := db.SqliteForTest()
	require.NoError(t, err)

	conf := &config.DownloaderConfig{
		Transmission: &config.TransmissionConfig{
			URL: serv.URL,
		},
		SeedingPolicy: &config.SeedingPolicy{
			IntervalInDays:    3,
			UploadAtLeastInMB: 1,
		},
	}

	client, err := New("test", conf, d)
	require.NoError(t, err)

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	threeDaysAgo := time.Now().AddDate(0, 0, -3).Format("2006-01-02")
	twoMonthsAgo := time.Now().AddDate(0, -2, 0).Format("2006-01-02")

	// r1 is current seeding, less then 3 day
	r1 := &db.DownloadStatus{
		ID: "test/1",
		UploadHistories: map[string]int64{
			yesterday: 0,
		},
	}
	require.NoError(t, d.Create(r1).Error)

	// r2 is current seeding, more then 3 day, latest upload - 3 day ago > 1MB, it will continue seeding
	r2 := &db.DownloadStatus{
		ID: "test/2",
		UploadHistories: map[string]int64{
			threeDaysAgo: 0,
			twoMonthsAgo: 0,
		},
	}
	require.NoError(t, d.Create(r2).Error)

	// r3 is current seeding, more then 3 day, latest upload - 3 day ago < 1MB, it will stop seeding
	r3 := &db.DownloadStatus{
		ID: "test/3",
		UploadHistories: map[string]int64{
			threeDaysAgo: 0,
		},
	}
	require.NoError(t, d.Create(r3).Error)

	// r4 is not current seeding, more then 30 day, record should be deleted
	r4 := &db.DownloadStatus{
		ID:        "test/4",
		UpdatedAt: time.Now().AddDate(0, 0, -1-db.StoreMaxDays),
		UploadHistories: map[string]int64{
			threeDaysAgo: 0,
		},
	}
	require.NoError(t, d.Create(r4).Error)

	fake.resp = []any{
		&torrentGetResults{
			Torrents: []transmissionrpc.Torrent{
				newTorrent(1, "1", transmissionrpc.TorrentStatusSeed, 100),
				newTorrent(2, "2", transmissionrpc.TorrentStatusSeed, 1025*1024),
				newTorrent(3, "3", transmissionrpc.TorrentStatusSeed, 1000*1024),
				newTorrent(4, "4", transmissionrpc.TorrentStatusStopped, 1000*1024),
				// r5 is brand new, insert to db
				newTorrent(5, "5", transmissionrpc.TorrentStatusSeed, 1000*1024),
			},
		},
		&struct{}{},
	}

	client.checkDailySeeding()

	assert.Len(t, fake.reqs, 2)
	assert.Equal(t, "torrent-get", fake.reqs[0].Method)
	assert.Equal(t, "torrent-stop", fake.reqs[1].Method)
	assert.Equal(t, map[string]interface{}{"ids": []interface{}{float64(3)}}, fake.reqs[1].Arguments)

	{
		// r1 new history item added
		r := &db.DownloadStatus{}
		require.NoError(t, d.First(r, "id = ?", "test/1").Error)
		assert.Equal(t, map[string]int64{
			yesterday: 0,
			today:     100,
		}, r.UploadHistories)
	}

	{
		// r2 new history item added and cleanup old item.
		r := &db.DownloadStatus{}
		require.NoError(t, d.First(r, "id = ?", "test/2").Error)
		assert.Equal(t, map[string]int64{
			threeDaysAgo: 0,
			today:        1025 * 1024,
		}, r.UploadHistories)
	}

	{
		// r3 new history item added.
		r := &db.DownloadStatus{}
		require.NoError(t, d.First(r, "id = ?", "test/3").Error)
		assert.Equal(t, map[string]int64{
			threeDaysAgo: 0,
			today:        1000 * 1024,
		}, r.UploadHistories)
	}

	{
		// r4 should be deleted
		r := &db.DownloadStatus{}
		err := d.First(r, "id = ?", "test/4").Error
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	}

	{
		// r5 should be inserted
		r := &db.DownloadStatus{}
		require.NoError(t, d.First(r, "id = ?", "test/5").Error)
		assert.Equal(t, map[string]int64{
			today: 1000 * 1024,
		}, r.UploadHistories)
	}
}
