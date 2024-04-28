package collector

import (
	"github.com/mylastgame/yp-metrics-service/internal/agent/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollector_Collect(t *testing.T) {
	type fields struct {
		storage storage.Storage
	}
	tests := []struct {
		name  string
		title string
	}{
		{
			name:  Alloc,
			title: Alloc,
		},
		{
			name:  GCSys,
			title: GCSys,
		},
		{
			name:  RandomValue,
			title: RandomValue,
		},
		{
			name:  TotalAlloc,
			title: TotalAlloc,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collector{
				storage: storage.NewMemStorage(),
			}
			c.Collect()

			assert.NotEqual(t, 0, len(c.storage.GetGauges()), tt.name)
			assert.NotEqual(t, 0, c.storage.GetGauges()[tt.title], tt.name)
			assert.NotEqual(t, 0, c.storage.GetCounters()[PollCount], tt.name)
		})
	}
}
