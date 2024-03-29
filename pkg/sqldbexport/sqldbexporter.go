package sqldbexport

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
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

func convertTime(timestamp uint64) time.Time {
	startTimeSec := int64(timestamp / 1000000)
	startTimeNsec := (int64(timestamp) - (startTimeSec * 1000000)) * 1000
	return time.Unix(startTimeSec, startTimeNsec)
}

func convertDuration(duration int64) time.Duration {
	return time.Duration(duration * 1000)
}
