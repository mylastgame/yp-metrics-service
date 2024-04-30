package sender

import (
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
)

type Sender interface {
	Send(metric metrics.Metrics) error
}
