package test

import (
	"github.com/go-chi/chi/v5"
	"github.com/mylastgame/yp-metrics-service/internal/server/app"
	"github.com/mylastgame/yp-metrics-service/internal/server/domain/counter"
	"github.com/mylastgame/yp-metrics-service/internal/server/domain/gauge"
	counterStrg "github.com/mylastgame/yp-metrics-service/internal/server/storage/counter"
	gaugeStrg "github.com/mylastgame/yp-metrics-service/internal/server/storage/gauge"
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
		{"case3", http.MethodGet, "/value/gauge/g1", "0.00001", true, http.StatusOK},
		{"case3", http.MethodGet, "/value/gauge/g3", "99.076511", true, http.StatusOK},
		{"case3", http.MethodGet, "/value/gauge/g4", "", false, http.StatusNotFound},
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
	r := chi.NewRouter()
	gaugeRepo := gaugeStrg.NewMemRepo()
	counterRepo := counterStrg.NewMemRepo()
	app.Setup(r, gaugeRepo, counterRepo)

	gaugeRepo.Save(&gauge.Gauge{Title: "g1", Val: gauge.ValType(0.00001)})
	gaugeRepo.Save(&gauge.Gauge{Title: "g2", Val: gauge.ValType(1)})
	gaugeRepo.Save(&gauge.Gauge{Title: "g3", Val: gauge.ValType(99.076511)})

	counterRepo.Add(&counter.Counter{Title: "c1", Val: counter.ValType(1)})
	counterRepo.Add(&counter.Counter{Title: "c2", Val: counter.ValType(1)})
	counterRepo.Add(&counter.Counter{Title: "c3", Val: counter.ValType(99)})

	gaugeHTML := "Gauges: <ol>"
	for _, g := range gaugeRepo.GetAll() {
		gaugeHTML += html.Tag("li", g)
	}
	gaugeHTML += "</ol>"

	counterHTML := "Counters: <ol>"
	for _, c := range counterRepo.GetAll() {
		counterHTML += html.Tag("li", c)
	}
	counterHTML += "</ol>"

	return r, gaugeHTML + counterHTML
}
