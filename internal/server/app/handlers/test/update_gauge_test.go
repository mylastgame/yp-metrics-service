package test

import (
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"github.com/mylastgame/yp-metrics-service/internal/core/test"
	"github.com/mylastgame/yp-metrics-service/internal/server/app"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_UpdateGaugeHandler(t *testing.T) {
	var testTable = []struct {
		name        string
		method      string
		url         string
		want        string
		wantSuccess bool
		status      int
	}{
		{"case1", http.MethodPost, "/update/gauge/c1/1", "1", true, http.StatusOK},
		{"case2", http.MethodPost, "/update/gauge/c1/1", "1", true, http.StatusOK},
		{"case3", http.MethodPost, "/update/gauge/c1/8.001", "8.001", true, http.StatusOK},
		{"case3", http.MethodPost, "/update/gauge/c2/5", "8.001", true, http.StatusOK},
		{"case4", http.MethodPost, "/update/gauge/c1/8,1", "", false, http.StatusBadRequest},
		{"case5", http.MethodPost, "/update/gauge/c1/8a", "", false, http.StatusBadRequest},
	}

	repo := storage.NewMemRepo()
	fileStorage := test.NewMockFileStorage(repo)
	r := app.NewRouter(repo, fileStorage)

	logger.Initialize("info")
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, v := range testTable {
		resp, _ := testRequest(t, ts, v.method, v.url)
		assert.Equal(t, v.status, resp.StatusCode, v.name)

		if !v.wantSuccess {
			continue
		}

		c, err := repo.Get(metrics.Gauge, "c1")

		require.NoError(t, err, v.name)
		assert.Equal(t, v.want, c, v.name)

		resp.Body.Close()
	}
}
