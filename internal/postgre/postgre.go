package postgre

import (
	"context"
	"database/sql"
	"time"

	"github.com/Mr-Punder/go-alerting-service/internal/interfaces"
	"github.com/Mr-Punder/go-alerting-service/internal/metrics"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgreDB struct {
	db  *sql.DB
	log interfaces.Logger
}

func NewPostgreDB(dsn string, log interfaces.Logger) (*PostgreDB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Errorf("Error opening postgre database %s", err)
		return nil, err
	}

	Pdb := PostgreDB{
		db:  db,
		log: log,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = Pdb.InitTable(ctx)
	if err != nil {
		log.Errorf("Error initializing table metric %s", err)
		return &Pdb, err // have to do it that way to pass tests in iter10 where dsn is wrong but server has to work and be able to ping smth
	}

	log.Info("table meric Initialized")

	return &Pdb, nil
}

func (db *PostgreDB) InitTable(ctx context.Context) error {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_name = 'metric'
		)
	`

	var existsGauge bool
	err := db.db.QueryRowContext(ctx, query).Scan(&existsGauge)
	if err != nil {
		db.log.Errorf("Error searching table %s", err)
		return err
	}

	if existsGauge {
		db.log.Info("Found table")
		return nil
	}
	query = `
		CREATE TABLE metric (
			m_name VARCHAR(50) PRIMARY KEY,
			m_type VARCHAR(50) NOT NULL,
			delta BIGINT,
			value DOUBLE PRECISION
		)
	`

	_, err = db.db.ExecContext(ctx, query)
	if err != nil {
		db.log.Errorf("Error creating table %s", err)
		return err
	}

	db.log.Info("Table has not found and then created")

	_, err = db.db.ExecContext(ctx, "CREATE INDEX IF NOT EXIST m_name ON metric (m_name)")
	if err != nil {
		db.log.Errorf("Error creating index %s", err)
		return err
	}

	db.log.Info("Index created")
	return nil

}

func (db *PostgreDB) Close() error {
	err := db.db.Close()
	if err != nil {
		db.log.Errorf("Error cloding PostgreDB %s", err)
		return err
	}
	return nil
}

func (db *PostgreDB) Ping() error {
	return db.db.Ping()
}

func (db *PostgreDB) GetAll(ctx context.Context) map[string]metrics.Metrics {
	query := `SELECT m_name, m_type, delta, value FROM metric`

	rows, err := db.db.QueryContext(ctx, query)
	if err != nil {
		db.log.Errorf("Error selecting all metrics %s", err)
		return make(map[string]metrics.Metrics)

	}

	if rows.Err() != nil {
		db.log.Errorf("Error selecting all metrics %s", err)
		return make(map[string]metrics.Metrics)

	}

	defer rows.Close()
	metricMap := make(map[string]metrics.Metrics)

	for rows.Next() {
		var metric metrics.Metrics

		err := rows.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)
		if err != nil {
			db.log.Errorf("Error scaning metric %s", err)
			return make(map[string]metrics.Metrics)
		}

		metricMap[metric.ID] = metric
	}

	return metricMap
}
func (db *PostgreDB) Get(ctx context.Context, metric metrics.Metrics) (metrics.Metrics, bool) {

	query := `
		SELECT m_name, m_type, delta, value
		FROM metric
		WHERE m_name = $1
	`
	id := metric.ID

	var resMetric metrics.Metrics

	err := db.db.QueryRowContext(ctx, query, id).Scan(&resMetric.ID, &resMetric.MType, &resMetric.Delta, &resMetric.Value)
	if err == sql.ErrNoRows {
		db.log.Infof("Not found metric with id %s", id)
		return metrics.Metrics{}, false
	}

	if err != nil {
		db.log.Errorf("Error getting metric %s with error %s", id, err)
		return metrics.Metrics{}, false
	}

	return resMetric, true

}

func (db *PostgreDB) Delete(ctx context.Context, metric metrics.Metrics) error {
	quary := "DELETE FROM matric WHERE m_name = $1"

	id := metric.ID

	_, err := db.db.ExecContext(ctx, quary, id)
	if err != nil {
		db.log.Errorf("Error deleting meric %s  error: %s", metric.ID, err)
		return err
	}
	return nil
}

func (db *PostgreDB) Set(ctx context.Context, metric metrics.Metrics) error {
	quary := `
		INSERT INTO metric (m_name, m_type, delta, value)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (m_name) DO update
		SET  m_type = EXCLUDED.m_type, delta = metric.delta + EXCLUDED.delta, value = EXCLUDED.value
	`

	var delta int64
	var value float64
	if metric.Delta == nil {
		delta = 0
	} else {
		delta = *metric.Delta
	}
	if metric.Value == nil {
		value = 0.0
	} else {
		value = *metric.Value
	}

	_, err := db.db.ExecContext(ctx, quary, metric.ID, metric.MType, delta, value)
	if err != nil {
		db.log.Errorf("Error updating metric %s  error: %s", metric.ID, err)
		return err
	}

	return nil
}

func (db *PostgreDB) SetAll(ctx context.Context, metrics []metrics.Metrics) error {
	tx, err := db.db.Begin()
	if err != nil {
		db.log.Errorf("Error creating transaction %s", err)
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `
	INSERT INTO metric (m_name, m_type, delta, value)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (m_name) DO update
	SET  m_type = EXCLUDED.m_type, delta = metric.delta + EXCLUDED.delta, value = EXCLUDED.value
`)

	if err != nil {
		db.log.Errorf("Error preparing query", err)
		return err
	}
	defer stmt.Close()

	for _, metric := range metrics {
		var delta int64
		var value float64
		if metric.Delta == nil {
			delta = 0
		} else {
			delta = *metric.Delta
		}
		if metric.Value == nil {
			value = 0.0
		} else {
			value = *metric.Value
		}

		_, err := stmt.ExecContext(ctx, metric.ID, metric.MType, delta, value)
		if err != nil {
			db.log.Errorf("Error updating metric %s  error: %s in transaction", metric.ID, err)
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
