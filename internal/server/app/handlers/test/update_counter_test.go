package test

import (
	"context"
	"fmt"
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

func TestUpdateCounterHandler(t *testing.T) {
	var testTable = []struct {
		method      string
		url         string
		want        string
		wantSuccess bool
		status      int
	}{
		{http.MethodPost, "/update/counter/c1/1", "1", true, http.StatusOK},
		{http.MethodPost, "/update/counter/c1/1", "2", true, http.StatusOK},
		{http.MethodPost, "/update/counter/c1/8", "10", true, http.StatusOK},
		{http.MethodPost, "/update/counter/c2/2", "10", true, http.StatusOK},
		{http.MethodPost, "/update/counter/c3/99", "10", true, http.StatusOK},
		{http.MethodPost, "/update/counter/c1/8.1", "0", false, http.StatusBadRequest},
		{http.MethodPost, "/update/counter/c1/8a", "0", false, http.StatusBadRequest},
	}

	log, err := logger.NewLogger("info")
	if err != nil {
		fmt.Printf("Error init logger: %v/n", err)
		panic(err)
	}
	repo := storage.NewMemRepo()
	fileStorage := test.NewMockFileStorage(repo)
	r := app.NewRouter(repo, fileStorage, log)

	ts := httptest.NewServer(r)
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		ts.Close()
		cancel()
	}()

	for _, v := range testTable {
		resp, _ := testRequest(t, ts, v.method, v.url)

		assert.Equal(t, v.status, resp.StatusCode)

		if !v.wantSuccess {
			continue
		}

		c, err := repo.Get(ctx, metrics.Counter, "c1")

		require.NoError(t, err)
		assert.Equal(t, v.want, c)

		resp.Body.Close()
	}
}
