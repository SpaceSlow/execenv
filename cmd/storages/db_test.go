package storages

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SpaceSlow/execenv/cmd/metrics"
)

type postgresContainer struct {
	pool     *dockertest.Pool
	resource *dockertest.Resource
	dsn      string
}

func newPostgresContainer() *postgresContainer {
	return &postgresContainer{}
}

func (c *postgresContainer) Start() error {
	var err error
	c.pool, err = dockertest.NewPool("")
	if err != nil {
		return fmt.Errorf("could not construct pool: %s", err)
	}

	err = c.pool.Client.Ping()
	if err != nil {
		return fmt.Errorf("could not connect to Docker: %s", err)
	}

	c.resource, err = c.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16",
		Env: []string{
			"POSTGRES_USER=test",
			"POSTGRES_PASSWORD=test",
			"POSTGRES_DATABASE=test",
			"listen_addresses = '*'",
		}}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return fmt.Errorf("could not start resource: %s", err)
	}

	c.dsn = fmt.Sprintf("postgres://test:test@localhost:%s/test?sslmode=disable", c.resource.GetPort("5432/tcp"))
	if err := c.pool.Retry(func() error {
		var err error
		db, err := sql.Open("pgx", c.dsn)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		return fmt.Errorf("could not connect to database: %s", err)
	}
	return nil
}

func (c *postgresContainer) Stop() error {
	return c.pool.Purge(c.resource)
}

type DBStorageWithDeleting struct {
	*DBStorage
}

func NewDBStorageWithDeleting(storage *DBStorage) *DBStorageWithDeleting {
	return &DBStorageWithDeleting{storage}
}

func (s DBStorageWithDeleting) DeleteMetrics() {
	s.db.ExecContext(s.ctx, "DELETE FROM metrics")
}

var storage *DBStorageWithDeleting

func TestMain(m *testing.M) {
	container := newPostgresContainer()
	err := container.Start()
	if err != nil {
		log.Fatalf("Could not start postgres container: %s", err)
	}
	defer func(container *postgresContainer) {
		err := container.Stop()
		if err != nil {
			log.Fatalf("Could not stop postgres container: %s", err)
		}
	}(container)
	s, err := NewDBStorage(context.Background(), container.dsn, []time.Duration{time.Second})
	if err != nil {
		log.Printf("Could not create DBStorage: %s", err)
		return
	}
	storage = NewDBStorageWithDeleting(s)

	if err := createMetricsTable(context.Background(), storage.db); err != nil {
		log.Printf("Could not create metrics table: %s", err)
		return
	}

	m.Run()
}

func TestDBStorage_AddGet(t *testing.T) {
	tests := []struct {
		name       string
		metric     *metrics.Metric
		wantMetric *metrics.Metric
		wantErr    error
	}{
		{
			name: "adding counter metric",
			metric: &metrics.Metric{
				Type:  metrics.Counter,
				Name:  "PollCount",
				Value: int64(5),
			},
			wantMetric: &metrics.Metric{
				Type:  metrics.Counter,
				Name:  "PollCount",
				Value: int64(5),
			},
			wantErr: nil,
		},
		{
			name: "adding counter metric (check incrementing PollCount)",
			metric: &metrics.Metric{
				Type:  metrics.Counter,
				Name:  "PollCount",
				Value: int64(5),
			},
			wantMetric: &metrics.Metric{
				Type:  metrics.Counter,
				Name:  "PollCount",
				Value: int64(10),
			},
			wantErr: nil,
		},
		{
			name: "adding gauge metric",
			metric: &metrics.Metric{
				Type:  metrics.Gauge,
				Name:  "RandomValue",
				Value: 6.54,
			},
			wantMetric: &metrics.Metric{
				Type:  metrics.Gauge,
				Name:  "RandomValue",
				Value: 6.54,
			},
			wantErr: nil,
		},
		{
			name: "adding gauge metric (check updating RandomValue)",
			metric: &metrics.Metric{
				Type:  metrics.Gauge,
				Name:  "RandomValue",
				Value: 10.11,
			},
			wantMetric: &metrics.Metric{
				Type:  metrics.Gauge,
				Name:  "RandomValue",
				Value: 10.11,
			},
			wantErr: nil,
		},
		{
			name: "adding metric with incorrect type",
			metric: &metrics.Metric{
				Type:  -1,
				Name:  "RandomValue",
				Value: 6.54,
			},
			wantMetric: nil,
			wantErr:    metrics.ErrIncorrectMetricTypeOrValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := storage.Add(tt.metric)
			require.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantMetric, got)

			stored, _ := storage.Get(tt.metric.Type, tt.metric.Name)
			assert.ObjectsAreEqual(tt.wantMetric, stored)
		})
	}
}

func TestDBStorage_BatchList(t *testing.T) {
	tests := []struct {
		name        string
		metricSlice []metrics.Metric
		wantErr     error
	}{
		{
			name: "batch one metric",
			metricSlice: []metrics.Metric{
				{
					Type:  metrics.Gauge,
					Name:  "GaugeMetric",
					Value: 0.11,
				},
			},
			wantErr: nil,
		},
		{
			name: "batch several metrics",
			metricSlice: []metrics.Metric{
				{
					Type:  metrics.Counter,
					Name:  "PollCount",
					Value: int64(5),
				},
				{
					Type:  metrics.Gauge,
					Name:  "RandomValue",
					Value: 10.11,
				},
			},
			wantErr: nil,
		},
		{
			name: "batch metrics with incorrect type",
			metricSlice: []metrics.Metric{
				{
					Type:  -1,
					Name:  "PollCount",
					Value: int64(15),
				},
				{
					Type:  metrics.Gauge,
					Name:  "RandomValue",
					Value: 3.22,
				},
			},
			wantErr: metrics.ErrIncorrectMetricTypeOrValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.DeleteMetrics()
			err := storage.Batch(tt.metricSlice)
			assert.ErrorIs(t, err, tt.wantErr)

			if err != nil {
				return
			}

			actualMetrics := storage.List()
			assert.ElementsMatch(t, tt.metricSlice, actualMetrics)
		})
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func getRandomMetric() metrics.Metric {
	mType := metrics.MetricType(rand.Intn(2) + 1)

	metric := metrics.Metric{
		Type: mType,
		Name: randStringBytes(rand.Intn(10) + 10),
	}

	switch mType {
	case metrics.Counter:
		metric.Value = rand.Int63n(1000)
	case metrics.Gauge:
		metric.Value = rand.Float64()
	}
	return metric
}

func BenchmarkDBStorage_Batch(b *testing.B) {
	container := newPostgresContainer()
	err := container.Start()
	if err != nil {
		log.Fatalf("Could not start postgres container: %s", err)
	}
	defer func(container *postgresContainer) {
		err := container.Stop()
		if err != nil {
			log.Fatalf("Could not stop postgres container: %s", err)
		}
	}(container)

	metricSlice := make([]metrics.Metric, 50)
	for i := range metricSlice {
		metricSlice[i] = getRandomMetric()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := storage.Batch(metricSlice)
		if err != nil {
			log.Printf("Error occured on batching metrics into storage: %s, metric: %v", err, metricSlice[i])
			return
		}
	}
}
