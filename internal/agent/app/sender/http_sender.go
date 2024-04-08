package sender

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/mylastgame/yp-metrics-service/internal/agent/metric"
	"time"
)

type httpSender struct {
	server string
	method string
}

func NewHTTPSender(server, method, prefix string) *httpSender {
	return &httpSender{server: server, method: method}
}

func (s *httpSender) Send(m metric.Metric) error {
	req := fmt.Sprintf("%s/%s/%s/%s", s.server, m.Mtype, m.Title, m.Val)

	client := resty.New()
	client.
		// устанавливаем количество повторений
		SetRetryCount(3).
		// длительность ожидания между попытками
		SetRetryWaitTime(100 * time.Millisecond).
		// длительность максимального ожидания
		SetRetryMaxWaitTime(301 * time.Millisecond)

	r, err := client.R().
		Post(req)

	//r, err := http.NewRequest(s.method, req, nil)
	//if err != nil {
	//	return err
	//}
	//
	//client := &http.Client{}
	//res, err := client.Do(r)
	if err != nil {
		return err
	}
	//defer res.Body.Close()

	return nil
}
