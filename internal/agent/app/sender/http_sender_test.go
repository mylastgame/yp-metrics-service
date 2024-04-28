package sender

import (
	"github.com/mylastgame/yp-metrics-service/internal/agent/metric"
	"github.com/mylastgame/yp-metrics-service/internal/server/app"
	"github.com/mylastgame/yp-metrics-service/internal/server/domain/metrics"
	"github.com/mylastgame/yp-metrics-service/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_httpSender_Send(t *testing.T) {
	tests := []struct {
		name       string
		m          metric.Metric
		wantStatus int
		want       any
		wantErr    bool
	}{
		{
			name:       "case 1",
			m:          metric.Metric{Mtype: "gauge", Title: "g1", Val: "3.006"},
			wantStatus: http.StatusOK,
			want:       "3.006",
			wantErr:    false,
		},
		{
			name:       "case 2",
			m:          metric.Metric{Mtype: "counter", Title: "c1", Val: "22"},
			wantStatus: http.StatusOK,
			want:       "22",
			wantErr:    false,
		},
	}

	repo := storage.NewMemRepo()
	r := app.NewRouter(repo)

	s := httptest.NewServer(r)
	defer s.Close()

	sender := httpSender{s.URL, http.MethodPost, "update"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sender.Send(tt.m)
			if tt.wantErr {
				require.Error(t, err, tt.name)
			}

			if tt.m.Mtype == "counter" {
				get, err := repo.Get(metrics.Counter, tt.m.Title)
				require.NoError(t, err, tt.name)
				assert.Equal(t, tt.want, get, tt.name)
			} else {
				get, err := repo.Get(metrics.Gauge, tt.m.Title)
				require.NoError(t, err, tt.name)
				assert.Equal(t, tt.want, get, tt.name)
			}
		})
	}
}
