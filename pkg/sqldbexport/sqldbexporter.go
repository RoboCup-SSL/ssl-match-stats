package sqldbexport

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"time"
)

type SqlDbExporter struct {
	db       *sql.DB
	ctx      context.Context
	autoRefs map[string]uuid.UUID
}

func (p *SqlDbExporter) Connect(driver string, dataSourceName string) error {
	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return err
	}
	p.db = db
	p.ctx = context.Background()
	p.autoRefs = map[string]uuid.UUID{}
	return nil
}

func convertTime(timestamp uint64) time.Time {
	startTimeSec := int64(timestamp / 1000000)
	startTimeNsec := (int64(timestamp) - (startTimeSec * 1000000)) * 1000
	return time.Unix(startTimeSec, startTimeNsec)
}

func convertDuration(duration int64) float64 {
	return time.Duration(duration * 1000).Seconds()
}
