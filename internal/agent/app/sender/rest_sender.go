package sender

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/mylastgame/yp-metrics-service/internal/agent/config"
	"github.com/mylastgame/yp-metrics-service/internal/core/hash"
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"github.com/mylastgame/yp-metrics-service/internal/service"
	"net/http"
	"time"
)

type RESTSender struct {
	endpoint string
	method   string
	path     string
	logger   *logger.Logger
	cfg      *config.Config
}

func NewRESTSender(e, m, p string, log *logger.Logger, cfg *config.Config) *RESTSender {
	return &RESTSender{endpoint: e, method: m, path: p, logger: log, cfg: cfg}
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
	var (
		hashCode []byte
		err      error
	)

	if s.cfg.Key != "" {
		//create request hash
		hashCode, err = hash.GetSHA256Hash(s.cfg.Key, body)

		if err != nil {
			return fmt.Errorf("error while calculate hash of metrics data: %s", err.Error())
		}

	}

	// сжимаем содержимое data
	bodyCompressed, err := service.Compress(body)
	if err != nil {
		return fmt.Errorf("error while compressing data: %s", err.Error())
	}

	req := fmt.Sprintf("%s/%s/", s.endpoint, s.path)
	client := resty.New()
	client.
		// устанавливаем количество повторений
		//SetRetryCount(3).
		// длительность ожидания между попытками
		//SetRetryWaitTime(100 * time.Millisecond).
		// длительность максимального ожидания
		SetRetryMaxWaitTime(301 * time.Millisecond)

	r := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip")

	if s.cfg.Key != "" {
		r = r.SetHeader("Hash", base64.URLEncoding.EncodeToString(hashCode))
	}

	res, err := r.SetBody(bodyCompressed).
		Post(req)

	if err != nil {
		s.logger.Log.Error("error sending metrics request: " + err.Error())
		return NewErrSendRequest(err)
	}

	if res.StatusCode() != http.StatusOK {
		s.logger.Log.Error("bad status on request: " + res.Status())
		return fmt.Errorf("response status code: %d, for url: %s", res.StatusCode(), req)
	}

	return nil
}
