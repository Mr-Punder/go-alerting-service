package postgre

import (
	"database/sql"

	"github.com/Mr-Punder/go-alerting-service/internal/interfaces"
	"github.com/Mr-Punder/go-alerting-service/internal/metrics"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgreDB struct {
	db  *sql.DB
	log interfaces.Logger
}

func NewPostgreDB(initstr string, log interfaces.Logger) (*PostgreDB, error) {
	db, err := sql.Open("pgx", initstr)
	if err != nil {
		log.Errorf("Error opening postgre database %s", err)
		return nil, err
	}

	return &PostgreDB{
		db:  db,
		log: log,
	}, nil
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

func (db *PostgreDB) GetAll() map[string]metrics.Metrics {
	return nil
}
func (db *PostgreDB) Get(metric metrics.Metrics) (metrics.Metrics, bool) {
	return metrics.Metrics{}, false
}

func (db *PostgreDB) Delete(metric metrics.Metrics) {}

func (db *PostgreDB) Set(metric metrics.Metrics) error {
	return nil
}
