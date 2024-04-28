package sender

import (
	"github.com/mylastgame/yp-metrics-service/internal/agent/metric"
)

type Sender interface {
	Send(metric metric.Metric) error
}
