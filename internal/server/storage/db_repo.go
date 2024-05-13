package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mylastgame/yp-metrics-service/internal/core/metrics"
	"github.com/mylastgame/yp-metrics-service/internal/service/convert"
)

const MetricsTable = "metrics"

type DBRepo struct {
	conn *sql.DB
}

func NewDBRepo(ctx context.Context, db *sql.DB) (*DBRepo, error) {
	repo := &DBRepo{conn: db}

	err := repo.Bootstrap(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing DB: %v", err)
	}

	return repo, nil
}

func (r *DBRepo) Close() error {
	err := r.conn.Close()
	if err != nil {
		return fmt.Errorf("error closing DB connection: %v", err)
	}

	return err
}

func (r *DBRepo) Bootstrap(ctx context.Context) error {
	tx, err := r.conn.BeginTx(ctx, nil)
	defer tx.Commit()

	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}

	//drop table if exists
	_, err = tx.ExecContext(ctx, `DROP TABLE IF EXISTS metrics`)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("error drop table: %v", err)
	}

	//create table
	_, err = tx.ExecContext(ctx, `
        CREATE TABLE metrics (
            id varchar(128) PRIMARY KEY NOT NULL,
            type varchar(128) NOT NULL,
            value double precision NULL,
            delta bigint NULL
        )
    `)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("error creating database: %v", err)
	}

	//create index?
	_, err = tx.ExecContext(ctx, `CREATE INDEX type_idx ON metrics (type)`)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("error creating index: %v", err)
	}

	return nil
}

func (r *DBRepo) Ping() error {
	return r.conn.Ping()
}

func (r *DBRepo) Get(ctx context.Context, t string, k string) (string, error) {
	if t == metrics.Gauge {
		val, err := r.GetGauge(ctx, k)
		if err != nil {
			if err != ErrorNotExists {
				return "", fmt.Errorf("error getting metrics: %v", err)
			}

			return "", err
		}

		return convert.GaugeToString(val), nil
	}

	if t == metrics.Counter {
		val, err := r.GetCounter(ctx, k)
		if err != nil {
			if err != ErrorNotExists {
				return "", fmt.Errorf("error getting metrics: %v", err)
			}

			return "", err
		}

		return convert.CounterToString(val), nil
	}

	return "", NewStorageError(BadMetricType, t, k)
}

func (r *DBRepo) Set(ctx context.Context, t string, id string, v string) error {
	if t == metrics.Gauge {
		gaugeVal, err := convert.StringToGauge(v)
		if err != nil {
			return NewStorageError(BadValue, t, v)
		}

		return r.SetGauge(ctx, id, gaugeVal)
	}

	if t == metrics.Counter {
		counterVal, err := convert.StringToCounter(v)
		if err != nil {
			return NewStorageError(BadValue, t, v)
		}

		return r.SetCounter(ctx, id, counterVal)
	}

	return NewStorageError(BadMetricType, t, id)
}

func (r *DBRepo) SetGauge(ctx context.Context, id string, v float64) error {
	tx, err := r.conn.BeginTx(ctx, nil)
	defer tx.Commit()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}

	row := r.conn.QueryRowContext(ctx,
		"SELECT id FROM metrics WHERE type = $1 AND id = $2",
		metrics.Gauge, id,
	)

	var existedID string

	err = row.Scan(&existedID)
	if err != nil {
		if err != sql.ErrNoRows {
			_ = tx.Rollback()
			return fmt.Errorf("error scanning metrics: %v", err)
		}

		//Insert new value
		_, err = r.conn.ExecContext(ctx, "INSERT INTO metrics (id ,type, value) VALUES ($1, $2, $3)", id, metrics.Gauge, v)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("error inserting metrics: %v", err)
		}

	} else {
		//update old value
		_, err = tx.ExecContext(ctx, "UPDATE metrics SET value = $1 WHERE id = $2 AND type = $3", v, id, metrics.Gauge)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("error updating metrics: %v", err)
		}
	}

	return nil
}

func (r *DBRepo) SetCounter(ctx context.Context, id string, v int64) error {
	tx, err := r.conn.BeginTx(ctx, nil)
	defer tx.Commit()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}

	row := r.conn.QueryRowContext(ctx,
		"SELECT id FROM metrics WHERE type = $1 AND id = $2",
		metrics.Counter, id,
	)

	var existedID string

	err = row.Scan(&existedID)
	if err != nil {
		if err != sql.ErrNoRows {
			_ = tx.Rollback()
			return fmt.Errorf("error scanning metrics: %v", err)
		}

		//Insert new value
		_, err = r.conn.ExecContext(ctx, "INSERT INTO metrics (id ,type, delta) VALUES ($1, $2, $3)", id, metrics.Counter, v)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("error inserting metrics: %v", err)
		}

	} else {
		//update old value
		_, err = tx.ExecContext(ctx, "UPDATE metrics SET delta = delta + $1 WHERE id = $2 AND type = $3", v, id, metrics.Counter)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("error updating metrics: %v", err)
		}
	}

	return nil
}

func (r *DBRepo) GetCounter(ctx context.Context, id string) (int64, error) {
	var counter int64
	row := r.conn.QueryRowContext(ctx, "SELECT delta FROM metrics WHERE type = $1 AND id = $2", metrics.Counter, id)
	err := row.Scan(&counter)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, fmt.Errorf("error scanning metrics: %v", err)
		}

		return 0, ErrorNotExists
	}

	return counter, nil
}

func (r *DBRepo) GetGauge(ctx context.Context, id string) (float64, error) {
	var gauge float64
	row := r.conn.QueryRowContext(ctx, "SELECT value FROM metrics WHERE type = $1 AND id = $2", metrics.Gauge, id)
	err := row.Scan(&gauge)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, fmt.Errorf("error scanning metrics: %v", err)
		}

		return 0, ErrorNotExists
	}

	return gauge, nil
}

func (r *DBRepo) GetGauges(ctx context.Context) (metrics.GaugeList, error) {
	res := metrics.GaugeList{}
	rows, err := r.conn.QueryContext(ctx, "SELECT id, value FROM metrics WHERE type = $1", metrics.Gauge)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("error scanning metrics: %v", err)
		}

		return res, ErrorNotExists
	}

	var (
		value float64
		id    string
	)
	for rows.Next() {
		err = rows.Scan(&id, &value)
		if err != nil {
			return nil, fmt.Errorf("error scanning metrics: %v", err)
		}

		res[id] = value
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return res, nil
}

func (r *DBRepo) GetCounters(ctx context.Context) (metrics.CounterList, error) {
	res := metrics.CounterList{}
	rows, err := r.conn.QueryContext(ctx, "SELECT id, delta FROM metrics WHERE type = $1", metrics.Counter)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("error scanning metrics: %v", err)
		}

		return res, ErrorNotExists
	}

	var (
		value int64
		id    string
	)
	for rows.Next() {
		err = rows.Scan(&id, &value)
		if err != nil {
			return nil, fmt.Errorf("error scanning metrics: %v", err)
		}

		res[id] = value
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return res, nil
}

func (r *DBRepo) SaveMetric(ctx context.Context, metric metrics.Metrics) error {
	if metric.MType == metrics.Gauge {
		err := r.SetGauge(ctx, metric.ID, *metric.Value)
		if err != nil {
			return err
		}
		return nil
	}

	if metric.MType == metrics.Counter {
		err := r.SetCounter(ctx, metric.ID, *metric.Delta)
		if err != nil {
			return err
		}
		return nil
	}

	return NewStorageError(BadMetricType, metric.MType, metric.ID)
}

func (r *DBRepo) GetMetric(ctx context.Context, mType string, id string) (metrics.Metrics, error) {
	metric := metrics.Metrics{}
	var (
		err  error
		gVal float64
		cVal int64
	)

	if mType == metrics.Gauge {
		gVal, err = r.GetGauge(ctx, id)
		if err == nil {
			metric.MType = mType
			metric.ID = id
			metric.Value = &gVal
		}

		return metric, err
	}

	if mType == metrics.Counter {
		cVal, err = r.GetCounter(ctx, id)
		if err == nil {
			metric.MType = mType
			metric.ID = id
			metric.Delta = &cVal
		}
		return metric, err
	}

	return metric, NewStorageError(BadMetricType, metric.MType, metric.ID)
}

func (r *DBRepo) SaveMetrics(ctx context.Context, list []metrics.Metrics) error {
	tx, err := r.conn.BeginTx(ctx, nil)
	defer tx.Commit()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}

	gaugeQuery := `
		INSERT INTO metrics (id ,type, value) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (id) DO UPDATE SET value = EXCLUDED.value
	`
	gaugeStmt, err := tx.PrepareContext(ctx, gaugeQuery)
	if err != nil {
		return fmt.Errorf("error preparing gauge query: %v", err)
	}
	defer gaugeStmt.Close()

	counterQuery := `
		INSERT INTO metrics (id ,type, delta) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (id) DO UPDATE SET delta = metrics.delta + EXCLUDED.delta
	`
	counterStmt, err := tx.PrepareContext(ctx, counterQuery)
	if err != nil {
		return fmt.Errorf("error preparing counter query: %v", err)
	}
	defer counterStmt.Close()

	for _, metric := range list {
		if metric.MType == metrics.Gauge {
			_, err = gaugeStmt.ExecContext(ctx, metric.ID, metric.MType, metric.Value)
			if err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("error updating metrics: %v", err)
			}
			continue
		}

		if metric.MType == metrics.Counter {
			_, err = counterStmt.ExecContext(ctx, metric.ID, metric.MType, metric.Delta)
			if err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("error updating metrics: %v", err)
			}
			continue
		}

		_ = tx.Rollback()
		return fmt.Errorf("unknown metric type: %v", metric.MType)
	}

	return nil
}
