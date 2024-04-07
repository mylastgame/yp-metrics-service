package sender

import (
	"fmt"
	"github.com/mylastgame/yp-metrics-service/internal/agent/metric"
	"net/http"
)

type httpSender struct {
	server string
	method string
	prefix string
}

func NewHttpSender(server, method, prefix string) *httpSender {
	return &httpSender{server: server, method: method, prefix: prefix}
}

func (s *httpSender) Send(m metric.Metric) error {
	req := fmt.Sprintf("%s/%s/%s/%s/%s", s.server, s.prefix, m.Mtype, m.Title, m.Val)
	r, err := http.NewRequest(s.method, req, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
