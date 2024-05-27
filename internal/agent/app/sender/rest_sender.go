package sender

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"github.com/mylastgame/yp-metrics-service/internal/service"
	"log"
	"net/http"
	"time"
)

type RESTSender struct {
	endpoint string
	method   string
	path     string
	logger   *logger.Logger
}

func NewRESTSender(e, m, p string, log *logger.Logger) *RESTSender {
	return &RESTSender{endpoint: e, method: m, path: p, logger: log}
}

func (s *RESTSender) Send(m metrics.Metrics) error {
	body, err := json.Marshal(m)
	if err != nil {
		s.logger.Log.Error("marshal metrics error: " + err.Error())
		return err
	}

	return s.sendData(body)
}

func (s *RESTSender) SendBatch(list []metrics.Metrics) error {
	if len(list) == 0 {
		return fmt.Errorf("empty metrics list, sending canceled")
	}

	body, err := json.Marshal(list)
	if err != nil {
		s.logger.Log.Error("marshal metrics error: " + err.Error())
		return err
	}

	return s.sendData(body)
}

func (s *RESTSender) sendData(body []byte) error {
	// сжимаем содержимое data
	bodyCompressed, err := service.Compress(body)
	if err != nil {
		log.Fatal(err)
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
		SetHeader("Content-Encoding", "gzip").
		SetBody(bodyCompressed).
		Post(req)

	if err != nil {
		s.logger.Log.Error("error sending metrics request: " + err.Error())
		return err
	}

	if res.StatusCode() != http.StatusOK {
		s.logger.Log.Error("bad status on request: " + res.Status())
		return fmt.Errorf("response status code: %d, for url: %s", res.StatusCode(), req)
	}

	return nil
}
