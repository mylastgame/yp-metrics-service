package sender

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"net/http"
	"time"
)

type RESTSender struct {
	endpoint string
	method   string
	path     string
}

func NewRESTSender(e, m, p string) *RESTSender {
	return &RESTSender{endpoint: e, method: m, path: p}
}

func (s *RESTSender) Send(m metrics.Metrics) error {
	body, err := json.Marshal(m)
	if err != nil {
		logger.Log.Error("marshal metrics error: " + err.Error())
		return err
	}

	req := fmt.Sprintf("%s/%s/", s.endpoint, s.path)
	client := resty.New()
	client.
		// устанавливаем количество повторений
		SetRetryCount(3).
		// длительность ожидания между попытками
		SetRetryWaitTime(100 * time.Millisecond).
		// длительность максимального ожидания
		SetRetryMaxWaitTime(301 * time.Millisecond)

	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(req)

	if err != nil {
		logger.Log.Error("error sending metrics request: " + err.Error())
		return err
	}

	if res.StatusCode() != http.StatusOK {
		logger.Log.Error("bad status on request: " + res.Status())
		return fmt.Errorf("response status code: %d, for url: %s", res.StatusCode(), req)
	}

	return nil
}
