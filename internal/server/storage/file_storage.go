package storage

import (
	"bufio"
	"encoding/json"
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"github.com/mylastgame/yp-metrics-service/internal/server/config"
	"os"
	"sync"
)

type PersistentStorage interface {
	Save() error
	Restore() error
}

type FileStorage struct {
	repo Repo
	m    *sync.Mutex
}

func NewFileStorage(repo Repo) *FileStorage {
	return &FileStorage{
		repo: repo,
		m:    &sync.Mutex{},
	}
}

func (s *FileStorage) Save() error {
	logger.Sugar.Info("Saving file to storage")
	file, err := os.OpenFile(config.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	s.m.Lock()
	defer func() {
		s.m.Unlock()
		err := file.Close()
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}()

	if err != nil {
		return err
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", "")

	gauges := s.repo.GetGauges()
	for k, v := range gauges {
		err = enc.Encode(metrics.Metrics{MType: metrics.Gauge, ID: k, Value: &v})
		if err != nil {
			logger.Log.Error(err.Error())
			return err
		}
	}

	counters := s.repo.GetCounters()
	for k, v := range counters {
		err = enc.Encode(metrics.Metrics{MType: metrics.Counter, ID: k, Delta: &v})
		if err != nil {
			logger.Log.Error(err.Error())
			return err
		}
	}

	return nil
}

func (s *FileStorage) Restore() error {
	logger.Sugar.Info("Restoring data to repository from file")
	s.m.Lock()
	file, err := os.OpenFile(config.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	defer func() {
		s.m.Unlock()
		err := file.Close()
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}()

	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)

	// одиночное сканирование до следующей строки
	for scanner.Scan() {
		data := scanner.Bytes()
		metric := metrics.Metrics{}

		err = json.Unmarshal(data, &metric)
		if err != nil {
			logger.Log.Error(err.Error())
			return err
		}

		err = s.repo.SaveMetric(metric)
		if err != nil {
			logger.Log.Error(err.Error())
			return err
		}
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	return nil
}
