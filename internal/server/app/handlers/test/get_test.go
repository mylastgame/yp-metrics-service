package test

import (
	"github.com/go-chi/chi/v5"
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/server/app"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage"
	"github.com/mylastgame/yp-metrics-service/internal/service/html"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHandler(t *testing.T) {
	r, _ := setup()
	logger.Initialize("info")
	ts := httptest.NewServer(r)
	defer ts.Close()

	var testTable = []struct {
		name        string
		method      string
		url         string
		want        string
		wantSuccess bool
		status      int
	}{
		{"case1", http.MethodGet, "/value/counter/c1", "1", true, http.StatusOK},
		{"case2", http.MethodGet, "/value/counter/c3", "99", true, http.StatusOK},
		{"case2", http.MethodGet, "/value/counter/c4", "", false, http.StatusNotFound},
		{"case4", http.MethodGet, "/value/gauge/g1", "0.00001", true, http.StatusOK},
		{"case5", http.MethodGet, "/value/gauge/g3", "99.076511", true, http.StatusOK},
		{"case6", http.MethodGet, "/value/gauge/g4", "", false, http.StatusNotFound},
		//{"case3", http.MethodGet, "/", getAllHtml, true, http.StatusOK},
	}

	for _, v := range testTable {
		resp, get := testRequest(t, ts, v.method, v.url)

		assert.Equal(t, v.status, resp.StatusCode, v.name)

		if !v.wantSuccess {
			continue
		}
		assert.Equal(t, v.want, get, v.name)

		resp.Body.Close()
	}
}

func setup() (chi.Router, string) {
	repo := storage.NewMemRepo()
	r := app.NewRouter(repo)

	repo.Set("gauge", "g1", "0.00001")
	repo.Set("gauge", "g2", "1")
	repo.Set("gauge", "g3", "99.076511")
	repo.Set("counter", "c1", "1")
	repo.Set("counter", "c2", "1")
	repo.Set("counter", "c3", "99")

	gaugeHTML := "Gauges: <ol>"
	for _, g := range repo.GetGauges() {
		gaugeHTML += html.Tag("li", g)
	}
	gaugeHTML += "</ol>"

	counterHTML := "Counters: <ol>"
	for _, c := range repo.GetCounters() {
		counterHTML += html.Tag("li", c)
	}
	counterHTML += "</ol>"

	return r, gaugeHTML + counterHTML
}
