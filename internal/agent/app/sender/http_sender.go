package sender

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"net/http"
	"time"
)

type httpSender struct {
	endpoint string
	method   string
	path     string
}

func NewHTTPSender(e, m, p string) *httpSender {
	return &httpSender{endpoint: e, method: m, path: p}
}

func (s *httpSender) Send(m metrics.Metrics) error {
	var val string
	if m.MType == metrics.Counter {
		val = fmt.Sprintf("%d", *m.Delta)
	} else if m.MType == metrics.Gauge {
		val = fmt.Sprintf("%f", *m.Value)
	}

	req := fmt.Sprintf("%s/%s/%s/%s/%s", s.endpoint, s.path, m.MType, m.ID, val)

	client := resty.New()
	client.
		// устанавливаем количество повторений
		SetRetryCount(3).
		// длительность ожидания между попытками
		SetRetryWaitTime(100 * time.Millisecond).
		// длительность максимального ожидания
		SetRetryMaxWaitTime(301 * time.Millisecond)

	res, err := client.R().
		Post(req)

	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("response status code: %d, for url: %s", res.StatusCode(), req)
	}

	return nil
}

func (s *httpSender) SendBatch(metrics []metrics.Metrics) error {
	return fmt.Errorf("not implemented")
}
