package sqldbexport

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"sync"
	"time"
)

type SqlDbExporter struct {
	db       *sql.DB
	ctx      context.Context
	autoRefs map[string]uuid.UUID
	mutex    sync.Mutex
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

func (p *SqlDbExporter) AutoRefId(name string) (uuid.UUID, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	id, ok := p.autoRefs[name]
	return id, ok
}

func (p *SqlDbExporter) PutAutoRef(name string, id uuid.UUID) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.autoRefs[name] = id
}

func convertTime(timestamp uint64) time.Time {
	startTimeSec := int64(timestamp / 1000000)
	startTimeNsec := (int64(timestamp) - (startTimeSec * 1000000)) * 1000
	return time.Unix(startTimeSec, startTimeNsec)
}

func convertDuration(duration int64) float64 {
	return time.Duration(duration * 1000).Seconds()
}
