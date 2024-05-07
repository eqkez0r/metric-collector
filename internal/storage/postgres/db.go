package postgres

import (
	"context"
	store "github.com/Eqke/metric-collector/internal/storage"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"time"
)

const (
	TYPE = "PostgreSQL database"

	queryCreateGauges   = `CREATE TABLE IF NOT EXISTS gauges(name text primary key, value double precision)`
	queryCreateCounters = `CREATE TABLE IF NOT EXISTS counters(name text primary key, value int)`

	queryGetGauge    = `SELECT value FROM gauges WHERE name = $1`
	queryGetAllGauge = `Select name, value FROM gauges`
	querySetGauge    = `INSERT INTO gauges(name, value) VALUES($1, $2) ON CONFLICT(name) DO UPDATE SET value = $2`

	queryGetCounter    = `SELECT value FROM counters WHERE name = $1`
	queryGetAllCounter = `SELECT name, value FROM counters`
	querySetCounter    = `INSERT INTO counters(name, value) VALUES($1, $2) ON CONFLICT(name) DO UPDATE SET value = $2`
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

func (p *PSQLStorage) SetValue(metricType, name, value string) error {
	switch metricType {
	case metric.TypeCounter.String():
		{
			_, err := p.conn.Exec(p.ctx, querySetCounter, name, value)
			if err != nil {
				p.logger.Errorf("Database exec error: %v", err)
				return err
			}
		}
	case metric.TypeGauge.String():
		{
			_, err := p.conn.Exec(p.ctx, querySetGauge, name, value)
			if err != nil {
				p.logger.Errorf("Database exec error: %v", err)
				return err
			}
		}
	default:
		{
			p.logger.Error(store.ErrPointSetValue, store.ErrIsUnknownType)
			return store.ErrIsUnknownType

		}
	}
	p.logger.Infof("metric was saved with type: %s, name: %s, value: %s",
		metricType, name, value)

	return nil
}

func (p *PSQLStorage) SetMetric(m metric.Metrics) error {
	switch m.MType {
	case metric.TypeCounter.String():
		{
			_, err := p.conn.Exec(p.ctx, querySetCounter, m.ID, *m.Delta)
			if err != nil {
				p.logger.Errorf("Database exec error: %v", err)
				return err
			}
		}
	case metric.TypeGauge.String():
		{
			_, err := p.conn.Exec(p.ctx, querySetGauge, m.ID, *m.Value)
			if err != nil {
				p.logger.Errorf("Database exec error: %v", err)
				return err
			}
		}
	default:
		{
			p.logger.Error(store.ErrPointSetValue, store.ErrIsUnknownType)
			return store.ErrIsUnknownType
		}
	}
	return nil
}

func (p *PSQLStorage) SetMetrics(m []metric.Metrics) error {
	for _, v := range m {
		err := p.SetMetric(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PSQLStorage) GetValue(metricType, name string) (string, error) {
	switch metricType {
	case metric.TypeCounter.String():
		{
			row := p.conn.QueryRow(p.ctx, queryGetCounter, name)
			var value string
			if err := row.Scan(&value); err != nil {
				p.logger.Errorf("Database scan error: %v", err)
				return "", err
			}
			return value, nil
		}
	case metric.TypeGauge.String():
		{
			row := p.conn.QueryRow(p.ctx, queryGetGauge, name)
			var value string
			if err := row.Scan(&value); err != nil {
				p.logger.Errorf("Database scan error: %v", err)
				return "", err
			}
			return value, nil
		}
	default:
		{
			p.logger.Error(store.ErrPointGetValue, store.ErrIsUnknownType)
			return "", store.ErrIsUnknownType
		}
	}
}

func (p *PSQLStorage) GetMetrics() ([]store.Metric, error) {
	metrics := make([]store.Metric, 0)

	rows, err := p.conn.Query(p.ctx, queryGetAllGauge)
	if err != nil {
		p.logger.Errorf("Database query error: %v", err)
		return nil, err
	}

	for rows.Next() {
		var m store.Metric
		if err = rows.Scan(m.Name, m.Value); err != nil {
			p.logger.Errorf("Database scan error: %v", err)
			return nil, err
		}
		metrics = append(metrics, m)
	}

	rows, err = p.conn.Query(p.ctx, queryGetAllCounter)
	if err != nil {
		p.logger.Errorf("Database query error: %v", err)
		return nil, err
	}

	for rows.Next() {
		var m store.Metric
		if err = rows.Scan(m.Name, m.Value); err != nil {
			p.logger.Errorf("Database scan error: %v", err)
			return nil, err
		}
		metrics = append(metrics, m)
	}

	return metrics, nil
}

func (p *PSQLStorage) GetMetric(m metric.Metrics) (metric.Metrics, error) {
	p.logger.Infof("get metric %v", m)
	switch m.MType {
	case metric.TypeCounter.String():
		{
			err := p.conn.QueryRow(p.ctx, queryGetCounter, m.ID).Scan(&m.Delta)
			if err != nil {
				p.logger.Errorf("Database scan error: %v", err)
				return m, err
			}
			p.logger.Infof("delta %v", *m.Delta)
		}
	case metric.TypeGauge.String():
		{
			err := p.conn.QueryRow(p.ctx, queryGetGauge, m.ID).Scan(&m.Value)
			if err != nil {
				p.logger.Errorf("Database scan error: %v", err)
				return m, err
			}
			p.logger.Infof("value %v", *m.Value)
		}
	default:
		{
			p.logger.Error(store.ErrPointGetMetric, store.ErrIsUnknownType)
			return m, store.ErrIsUnknownType
		}
	}
	p.logger.Infof("metric %v", m)
	return m, nil
}

func (p *PSQLStorage) ToJSON() ([]byte, error) {
	return nil, nil
}

func (p *PSQLStorage) FromJSON(bytes []byte) error {
	return nil
}

func (p *PSQLStorage) ToFile(s string) error {
	return nil
}

func (p *PSQLStorage) FromFile(s string) error {
	return nil
}

func (p *PSQLStorage) Close() error {
	return p.conn.Close(context.Background())
}

func (p *PSQLStorage) Type() string {
	return TYPE
}
