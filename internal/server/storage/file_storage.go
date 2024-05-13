package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"github.com/mylastgame/yp-metrics-service/internal/core/logger"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"github.com/mylastgame/yp-metrics-service/internal/server/config"
	"os"
	"sync"
)

type PersistentStorage interface {
	Save(context.Context) error
	Restore(context.Context) error
}

type FileStorage struct {
	repo   Repo
	m      *sync.Mutex
	logger *logger.Logger
}

func NewFileStorage(repo Repo, log *logger.Logger) *FileStorage {
	return &FileStorage{
		repo:   repo,
		m:      &sync.Mutex{},
		logger: log,
	}
}

func (s *FileStorage) Save(ctx context.Context) error {
	s.logger.Sugar.Info("Saving file to storage")
	var (
		err      error
		gauges   metrics.GaugeList
		counters metrics.CounterList
	)
	file, err := os.OpenFile(config.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	//s.m.Lock()
	defer func() {
		//s.m.Unlock()
		err := file.Close()
		if err != nil {
			s.logger.Log.Error(err.Error())
		}
	}()

	if err != nil {
		return err
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", "")

	gauges, err = s.repo.GetGauges(ctx)
	if err != nil && !errors.Is(err, ErrorNotExists) {
		return err
	}
	for k, v := range gauges {
		err = enc.Encode(metrics.Metrics{MType: metrics.Gauge, ID: k, Value: &v})
		if err != nil {
			s.logger.Log.Error(err.Error())
			return err
		}
	}

	counters, err = s.repo.GetCounters(ctx)
	if err != nil && !errors.Is(err, ErrorNotExists) {
		return err
	}
	for k, v := range counters {
		err = enc.Encode(metrics.Metrics{MType: metrics.Counter, ID: k, Delta: &v})
		if err != nil {
			s.logger.Log.Error(err.Error())
			return err
		}
	}

	return nil
}

func (s *FileStorage) Restore(ctx context.Context) error {
	s.logger.Sugar.Info("Restoring data to repository from file")
	file, err := os.OpenFile(config.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	defer func() {
		err := file.Close()
		if err != nil {
			s.logger.Log.Error(err.Error())
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
			s.logger.Log.Error(err.Error())
			return err
		}

		err = s.repo.SaveMetric(ctx, metric)
		if err != nil {
			s.logger.Log.Error(err.Error())
			return err
		}
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	return nil
}
