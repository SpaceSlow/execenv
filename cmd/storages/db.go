package storages

import (
	"context"
	"database/sql"
	"github.com/SpaceSlow/execenv/cmd/metrics"
	_ "github.com/jackc/pgx/v5"
)

type DBStorage struct {
	ctx context.Context
	db  *sql.DB
}

func (s DBStorage) Add(metric *metrics.Metric) error {
	//TODO implement me
	panic("implement me")
}

func (s DBStorage) Get(metricType metrics.MetricType, name string) (*metrics.Metric, bool) {
	rows, err := s.db.QueryContext(s.ctx, "SELECT is_gauge, delta, value FROM metrics WHERE name = $1 LIMIT 1;", name)
	if err != nil {
		return nil, false
	}

	var (
		isGauge bool
		delta   float64
		value   int64
	)
	if err := rows.Scan(&isGauge, &delta, &value); err != nil {
		return nil, false
	}

	if isGauge && metricType == metrics.Gauge {
		return &metrics.Metric{
			Type:  metricType,
			Name:  name,
			Value: delta,
		}, true
	} else if !isGauge && metricType == metrics.Counter {
		return &metrics.Metric{
			Type:  metricType,
			Name:  name,
			Value: value,
		}, true
	}
	return nil, false
}

func (s DBStorage) List() []metrics.Metric {
	rows, err := s.db.QueryContext(s.ctx, "SELECT name, is_gauge, delta, value FROM metrics;")
	if err != nil {
		return make([]metrics.Metric, 0)
	}

	metricSlice := make([]metrics.Metric, 0)
	var (
		name    string
		isGauge bool
		delta   float64
		value   int64
	)
	for rows.Next() {
		m := metrics.Metric{}

		if err := rows.Scan(&name, &isGauge, &delta, &value); err != nil {
			return make([]metrics.Metric, 0)
		}

		if isGauge {
			m.Type = metrics.Gauge
			m.Value = delta
		} else {
			m.Type = metrics.Counter
			m.Value = value
		}
		m.Name = name

		metricSlice = append(metricSlice, m)
	}

	if rows.Err() != nil {
		return make([]metrics.Metric, 0)
	}

	return metricSlice
}

func (s DBStorage) Close() error {
	return s.db.Close()
}

func (s DBStorage) CheckConnection() bool {
	return s.db.PingContext(s.ctx) == nil
}

func NewDBStorage(ctx context.Context, dsn string) (*DBStorage, error) {
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	ok, err := checkExistMetricTable(ctx, db)
	if err != nil {
		return nil, err
	}

	if !ok {
		if err := createMetricsTable(ctx, db); err != nil {
			return nil, err
		}
	}

	return &DBStorage{
		ctx: ctx,
		db:  db,
	}, nil
}

func checkExistMetricTable(ctx context.Context, db *sql.DB) (bool, error) {
	rows, err := db.QueryContext(ctx, "SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename  = 'metrics');")
	if err != nil {
		return false, err
	}

	var tableExist bool
	if err := rows.Scan(&tableExist); err != nil {
		return false, err
	}
	return true, nil
}

func createMetricsTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE metrics 
		(
			id 			SERIAL PRIMARY KEY,
			name 		VARCHAR(30),
			is_gauge 	BOOLEAN,
			delta 		DOUBLE PRECISION,
			value		INTEGER,
		);
		`)

	return err
}
