package sqldbexport

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
)

type SqlDbExporter struct {
	db  *sql.DB
	ctx context.Context
}

func (p *SqlDbExporter) Connect(driver string, dataSourceName string) error {
	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return err
	}
	p.db = db
	p.ctx = context.Background()
	return nil
}
