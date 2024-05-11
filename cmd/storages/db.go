package storages

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SpaceSlow/execenv/cmd/metrics"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	ctx context.Context
	db  *sql.DB
}

func (s DBStorage) Add(metric *metrics.Metric) (*metrics.Metric, error) {
	var (
		updMetric *metrics.Metric
		err       error
	)
	switch metric.Type {
	case metrics.Gauge:
		_, err = s.db.ExecContext(s.ctx, "INSERT INTO metrics (name, is_gauge, value) VALUES ($1, TRUE, $2) ON CONFLICT (name) DO UPDATE SET value=excluded.value;", metric.Name, metric.Value.(float64))
		if err != nil {
			return nil, err
		}
		updMetric = metric.Copy()
	case metrics.Counter:
		row := s.db.QueryRowContext(s.ctx, "SELECT delta FROM metrics WHERE (name=$1 AND is_gauge=FALSE) LIMIT 1;", metric.Name)
		var prevValue int64
		err := row.Scan(&prevValue)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		updValue := metric.Value.(int64) + prevValue
		_, err = s.db.ExecContext(s.ctx, "INSERT INTO metrics (name, is_gauge, delta) VALUES ($1, FALSE, $2) ON CONFLICT (name) DO UPDATE SET delta=excluded.delta;", metric.Name, updValue)

		updMetric = metric.Copy()
		updMetric.Value = updValue
		if err != nil {
			return nil, err
		}
	default:
		err = &metrics.IncorrectMetricTypeOrValueError{}
	}
	return updMetric, nil
}

func (s DBStorage) Batch(metricSlice []metrics.Metric) ([]metrics.Metric, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	updMetrics := make([]metrics.Metric, 0, len(metricSlice))
	for _, metric := range metricSlice {
		var (
			updMetric *metrics.Metric
			err       error
		)
		switch metric.Type {
		case metrics.Gauge:
			_, err = tx.ExecContext(s.ctx, "INSERT INTO metrics (name, is_gauge, value) VALUES ($1, TRUE, $2) ON CONFLICT (name) DO UPDATE SET value=excluded.value;", metric.Name, metric.Value.(float64))
			if err != nil {
				return nil, err
			}
			updMetric = metric.Copy()
		case metrics.Counter:
			row := tx.QueryRowContext(s.ctx, "SELECT delta FROM metrics WHERE (name=$1 AND is_gauge=FALSE) LIMIT 1;", metric.Name)
			var prevValue int64
			err := row.Scan(&prevValue)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
			updValue := metric.Value.(int64) + prevValue
			_, err = tx.ExecContext(s.ctx, "INSERT INTO metrics (name, is_gauge, delta) VALUES ($1, FALSE, $2) ON CONFLICT (name) DO UPDATE SET delta=excluded.delta;", metric.Name, updValue)

			updMetric = metric.Copy()
			updMetric.Value = updValue
			if err != nil {
				return nil, err
			}
		default:
			err = &metrics.IncorrectMetricTypeOrValueError{}
		}
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		updMetrics = append(updMetrics, *updMetric)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return updMetrics, nil
}

func (s DBStorage) Get(metricType metrics.MetricType, name string) (*metrics.Metric, bool) {
	switch metricType {
	case metrics.Gauge:
		var value float64
		row := s.db.QueryRowContext(s.ctx, "SELECT value FROM metrics WHERE (name = $1 AND is_gauge=TRUE) LIMIT 1;", name)
		if err := row.Scan(&value); err != nil {
			return nil, false
		}

		return &metrics.Metric{
			Type:  metricType,
			Name:  name,
			Value: value,
		}, true
	case metrics.Counter:
		var delta int64
		row := s.db.QueryRowContext(s.ctx, "SELECT delta FROM metrics WHERE (name = $1 AND is_gauge=FALSE) LIMIT 1;", name)
		if err := row.Scan(&delta); err != nil {
			return nil, false
		}

		return &metrics.Metric{
			Type:  metricType,
			Name:  name,
			Value: delta,
		}, true
	default:
		return nil, false
	}
}

func (s DBStorage) List() []metrics.Metric {
	rows, err := s.db.QueryContext(s.ctx, "SELECT name, is_gauge, delta, value FROM metrics;")
	if err != nil {
		return make([]metrics.Metric, 0)
	}

	defer rows.Close()

	metricSlice := make([]metrics.Metric, 0)
	var (
		name    string
		isGauge bool
		delta   *int64
		value   *float64
	)
	for rows.Next() {
		m := metrics.Metric{}

		if err := rows.Scan(&name, &isGauge, &delta, &value); err != nil {
			return make([]metrics.Metric, 0)
		}

		if isGauge {
			m.Type = metrics.Gauge
			m.Value = *value
		} else {
			m.Type = metrics.Counter
			m.Value = *delta
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
	db, err := sql.Open("pgx", dsn)

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
	row := db.QueryRowContext(ctx, "SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename  = 'metrics');")

	var tableExist bool
	if err := row.Scan(&tableExist); err != nil {
		return false, err
	}
	return tableExist, nil
}

func createMetricsTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE metrics (
			id 			SERIAL PRIMARY KEY,
			name 		VARCHAR(30) UNIQUE NOT NULL,
			is_gauge 	BOOLEAN NOT NULL,
			delta 		INTEGER,
			value		DOUBLE PRECISION
		);
		`)

	return err
}
