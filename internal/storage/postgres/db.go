package postgres

import (
	"context"
	store "github.com/Eqke/metric-collector/internal/storage"
	e "github.com/Eqke/metric-collector/pkg/error"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"time"
)

const (
	TYPE = "PostgreSQL database"

	queryCreateGauges   = `CREATE TABLE IF NOT EXISTS gauges(name text primary key, value double precision)`
	queryCreateCounters = `CREATE TABLE IF NOT EXISTS counters(name text primary key, value int)`

	queryGetGauge = `SELECT value FROM gauges WHERE name = $1`
	querySetGauge = `INSERT INTO gauges(name, value) VALUES($1, $2) ON CONFLICT(name) DO UPDATE SET value = $2`

	queryGetCounter = `SELECT value FROM counters WHERE name = $1`
	querySetCounter = `INSERT INTO counters(name, value) VALUES($1, $2) ON CONFLICT(name) DO UPDATE SET value = $2`
)

type PSQLStorage struct {
	ctx    context.Context
	conn   *pgx.Conn
	logger *zap.SugaredLogger
}

func New(ctx context.Context, logger *zap.SugaredLogger, conn string) (*PSQLStorage, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	db, err := pgx.Connect(ctx, conn)
	if err != nil {
		logger.Errorf("Database connection error: %v", err)
		return nil, err
	}

	err = db.Ping(ctxT)
	if err != nil {
		logger.Errorf("Database ping error: %v", err)
		return nil, err
	}

	_, err = db.Exec(ctxT, queryCreateGauges)
	if err != nil {
		logger.Errorf("Database exec error: %v", err)
		return nil, err
	}

	_, err = db.Exec(ctxT, queryCreateCounters)
	if err != nil {
		logger.Errorf("Database exec error: %v", err)
		return nil, err
	}

	return &PSQLStorage{
		ctx:    ctx,
		conn:   db,
		logger: logger,
	}, nil
}

func (P *PSQLStorage) SetValue(metricType, name, value string) error {
	switch metricType {
	case metric.TypeCounter.String():
		{
			_, err := P.conn.Exec(P.ctx, querySetCounter, name, value)
			if err != nil {
				P.logger.Errorf("Database exec error: %v", err)
				return e.WrapError(store.ErrPointSetValue, err)
			}
		}
	case metric.TypeGauge.String():
		{
			_, err := P.conn.Exec(P.ctx, querySetGauge, name, value)
			if err != nil {
				P.logger.Errorf("Database exec error: %v", err)
				return e.WrapError(store.ErrPointSetValue, err)
			}
		}
	default:
		{
			P.logger.Error(store.ErrPointSetValue, store.ErrIsUnknownType)
			return e.WrapError(store.ErrPointSetValue, store.ErrIsUnknownType)

		}
	}
	P.logger.Infof("metric was saved with type: %s, name: %s, value: %s",
		metricType, name, value)

	return nil
}

func (P *PSQLStorage) SetMetric(m metric.Metrics) error {
	switch m.MType {
	case metric.TypeCounter.String():
		{
			_, err := P.conn.Exec(P.ctx, querySetCounter, m.ID, m.Delta)
			if err != nil {
				P.logger.Errorf("Database exec error: %v", err)
				return e.WrapError(store.ErrPointSetValue, err)
			}
		}
	case metric.TypeGauge.String():
		{
			_, err := P.conn.Exec(P.ctx, querySetGauge, m.ID, m.Value)
			if err != nil {
				P.logger.Errorf("Database exec error: %v", err)
				return e.WrapError(store.ErrPointSetValue, err)
			}
		}
	default:
		{
			P.logger.Error(store.ErrPointSetValue, store.ErrIsUnknownType)
			return e.WrapError(store.ErrPointSetValue, store.ErrIsUnknownType)

		}
	}
	return nil
}

func (P *PSQLStorage) GetValue(metricType, name string) (string, error) {
	switch metricType {
	case metric.TypeCounter.String():
		{
			row := P.conn.QueryRow(P.ctx, queryGetCounter, name)
			var value string
			if err := row.Scan(&value); err != nil {
				P.logger.Errorf("Database scan error: %v", err)
				return "", e.WrapError(store.ErrPointGetValue, err)
			}
			return value, nil
		}
	case metric.TypeGauge.String():
		{
			row := P.conn.QueryRow(P.ctx, queryGetGauge, name)
			var value string
			if err := row.Scan(&value); err != nil {
				P.logger.Errorf("Database scan error: %v", err)
				return "", e.WrapError(store.ErrPointGetValue, err)
			}
			return value, nil
		}
	default:
		{
			P.logger.Error(store.ErrPointGetValue, store.ErrIsUnknownType)
			return "", e.WrapError(store.ErrPointGetValue, store.ErrIsUnknownType)
		}
	}
}

func (P *PSQLStorage) GetMetrics() ([]store.Metric, error) {
	metrics := make([]store.Metric, 0)

	rows, err := P.conn.Query(P.ctx, queryCreateCounters)
	if err != nil {
		P.logger.Errorf("Database query error: %v", err)
		return nil, e.WrapError(store.ErrPointGetMetrics, err)
	}

	for rows.Next() {
		var m store.Metric
		if err = rows.Scan(m.Name, m.Value); err != nil {
			P.logger.Errorf("Database scan error: %v", err)
			return nil, e.WrapError(store.ErrPointGetMetrics, err)
		}
		metrics = append(metrics, m)
	}

	rows, err = P.conn.Query(P.ctx, queryCreateGauges)
	if err != nil {
		P.logger.Errorf("Database query error: %v", err)
		return nil, e.WrapError(store.ErrPointGetMetrics, err)
	}

	for rows.Next() {
		var m store.Metric
		if err = rows.Scan(m.Name, m.Value); err != nil {
			P.logger.Errorf("Database scan error: %v", err)
			return nil, e.WrapError(store.ErrPointGetMetrics, err)
		}
		metrics = append(metrics, m)
	}

	return metrics, nil
}

func (P *PSQLStorage) GetMetric(m metric.Metrics) (metric.Metrics, error) {

	var met metric.Metrics
	switch m.MType {
	case metric.TypeCounter.String():
		{
			var val metric.Counter
			err := P.conn.QueryRow(P.ctx, queryGetCounter, m.ID).Scan(&val)
			if err != nil {
				P.logger.Errorf("Database scan error: %v", err)
				return met, e.WrapError(store.ErrPointGetMetric, err)
			}
		}
	case metric.TypeGauge.String():
		{
			var val metric.Gauge
			err := P.conn.QueryRow(P.ctx, queryGetGauge, m.ID).Scan(&val)
			if err != nil {
				P.logger.Errorf("Database scan error: %v", err)
				return met, e.WrapError(store.ErrPointGetMetric, err)
			}
		}
	default:
		{
			P.logger.Error(store.ErrPointGetMetric, store.ErrIsUnknownType)
			return met, e.WrapError(store.ErrPointGetMetric, store.ErrIsUnknownType)
		}
	}
	return met, nil
}

func (P *PSQLStorage) ToJSON() ([]byte, error) {
	return nil, nil
}

func (P *PSQLStorage) FromJSON(bytes []byte) error {
	return nil
}

func (P *PSQLStorage) ToFile(s string) error {
	return nil
}

func (P *PSQLStorage) FromFile(s string) error {
	return nil
}

func (P *PSQLStorage) Close() error {
	return P.conn.Close(context.Background())
}

func (P *PSQLStorage) Type() string {
	return TYPE
}
