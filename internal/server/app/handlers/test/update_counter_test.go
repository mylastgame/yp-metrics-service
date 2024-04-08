package test

import (
	"github.com/go-chi/chi/v5"
	"github.com/mylastgame/yp-metrics-service/internal/server/app"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage/counter"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage/gauge"
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
		want        int64
		wantSuccess bool
		status      int
	}{
		{http.MethodPost, "/update/counter/c1/1", 1, true, http.StatusOK},
		{http.MethodPost, "/update/counter/c1/1", 2, true, http.StatusOK},
		{http.MethodPost, "/update/counter/c1/8", 10, true, http.StatusOK},
		{http.MethodPost, "/update/counter/c2/2", 10, true, http.StatusOK},
		{http.MethodPost, "/update/counter/c3/99", 10, true, http.StatusOK},
		{http.MethodPost, "/update/counter/c1/8.1", 0, false, http.StatusBadRequest},
		{http.MethodPost, "/update/counter/c1/8a", 0, false, http.StatusBadRequest},
	}

	r := chi.NewRouter()
	repo := counter.NewMemRepo()
	app.Setup(r, gauge.NewMemRepo(), repo)

	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, v := range testTable {
		resp, _ := testRequest(t, ts, v.method, v.url)

		assert.Equal(t, v.status, resp.StatusCode)

		if !v.wantSuccess {
			continue
		}

		c, ok := repo.Get("c1")

		require.True(t, ok)
		assert.Equal(t, v.want, int64(c.Val))

		resp.Body.Close()
	}
}
