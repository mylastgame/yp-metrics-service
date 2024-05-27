package sender

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

func Test_httpSender_Send(t *testing.T) {
	val1 := 3.006
	var val2 int64 = 22

	tests := []struct {
		name       string
		m          metrics.Metrics
		wantStatus int
		want       any
		wantErr    bool
	}{
		{
			name:       "case 1",
			m:          metrics.Metrics{MType: "gauge", ID: "g1", Value: &val1},
			wantStatus: http.StatusOK,
			want:       "3.006",
			wantErr:    false,
		},
		{
			name:       "case 2",
			m:          metrics.Metrics{MType: "counter", ID: "c1", Delta: &val2},
			wantStatus: http.StatusOK,
			want:       "22",
			wantErr:    false,
		},
	}

	log, err := logger.NewLogger("info")
	if err != nil {
		fmt.Printf("Error init logger: %v/n", err)
		panic(err)
	}
	repo := storage.NewMemRepo()
	fileStorage := test.NewMockFileStorage(repo)
	r := app.NewRouter(repo, fileStorage, log)

	s := httptest.NewServer(r)
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		s.Close()
		cancel()
	}()

	sender := httpSender{s.URL, http.MethodPost, "update"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sender.Send(tt.m)
			if tt.wantErr {
				require.Error(t, err, tt.name)
			}

			if tt.m.MType == "counter" {
				get, err := repo.Get(ctx, metrics.Counter, tt.m.ID)
				require.NoError(t, err, tt.name)
				assert.Equal(t, tt.want, get, tt.name)
			} else {
				get, err := repo.Get(ctx, metrics.Gauge, tt.m.ID)
				require.NoError(t, err, tt.name)
				assert.Equal(t, tt.want, get, tt.name)
			}
		})
	}
}
