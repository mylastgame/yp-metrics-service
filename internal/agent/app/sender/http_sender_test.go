package sender

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mylastgame/yp-metrics-service/internal/agent/metric"
	"github.com/mylastgame/yp-metrics-service/internal/server/app"
	counterStrg "github.com/mylastgame/yp-metrics-service/internal/server/storage/counter"
	gaugeStrg "github.com/mylastgame/yp-metrics-service/internal/server/storage/gauge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_httpSender_Send(t *testing.T) {
	type fields struct {
		endpoint string
		method   string
	}
	type args struct {
		m metric.Metric
	}
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
			want:       3.006,
			wantErr:    false,
		},
		{
			name:       "case 2",
			m:          metric.Metric{Mtype: "counter", Title: "c1", Val: "22"},
			wantStatus: http.StatusOK,
			want:       22,
			wantErr:    false,
		},
	}

	r := chi.NewRouter()
	gaugeRepo := gaugeStrg.NewMemRepo()
	counterRepo := counterStrg.NewMemRepo()
	app.Setup(r, gaugeRepo, counterRepo)

	s := httptest.NewServer(r)
	defer s.Close()

	sender := httpSender{fmt.Sprintf("%s/update", s.URL), http.MethodPost}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sender.Send(tt.m)
			if tt.wantErr {
				require.Error(t, err, tt.name)
			}

			if tt.m.Mtype == "counter" {
				get, ok := counterRepo.Get(tt.m.Title)
				require.True(t, ok, tt.name)
				assert.Equal(t, tt.want, int(get.Val), tt.name)
			} else {
				get, ok := gaugeRepo.Get(tt.m.Title)
				require.True(t, ok, tt.name)
				assert.Equal(t, tt.want, float64(get.Val), tt.name)
			}
		})
	}
}
