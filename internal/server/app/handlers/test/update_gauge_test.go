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

func TestHandler_UpdateGaugeHandler(t *testing.T) {
	var testTable = []struct {
		name        string
		method      string
		url         string
		want        float64
		wantSuccess bool
		status      int
	}{
		{"case1", http.MethodPost, "/update/gauge/c1/1", 1, true, http.StatusOK},
		{"case2", http.MethodPost, "/update/gauge/c1/1", 1, true, http.StatusOK},
		{"case3", http.MethodPost, "/update/gauge/c1/8.001", 8.001, true, http.StatusOK},
		{"case3", http.MethodPost, "/update/gauge/c2/5", 8.001, true, http.StatusOK},
		{"case4", http.MethodPost, "/update/gauge/c1/8,1", 0, false, http.StatusBadRequest},
		{"case5", http.MethodPost, "/update/gauge/c1/8a", 0, false, http.StatusBadRequest},
	}

	r := chi.NewRouter()
	repo := gauge.NewMemRepo()
	app.Setup(r, repo, counter.NewMemRepo())

	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, v := range testTable {
		resp, _ := testRequest(t, ts, v.method, v.url)
		assert.Equal(t, v.status, resp.StatusCode, v.name)

		if !v.wantSuccess {
			continue
		}

		c, ok := repo.Get("c1")

		require.True(t, ok, v.name)
		assert.Equal(t, v.want, float64(c.Val), v.name)

		resp.Body.Close()
	}
}
