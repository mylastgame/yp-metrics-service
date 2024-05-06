package storage

import (
	"encoding/json"
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"github.com/mylastgame/yp-metrics-service/internal/server/config"
	repoStorage "github.com/mylastgame/yp-metrics-service/internal/server/storage"
	"os"
	"sync"
)

type FileStorage struct {
	repo repoStorage.Repo
	m    *sync.Mutex
}

var storage *FileStorage

func InitStorage(repo repoStorage.Repo) {
	storage = &FileStorage{
		repo: repo,
		m:    &sync.Mutex{},
	}
}

func Save() error {
	file, err := os.OpenFile(config.FileStoragePath, os.O_WRONLY|os.O_CREATE, 0666)
	defer func() {
		err := file.Close()
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}()

	if err != nil {
		return err
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	gauges := storage.repo.GetGauges()
	for k, v := range gauges {
		storage.m.Lock()
		err = enc.Encode(metrics.Metrics{MType: metrics.Gauge, ID: k, Value: &v})
		storage.m.Unlock()
		if err != nil {
			logger.Log.Error(err.Error())
			return err
		}
	}

	counters := storage.repo.GetCounters()
	for k, v := range counters {
		storage.m.Lock()
		err = enc.Encode(metrics.Metrics{MType: metrics.Counter, ID: k, Delta: &v})
		storage.m.Unlock()
		if err != nil {
			logger.Log.Error(err.Error())
			return err
		}
	}

	return nil
}
