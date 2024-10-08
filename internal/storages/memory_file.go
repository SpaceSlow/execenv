package storages

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/SpaceSlow/execenv/internal/logger"
	"github.com/SpaceSlow/execenv/internal/metrics"
)

var _ MetricStorage = (*MemFileStorage)(nil)

var ErrNoSpecifyFile = errors.New("no file specified")

// MemFileStorage хранит метрики и в памяти и в файле, поддерживает синхронизацию памяти с файлом.
type MemFileStorage struct {
	*MemStorage
	ctx         context.Context
	f           *os.File
	isSyncStore bool
}

func NewMemFileStorage(ctx context.Context, filename string, duration time.Duration, neededRestore bool) (*MemFileStorage, error) {
	storage := &MemFileStorage{
		ctx:        ctx,
		MemStorage: NewMemStorage(),
		f:          nil,
	}

	if filename == "" {
		return storage, nil
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	storage.f = file

	if neededRestore {
		err = storage.LoadMetricsFromFile()
	} else {
		err = storage.f.Truncate(0)
	}
	if err != nil {
		return nil, err
	}

	storage.isSyncStore = duration == 0
	if !storage.isSyncStore {
		go storage.startStoreMetricsPerSecondsTask(duration)
	}

	return storage, nil
}

func (s *MemFileStorage) Add(metric *metrics.Metric) (*metrics.Metric, error) {
	updMetric, err := s.MemStorage.Add(metric)
	if err != nil {
		return nil, err
	}
	if !s.isSyncStore {
		return updMetric, nil
	}
	return updMetric, s.SaveMetricsToFile()
}

func (s *MemFileStorage) Batch(metricSlice []metrics.Metric) error {
	err := s.MemStorage.Batch(metricSlice)
	if err != nil || !s.isSyncStore {
		return err
	}
	return s.SaveMetricsToFile()
}

func (s *MemFileStorage) Close() error {
	if s.f == nil {
		return nil
	}

	err := s.SaveMetricsToFile()
	if err != nil {
		logger.Log.Error("not saved metrics")
	}

	return s.f.Close()
}

func (s *MemFileStorage) SaveMetricsToFile() error {
	if s.f == nil {
		return ErrNoSpecifyFile
	}
	logger.Log.Info("saving metrics...")

	data, err := json.MarshalIndent(s.List(), "", "    ")
	if err != nil {
		return err
	}
	s.mu.Lock()
	_, err = s.f.WriteAt(data, 0)
	s.mu.Unlock()

	return err
}

func (s *MemFileStorage) LoadMetricsFromFile() error {
	if s.f == nil {
		return ErrNoSpecifyFile
	}

	s.mu.Lock()
	data, err := os.ReadFile(s.f.Name())
	s.mu.Unlock()
	if err != nil || len(data) == 0 {
		return err
	}
	var metricSlice []metrics.Metric
	err = json.Unmarshal(data, &metricSlice)
	if err != nil {
		return err
	}

	for _, metric := range metricSlice {
		if _, err = s.MemStorage.Add(&metric); err != nil {
			return err
		}
	}

	return nil
}

func (s *MemFileStorage) startStoreMetricsPerSecondsTask(duration time.Duration) {
	logger.Log.Info("start store metrics task")
	storeTicker := time.NewTicker(duration)
	defer storeTicker.Stop()
	for {
		select {
		case <-s.ctx.Done():
			logger.Log.Info("store metrics task has been finished")
			return
		case <-storeTicker.C:
			err := s.SaveMetricsToFile()
			if err != nil {
				logger.Log.Error("not saved metrics")
			}
		}
	}
}
