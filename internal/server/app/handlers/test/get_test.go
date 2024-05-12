package test

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/core/test"
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
	log, err := logger.NewLogger("info")
	if err != nil {
		fmt.Printf("Error init logger: %v/n", err)
		panic(err)
	}
	repo := storage.NewMemRepo()
	fileStorage := test.NewMockFileStorage(repo)
	r := app.NewRouter(repo, fileStorage, log)
	ctx := context.Background()

	repo.Set(ctx, "gauge", "g1", "0.00001")
	repo.Set(ctx, "gauge", "g2", "1")
	repo.Set(ctx, "gauge", "g3", "99.076511")
	repo.Set(ctx, "counter", "c1", "1")
	repo.Set(ctx, "counter", "c2", "1")
	repo.Set(ctx, "counter", "c3", "99")

	gaugeHTML := "Gauges: <ol>"
	//html.SliceToOlLi("Gauges", gauges)
	gauges, err := repo.GetGauges(ctx)
	if err != nil {
		panic(err)
	}
	for k, g := range gauges {
		gaugeHTML += html.Tag("li", fmt.Sprintf("%s: %f", k, g))
	}
	gaugeHTML += "</ol>"

	counterHTML := "Counters: <ol>"
	counters, err := repo.GetCounters(ctx)
	if err != nil {
		panic(err)
	}
	for k, c := range counters {
		counterHTML += html.Tag("li", fmt.Sprintf("%s: %d", k, c))
	}
	counterHTML += "</ol>"

	return r, gaugeHTML + counterHTML
}
