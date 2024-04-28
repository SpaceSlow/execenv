package storages

import (
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SpaceSlow/execenv/cmd/metrics"
	"github.com/SpaceSlow/execenv/cmd/middlewares"
)

type MemFileStorage struct {
	*MemStorage
	f           *os.File
	isSyncStore bool
}

func (s *MemFileStorage) Add(metric *metrics.Metric) error {
	err := s.MemStorage.Add(metric)
	if err != nil || !s.isSyncStore {
		return err
	}
	return s.SaveMetricsToFile()
}

func (s *MemFileStorage) Close() error {
	if s.f == nil {
		return nil
	}
	return s.f.Close()
}

func (s *MemFileStorage) SaveMetricsToFile() error {
	if s.f == nil {
		return errors.New("no file specified")
	}
	middlewares.Log.Info("saving metrics...")

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
		return errors.New("no file specified")
	}

	s.mu.Lock()
	data, err := os.ReadFile(s.f.Name())
	s.mu.Unlock()
	if err != nil || len(data) == 0 {
		return err
	}
	var metricSlice []*metrics.Metric
	err = json.Unmarshal(data, &metricSlice)
	if err != nil {
		return err
	}

	for _, metric := range metricSlice {
		err = s.MemStorage.Add(metric)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *MemFileStorage) startStoreMetricsPerSecondsTask(secs uint) {
	middlewares.Log.Info("start store metrics task")
	closed := make(chan os.Signal, 1)
	defer close(closed)
	signal.Notify(closed, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	interval := time.Duration(secs) * time.Second
	for {
		select {
		case <-closed:
			middlewares.Log.Info("finish store metrics task")
			os.Exit(1)
		case <-time.After(interval):
			err := s.SaveMetricsToFile()
			if err != nil {
				middlewares.Log.Error("not saved metrics")
			}
		}
	}
}

func NewMemFileStorage(filename string, storePerSeconds uint, neededRestore bool) (*MemFileStorage, error) {
	storage := &MemFileStorage{
		MemStorage: &MemStorage{counters: make(map[string]int64), gauges: make(map[string]float64)},
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

	storage.isSyncStore = storePerSeconds == 0
	if !storage.isSyncStore {
		go storage.startStoreMetricsPerSecondsTask(storePerSeconds)
	}

	return storage, nil
}
