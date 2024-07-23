package storages

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/SpaceSlow/execenv/cmd/metrics"
)

var _ MetricStorage = (*DBStorage)(nil)

type RetryDB struct {
	*sql.DB
	delays []time.Duration
}

func (db *RetryDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	rowCh := make(chan *sql.Row, 1)
	defer close(rowCh)
	queryRowContext := func() error {
		row := db.DB.QueryRowContext(ctx, query, args...)
		var pgConErr *pgconn.ConnectError
		if row.Err() != nil && errors.As(row.Err(), &pgConErr) {
			if len(rowCh) > 0 {
				<-rowCh
			}
			rowCh <- row
			return row.Err()
		}
		if len(rowCh) > 0 {
			<-rowCh
		}
		rowCh <- row
		return nil
	}
	<-metrics.RetryFunc(queryRowContext, db.delays)

	return <-rowCh
}

func (db *RetryDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	resultCh := make(chan sql.Result, 1)
	defer close(resultCh)
	execContext := func() error {
		res, err := db.DB.ExecContext(ctx, query, args...)
		var pgConErr *pgconn.ConnectError
		if err != nil && errors.As(err, &pgConErr) {
			if len(resultCh) > 0 {
				<-resultCh
			}
			resultCh <- res
			return err
		}
		if len(resultCh) > 0 {
			<-resultCh
		}
		resultCh <- res
		return nil
	}
	err := <-metrics.RetryFunc(execContext, db.delays)

	return <-resultCh, err
}

type DBStorage struct {
	ctx context.Context
	db  RetryDB
}

func NewDBStorage(ctx context.Context, dsn string, delays []time.Duration) (*DBStorage, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}
	rdb := RetryDB{db, delays}

	ok, err := checkExistMetricTable(ctx, rdb)
	if err != nil {
		return nil, err
	}

	if !ok {
		if err := createMetricsTable(ctx, rdb); err != nil {
			return nil, err
		}
	}

	return &DBStorage{
		ctx: ctx,
		db:  rdb,
	}, nil
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
		err = metrics.ErrIncorrectMetricTypeOrValue
	}
	return updMetric, err
}

func (s DBStorage) Batch(metricSlice []metrics.Metric) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	for i := range metricSlice {
		switch metricSlice[i].Type {
		case metrics.Gauge:
			_, err = tx.ExecContext(s.ctx, "INSERT INTO metrics (name, is_gauge, value) VALUES ($1, TRUE, $2) ON CONFLICT (name) DO UPDATE SET value=excluded.value;", metricSlice[i].Name, metricSlice[i].Value.(float64))
		case metrics.Counter:
			_, err = tx.ExecContext(s.ctx, "INSERT INTO metrics (name, is_gauge, delta) VALUES ($1, FALSE, $2) ON CONFLICT (name) DO UPDATE SET delta=(excluded.delta + (SELECT delta FROM metrics WHERE (name=$1 AND is_gauge=FALSE) LIMIT 1));", metricSlice[i].Name, metricSlice[i].Value.(int64))
		default:
			err = metrics.ErrIncorrectMetricTypeOrValue
		}
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s DBStorage) Get(metricType metrics.MetricType, name string) (*metrics.Metric, bool) {
	switch metricType {
	case metrics.Gauge:
		var value float64
		row := s.db.QueryRowContext(s.ctx, "SELECT value FROM metrics WHERE (name=$1 AND is_gauge=TRUE) LIMIT 1;", name)
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
		row := s.db.QueryRowContext(s.ctx, "SELECT delta FROM metrics WHERE (name=$1 AND is_gauge=FALSE) LIMIT 1;", name)
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

func checkExistMetricTable(ctx context.Context, db RetryDB) (bool, error) {
	row := db.QueryRowContext(ctx, "SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename  = 'metrics');")

	var tableExist bool
	if err := row.Scan(&tableExist); err != nil {
		return false, err
	}
	return tableExist, nil
}

func createMetricsTable(ctx context.Context, db RetryDB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE metrics (
			id 			SERIAL PRIMARY KEY,
			name 		VARCHAR(30) UNIQUE NOT NULL,
			is_gauge 	BOOLEAN NOT NULL,
			delta 		BIGINT,
			value		DOUBLE PRECISION
		);
		`)

	return err
}
