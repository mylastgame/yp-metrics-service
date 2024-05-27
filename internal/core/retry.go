package core

import (
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"time"
)

func Retry(label string, tries int, fn func() error, log *logger.Logger) error {
	err := fn()

	if err != nil {
		for i := 0; i < tries; i++ {
			delay := i*2 + 1
			time.Sleep(time.Duration(delay) * time.Second)
			log.Sugar.Infof("%s: retry with delay %ds", label, delay)

			err = fn()
			if err == nil {
				return nil
			}
		}
	}

	return err
}
